package result

import (
	"code.cloudfoundry.org/garden"
	"github.com/teddyking/ladybug/sys"
)

type Containers map[string]CInfo

type CInfo struct {
	Ip           string
	ProcessName  string
	CreatedAt    string
	PortMappings []garden.PortMapping
}

type resultModifier func(c Containers) (Containers, error)

func WithHandles(gdnContainers []garden.Container) resultModifier {
	return func(c Containers) (Containers, error) {
		for _, gdnContainer := range gdnContainers {
			c[gdnContainer.Handle()] = CInfo{}
		}

		return c, nil
	}
}

func WithIPs(gdnContainers []garden.Container) resultModifier {
	return func(c Containers) (Containers, error) {
		for _, gdnContainer := range gdnContainers {
			handle := gdnContainer.Handle()
			containerInfo, err := gdnContainer.Info()
			if err != nil {
				return nil, err
			}

			currentCInfo := c[handle]
			currentCInfo.Ip = containerInfo.ContainerIP
			c[handle] = currentCInfo
		}

		return c, nil
	}
}

func WithProcessNames(gdnContainers []garden.Container, host sys.Host) resultModifier {
	return func(c Containers) (Containers, error) {
		for _, gdnContainer := range gdnContainers {
			handle := gdnContainer.Handle()

			containerPids, err := host.ContainerPids(handle)
			if err != nil {
				return nil, err
			}

			containerProcessName := "N/A"
			if len(containerPids) > 0 {
				var err error
				containerProcessName, err = host.ContainerProcessName(containerPids[0])
				if err != nil {
					return nil, err
				}
			}

			currentCInfo := c[handle]
			currentCInfo.ProcessName = containerProcessName
			c[handle] = currentCInfo
		}

		return c, nil
	}
}

func WithCreatedAtTimes(gdnContainers []garden.Container, host sys.Host) resultModifier {
	return func(c Containers) (Containers, error) {
		for _, gdnContainer := range gdnContainers {
			handle := gdnContainer.Handle()

			containerCreationTime, err := host.ContainerCreationTime(handle)
			if err != nil {
				return nil, err
			}

			currentCInfo := c[handle]
			currentCInfo.CreatedAt = containerCreationTime
			c[handle] = currentCInfo
		}

		return c, nil
	}
}

func WithPortMappings(gdnContainers []garden.Container) resultModifier {
	return func(c Containers) (Containers, error) {
		for _, gdnContainer := range gdnContainers {
			handle := gdnContainer.Handle()
			containerInfo, err := gdnContainer.Info()
			if err != nil {
				return nil, err
			}

			currentCInfo := c[handle]
			currentCInfo.PortMappings = containerInfo.MappedPorts
			c[handle] = currentCInfo
		}

		return c, nil
	}
}

func (c Containers) Generate(modifiers ...resultModifier) error {
	for _, mod := range modifiers {
		result, err := mod(c)
		if err != nil {
			return err
		}
		c = result
	}

	return nil
}
