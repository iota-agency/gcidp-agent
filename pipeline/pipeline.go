package pipeline

import (
	"github.com/docker/docker/client"
)

type StageContext struct {
	Client     *client.Client
	WorkingDir string
	Branch     string
}

type Stage interface {
	Run(context *StageContext) error
}

type PipeLine struct {
	stages []Stage
}

func (p *PipeLine) Run(context *StageContext) error {
	for _, s := range p.stages {
		if err := s.Run(context); err != nil {
			return err
		}
	}
	return nil
}
