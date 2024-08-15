package handlers

import (
	"elvy-loadbalancer/loadbalancer"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func AddServerHandler(lb *loadbalancer.LoadBalancerServers) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var backend loadbalancer.Backend
		err := json.NewDecoder(req.Body).Decode(&backend)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		server := loadbalancer.CreateServer(backend)
		lb.AddServer(server)
		rw.WriteHeader(http.StatusCreated)
	}
}

func RemoveServerHandler(lb *loadbalancer.LoadBalancerServers) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		address, err := ioutil.ReadAll(req.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		lb.RemoveServer(string(address))
		rw.WriteHeader(http.StatusNoContent)
	}
}
