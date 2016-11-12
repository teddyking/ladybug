package system

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type LinuxHost struct {
	DepotDir string
	Proc     string
	RunDir   string
}

type RuncContainerState struct {
	Created string `json:"created"`
}

func (lh *LinuxHost) ContainerPids(handle string) ([]string, error) {
	bundleDir := filepath.Join(lh.DepotDir, handle)

	processesPathsPattern := filepath.Join(lh.DepotDir, handle, "processes", "*")
	pidfilePathsPattern := filepath.Join(lh.DepotDir, handle, "processes", "*", "pidfile")

	processesPaths, _ := filepath.Glob(processesPathsPattern)
	pidfilePaths, _ := filepath.Glob(pidfilePathsPattern)

	if _, err := os.Stat(lh.DepotDir); os.IsNotExist(err) {
		// depot dir not found
		return nil, fmt.Errorf("Depot directory at '%s' not found", lh.DepotDir)
	}

	if _, err := os.Stat(bundleDir); os.IsNotExist(err) {
		// bundle dir not found - container doesn't exist
		return nil, fmt.Errorf("Container with handle '%s' not found", handle)
	}

	if len(processesPaths) != len(pidfilePaths) {
		// a pidfile is missing from one of a container's process dirs
		return nil, fmt.Errorf("One or more pidfiles are missing")
	}

	var pids []string
	for _, pidfilePath := range pidfilePaths {
		pid, err := ioutil.ReadFile(pidfilePath)
		if err != nil {
			return nil, err
		}

		pids = append(pids, string(pid))
	}

	return pids, nil
}

func (lh *LinuxHost) ContainerProcessName(pid string) (string, error) {
	statusfilePath := filepath.Join(lh.Proc, pid, "status")
	statusfile, err := os.Open(statusfilePath)
	defer statusfile.Close()
	if err != nil {
		return "", fmt.Errorf("Unable to open %s", statusfilePath)
	}
	scanner := bufio.NewScanner(statusfile)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Name:") {
			return strings.Fields(line)[1], nil
		}
	}

	return "N/A", nil
}

func (lh *LinuxHost) ContainerCreationTime(handle string) (string, error) {
	statefilePath := filepath.Join(lh.RunDir, handle, "state.json")
	statefile, err := os.Open(statefilePath)
	defer statefile.Close()
	if err != nil {
		return "", fmt.Errorf("Unable to open %s", statefilePath)
	}

	var state RuncContainerState
	if err := json.NewDecoder(statefile).Decode(&state); err != nil {
		return "", err
	}

	createdAt := state.Created
	if createdAt == "" {
		createdAt = "N/A"
	}

	return createdAt, nil
}
