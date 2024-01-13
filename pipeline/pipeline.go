package pipeline

import (
	"github.com/docker/docker/client"
)

type StageContext struct {
	Client          *client.Client
	Logger          Logger
	InternalNetwork string
	WorkingDir      string
	Branch          string
	Repo            string
	Secrets         SecretsStore
	Meta            Meta
}

type Stage interface {
	Run(ctx *StageContext) error
}

type PipeLine struct {
	stages []Stage
}

func (p *PipeLine) Run(ctx *StageContext) error {
	for _, s := range p.stages {
		if err := s.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}
