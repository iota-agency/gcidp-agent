package docker

import (
	"context"
	"github.com/apollo-studios/gcidp-agent/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type BuildCommand struct {
	target  string
	image   string
	context string
}

func Build(image, context string) *BuildCommand {
	return &BuildCommand{image: image, context: context}
}

func (d *BuildCommand) Run(cli *client.Client) error {
	tar, err := archive.TarWithOptions(d.context, &archive.TarOptions{})
	if err != nil {
		return err
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // TODO: Make this configurable
		Context:    tar,
		Target:     d.target,
		Tags:       []string{d.image},
	}
	build, err := cli.ImageBuild(context.Background(), tar, opts)
	if err != nil {
		return err
	}
	defer build.Body.Close()
	if err := utils.Print(build.Body); err != nil {
		return err
	}
	return nil
}

func (d *BuildCommand) Target(t string) *BuildCommand {
	d.target = t
	return d
}
