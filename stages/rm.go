package stages

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerRm struct {
	name  string
	image bool
}

func NewDockerRm() *DockerRm {
	return &DockerRm{}
}

func (d *DockerRm) Container(name string) *DockerRm {
	d.name = name
	return d
}

func (d *DockerRm) Image(name string) *DockerRm {
	d.name = name
	d.image = true
	return d
}

func (d *DockerRm) Run(cli *client.Client) error {
	_, err := cli.ImageRemove(context.Background(), d.name, types.ImageRemoveOptions{Force: true})
	if err != nil {
		return err
	}
	return nil
}
