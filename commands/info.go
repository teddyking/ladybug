package commands

import (
	"code.cloudfoundry.org/garden"
	"github.com/teddyking/ladybug/print"
)

type Info struct {
	Client  garden.Client
	Printer print.Printer
}

func (command *Info) Execute(args []string) error {
	var result print.InfoResult

	containers, err := command.Client.Containers(garden.Properties{})
	if err != nil {
		return err
	}

	result.ContainersCount = len(containers)

	return command.Printer.PrintInfo(result)
}
