package docker

import (
	"context"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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

func (d *PruneCommand) Run(ctx *pipeline.StageContext) error {
	if _, err := ctx.Client.ImagesPrune(context.Background(), filters.Args{}); err != nil {
		return err
	}
	if _, err := ctx.Client.ContainersPrune(context.Background(), filters.Args{}); err != nil {
		return err
	}
	if _, err := ctx.Client.VolumesPrune(context.Background(), filters.Args{}); err != nil {
		return err
	}
	if _, err := ctx.Client.BuildCachePrune(context.Background(), types.BuildCachePruneOptions{}); err != nil {
		return err
	}
	if _, err := ctx.Client.NetworksPrune(context.Background(), filters.Args{}); err != nil {
		return err
	}
	return nil
}

func (d *RmContainerCommand) Run(ctx *pipeline.StageContext) error {
	err := ctx.Client.ContainerRemove(context.Background(), d.name, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		if d.silent {
			log.Println(err)
		} else {
			return err
		}
	}
	return nil
}

func (d *RmImageCommand) Run(ctx *pipeline.StageContext) error {
	_, err := ctx.Client.ImageRemove(context.Background(), d.name, types.ImageRemoveOptions{Force: true})
	if err != nil {
		if d.silent {
			log.Println(err)
		} else {
			return err
		}
	}
	return nil
}
