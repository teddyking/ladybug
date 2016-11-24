package commands

import (
	"code.cloudfoundry.org/garden"
	"github.com/teddyking/ladybug/output"
	"github.com/teddyking/ladybug/result"
	"github.com/teddyking/ladybug/sys"
)

type Containers struct {
	Client  garden.Client
	Host    sys.Host
	Printer output.Printer
}

func (command *Containers) Execute(args []string) error {
	containers, err := command.Client.Containers(garden.Properties{})
	if err != nil {
		return err
	}

	containersResult := make(result.Containers, len(containers))

	err = containersResult.Generate(
		result.WithHandles(containers),
		result.WithIPs(containers),
		result.WithProcessNames(containers, command.Host),
		result.WithCreatedAtTimes(containers, command.Host),
		result.WithPortMappings(containers),
	)
	if err != nil {
		return err
	}

	return command.Printer.PrintContainers(containersResult)
}
