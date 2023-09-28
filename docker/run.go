package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"strings"
)

type RunCommand struct {
	env       map[string]string
	container string
	image     string
	network   string
	labels    map[string]string
}

func Run(container, image string) *RunCommand {
	return &RunCommand{
		container: container,
		image:     image,
		labels:    map[string]string{},
		env:       map[string]string{},
	}
}

func (d *RunCommand) envList() []string {
	var env []string
	for k, v := range d.env {
		env = append(env, k+"="+v)
	}
	return env
}

func (d *RunCommand) Run(cli *client.Client) error {
	var networkConfig *network.NetworkingConfig
	if d.network != "" {
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				d.network: {
					NetworkID: d.network,
				},
			},
		}
	}
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:  d.image,
		Tty:    false,
		Labels: d.labels,
		Env:    d.envList(),
	}, nil, networkConfig, nil, d.container)
	if err != nil {
		return err
	}
	return cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
}

func (d *RunCommand) Label(key, value string) *RunCommand {
	d.labels[key] = value
	return d
}

func (d *RunCommand) LabelString(label string) *RunCommand {
	v := strings.Split(label, "=")
	d.labels[v[0]] = v[1]
	return d
}

func (d *RunCommand) Env(key, value string) *RunCommand {
	d.env[key] = value
	return d
}

func (d *RunCommand) Network(network string) *RunCommand {
	d.network = network
	return d
}
