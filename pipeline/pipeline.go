package pipeline

import (
	"github.com/docker/docker/client"
)

type Stage interface {
	Run(cli *client.Client) error
}

type PipeLine struct {
	stages []Stage
}

func (p *PipeLine) Run(cli *client.Client) error {
	for _, s := range p.stages {
		if err := s.Run(cli); err != nil {
			return err
		}
	}
	return nil
}
