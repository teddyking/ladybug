package commands

import (
	"fmt"

	"github.com/concourse/fly/commands/internal/displayhelpers"
	"github.com/concourse/fly/rc"
)

type UnpausePipelineCommand struct {
	Pipeline string `short:"p" long:"pipeline" required:"true" description:"Pipeline to unpause"`
}

func (command *UnpausePipelineCommand) Execute(args []string) error {
	pipelineName := command.Pipeline

	target, err := rc.LoadTarget(Fly.Target)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	found, err := target.Team().UnpausePipeline(pipelineName)
	if err != nil {
		return err
	}

	if found {
		fmt.Printf("unpaused '%s'\n", pipelineName)
	} else {
		displayhelpers.Failf("pipeline '%s' not found\n", pipelineName)
	}

	return nil
}
