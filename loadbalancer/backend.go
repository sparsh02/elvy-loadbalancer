package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
	"time"
)

// Pointers needs to be passed to the function
type Backend struct {
	Address string
	Alive   bool
}

// Interface already holds the reference, so no need to pass pointer to the function
type Server interface {
	// Address returns the address with which to access the server
	Address() string

	// IsAlive returns true if the server is alive and able to serve requests
	IsAlive() bool

	// Serve uses this server to process the request
	Serve(rw http.ResponseWriter, req *http.Request)

	ActiveConnections() int
}

type simpleServer struct {
	addr           string
	alive          bool
	proxy          *httputil.ReverseProxy
	activeRequests int64
}

func (s *simpleServer) Address() string { return s.addr }

func (s *simpleServer) IsAlive() bool { return s.alive }

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Proxying request to %s\n", s.addr)
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)
	s.proxy.ServeHTTP(rw, req)
}

func (s *simpleServer) ActiveConnections() int {
	return int(atomic.LoadInt64(&s.activeRequests))
}

func (s *simpleServer) checkHealth() {
	for {
		resp, errs := http.Get(s.addr)
		if errs != nil || resp.StatusCode >= 400 {
			s.alive = false
		} else {
			s.alive = true
		}
		time.Sleep(30 * time.Second)
	}
}
func CreateServer(backend Backend) Server {
	addr := backend.Address
	alive := backend.Alive
	serverUrl, err := url.Parse(addr)

	handleErr(err)
	server := &simpleServer{
		addr:  addr,
		alive: alive,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}

	// starting healcheck in separate goroutine
	// checkHealth will update the alive status of the server
	go server.checkHealth()
	return server
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
