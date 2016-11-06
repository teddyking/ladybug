package commands

import (
	"fmt"
	"io"

	"code.cloudfoundry.org/garden"
	"github.com/teddyking/ladybug/system"
)

type Containers struct {
	Client garden.Client
	Host   system.Host
	Out    io.Writer
}

func (command *Containers) Execute(args []string) error {
	containers, err := command.Client.Containers(garden.Properties{})
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		command.Out.Write([]byte("0 running containers found on this host\n"))
		return nil
	}

	for _, container := range containers {
		handle := container.Handle()

		containerInfo, err := container.Info()
		if err != nil {
			return err
		}

		containerPids, err := command.Host.ContainerPids(handle)
		if err != nil {
			return err
		}

		containerProcessName := "N/A"
		if len(containerPids) > 0 {
			containerProcessName, err = command.Host.ContainerProcessName(containerPids[0])
			if err != nil {
				return err
			}
		}

		detailedContainerInfo := fmt.Sprintf("%s - %s - %s\n", handle, containerInfo.ContainerIP, containerProcessName)
		command.Out.Write([]byte(detailedContainerInfo))
	}

	return nil
}
