package print

import "fmt"

type InfoResult struct {
	ContainersCount int
}

func (r *ResultPrinter) PrintInfo(result InfoResult) {
	r.Out.Write([]byte(fmt.Sprintf("Running containers: %d\n", result.ContainersCount)))
}
