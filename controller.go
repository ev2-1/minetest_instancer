package main

import (
	"fmt"
	"os"
)

var pwd string

var srvNetwork = "ev21_default"
var minetestContainer = "minetest_5.4.1"
var worldPath = "/../worlds/"
var configPath = "/../config/"
var gamePath = "/../games/"

var c *Controller
var servers map[string]*minetestServer

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: minetest <port>")
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		worldPath = pwd + worldPath
		configPath = pwd + configPath
		gamePath = pwd + gamePath
	}

	c, err = NewController()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = telnetServer(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
}

func handle(prefix string, tw *TelnetWriter, s []string) string {
	switch s[0] {
	case "srv_create":
		if len(s) != 4 {
			return "err, usage: srv_create <name> <world> <game>"
		}

		srv := &minetestServer{
			name:   s[1],
			world:  s[2],
			game:   s[3],
			config: "conf",
		}

		// create minetest server
		id, err := c.MinetestRun(srv)
		if err != nil {
			return "err, " + err.Error()
		}

		srv.id = id
		servers[srv.id] = srv

		return "OK, " + id

	case "srv_connect":
		if len(s) != 3 {
			return "err, usage: srv_connect <container> <network>"
		}

		if _, exist := servers[s[1]]; !exist {
			return "err, srv dosnt exist"
		}

		if s[2] == "default" {
			s[2] = srvNetwork
		}

		err := c.NetworkConnect(s[1], s[2])
		if err != nil {
			return "err, " + err.Error()
		}

		servers[s[1]].net = s[2]

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

	default:
		return "err"
	}
}
