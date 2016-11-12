package print

import (
	"io"

	"github.com/concourse/fly/ui"
)

type Printer interface {
	PrintContainers(result ContainersResult) error
}

type ResultPrinter struct {
	Out io.Writer
}

func NewResultPrinter(out io.Writer) *ResultPrinter {
	return &ResultPrinter{
		Out: out,
	}
}

type ContainerInfo struct {
	Handle      string
	Ip          string
	ProcessName string
}

type ContainersResult struct {
	ContainerInfos []ContainerInfo
}

func (r *ResultPrinter) PrintContainers(result ContainersResult) error {
	table := ui.Table{
		Headers: ui.TableRow{
			{Contents: "Handle"},
			{Contents: "IP Address"},
			{Contents: "Process Name"},
		},
	}

	for _, containerInfo := range result.ContainerInfos {
		row := ui.TableRow{
			{Contents: containerInfo.Handle},
			{Contents: containerInfo.Ip},
			{Contents: containerInfo.ProcessName},
		}

		table.Data = append(table.Data, row)
	}

	return table.Render(r.Out)
}
