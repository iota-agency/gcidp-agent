package pipeline

import (
	"fmt"
	"github.com/docker/docker/client"
	"os"
)

type Stage interface {
	Run(cli *client.Client) error
}

type PipeLine struct {
	Client *client.Client
	Stages []Stage
}

func New() *PipeLine {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return &PipeLine{
		Client: cli,
	}
}

func (p *PipeLine) Run() {
	for _, s := range p.Stages {
		if err := s.Run(p.Client); err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}

func (p *PipeLine) Stage(s Stage) *PipeLine {
	p.Stages = append(p.Stages, s)
	return p
}

func (p *PipeLine) Branch() string {
	return os.Getenv("GITHUB_REF_NAME")
}
