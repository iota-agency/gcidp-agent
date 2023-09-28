package main

import (
	"fmt"
	"github.com/apollo-studios/gcidp-agent/pipeline"
)

const projectName = "website"

func Cleanup(pl *pipeline.PipeLine, branch string) {
	fmt.Println("No Cleanup function is registered")
}

func Build(pl *pipeline.PipeLine, branch string) {
	fmt.Println("No Build function is registered")
}
