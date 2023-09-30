package pipeline

import (
	"github.com/docker/docker/client"
	"log"
	"sync"
)

type Runner struct {
	pipelines  []*PipeLine
	Client     *client.Client
	WorkingDir string
	Branch     string
}

func NewRunner(dir, branch string) *Runner {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return &Runner{
		Client:     cli,
		WorkingDir: dir,
		Branch:     branch,
	}
}

func (r *Runner) Run() error {
	var wg sync.WaitGroup
	wg.Add(len(r.pipelines))
	context := &StageContext{
		Client:     r.Client,
		Branch:     r.Branch,
		WorkingDir: r.WorkingDir,
	}
	for _, pl := range r.pipelines {
		go func(p *PipeLine) {
			if err := p.Run(context); err != nil {
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
