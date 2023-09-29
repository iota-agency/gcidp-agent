package pipeline

import (
	"github.com/docker/docker/client"
	"os"
	"sync"
)

type Runner struct {
	pipelines []*PipeLine
	Client    *client.Client
}

func NewRunner() *Runner {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return &Runner{Client: cli}
}

func (r *Runner) Run() error {
	var wg sync.WaitGroup
	wg.Add(len(r.pipelines))
	var errors chan error
	for _, pl := range r.pipelines {
		go func(p *PipeLine) {
			if err := p.Run(r.Client); err != nil {
				errors <- err
			}
			defer wg.Done()
		}(pl)
	}
	wg.Wait()
	return nil
}

func (r *Runner) Pipeline() *PipeLine {
	p := &PipeLine{}
	r.pipelines = append(r.pipelines, p)
	return p
}

func (r *Runner) Branch() string {
	return os.Getenv("GITHUB_REF_NAME")
}
