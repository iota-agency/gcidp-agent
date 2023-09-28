package stages

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"strings"
)

type DockerRun struct {
	env       map[string]string
	container string
	image     string
	network   string
	labels    map[string]string
}

func NewDockerRun(container, image string) *DockerRun {
	return &DockerRun{
		container: container,
		image:     image,
		labels:    map[string]string{},
		env:       map[string]string{},
	}
}

func (d *DockerRun) envList() []string {
	var env []string
	for k, v := range d.env {
		env = append(env, k+"="+v)
	}
	return env
}

func (d *DockerRun) Run(cli *client.Client) error {
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			d.network: {
				NetworkID: d.network,
			},
		},
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

func (d *DockerRun) Label(key, value string) *DockerRun {
	d.labels[key] = value
	return d
}

func (d *DockerRun) LabelString(label string) *DockerRun {
	v := strings.Split(label, "=")
	d.labels[v[0]] = v[1]
	return d
}

func (d *DockerRun) Env(key, value string) *DockerRun {
	d.env[key] = value
	return d
}

func (d *DockerRun) Network(network string) *DockerRun {
	d.network = network
	return d
}
