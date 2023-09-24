package stages

import (
	"fmt"
	"os/exec"
)

type EnvironmentVariable struct {
	Key   string
	Value string
}

type DockerRun struct {
	envVariables  []EnvironmentVariable
	containerName string
	imageName     string
	network       string
	labels        []string
}

func NewDockerRun(containerName, imageName string) *DockerRun {
	return &DockerRun{containerName: containerName, imageName: imageName}
}

func (d *DockerRun) Run() error {
	args := []string{
		"run",
		"-d",
		"--name",
		d.containerName,
	}
	for _, l := range d.labels {
		args = append(args, "--label", l)
	}
	for _, e := range d.envVariables {
		args = append(args, "-e", fmt.Sprintf("%s=%s", e.Key, e.Value))
	}
	if d.network != "" {
		args = append(args, "--network", d.network)
	}
	args = append(args, d.imageName)
	cmd := exec.Command("docker", args...)
	return runCmd(cmd)
}

func (d *DockerRun) Label(label string) *DockerRun {
	d.labels = append(d.labels, label)
	return d
}

func (d *DockerRun) Env(key, value string) *DockerRun {
	d.envVariables = append(d.envVariables, EnvironmentVariable{Key: key, Value: value})
	return d
}

func (d *DockerRun) Network(network string) *DockerRun {
	d.network = network
	return d
}
