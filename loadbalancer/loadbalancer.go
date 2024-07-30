package loadbalancer

import (
	"fmt"
	"hash/fnv"
	"net"
	"net/http"
	"sync"
)

type LoadBalancerServers struct {
	port        string
	servers     []Server
	current     int
	algo        string
	sticky      bool
	stickyTable map[string]Server
	stickyLock  sync.Mutex
}

func NewLoadBalancerServers(port string, servers []Server, algo string, stickySession bool) *LoadBalancerServers {
	return &LoadBalancerServers{
		port:        port,
		servers:     servers,
		current:     0,
		algo:        algo,
		sticky:      stickySession,
		stickyTable: make(map[string]Server),
	}
}

// Getting next server based on the algorithm

func (lb *LoadBalancerServers) GetNextServer(req *http.Request) Server {

	fmt.Printf("Using %s algorithm\n", lb.algo)
	userIp := getClientIP(req)

	if lb.sticky {
		lb.stickyLock.Lock()
		if server, exists := lb.stickyTable[userIp]; exists && server.IsAlive() {
			lb.stickyLock.Unlock()
			return server
		}
		lb.stickyLock.Unlock()
	}

	var server Server

	switch lb.algo {
	case "round_robin":
		server = lb.roundRobin()
	case "least_conn":
		server = lb.leastConn()
	case "ip_hash":
		server = lb.iphash(userIp)
	default:
		server = lb.roundRobin()
	}

	if lb.sticky {
		lb.stickyLock.Lock()
		lb.stickyTable[userIp] = server
		lb.stickyLock.Unlock()
	}
	return server
}

func (lb *LoadBalancerServers) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	server := lb.GetNextServer(req)
	server.Serve(rw, req)
}

func (lb *LoadBalancerServers) Port() string {
	return lb.port
}

// Algorithm implementations

func (lb *LoadBalancerServers) roundRobin() Server {
	server := lb.servers[lb.current%len(lb.servers)]
	fmt.Println("Current server is ", server)

	for !server.IsAlive() {
		lb.current++
		server = lb.servers[lb.current%len(lb.servers)]
	}
	lb.current++
	return server
}

func (lb *LoadBalancerServers) leastConn() Server {

	var selected Server

	minConnections := int(^uint(0) >> 1) // Max int value

	for _, server := range lb.servers {
		if server.IsAlive() {
			connections := server.ActiveConnections()
			if connections < minConnections {
				minConnections = connections
				selected = server
			}
		}
	}
	return selected
}
func (lb *LoadBalancerServers) iphash(ip string) Server {

	hash := fnv.New32a()
	hash.Write([]byte(ip))

	index := hash.Sum32() % uint32(len(lb.servers))
	return lb.servers[index]

}

// Utilities

func getClientIP(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(req.RemoteAddr)
	}
	userIp := net.ParseIP(ip)
	return userIp.String()
}
