package main

import (
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"os"
	"os/exec"
	"plugin"
)

func main() {
	if err := exec.Command("go", "build", "-buildmode=plugin", "-o", "plugin.so", "./plugin").Run(); err != nil {
		panic(err)
	}
	pl := pipeline.New()
	branch := pl.BranchNormalized()
	p, err := plugin.Open("./plugin.so")
	if err != nil {
		panic(err)
	}
	if os.Getenv("GCIDP_CLEANUP") == "True" {
		cleanup, err := p.Lookup("Cleanup")
		if err != nil {
			panic(err)
		}
		cleanup.(func(pl *pipeline.PipeLine, branch string))(pl, branch)
	} else {
		build, err := p.Lookup("Build")
		if err != nil {
			panic(err)
		}
		build.(func(pl *pipeline.PipeLine, branch string))(pl, branch)
	}
	pl.Run()
}
