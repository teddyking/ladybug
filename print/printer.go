package print

import "io"

// go:generate counterfeiter . Printer
type Printer interface {
	PrintContainers(result ContainersResult) error
	PrintInfo(result InfoResult) error
}

type ResultPrinter struct {
	Out io.Writer
}

func NewResultPrinter(out io.Writer) *ResultPrinter {
	return &ResultPrinter{
		Out: out,
	}
}
