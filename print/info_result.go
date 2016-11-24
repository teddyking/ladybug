package print

import (
	"fmt"

	"github.com/teddyking/ladybug/result"
)

func (r *ResultPrinter) PrintInfo(infoResult result.InfoResult) error {
	_, err := r.Out.Write([]byte(fmt.Sprintf("Running containers: %d\n", infoResult.ContainersCount)))
	return err
}
