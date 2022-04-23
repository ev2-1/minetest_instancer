package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

func (c *Controller) ContainerRun(image string, hostname string, mounts []mount.Mount) (id string, err error) {
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

	err = c.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *Controller) NetworkConnect(ID, network string) error {
	return c.cli.NetworkConnect(context.Background(), network, ID, nil)
}

// create minetest container
func (c *Controller) MinetestRun(mtSrv *minetestServer) (id string, err error) {
	return c.ContainerRun(minetestContainer, mtSrv.name, []mount.Mount{{
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
	}})
}

// delete container:
func (c *Controller) DeleteContainer(id string) error {
	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	return c.cli.ContainerRemove(context.Background(), id, removeOptions)
	
}
