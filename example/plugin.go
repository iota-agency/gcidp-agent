package main

import (
	"fmt"
	"github.com/apollo-studios/gcidp-agent/docker"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/apollo-studios/gcidp-agent/traefik"
)

var BuildServerVersion = "0.1.6"

const projectName = "website"

func Cleanup(runner *pipeline.Runner, branch string) {
	containerName := fmt.Sprintf("%s-front-%s", projectName, branch)
	imageName := fmt.Sprintf("%s-front:%s", projectName, branch)
	runner.Pipeline(
		docker.RmContainer(containerName, true),
		docker.RmImage(imageName, true),
	)
}

func Build(runner *pipeline.Runner, branch string) {
	containerName := fmt.Sprintf("%s-front-%s", projectName, branch)
	imageName := fmt.Sprintf("%s-front:%s", projectName, branch)
	routerName := fmt.Sprintf("%s-%s-front", projectName, branch)
	runner.Pipeline(
		docker.Build(imageName, "./context/front").Target("prod"),
		docker.RmContainer(containerName, true),
		docker.Run(containerName, imageName).Config(
			docker.Label("gcidp.branch", branch),
			docker.Label(traefik.Enable, traefik.True),
			docker.Label(traefik.TLS(routerName), traefik.True),
			docker.Label(traefik.TLSResolver(routerName), "letsencrypt"),
			docker.Label(traefik.Rule(routerName), traefik.Host(fmt.Sprintf("%s.%s.apollos.studio", branch, projectName))),
			docker.Label(traefik.LoadBalancerPort(routerName), "80"),
			docker.Env("NUXT_PUBLIC_API_URL", "https://api.apollos.studio"),
			docker.Env("NUXT_PUBLIC_SSR_API_URL", "http://back:3030"),
			docker.Network("app"),
		),
	)

	containerName = fmt.Sprintf("%s-houston-%s", projectName, branch)
	imageName = fmt.Sprintf("%s-houston:%s", projectName, branch)
	routerName = fmt.Sprintf("%s-%s-houston", projectName, branch)
	runner.Pipeline(
		docker.Build(imageName, "./context/admin").Target("prod"),
		docker.RmContainer(containerName, true),
		docker.Run(containerName, imageName).Config(
			docker.Label("gcidp.branch", branch),
			docker.Label(traefik.Enable, traefik.True),
			docker.Label(traefik.TLS(routerName), traefik.True),
			docker.Label(traefik.TLSResolver(routerName), "letsencrypt"),
			docker.Label(traefik.Rule(routerName), traefik.Host(fmt.Sprintf("houston.%s.%s.apollos.studio", branch, projectName))),
			docker.Label(traefik.LoadBalancerPort(routerName), "80"),
			docker.Env("NUXT_PUBLIC_API_URL", "https://api.apollos.studio"),
			docker.Env("NUXT_PUBLIC_SSR_API_URL", "http://back:3030"),
			docker.Network("app"),
		),
	)
}
