package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type RmContainerCommand struct {
	name   string
	silent bool
}

type RmImageCommand struct {
	name   string
	silent bool
}

func RmContainer(name string, silent bool) *RmContainerCommand {
	return &RmContainerCommand{name: name, silent: silent}
}

func RmImage(name string, silent bool) *RmImageCommand {
	return &RmImageCommand{name: name, silent: silent}
}

func (d *RmContainerCommand) Run(cli *client.Client) error {
	err := cli.ContainerRemove(context.Background(), d.name, types.ContainerRemoveOptions{Force: true})
	if !d.silent && err != nil {
		return err
	}
	return nil
}

func (d *RmImageCommand) Run(cli *client.Client) error {
	_, err := cli.ImageRemove(context.Background(), d.name, types.ImageRemoveOptions{Force: true})
	if !d.silent && err != nil {
		return err
	}
	return nil
}
