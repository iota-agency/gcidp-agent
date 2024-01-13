package docker

import (
	"context"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/docker/docker/api/types"
)

type RmImageCommand struct {
	name   string
	silent bool
}

func RmImage(name string, silent bool) *RmImageCommand {
	return &RmImageCommand{name: name, silent: silent}
}

func (d *RmImageCommand) Run(ctx *pipeline.StageContext) error {
	_, err := ctx.Client.ImageRemove(context.Background(), d.name, types.ImageRemoveOptions{Force: true})
	if err != nil {
		if d.silent {
			ctx.Logger.Warn(err.Error())
		} else {
			return err
		}
	}
	return nil
}
