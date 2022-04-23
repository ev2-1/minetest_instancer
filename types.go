package main

import (
	//"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Controller struct {
	cli *client.Client
}

type minetestServer struct {
	name   string
	world  string
	game   string
	config string
	id     string
	net    string
}

func NewController() (c *Controller, err error) {
	c = new(Controller)

	c.cli, err = client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}
	return c, nil
}
