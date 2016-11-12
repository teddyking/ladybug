package commands

import (
	"fmt"

	"github.com/concourse/fly/commands/internal/displayhelpers"
	"github.com/concourse/fly/rc"
)

type ExposePipelineCommand struct {
	Pipeline string `short:"p" long:"pipeline" required:"true" description:"Pipeline to expose"`
}

func (command *ExposePipelineCommand) Execute(args []string) error {
	pipelineName := command.Pipeline

	target, err := rc.LoadTarget(Fly.Target)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	found, err := target.Team().ExposePipeline(pipelineName)
	if err != nil {
		return err
	}

	if found {
		fmt.Printf("exposed '%s'\n", pipelineName)
	} else {
		displayhelpers.Failf("pipeline '%s' not found\n", pipelineName)
	}

	return nil
}
