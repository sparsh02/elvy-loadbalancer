package loadbalancer

import "sync/atomic"

type LoadBalancer struct {
    backends []*Backend
    current  uint64
}

func NewLoadBalancer(backends []*Backend) *LoadBalancer {
    return &LoadBalancer{
        backends: backends,
        current:  0,
    }
}

func (lb *LoadBalancer) GetNextBackend() *Backend {
    next := atomic.AddUint64(&lb.current, 1)
    return lb.backends[next%uint64(len(lb.backends))]
}
