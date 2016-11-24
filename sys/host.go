package sys

// go:generate counterfeiter . Host
type Host interface {
	ContainerPids(handle string) ([]string, error)
	ContainerProcessName(pid string) (string, error)
	ContainerCreationTime(handle string) (string, error)
}
