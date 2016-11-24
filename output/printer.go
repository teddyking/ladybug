package output

import (
	"io"

	"github.com/teddyking/ladybug/result"
)

// go:generate counterfeiter . Printer
type Printer interface {
	PrintContainers(containersResult result.Containers) error
	PrintInfo(infoResult result.Info) error
}

type ResultPrinter struct {
	Out io.Writer
}

func NewResultPrinter(out io.Writer) *ResultPrinter {
	return &ResultPrinter{
		Out: out,
	}
}
