package docker

import (
	"context"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerNetwork "github.com/docker/docker/api/types/network"
)

type RunCommand struct {
	cName         string
	config        *container.Config
	hostConfig    *container.HostConfig
	networkConfig *dockerNetwork.NetworkingConfig
	confs         []Conf
}

func Run(cName, image string) *RunCommand {
	return &RunCommand{
		cName: cName,
		config: &container.Config{
			Image:  image,
			Tty:    false,
			Labels: map[string]string{},
			Env:    []string{},
		},
		hostConfig: &container.HostConfig{},
	}
}

func (d *RunCommand) Run(ctx *pipeline.StageContext) error {
	confs := append(d.confs, Network(ctx.InternalNetwork))
	for _, c := range confs {
		err := c.apply(d)
		if err != nil {
			return err
		}
	}
	var endpointsConfig map[string]*dockerNetwork.EndpointSettings
	if d.networkConfig != nil && d.networkConfig.EndpointsConfig != nil {
		endpointsConfig = d.networkConfig.EndpointsConfig
	}
	resp, err := ctx.Client.ContainerCreate(context.Background(), d.config, d.hostConfig, d.networkConfig, nil, d.cName)
	if err != nil {
		return err
	}
	if err := ctx.Client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	for _, net := range endpointsConfig {
		err = ctx.Client.NetworkConnect(context.Background(), net.NetworkID, resp.ID, nil)
		if err != nil {
			return err
		}
	}
	return ctx.Meta.Add("container", resp.ID)
}

func (d *RunCommand) Config(confs ...Conf) *RunCommand {
	d.confs = append(d.confs, confs...)
	return d
}
