package loader

import (
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/apollo-studios/gcidp-agent/utils"
	"os/exec"
	"plugin"
)

type Plugin struct {
	BuildVersion string
	Build        func(runner *pipeline.Runner, branch string)
	Cleanup      func(runner *pipeline.Runner, branch string)
}

func Load(fn string) (*Plugin, error) {
	p, err := plugin.Open(fn)
	if err != nil {
		return nil, err
	}
	cleanup, err := p.Lookup("Cleanup")
	if err != nil {
		return nil, err
	}
	build, err := p.Lookup("Build")
	if err != nil {
		return nil, err
	}
	v, err := p.Lookup("BuildServerVersion")
	if err != nil {
		return nil, err
	}
	return &Plugin{
		Cleanup:      cleanup.(func(runner *pipeline.Runner, branch string)),
		Build:        build.(func(runner *pipeline.Runner, branch string)),
		BuildVersion: *v.(*string),
	}, nil
}

func BuildPlugin(dst string, src string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", dst, src)
	return utils.RunCmd(cmd)
}
