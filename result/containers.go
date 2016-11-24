package result

import (
	"code.cloudfoundry.org/garden"
	"github.com/teddyking/ladybug/system"
)

type ContainersResult map[string]CInfo

type CInfo struct {
	Ip           string
	ProcessName  string
	CreatedAt    string
	PortMappings []garden.PortMapping
}

type resultModifier func(r ContainersResult) (ContainersResult, error)

func WithHandles(containers []garden.Container) resultModifier {
	return func(r ContainersResult) (ContainersResult, error) {
		for _, container := range containers {
			r[container.Handle()] = CInfo{}
		}

		return r, nil
	}
}

func WithIPs(containers []garden.Container) resultModifier {
	return func(r ContainersResult) (ContainersResult, error) {
		for _, container := range containers {
			handle := container.Handle()
			containerInfo, err := container.Info()
			if err != nil {
				return nil, err
			}

			currentCInfo := r[handle]
			currentCInfo.Ip = containerInfo.ContainerIP
			r[handle] = currentCInfo
		}

		return r, nil
	}
}

func WithProcessNames(containers []garden.Container, host system.Host) resultModifier {
	return func(r ContainersResult) (ContainersResult, error) {
		for _, container := range containers {
			handle := container.Handle()

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

			currentCInfo := r[handle]
			currentCInfo.ProcessName = containerProcessName
			r[handle] = currentCInfo
		}

		return r, nil
	}
}

func WithCreatedAtTimes(containers []garden.Container, host system.Host) resultModifier {
	return func(r ContainersResult) (ContainersResult, error) {
		for _, container := range containers {
			handle := container.Handle()

			containerCreationTime, err := host.ContainerCreationTime(handle)
			if err != nil {
				return nil, err
			}

			currentCInfo := r[handle]
			currentCInfo.CreatedAt = containerCreationTime
			r[handle] = currentCInfo
		}

		return r, nil
	}
}

func WithPortMappings(containers []garden.Container) resultModifier {
	return func(r ContainersResult) (ContainersResult, error) {
		for _, container := range containers {
			handle := container.Handle()
			containerInfo, err := container.Info()
			if err != nil {
				return nil, err
			}

			currentCInfo := r[handle]
			currentCInfo.PortMappings = containerInfo.MappedPorts
			r[handle] = currentCInfo
		}

		return r, nil
	}
}

func (r ContainersResult) Generate(modifiers ...resultModifier) error {
	for _, mod := range modifiers {
		result, err := mod(r)
		if err != nil {
			return err
		}
		r = result
	}

	return nil
}
