package main

import (
	"fmt"
	"net/http"
)

func handler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Server 1!")
}

func main() {
	http.HandleFunc("/", handler1)
	http.ListenAndServe(":8081", nil)
}
