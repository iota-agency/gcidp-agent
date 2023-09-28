package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"log"
)

type RmContainerCommand struct {
	name   string
	silent bool
}

type RmImageCommand struct {
	name   string
	silent bool
}

type PruneCommand struct {
}

func RmContainer(name string, silent bool) *RmContainerCommand {
	return &RmContainerCommand{name: name, silent: silent}
}

func RmImage(name string, silent bool) *RmImageCommand {
	return &RmImageCommand{name: name, silent: silent}
}

func Prune() *PruneCommand {
	return &PruneCommand{}
}

func (d *PruneCommand) Run(cli *client.Client) error {
	if _, err := cli.ImagesPrune(context.Background(), filters.Args{}); err != nil {
		return err
	}
	if _, err := cli.ContainersPrune(context.Background(), filters.Args{}); err != nil {
		return err
	}
	if _, err := cli.VolumesPrune(context.Background(), filters.Args{}); err != nil {
		return err
	}
	if _, err := cli.BuildCachePrune(context.Background(), types.BuildCachePruneOptions{}); err != nil {
		return err
	}
	if _, err := cli.NetworksPrune(context.Background(), filters.Args{}); err != nil {
		return err
	}
	return nil
}

func (d *RmContainerCommand) Run(cli *client.Client) error {
	err := cli.ContainerRemove(context.Background(), d.name, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		if d.silent {
			log.Println(err)
		} else {
			return err
		}
	}
	return nil
}

func (d *RmImageCommand) Run(cli *client.Client) error {
	_, err := cli.ImageRemove(context.Background(), d.name, types.ImageRemoveOptions{Force: true})
	if err != nil {
		if d.silent {
			log.Println(err)
		} else {
			return err
		}
	}
	return nil
}
