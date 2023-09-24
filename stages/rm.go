package stages

import (
	"os/exec"
)

type DockerRm struct {
	name  string
	image bool
}

func NewDockerRm() *DockerRm {
	return &DockerRm{}
}

func (d *DockerRm) Container(name string) *DockerRm {
	d.name = name
	return d
}

func (d *DockerRm) Image(name string) *DockerRm {
	d.name = name
	d.image = true
	return d
}

func (d *DockerRm) Run() error {
	var args []string
	if d.image {
		args = append(args, "image")
	}
	args = append(args, "rm", "-f", d.name)
	cmd := exec.Command("docker", args...)
	return runCmd(cmd)
}
