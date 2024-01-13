package main

import (
	"fmt"
	"github.com/apollo-studios/gcidp-agent/docker"
	"github.com/apollo-studios/gcidp-agent/pipeline"
)

var BuildServerVersion = "0.2.1"

const projectName = "website"

func Cleanup(runner *pipeline.Runner) {
	branch := runner.Branch
	containerName := fmt.Sprintf("%s-app-%s", projectName, branch)
	imageName := fmt.Sprintf("%s-app:%s", projectName, branch)
	runner.Pipeline(
		docker.RmContainer(containerName, true),
		docker.RmImage(imageName, true),
	)
}

func Build(runner *pipeline.Runner) {
	branch := runner.Branch
	containerName := fmt.Sprintf("%s-app-%s", projectName, branch)
	imageName := fmt.Sprintf("%s-app:%s", projectName, branch)
	runner.Pipeline(
		docker.Build(imageName, "./app").Target("prod"),
		docker.RmContainer(containerName, true),
		docker.Run(containerName, imageName).Config(
			docker.PortBinding("8080", "8080"),
			docker.Env("SOME_VAR", "some_value"),
		),
	)
}
