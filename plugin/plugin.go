package main

import (
	"fmt"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/apollo-studios/gcidp-agent/stages"
	"github.com/apollo-studios/gcidp-agent/traefik"
)

const projectName = "website"

func Cleanup(pl *pipeline.PipeLine, branch string) {
	containerName := fmt.Sprintf("%s-front-%s", projectName, branch)
	imageName := fmt.Sprintf("%s-front:%s", projectName, branch)
	pl.Stage(stages.NewDockerRm().Container(containerName))
	pl.Stage(stages.NewDockerRm().Image(imageName))
	pl.Run()
}

func Build(pl *pipeline.PipeLine, branch string) {
	containerName := fmt.Sprintf("%s-front-%s", projectName, branch)
	imageName := fmt.Sprintf("%s-front:%s", projectName, branch)
	routerName := fmt.Sprintf("%s-%s-front", projectName, branch)

	pl.Stage(stages.NewDockerRm().Container(containerName))
	pl.Stage(stages.NewDockerBuild(imageName, "./front").Target("prod"))
	pl.Stage(
		stages.NewDockerRun(containerName, imageName).
			Label(traefik.Enable, traefik.True).
			Label(traefik.Host(routerName), fmt.Sprintf("%s.%s.apollos.studio", branch, projectName)).
			Label(traefik.TLS(routerName), traefik.True).
			Label(traefik.TLSResolver(routerName), "letsencrypt").
			Label(traefik.LoadBalancerPort(routerName), "80").
			Env("NUXT_PUBLIC_API_URL", "https://api.apollos.studio").
			Env("NUXT_PUBLIC_SSR_API_URL", "http://back:3030").
			Network("app"),
	)
}
