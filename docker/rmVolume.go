package docker

import (
	"context"
	"github.com/apollo-studios/gcidp-agent/pipeline"
)

type RmVolumeCommand struct {
	Volume string
}

func RmVolume(volume string) *RmVolumeCommand {
	return &RmVolumeCommand{
		Volume: volume,
	}
}

func (d *RmVolumeCommand) Run(ctx *pipeline.StageContext) error {
	_ = ctx.Client.VolumeRemove(context.Background(), d.Volume, true)
	return nil
}
