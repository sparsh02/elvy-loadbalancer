package main

import (
	"elvy-loadbalancer/config"
	"fmt"
	"log"
	"net/http"

	"elvy-loadbalancer/loadbalancer"
)

func main() {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	var servers []loadbalancer.Server

	for _, backend := range config.Backends {
		servers = append(servers, loadbalancer.CreateServer(backend))
	}

	lb := loadbalancer.NewLoadBalancerServers(config.Port, servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("Received request: %s %s\n", req.Method, req.URL)
		lb.ServeProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)
	fmt.Printf("forwarding request to address %q\n", lb.Port())

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}

}
