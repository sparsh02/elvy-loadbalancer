package loadbalancer

import (
    "io"
    "log"
    "net"
)

type Handler struct {
    LoadBalancer *LoadBalancer
}

func NewHandler(lb *LoadBalancer) *Handler {
    return &Handler{
        LoadBalancer: lb,
    }
}

func (h *Handler) HandleConnection(conn net.Conn) {
    backend := h.LoadBalancer.GetNextBackend()
    if backend.Alive {
        backendConn, err := net.Dial("tcp", backend.Address)
        if err != nil {
            log.Println("Error connecting to backend:", err)
            conn.Close()
            return
        }
        defer backendConn.Close()

        go io.Copy(backendConn, conn)
        io.Copy(conn, backendConn)
    } else {
        conn.Close()
    }
}
