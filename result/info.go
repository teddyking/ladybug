package result

import "code.cloudfoundry.org/garden"

type Info struct {
	ContainersCount int
}

func (i *Info) Generate(containers []garden.Container) {
	i.ContainersCount = len(containers)
}
