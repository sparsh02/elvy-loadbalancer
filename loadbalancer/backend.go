package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

type Backend struct {
	Address   string
	Alive     bool
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
}

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
	ActiveConnections() int
}

type simpleServer struct {
	addr           string
	alive          bool
	proxy          *httputil.ReverseProxy
	activeRequests int64
	limiter        *rate.Limiter
}

func (s *simpleServer) Address() string { return s.addr }

func (s *simpleServer) IsAlive() bool { return s.alive }

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	if s.limiter != nil && !s.limiter.Allow() {
		http.Error(rw, "Too Many Requests", http.StatusTooManyRequests)
		return
	}
	atomic.AddInt64(&s.activeRequests, 1)
	defer atomic.AddInt64(&s.activeRequests, -1)

	s.proxy.ServeHTTP(rw, req)
}

func (s *simpleServer) ActiveConnections() int {
	return int(atomic.LoadInt64(&s.activeRequests))
}

func (s *simpleServer) checkHealth() {
	for {
		resp, err := http.Get(s.addr)
		if err != nil || resp.StatusCode >= 400 {
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

	var limiter *rate.Limiter
	if backend.RateLimit.Enabled {

		limiter = rate.NewLimiter(rate.Limit(float64(backend.RateLimit.RequestsPerMinute)/60.0), backend.RateLimit.RequestsPerMinute)
	}

	server := &simpleServer{
		addr:    addr,
		alive:   alive,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
		limiter: limiter,
	}

	go server.checkHealth()
	return server
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
