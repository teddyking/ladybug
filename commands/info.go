package commands

import (
	"fmt"
	"io"

	"code.cloudfoundry.org/garden"
)

type Info struct {
	Client garden.Client
	Out    io.Writer
}

func (command *Info) Execute(args []string) error {
	containers, err := command.Client.Containers(garden.Properties{})
	if err != nil {
		return err
	}

	command.Out.Write([]byte(fmt.Sprintf("Running containers: %d\n", len(containers))))
	return nil
}
