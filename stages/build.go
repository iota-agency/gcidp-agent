package stages

import (
	"os/exec"
)

type DockerBuild struct {
	target    string
	imageName string
	context   string
}

func NewDockerBuild(imageName, context string) *DockerBuild {
	return &DockerBuild{imageName: imageName, context: context}
}

func (d *DockerBuild) Run() error {
	args := []string{
		"build",
	}
	if d.target != "" {
		args = append(args, "--target "+d.target)
	}
	args = append(args, "--tag "+d.imageName, d.context)
	cmd := exec.Command("docker", args...)
	return cmd.Run()
}

func (d *DockerBuild) Target(t string) *DockerBuild {
	d.target = t
	return d
}
