package main

import (
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
	runner := pipeline.NewRunner()
	branch := runner.Branch()
	p, err := plugin.Open("./plugin.so")
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("GCIDP_CLEANUP") == "True" {
		cleanup, err := p.Lookup("Cleanup")
		if err != nil {
			log.Fatal(err)
		}
		cleanup.(func(runner *pipeline.Runner, branch string))(runner, branch)
	} else {
		build, err := p.Lookup("Build")
		if err != nil {
			log.Fatal(err)
		}
		build.(func(runner *pipeline.Runner, branch string))(runner, branch)
	}
	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}
}
