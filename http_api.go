package main

import (
	"fmt"
	"net/http"

	"strings"
)

func listen_http(addr string) {
	http.HandleFunc("/ready", readyFunc)

	http.HandleFunc("/send", sendFunc)

	http.HandleFunc("/pkexec", pkexecFunc) // because sudo is bad /s
	
	http.HandleFunc("/kick", kickFunc)

	go func() {
		http.ListenAndServe(addr, nil)
	}()
}

func pkexecFunc(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query()
	if param["name"] == nil || param["cmd"] == nil {
		w.WriteHeader(400)
		return
	}

	fmt.Println("[HTTP] pkexec player", param["name"][0], "cmd", param["cmd"][0])

	broadcast("pkexec "+param["name"][0]+" "+param["cmd"][0], "multiserver")	
}

func sendFunc(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query()
	if param["name"] == nil || param["dest"] == nil {
		w.WriteHeader(400)
		return
	}

	fmt.Println("[HTTP] sending player", param["name"][0], "to", param["dest"][0])

	broadcast("send "+param["name"][0]+" "+param["dest"][0], "multiserver")
}

func kickFunc(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query()
	if param["name"] == nil || param["reason"] == nil {
		w.WriteHeader(400)
		return
	}

	fmt.Println("[HTTP] kicking player", param["name"][0], "because", strings.ReplaceAll(param["reason"][0], "\n", "\\n"))

	broadcast("kick "+param["name"][0]+" "+param["reason"][0], "multiserver")
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
