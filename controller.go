package main

import (
	"fmt"
	"os"
	"strings"

	"net/http"
)

var pwd string

var c *Controller
var ip2ms map[string]*minetestServer
var id2ms map[string]*minetestServer
var name2ms map[string]*minetestServer
var waits map[string]chan struct{}

func main() {
	ip2ms = make(map[string]*minetestServer)
	id2ms = make(map[string]*minetestServer)
	name2ms = make(map[string]*minetestServer)
	waits = make(map[string]chan struct{})

	// read config
	err := ReadConfig()
	if err != nil {
		fmt.Println("[MAIN] couldn't read config", err)
		os.Exit(1)
	}

	c, err = NewController()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// web server
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
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
	})

	go func() {
		http.ListenAndServe("[::]:80", nil)
	}()

	// start telnet server
	err = telnetServer("8888")
	if err != nil {
		fmt.Println(err)
	}
}

func handle(prefix string, tw *TelnetWriter, s []string) string {
	switch s[0] {
	case "get_ip":
		if len(s) != 2 {
			return "err, usage: get_ip <container>"
		}

		ip, err := c.GetIp(s[1])
		if err != nil {
			return "err, " + err.Error()
		}

		return "OK, " + ip

	case "srv_create":
		if len(s) != 4 && len(s) != 5 {
			return "err, usage: srv_create <name> <world> <game> [net]"
		}

		if len(s) == 4 {
			s[4] = ""
		}

		// create world dir if not exist:
		os.Mkdir(worldMount+s[2], 0755)

		srv := &minetestServer{
			name:   s[1],
			world:  s[2],
			game:   s[3],
			config: "conf",
			net:    s[4],
		}

		name2ms[srv.name] = srv

		id, err := c.MinetestCreate(srv)
		if err != nil {
			return "err, (create), " + err.Error()
		}
		srv.id = id
		id2ms[id[0:11]] = srv

		// start server
		err = c.ContainerStart(srv.id)
		if err != nil {
			return "err, (start), " + err.Error()
		}

		// wait for get ip
		var ip string
		for len(ip) == 0 {
			ip, err = c.GetIp(id[0:11])
			if err != nil {
				return "err, (ip), " + err.Error()
			}
			fmt.Println("got ip", ip, "for host", id[0:11])
		}
		srv.ip = ip
		if waits[ip] == nil {
			waits[ip] = make(chan struct{})
		}

		// wait for server online
		<-waits[ip]

		return "OK, " + id

	case "srv_connect":
		if len(s) != 3 {
			return "err, usage: srv_connect <container> <network>"
		}

		if s[2] == "default" {
			s[2] = srvNetwork
		}

		if id2ms[s[1]] == nil {
			return "err, srv not registerd"
		}
		id2ms[s[1]].net = s[2]

		err := c.NetworkConnect(s[1], s[2])
		if err != nil {
			return "err, " + err.Error()
		}

		return "OK"

	case "srv_delete":
		if len(s) != 2 {
			return "err, usage: srv_delete <container>"
		}

		err := c.DeleteContainer(s[1])
		if err != nil {
			return "err, " + err.Error()
		}

		return "OK"

	case "srv_state":
		if len(s) != 2 {
			return "err, usage: srv_exists <container>"
		}

		resp, err := c.Inspect(s[1])
		if err != nil {
			return "OK, false"
		} else {
			var str string
			if resp.State.Running {
				str = "Running"
			} else {
				str = "Stopped"
			}
			return "OK, " + resp.State.Status + ", " + str
		}

	case "servers":
		var resp string

		for id, srv := range id2ms {
			var s string
			if srv.ready {
				s = "true"
			} else {
				s = "false"
			}

			resp += id + " - " + srv.name + "|" + srv.world + "|" + srv.game + "|" + srv.config + "|" + srv.id + "|" + srv.net + "|" + srv.ip + "|" + s + "; "
		}

		return "OK, " + resp

	case "srv_start":
		if len(s) != 2 {
			return "err, usage: srv_start <container>"
		}

		err := c.ContainerStart(s[1])
		if err != nil {
			return "err, " + err.Error()
		}

		return "OK"

	default:
		return "err"
	}
}
