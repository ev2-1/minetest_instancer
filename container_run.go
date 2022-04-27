package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

func (c *Controller) ContainerCreate(image string, hostname string, mounts []mount.Mount, network string) (id string, err error) {
	hostConfig := container.HostConfig{}

	hostConfig.Mounts = mounts

	resp, err := c.cli.ContainerCreate(context.Background(), &container.Config{
		Tty:      true,
		Image:    image,
		Hostname: hostname,
	}, &hostConfig, nil, nil, hostname)

	if err != nil {
		return "", err
	}

	if len(network) != 0 {
		err = c.NetworkConnect(resp.ID, network)
		if err != nil {
			return "", err
		}
	}

	return resp.ID, nil
}

// connect a docker container to network
func (c *Controller) NetworkConnect(ID, network string) error {
	return c.cli.NetworkConnect(context.Background(), network, ID, nil)
}

// start a docker container
func (c *Controller) ContainerStart(id string) error {
	return c.cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}

// create minetest container
func (c *Controller) MinetestCreate(mtSrv *minetestServer) (id string, err error) {
	return c.ContainerCreate(minetestContainer, mtSrv.name, []mount.Mount{{
		Type:   mount.TypeBind,
		Source: worldPath + mtSrv.world,
		Target: "/minetest/worlds/world",
	}, {
		Type:   mount.TypeBind,
		Source: gamePath + mtSrv.game,
		Target: "/minetest/games/game",
	}, {
		Type:   mount.TypeBind,
		Source: configPath + mtSrv.config,
		Target: "/config/config.yml",
	}}, mtSrv.net)
}

// delete container:
func (c *Controller) DeleteContainer(id string) error {
	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	return c.cli.ContainerRemove(context.Background(), id, removeOptions)

}

func (c *Controller) Inspect(id string) (types.ContainerJSON, error) {
	return c.cli.ContainerInspect(context.Background(), id)
}

// get ip of container:
func (c *Controller) GetIp(id string) (string, error) {
	resp, err := c.Inspect(id)
	if err != nil {
		return "", err
	}

	nets := resp.NetworkSettings.Networks

	if nets[srvNetwork] == nil { // id not in server network
		return "", nil
	}

	return nets[srvNetwork].IPAddress, nil
}
