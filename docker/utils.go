package docker

import (
	"fmt"
	"github.com/apollo-studios/gcidp-agent/utils"
	"github.com/docker/docker/api/types/mount"
	dockerNetwork "github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"os"
	"path/filepath"
	"strings"
)

type Conf interface {
	apply(d *RunCommand) error
}

type label struct {
	key, value string
}

func (l *label) apply(d *RunCommand) error {
	d.config.Labels[l.key] = l.value
	return nil
}

func Label(key, value string) Conf {
	return &label{key, value}
}

type env struct {
	key, value string
}

func (e *env) apply(d *RunCommand) error {
	d.config.Env = append(d.config.Env, e.key+"="+e.value)
	return nil
}

func Env(key, value string) Conf {
	return &env{key, value}
}

type network struct {
	name string
}

func (n *network) apply(d *RunCommand) error {
	d.networkConfig = &dockerNetwork.NetworkingConfig{
		EndpointsConfig: map[string]*dockerNetwork.EndpointSettings{
			n.name: {
				NetworkID: n.name,
			},
		},
	}
	return nil
}

func Network(name string) Conf {
	return &network{name}
}

type portBinding struct {
	hostPort, containerPort string
}

func (p *portBinding) apply(d *RunCommand) error {
	d.hostConfig.PortBindings = nat.PortMap{
		nat.Port(p.containerPort): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: p.hostPort,
			},
		},
	}
	return nil
}

func PortBinding(hostPort, containerPort string) Conf {
	return &portBinding{hostPort, containerPort}
}

type volume struct {
	source, target string
}

func (v *volume) apply(d *RunCommand) error {
	fmt.Println("source", v.source)
	if err := utils.MkDirIfNone(v.source); err != nil {
		return err
	}
	if strings.HasPrefix(v.source, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		v.source = filepath.Join(home, v.source[2:])
	}
	d.hostConfig.Mounts = append(d.hostConfig.Mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: v.source,
		Target: v.target,
	})
	return nil
}

func Volume(source, target string) Conf {
	return &volume{source, target}
}

func Hostname(name string) Conf {
	return &hostname{name}
}

type hostname struct {
	name string
}

func (h *hostname) apply(d *RunCommand) error {
	d.config.Hostname = h.name
	return nil
}
