package result

import "code.cloudfoundry.org/garden"

type InfoResult struct {
	ContainersCount int
}

func (r *InfoResult) Generate(containers []garden.Container) {
	r.ContainersCount = len(containers)
}
