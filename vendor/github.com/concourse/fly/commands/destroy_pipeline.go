package commands

import (
	"fmt"

	"github.com/concourse/fly/rc"
	"github.com/vito/go-interact/interact"
)

type DestroyPipelineCommand struct {
	Pipeline        string `short:"p"  long:"pipeline" required:"true" description:"Pipeline to destroy"`
	SkipInteractive bool   `short:"n"  long:"non-interactive"          description:"Destroy the pipeline without confirmation"`
}

func (command *DestroyPipelineCommand) Execute(args []string) error {
	target, err := rc.LoadTarget(Fly.Target)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	pipelineName := command.Pipeline
	fmt.Printf("!!! this will remove all data for pipeline `%s`\n\n", pipelineName)

	confirm := command.SkipInteractive
	if !confirm {
		err := interact.NewInteraction("are you sure?").Resolve(&confirm)
		if err != nil || !confirm {
			fmt.Println("bailing out")
			return err
		}
	}

	found, err := target.Team().DeletePipeline(pipelineName)
	if err != nil {
		return err
	}

	if !found {
		fmt.Printf("`%s` does not exist\n", pipelineName)
	} else {
		fmt.Printf("`%s` deleted\n", pipelineName)
	}

	return nil
}
