package stages

import (
	"os/exec"
)

type DockerRm struct {
	name string
}

func NewDockerRm(name string) *DockerRm {
	return &DockerRm{name: name}
}

func (d *DockerRm) Run() error {
	args := []string{
		"rm",
		"-f",
		d.name,
	}
	cmd := exec.Command("docker", args...)
	return cmd.Run()
}
