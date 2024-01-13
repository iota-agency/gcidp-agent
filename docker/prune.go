package docker

import (
	"context"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type PruneCommand struct {
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
