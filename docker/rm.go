package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type RmCommand struct {
	name   string
	image  bool
	silent bool
}

func Rm(silent bool) *RmCommand {
	return &RmCommand{silent: silent}
}

func (d *RmCommand) Container(name string) *RmCommand {
	d.name = name
	return d
}

func (d *RmCommand) Image(name string) *RmCommand {
	d.name = name
	d.image = true
	return d
}

func (d *RmCommand) Run(cli *client.Client) error {
	_, err := cli.ImageRemove(context.Background(), d.name, types.ImageRemoveOptions{Force: true})
	if !d.silent && err != nil {
		return err
	}
	return nil
}
