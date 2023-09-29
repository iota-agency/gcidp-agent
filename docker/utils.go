package docker

import (
	dockerNetwork "github.com/docker/docker/api/types/network"
)

type Conf interface {
	apply(d *RunCommand)
}

type label struct {
	key, value string
}

func (l *label) apply(d *RunCommand) {
	d.config.Labels[l.key] = l.value
}

func Label(key, value string) Conf {
	return &label{key, value}
}

type env struct {
	key, value string
}

func (e *env) apply(d *RunCommand) {
	d.config.Env = append(d.config.Env, e.key+"="+e.value)
}

func Env(key, value string) Conf {
	return &env{key, value}
}

type network struct {
	name string
}

func (n *network) apply(d *RunCommand) {
	d.networkConfig = &dockerNetwork.NetworkingConfig{
		EndpointsConfig: map[string]*dockerNetwork.EndpointSettings{
			n.name: {
				NetworkID: n.name,
			},
		},
	}
}

func Network(name string) Conf {
	return &network{name}
}
