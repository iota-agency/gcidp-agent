package main

import (
	"fmt"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/apollo-studios/gcidp-agent/stages"
	"github.com/apollo-studios/gcidp-agent/traefik"
	"os"
)

func main() {
	projectName := "website"
	pl := pipeline.New()
	branch := pl.BranchNormalized()
	containerName := fmt.Sprintf("%s-front-%s", projectName, branch)
	imageName := fmt.Sprintf("%s-front:%s", projectName, branch)
	routerName := fmt.Sprintf("%s-%s-front", projectName, branch)

	if os.Getenv("GCIDP_CLEANUP") == "True" {
		pl.Stage(stages.NewDockerRm().Container(containerName))
		pl.Stage(stages.NewDockerRm().Image(imageName))
		pl.Run()
		return
	}

	pl.Stage(stages.NewDockerRm().Container(containerName))
	pl.Stage(stages.NewDockerBuild(imageName, "./front").Target("prod"))
	pl.Stage(
		stages.NewDockerRun(containerName, imageName).
			Label(traefik.Enable()).
			Label(traefik.Host(routerName, fmt.Sprintf("%s.%s.apollos.studio", branch, projectName))).
			Label(traefik.EnableTls(routerName)).
			Label(traefik.TslResolver(routerName)).
			Label(traefik.LoadBalancerPort(routerName, 80)).
			Env("NUXT_PUBLIC_API_URL", "https://api.apollos.studio").
			Env("NUXT_PUBLIC_SSR_API_URL", "http://back:3030"),
	)
	pl.Run()
}

// GOOS=linux GOARCH=amd64 go build -o build .
