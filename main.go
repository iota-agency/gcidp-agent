package main

import (
	"github.com/apollo-studios/gcidp-agent/loader"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"log"
	"os"
)

func main() {
	pluginFile := "./example/plugin.so"
	if err := loader.BuildPlugin(pluginFile, "./example"); err != nil {
		log.Fatal(err)
	}
	runner := pipeline.NewRunner()
	branch := runner.Branch()
	if branch == "" {
		panic("branch is empty")
	}
	p, err := loader.Load(pluginFile)
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("GCIDP_CLEANUP") == "True" {
		p.Cleanup(runner, branch)
	} else {
		p.Build(runner, branch)
	}
	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}
}
