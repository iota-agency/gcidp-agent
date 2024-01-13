package pipeline

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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
	Cleanup    bool
}

type SecretsStore interface {
	Get(key string) (string, error)
}

type Logger interface {
	Debug(log string)
	Info(log string)
	Error(log string)
	Warn(log string)
}

type RunnerOptions struct {
	WorkingDir string
	Branch     string
	Repo       string
	Secrets    SecretsStore
	Logger     Logger
	Cleanup    bool
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
		Cleanup:    opts.Cleanup,
	}
}

func (r *Runner) Run() error {
	internalNetwork := fmt.Sprintf("%s-%s-internal", r.Repo, r.Branch)
	if !r.Cleanup {
		filter := filters.NewArgs()
		filter.Add("name", internalNetwork)
		networks, err := r.Client.NetworkList(context.Background(), types.NetworkListOptions{Filters: filter})
		if err != nil {
			return err
		}
		if len(networks) == 0 {
			_, err := r.Client.NetworkCreate(context.Background(), internalNetwork, types.NetworkCreate{})
			if err != nil {
				return err
			}
		}
	}
	ctx := &StageContext{
		Client:          r.Client,
		Logger:          r.Logger,
		Branch:          r.Branch,
		WorkingDir:      r.WorkingDir,
		Repo:            r.Repo,
		InternalNetwork: internalNetwork,
	}
	var wg sync.WaitGroup
	wg.Add(len(r.pipelines))
	for _, pl := range r.pipelines {
		go func(p *PipeLine) {
			if err := p.Run(ctx); err != nil {
				r.Logger.Error(err.Error())
			}
			defer wg.Done()
		}(pl)
	}
	wg.Wait()
	if r.Cleanup {
		err := r.Client.NetworkRemove(context.Background(), internalNetwork)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) Pipeline(stages ...Stage) *PipeLine {
	p := &PipeLine{stages: stages}
	r.pipelines = append(r.pipelines, p)
	return p
}
