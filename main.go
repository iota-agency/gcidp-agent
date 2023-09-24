package gcidp_agent

import (
	"github.com/apollo-studios/gcidp-agent/pipeline"
	"github.com/apollo-studios/gcidp-agent/stages"
)

func main() {
	pl := pipeline.New()
	build := stages.NewDockerBuild()
	build.Label("")
	pl.Stage(build)

	run := stages.NewDockerRun()
	run.Label("")
	pl.Stage(run)
}
