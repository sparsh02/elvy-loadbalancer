package main

import (
	"context"
	"elvy-loadbalancer/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

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

	server := &http.Server{
		Addr: ":" + config.Port,
		// Handler: nil,
	}

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("Received request: %s %s\n", req.Method, req.URL)
		lb.ServeProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)

	// this function will start the server in a goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	fmt.Printf("forwarding request to address %q\n", lb.Port())

	// creating channelm to listen for os signals

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	// Block until a signal is received
	<-c
	fmt.Println("Shutting down server...")
	error := server.Shutdown(context.Background())
	if error != nil {
		log.Fatalf("Error shutting down server: %v", error)
	}
	fmt.Println("Server gracefully stopped")

}
