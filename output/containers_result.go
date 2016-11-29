package output

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"code.cloudfoundry.org/garden"
	"github.com/concourse/fly/ui"
	"github.com/teddyking/ladybug/result"
)

func (r *ResultPrinter) PrintContainers(containersResult result.Containers) error {
	table := ui.Table{
		Headers: ui.TableRow{
			{Contents: "Handle"},
			{Contents: "IP Address"},
			{Contents: "Process Name"},
			{Contents: "Created At"},
			{Contents: "Port Mappings"},
		},
	}

	containersToOutput := mapToSlice(containersResult)
	sort.Sort(containersByCreatedAt(containersToOutput))

	for _, c := range containersToOutput {
		row := ui.TableRow{
			{Contents: c.Handle},
			{Contents: c.Ip},
			{Contents: c.ProcessName},
			{Contents: trimTime(c.CreatedAt)},
		}

		var mappedPortsResult string
		if len(c.PortMappings) > 0 {
			for _, portMapping := range c.PortMappings {
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

type containerInfo struct {
	Handle       string
	Ip           string
	ProcessName  string
	CreatedAt    string
	PortMappings []garden.PortMapping
}

type containersByCreatedAt []containerInfo

func (c containersByCreatedAt) Len() int {
	return len(c)
}

func (c containersByCreatedAt) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c containersByCreatedAt) Less(i, j int) bool {
	layout := "2006-01-02T15:04:05.000000000Z"

	t1, err := time.Parse(layout, c[i].CreatedAt)
	if err != nil {
		panic(err)
	}
	t2, err := time.Parse(layout, c[j].CreatedAt)
	if err != nil {
		panic(err)
	}

	return t1.After(t2)
}

func trimTime(t string) string {
	// expects a string of the format 2016-11-15T06:48:15.137799416Z
	// and returns a string of the format 2016-11-15 06:48:15
	return strings.Replace(strings.Split(t, ".")[0], "T", " ", 1)
}

func mapToSlice(containersResult result.Containers) []containerInfo {
	var containers []containerInfo

	for handle, info := range containersResult {
		c := containerInfo{
			Handle:       handle,
			Ip:           info.Ip,
			ProcessName:  info.ProcessName,
			CreatedAt:    info.CreatedAt,
			PortMappings: info.PortMappings,
		}

		containers = append(containers, c)
	}

	return containers
}
