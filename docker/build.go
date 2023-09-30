package docker

import (
	"context"
	"fmt"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/apollo-studios/gcidp-agent/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"os"
	"path/filepath"
	"strings"
)

type BuildCommand struct {
	target  string
	image   string
	context string
	exclude []string
}

func Build(image, context string) *BuildCommand {
	return &BuildCommand{image: image, context: context}
}

func ReadIgnore(f string) []string {
	file, err := os.ReadFile(f)
	if err != nil {
		return []string{}
	}
	var result []string
	entries := strings.Split(string(file), "\n")
	for _, e := range entries {
		trimmed := strings.Trim(e, " ")
		if len(trimmed) > 0 {
			result = append(result, trimmed)
		}
	}
	return result
}

func (d *BuildCommand) Run(ctx *pipeline.StageContext) error {
	exclude := d.exclude
	dockerCtx := filepath.Join(ctx.WorkingDir, d.context)
	ignoreFile := filepath.Join(dockerCtx, ".dockerignore")
	if utils.FileExists(ignoreFile) {
		exclude = append(exclude, ReadIgnore(ignoreFile)...)
	}

	tar, err := archive.TarWithOptions(dockerCtx, &archive.TarOptions{
		ExcludePatterns: exclude,
	})
	if err != nil {
		return err
	}

	opts := types.ImageBuildOptions{
		Dockerfile:  "Dockerfile", // TODO: Make this configurable
		Context:     tar,
		Target:      d.target,
		Tags:        []string{d.image},
		Remove:      true,
		ForceRemove: true,
	}
	fmt.Println(opts.Tags, dockerCtx)
	build, err := ctx.Client.ImageBuild(context.Background(), tar, opts)
	if err != nil {
		return err
	}
	defer build.Body.Close()
	if err := utils.Print(build.Body); err != nil {
		return err
	}
	return nil
}

func (d *BuildCommand) Exclude(files []string) *BuildCommand {
	d.exclude = files
	return d
}

func (d *BuildCommand) Target(t string) *BuildCommand {
	d.target = t
	return d
}
