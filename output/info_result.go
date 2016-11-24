package output

import (
	"fmt"

	"github.com/teddyking/ladybug/result"
)

func (r *ResultPrinter) PrintInfo(infoResult result.Info) error {
	_, err := r.Out.Write([]byte(fmt.Sprintf("Running containers: %d\n", infoResult.ContainersCount)))
	return err
}
