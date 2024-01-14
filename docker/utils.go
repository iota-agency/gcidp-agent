package docker

import (
	"fmt"
	"github.com/apollo-studios/gcidp-agent/traefik"
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

type expose struct {
	host, port string
}

func (e *expose) apply(d *RunCommand) error {
	routerName := utils.RandStringLowerCharSet(10)
	isDevEnv := os.Getenv("GO_APP_ENV") == "development"
	confs := []Conf{
		Label(traefik.Enable, "true"),
		Label(traefik.DefineService(routerName), routerName),
		Label(traefik.LoadBalancerPort(routerName), e.port),
	}
	if isDevEnv {
		subdomain := strings.Split(e.host, ".")[0]
		confs = append(confs,
			Label(traefik.Rule(routerName), traefik.Host(subdomain+".localhost")),
		)
	} else {
		confs = append(confs,
			Label(traefik.TLS(routerName), traefik.True),
			Label(traefik.TLSResolver(routerName), "letsencrypt"),
			Label(traefik.Rule(routerName), traefik.Host(e.host)),
			Label(traefik.Wildcard(routerName, "main"), "apollos.studio"),
			Label(traefik.Wildcard(routerName, "sans"), "*.apollos.studio"),
		)
	}
	confs = append(confs,
		Label(traefik.Network, "app"),
		Network("app"),
	)
	for _, conf := range confs {
		if err := conf.apply(d); err != nil {
			return err
		}
	}
	return nil
}

func Expose(host, port string) Conf {
	return &expose{host, port}
}

type network struct {
	name string
}

func (n *network) apply(d *RunCommand) error {
	if d.networkConfig == nil {
		d.networkConfig = &dockerNetwork.NetworkingConfig{
			EndpointsConfig: map[string]*dockerNetwork.EndpointSettings{},
		}
	}
	d.networkConfig.EndpointsConfig[n.name] = &dockerNetwork.EndpointSettings{
		NetworkID: n.name,
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
	d.config.ExposedPorts = nat.PortSet{
		nat.Port(fmt.Sprintf("%s/tcp", p.containerPort)): struct{}{},
	}
	d.hostConfig.PortBindings = nat.PortMap{
		nat.Port(fmt.Sprintf("%s/tcp", p.containerPort)): []nat.PortBinding{
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

type mountVolume struct {
	source, target string
}

func (v *mountVolume) apply(d *RunCommand) error {
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

func BindVolume(source, target string) Conf {
	return &mountVolume{source, target}
}

type volume struct {
	name, target string
}

func (v *volume) apply(d *RunCommand) error {
	d.hostConfig.Mounts = append(d.hostConfig.Mounts, mount.Mount{
		Type:   mount.TypeVolume,
		Source: v.name,
		Target: v.target,
	})
	return nil
}

func Volume(name, target string) Conf {
	return &volume{name, target}
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
