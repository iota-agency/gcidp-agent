package docker

import (
	"context"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerNetwork "github.com/docker/docker/api/types/network"
	"log"
)

type RunCommand struct {
	cName         string
	config        *container.Config
	hostConfig    *container.HostConfig
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
		hostConfig: &container.HostConfig{},
	}
}

func (d *RunCommand) Run(ctx *pipeline.StageContext) error {
	d.config.Labels["gcidp.branch"] = ctx.Branch
	d.config.Labels["gcidp.repo"] = ctx.Repo
	resp, err := ctx.Client.ContainerCreate(context.Background(), d.config, d.hostConfig, d.networkConfig, nil, d.cName)
	if err != nil {
		return err
	}
	return ctx.Client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
}

func (d *RunCommand) Config(confs ...Conf) *RunCommand {
	for _, c := range confs {
		err := c.apply(d)
		if err != nil {
			log.Print(err)
		}
	}
	return d
}
