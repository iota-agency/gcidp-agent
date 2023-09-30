package pipeline

import (
	"github.com/docker/docker/client"
	"log"
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
	for _, pl := range r.pipelines {
		go func(p *PipeLine) {
			if err := p.Run(r.Client); err != nil {
				log.Println(err)
			}
			defer wg.Done()
		}(pl)
	}
	wg.Wait()
	return nil
}

func (r *Runner) Pipeline(stages ...Stage) *PipeLine {
	p := &PipeLine{stages: stages}
	r.pipelines = append(r.pipelines, p)
	return p
}

func (r *Runner) Branch() string {
	return os.Getenv("GITHUB_REF_NAME")
}
