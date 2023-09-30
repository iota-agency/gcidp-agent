package loader

import (
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/apollo-studios/gcidp-agent/utils"
	"os/exec"
	"plugin"
)

type Plugin struct {
	BuildVersion string
	Build        func(runner *pipeline.Runner)
	Cleanup      func(runner *pipeline.Runner)
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
		Cleanup:      cleanup.(func(runner *pipeline.Runner)),
		Build:        build.(func(runner *pipeline.Runner)),
		BuildVersion: *v.(*string),
	}, nil
}

func BuildPlugin(src string, dst string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", dst, src)
	return utils.RunCmd(cmd)
}
