package pipeline

import (
	"os"
	"strings"
)

type Stage interface {
	Run() error
}

type PipeLine struct {
	Stages []Stage
}

func New() *PipeLine {
	return &PipeLine{}
}

func (p *PipeLine) Run() {
	for _, s := range p.Stages {
		if err := s.Run(); err != nil {
			panic(err)
		}
	}
}

func (p *PipeLine) Stage(s Stage) {
	p.Stages = append(p.Stages, s)
}

func (p *PipeLine) Branch() string {
	return os.Getenv("GITHUB_REF_NAME")
}

func (p *PipeLine) BranchNormalized() string {
	b := strings.Replace(p.Branch(), "-", "", -1)
	b = strings.Replace(b, ".", "", -1)
	return strings.ToLower(b)
}
