package print

import (
	"fmt"
	"strings"

	"code.cloudfoundry.org/garden"
	"github.com/concourse/fly/ui"
)

type ContainersResult struct {
	ContainerInfos []ContainerInfo
}

type ContainerInfo struct {
	Handle       string
	Ip           string
	ProcessName  string
	CreatedAt    string
	PortMappings []garden.PortMapping
}

func (r *ResultPrinter) PrintContainers(result ContainersResult) error {
	table := ui.Table{
		Headers: ui.TableRow{
			{Contents: "Handle"},
			{Contents: "IP Address"},
			{Contents: "Process Name"},
			{Contents: "Created At"},
			{Contents: "Port Mappings"},
		},
	}

	for _, containerInfo := range result.ContainerInfos {
		row := ui.TableRow{
			{Contents: containerInfo.Handle},
			{Contents: containerInfo.Ip},
			{Contents: containerInfo.ProcessName},
			{Contents: trimTime(containerInfo.CreatedAt)},
		}

		var mappedPortsResult string
		if len(containerInfo.PortMappings) > 0 {
			for _, portMapping := range containerInfo.PortMappings {
				mappedPortsResult = fmt.Sprintf("%s%d->%d, ", mappedPortsResult, portMapping.HostPort, portMapping.ContainerPort)
			}
		}
		if mappedPortsResult == "" {
			mappedPortsResult = "N/A"
		}
		row = append(row, ui.TableCell{Contents: strings.Trim(mappedPortsResult, ", ")})

		table.Data = append(table.Data, row)
	}

	return table.Render(r.Out)
}

func trimTime(t string) string {
	// expects a string of the format 2016-11-15T06:48:15.137799416Z
	// and returns a string of the format 2016-11-15 06:48:15
	return strings.Replace(strings.Split(t, ".")[0], "T", " ", 1)
}