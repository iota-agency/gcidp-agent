package tests

import (
	"encoding/json"
	"fmt"
	"github.com/apollo-studios/gcidp-agent/docker"
	"github.com/apollo-studios/gcidp-agent/loader"
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"net/http"
	"os"
	"strings"
	"testing"
)

type MetaResponse struct {
	Env []string `json:"env"`
}

type Logger struct {
}

func (l *Logger) Debug(msg string) {
	fmt.Println("DEBUG: ", msg)
}

func (l *Logger) Info(msg string) {
	fmt.Println("INFO: ", msg)
}

func (l *Logger) Warn(msg string) {
	fmt.Println("WARN: ", msg)
}

func (l *Logger) Error(msg string) {
	fmt.Println("ERROR: ", msg)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestFullBuildCycle(t *testing.T) {
	buildRunner := pipeline.NewRunner(pipeline.RunnerOptions{
		WorkingDir: "test",
		Branch:     "test",
		Repo:       "test",
		Logger:     &Logger{},
	})
	var pluginFile = "./test/.gcidp/plugin.so"
	if err := loader.BuildPlugin("./test/.gcidp/build.go", pluginFile); err != nil {
		t.Fatal(err)
	}
	p, err := loader.Load(pluginFile)
	if err != nil {
		t.Fatal(err)
	}
	p.Build(buildRunner)
	if err := buildRunner.Run(); err != nil {
		t.Fatal("Load failed: ", err)
	}
	cleanupRunner := pipeline.NewRunner(pipeline.RunnerOptions{
		WorkingDir: "test",
		Branch:     "test",
		Repo:       "test",
		Logger:     &Logger{},
	})
	p.Cleanup(cleanupRunner)
	if err := cleanupRunner.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestEnvironmentVariables(t *testing.T) {
	buildRunner := pipeline.NewRunner(pipeline.RunnerOptions{
		WorkingDir: "test",
		Branch:     "test",
		Repo:       "test",
		Logger:     &Logger{},
	})
	buildRunner.Pipeline(
		docker.Build("test-app:test", "./app").Target("prod"),
		docker.Run("test-test-internal", "test-app:test").Config(
			docker.PortBinding("8080", "8080"),
			docker.Env("SOME_VAR", "some_value"),
		),
	)
	if err := buildRunner.Run(); err != nil {
		t.Fatal("Load failed: ", err)
	}
	resp, err := http.Get("http://localhost:8080/meta")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Status code: %d", resp.StatusCode)
	}
	meta := &MetaResponse{}
	if err := json.NewDecoder(resp.Body).Decode(meta); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(meta.Env, " "), "SOME_VAR=some_value") {
		t.Fatal("Env var not found")
	}
	cleanupRunner := pipeline.NewRunner(pipeline.RunnerOptions{
		WorkingDir: "test",
		Branch:     "test",
		Repo:       "test",
		Logger:     &Logger{},
		Cleanup:    true,
	})
	cleanupRunner.Pipeline(
		docker.RmContainer("test-test-internal", true),
		docker.RmImage("test-app:test", true),
	)
	if err := cleanupRunner.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestNamedVolumes(t *testing.T) {
	buildRunner := pipeline.NewRunner(pipeline.RunnerOptions{
		WorkingDir: "test",
		Branch:     "test",
		Repo:       "test",
		Logger:     &Logger{},
	})
	buildRunner.Pipeline(
		docker.Build("test-app:test", "./app").Target("prod"),
		docker.Run("test-test-internal", "test-app:test").Config(
			docker.PortBinding("8080", "8080"),
			docker.Env("SOME_VAR", "some_value"),
			docker.Volume("test-db-data", "/app/data"),
		),
	)
	if err := buildRunner.Run(); err != nil {
		t.Fatal("Run failed: ", err)
	}
	cleanupRunner := pipeline.NewRunner(pipeline.RunnerOptions{
		WorkingDir: "test",
		Branch:     "test",
		Repo:       "test",
		Logger:     &Logger{},
		Cleanup:    true,
	})
	cleanupRunner.Pipeline(
		docker.RmContainer("test-test-internal", true),
		docker.RmImage("test-app:test", true),
		docker.RmVolume("test-db-data"),
	)
	if err := cleanupRunner.Run(); err != nil {
		t.Fatal(err)
	}
}
