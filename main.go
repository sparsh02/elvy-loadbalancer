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

	// servers := []loadbalancer.Server{
	// 	loadbalancer.CreateServer("https://google.com"),
	// 	loadbalancer.CreateServer("https://yahoo.com"),
	// 	loadbalancer.CreateServer("https://bing.com")}

	// [] means, elements of the array are of type Backend struct from loadbalancer/backend.go file.
	// backends := []*loadbalancer.Backend{
	// 	{Address: "127.0.0.1:8081", Alive: true},
	// 	{Address: "127.0.0.1:8082", Alive: true},
	// }

	// lb := loadbalancer.NewLoadBalancer("8000",backends)  // This is the original code
	// handler := loadbalancer.NewHandler(lb)
	// lb := loadbalancer.NewLoadBalancerServers("8000", servers) // This is the modified code

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

	// listener, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	log.Fatal("Error starting TCP server:", err)
	// }
	// defer listener.Close()

	// log.Println("Load balancer started on :8080")

	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		log.Println("Error accepting connection:", err)
	// 		continue
	// 	}
	// 	go handler.HandleConnection(conn)
	// }
}
