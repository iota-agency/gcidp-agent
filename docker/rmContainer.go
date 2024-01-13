package docker

import (
	"context"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/docker/docker/api/types"
)

type RmContainerCommand struct {
	name   string
	silent bool
}

func RmContainer(name string, silent bool) *RmContainerCommand {
	return &RmContainerCommand{name: name, silent: silent}
}

func (d *RmContainerCommand) Run(ctx *pipeline.StageContext) error {
	err := ctx.Client.ContainerRemove(context.Background(), d.name, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		if d.silent {
			ctx.Logger.Warn(err.Error())
		} else {
			return err
		}
	}
	return nil
}
