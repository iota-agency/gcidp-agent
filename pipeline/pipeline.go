package pipeline

import "os"

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

}

func (p *PipeLine) Stage(s Stage) {
	p.Stages = append(p.Stages, s)
}

func (p *PipeLine) Branch() string {
	return os.Getenv("GITHUB_REF_NAME")
}
