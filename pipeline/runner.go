package pipeline

import (
	"github.com/docker/docker/client"
	"sync"
)

type Runner struct {
	pipelines  []*PipeLine
	Client     *client.Client
	Logger     Logger
	Secrets    SecretsStore
	WorkingDir string
	Branch     string
	Repo       string
}

type SecretsStore interface {
	Get(key string) (string, error)
}

type Logger interface {
	Debug(log string)
	Info(log string)
	Error(log string)
}

type RunnerOptions struct {
	WorkingDir string
	Branch     string
	Repo       string
	Secrets    SecretsStore
	Logger     Logger
}

func NewRunner(opts RunnerOptions) *Runner {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return &Runner{
		Client:     cli,
		WorkingDir: opts.WorkingDir,
		Branch:     opts.Branch,
		Repo:       opts.Repo,
		Logger:     opts.Logger,
		Secrets:    opts.Secrets,
	}
}

func (r *Runner) Run() error {
	var wg sync.WaitGroup
	wg.Add(len(r.pipelines))
	context := &StageContext{
		Client:     r.Client,
		Logger:     r.Logger,
		Branch:     r.Branch,
		WorkingDir: r.WorkingDir,
		Repo:       r.Repo,
	}
	for _, pl := range r.pipelines {
		go func(p *PipeLine) {
			if err := p.Run(context); err != nil {
				r.Logger.Error(err.Error())
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
