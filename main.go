package main

import (
	"github.com/apollo-studios/gcidp-agent/docker"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/apollo-studios/gcidp-agent/utils"
	"log"
	"os"
	"os/exec"
	"plugin"
)

func main() {
	if err := utils.RunCmd(exec.Command("go", "build", "-buildmode=plugin", "-o", "plugin.so", "./plugin")); err != nil {
		log.Fatal(err)
	}
	pl := pipeline.New()
	branch := pl.BranchNormalized()
	p, err := plugin.Open("./plugin.so")
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("GCIDP_CLEANUP") == "True" {
		cleanup, err := p.Lookup("Cleanup")
		if err != nil {
			log.Fatal(err)
		}
		cleanup.(func(pl *pipeline.PipeLine, branch string))(pl, branch)
	} else {
		build, err := p.Lookup("Build")
		if err != nil {
			log.Fatal(err)
		}
		build.(func(pl *pipeline.PipeLine, branch string))(pl, branch)
	}
	pl.Stage(docker.Prune())
	pl.Run()
}
