package main

import (
	"fmt"
	"net/http"

	"strings"
)

func listen_http(addr string) {
	http.HandleFunc("/ready", readyFunc)

	http.HandleFunc("/kick", kickFunc)

	go func() {
		http.ListenAndServe(addr, nil)
	}()
}

func kickFunc(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query()
	if param["name"] == nil || param["reason"] == nil {
		w.WriteHeader(400)
		return
	}

	fmt.Println("[HTTP] kicking player", param["name"][0], "because", param["reason"][0])
}

func readyFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK!")

	ipS := strings.Split(r.RemoteAddr, ":")
	ip := strings.Join(ipS[:len(ipS)-1], ":")

	fmt.Println("[HTTP] /ready by", ip)

	// mark server as online:
	srv := ip2ms[ip]
	if srv != nil {
		srv.ready = true
	} else {
		fmt.Println("[HTTP] info: no server with ip", ip, "registerd")
	}

	if waits[ip] != nil {
		select {
		case <-waits[ip]:
		default:
			close(waits[ip])
		}
	} else {
		waits[ip] = make(chan struct{})
		close(waits[ip])
	}
}
