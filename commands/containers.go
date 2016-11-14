package commands

import (
	"code.cloudfoundry.org/garden"
	"github.com/teddyking/ladybug/print"
	"github.com/teddyking/ladybug/system"
)

type Containers struct {
	Client  garden.Client
	Host    system.Host
	Printer print.Printer
}

func (command *Containers) Execute(args []string) error {
	var result print.ContainersResult
	var containerInfos []print.ContainerInfo

	containers, err := command.Client.Containers(garden.Properties{})
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		return command.Printer.PrintContainers(result)
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

		containerCreationTime, err := command.Host.ContainerCreationTime(handle)
		if err != nil {
			return err
		}

		containerInfos = append(containerInfos, print.ContainerInfo{
			Handle:       handle,
			Ip:           containerInfo.ContainerIP,
			ProcessName:  containerProcessName,
			CreatedAt:    containerCreationTime,
			PortMappings: containerInfo.MappedPorts,
		})
	}

	result.ContainerInfos = containerInfos
	return command.Printer.PrintContainers(result)
}
