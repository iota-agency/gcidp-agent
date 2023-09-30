package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerNetwork "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type RunCommand struct {
	cName         string
	config        *container.Config
	networkConfig *dockerNetwork.NetworkingConfig
}

func Run(cName, image string) *RunCommand {
	return &RunCommand{
		cName: cName,
		config: &container.Config{
			Image: image,
			Tty:   false,
			Labels: map[string]string{
				"gcidp.enable": "true",
			},
			Env: []string{},
		},
	}
}

func (d *RunCommand) Run(cli *client.Client) error {
	resp, err := cli.ContainerCreate(context.Background(), d.config, nil, d.networkConfig, nil, d.cName)
	if err != nil {
		return err
	}
	return cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
}

func (d *RunCommand) Config(confs ...Conf) *RunCommand {
	for _, c := range confs {
		c.apply(d)
	}
	return d
}
