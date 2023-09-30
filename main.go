package main

import (
	"github.com/apollo-studios/gcidp-agent/loader"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"log"
	"os"
)

func main() {
	pluginFile := "./example/plugin.so"
	branch := "APL-49"
	if err := loader.BuildPlugin(pluginFile, "./example"); err != nil {
		log.Fatal(err)
	}
	runner := pipeline.NewRunner("./context", branch)
	p, err := loader.Load(pluginFile)
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("GCIDP_CLEANUP") == "True" {
		p.Cleanup(runner)
	} else {
		p.Build(runner)
	}
	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}
}
