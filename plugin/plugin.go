package main

import (
	"fmt"
	"github.com/apollo-studios/gcidp-agent/docker"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/apollo-studios/gcidp-agent/traefik"
)

const projectName = "website"

func Cleanup(pl *pipeline.PipeLine, branch string) {
	containerName := fmt.Sprintf("%s-front-%s", projectName, branch)
	imageName := fmt.Sprintf("%s-front:%s", projectName, branch)
	pl.Stage(docker.RmContainer(containerName, true))
	pl.Stage(docker.RmImage(imageName, true))
	pl.Run()
}

func Build(pl *pipeline.PipeLine, branch string) {
	containerName := fmt.Sprintf("%s-front-%s", projectName, branch)
	imageName := fmt.Sprintf("%s-front:%s", projectName, branch)
	routerName := fmt.Sprintf("%s-%s-front", projectName, branch)

	pl.Stage(docker.RmImage(imageName, true))
	pl.Stage(docker.RmContainer(containerName, true))
	pl.Stage(docker.Build(imageName, "./context/front").Target("prod"))
	pl.Stage(
		docker.Run(containerName, imageName).
			Label("gcidp.branch", branch).
			Label(traefik.Enable, traefik.True).
			Label(traefik.TLS(routerName), traefik.True).
			Label(traefik.TLSResolver(routerName), "letsencrypt").
			Label(traefik.Rule(routerName), traefik.Host(fmt.Sprintf("%s.%s.apollos.studio", branch, projectName))).
			Label(traefik.LoadBalancerPort(routerName), "80").
			Env("NUXT_PUBLIC_API_URL", "https://api.apollos.studio").
			Env("NUXT_PUBLIC_SSR_API_URL", "http://back:3030").
			Network("app"),
	)
	pl.Stage(docker.Prune())
}
