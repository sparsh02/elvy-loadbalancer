package loadbalancer

import (
	"fmt"
	"net/http"
	// "sync/atomic"
)

// type LoadBalancer struct {
// 	port     string
// 	backends []*Backend
// 	current  uint64
// }

type LoadBalancerServers struct {
	port    string
	servers []Server
	current int
	algo    string
}

func NewLoadBalancerServers(port string, servers []Server, algo string) *LoadBalancerServers {
	return &LoadBalancerServers{
		port:    port,
		servers: servers,
		current: 0,
		algo:    algo,
	}
}

// func NewLoadBalancer(port string, backends []*Backend) *LoadBalancer {
// 	return &LoadBalancer{
// 		port:     port,
// 		backends: backends,
// 		current:  0,
// 	}
// }

// Round Robin algorithm

// func (lb *LoadBalancer) GetNextBackend() *Backend {
// 	next := atomic.AddUint64(&lb.current, 1)
// 	return lb.backends[next%uint64(len(lb.backends))]
// }

// Getting next server based on the algorithm

func (lb *LoadBalancerServers) GetNextServer() Server {

	fmt.Printf("Using %s algorithm\n", lb.algo)

	switch lb.algo {
	case "round_robin":
		return lb.roundRobin()
	case "least_conn":
		return lb.leastConn()
	case "ip_hash":
		return lb.iphash()
	default:
		return lb.roundRobin()
	}

	// server := lb.servers[lb.current%len(lb.servers)]
	// fmt.Println("Current server is ", server)

	// for !server.IsAlive() {
	// 	lb.current++
	// 	server = lb.servers[lb.current%len(lb.servers)]
	// }
	// lb.current++
	// return server
}

func (lb *LoadBalancerServers) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	server := lb.GetNextServer()
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
	// Implementing least connection algorithm

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
func (lb *LoadBalancerServers) iphash() Server {
	// Implementing ip hash algorithm

	// Get the client IP address
	// ip:= req

	return lb.roundRobin()
}

