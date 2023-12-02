package main

import (
	"github.com/apollo-studios/gcidp-agent/loader"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"log"
	"os"
)

type Logger struct {
}

func (l *Logger) Debug(msg string) {
	log.Println("DEBUG: ", msg)
}

func (l *Logger) Info(msg string) {
	log.Println("INFO: ", msg)
}

func (l *Logger) Error(msg string) {
	log.Println("ERROR: ", msg)
}

func main() {
	pluginFile := "./context/.gcidp/plugin.so"
	branch := "APL-49"
	if err := loader.BuildPlugin("./context/.gcidp/build.go", pluginFile); err != nil {
		log.Fatal("Build failed: ", err)
	}
	runner := pipeline.NewRunner(pipeline.RunnerOptions{
		WorkingDir: "./context",
		Branch:     branch,
		Repo:       "website",
		Logger:     &Logger{},
	})
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
