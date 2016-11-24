package commands

import (
	"code.cloudfoundry.org/garden"
	"github.com/teddyking/ladybug/print"
	"github.com/teddyking/ladybug/result"
)

type Info struct {
	Client  garden.Client
	Printer print.Printer
}

func (command *Info) Execute(args []string) error {
	var infoResult result.InfoResult

	containers, err := command.Client.Containers(garden.Properties{})
	if err != nil {
		return err
	}

	infoResult.Generate(containers)

	return command.Printer.PrintInfo(infoResult)
}
