package main

import (
	"github.com/apollo-studios/gcidp-agent/loader"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"log"
	"os"
)

func main() {
	pluginFile := "./context/.gcidp/plugin.so"
	branch := "APL-49"
	if err := loader.BuildPlugin("./context/.gcidp/build.go", pluginFile); err != nil {
		log.Fatal("Build failed: ", err)
	}
	runner := pipeline.NewRunner("./context", "website", branch)
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
