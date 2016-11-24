package commands

import (
	"code.cloudfoundry.org/garden"
	"github.com/teddyking/ladybug/output"
	"github.com/teddyking/ladybug/result"
)

type Info struct {
	Client  garden.Client
	Printer output.Printer
}

func (command *Info) Execute(args []string) error {
	var infoResult result.Info

	containers, err := command.Client.Containers(garden.Properties{})
	if err != nil {
		return err
	}

	infoResult.Generate(containers)

	return command.Printer.PrintInfo(infoResult)
}
