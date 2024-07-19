package loadbalancer

import (
	"net/http"
	"sync/atomic"
)

type LoadBalancer struct {
	port     string
	backends []*Backend
	current  uint64
}

type LoadBalancerServers struct {
	port    string
	servers  []Server
	current int
}
func NewLoadBalancerServers(port string, servers []Server) *LoadBalancerServers {
	return &LoadBalancerServers{
		port:     port,
		servers: servers,
		current:  0,
	}
}

func NewLoadBalancer(port string, backends []*Backend) *LoadBalancer {
	return &LoadBalancer{
		port:     port,
		backends: backends,
		current:  0,
	}
}

// Round Robin algorithm

func (lb *LoadBalancer) GetNextBackend() *Backend {
	next := atomic.AddUint64(&lb.current, 1)
	return lb.backends[next%uint64(len(lb.backends))]
}



// Round R

func (lb *LoadBalancerServers) GetNextServer() Server {
	server:= lb.servers[lb.current%len(lb.servers)]
	
	for !server.IsAlive(){
		lb.current++
		server = lb.servers[lb.current%len(lb.servers)]
	}
	lb.current++
	return server
}

func (lb *LoadBalancerServers) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	server := lb.GetNextServer()
	server.Serve(rw, req)
}

func (lb *LoadBalancerServers) Port() string{
	return lb.port
}