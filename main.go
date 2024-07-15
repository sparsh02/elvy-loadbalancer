package main

import (
    "log"
    "net"

    "elvy-loadbalancer/loadbalancer"
)

func main() {
    backends := []*loadbalancer.Backend{
        {Address: "127.0.0.1:8081", Alive: true},
        {Address: "127.0.0.1:8082", Alive: true},
    }

    lb := loadbalancer.NewLoadBalancer(backends)
    handler := loadbalancer.NewHandler(lb)

    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal("Error starting TCP server:", err)
    }
    defer listener.Close()

    log.Println("Load balancer started on :8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Error accepting connection:", err)
            continue
        }
        go handler.HandleConnection(conn)
    }
}
