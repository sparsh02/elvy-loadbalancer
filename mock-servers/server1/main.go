package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"strconv"
)

// Global variables to track connections
var (
	activeConnections int
	mu                 sync.Mutex
)

// Handler for simulating request processing with delay
func handler(w http.ResponseWriter, r *http.Request) {
	// Simulate varying processing times
	delay := r.URL.Query().Get("delay")
	if delay == "" {
		delay = "5000" // default 5 seconds
	}

	delayDuration, err := strconv.Atoi(delay)
	if err != nil {
		delayDuration = 5000 // default 5 seconds
	}

	mu.Lock()
	activeConnections++
	mu.Unlock()

	defer func() {
		mu.Lock()
		activeConnections--
		mu.Unlock()
	}()

	time.Sleep(time.Duration(delayDuration) * time.Millisecond)
	fmt.Fprintf(w, "Hello from Server! Active connections: %d", getActiveConnections())
}

// Endpoint to get the number of active connections
func connectionsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Active connections: %d", getActiveConnections())
}

// Function to get the number of active connections
func getActiveConnections() int {
	mu.Lock()
	defer mu.Unlock()
	return activeConnections
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/connections", connectionsHandler)
	fmt.Println("Starting server on port 8081...")
	http.ListenAndServe(":8081", nil)
}
