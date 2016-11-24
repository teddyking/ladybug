package sys_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/sys"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("Host", func() {
	var (
		fakeDepotDir string
		fakeProc     string
		fakeRunDir   string
		linuxHost    *LinuxHost
	)

	BeforeEach(func() {
		var err error

		fakeDepotDir, err = ioutil.TempDir("", "fake-depot-dir")
		Expect(err).NotTo(HaveOccurred())
		fakeProc, err = ioutil.TempDir("", "fake-proc")
		Expect(err).NotTo(HaveOccurred())
		fakeRunDir, err = ioutil.TempDir("", "fake-run-dir")
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		linuxHost = &LinuxHost{DepotDir: fakeDepotDir, Proc: fakeProc, RunDir: fakeRunDir}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(fakeDepotDir)).To(Succeed())
		Expect(os.RemoveAll(fakeProc)).To(Succeed())
		Expect(os.RemoveAll(fakeRunDir)).To(Succeed())
	})

	Describe("ContainerPids", func() {
		Context("when the container has 0 running processes", func() {
			BeforeEach(func() {
				bundlePath := filepath.Join(fakeDepotDir, "test-container")

				Expect(os.MkdirAll(bundlePath, 0755)).To(Succeed())
			})

			It("doesn't return any pids", func() {
				pids, err := linuxHost.ContainerPids("test-container")
				Expect(err).NotTo(HaveOccurred())

				Expect(len(pids)).To(Equal(0))
			})
		})

		Context("when the container has 1 running process", func() {
			BeforeEach(func() {
				createFakeContainerProcess(fakeDepotDir, "test-container", "fake-process-id", "100")
			})

			It("returns the pid of the process", func() {
				pids, err := linuxHost.ContainerPids("test-container")
				Expect(err).NotTo(HaveOccurred())

				Expect(len(pids)).To(Equal(1))
				Expect(pids[0]).To(Equal("100"))
			})
		})

		Context("when the container has > 1 running process", func() {
			BeforeEach(func() {
				createFakeContainerProcess(fakeDepotDir, "test-container", "fake-process-id", "100")
				createFakeContainerProcess(fakeDepotDir, "test-container", "fake-process-id2", "101")
			})

			It("returns the pid of the process", func() {
				pids, err := linuxHost.ContainerPids("test-container")
				Expect(err).NotTo(HaveOccurred())

				Expect(len(pids)).To(Equal(2))
				Expect(pids[0]).To(Equal("100"))
				Expect(pids[1]).To(Equal("101"))
			})
		})

		Context("when the depot dir doesn't exist", func() {
			BeforeEach(func() {
				fakeDepotDir = "/does/not/exist"
			})

			It("returns a meaningful error", func() {
				_, err := linuxHost.ContainerPids("test-container")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(Equal("Depot directory at '/does/not/exist' not found"))
			})
		})

		Context("when the container doesn't exist", func() {
			It("returns a meaningful error", func() {
				_, err := linuxHost.ContainerPids("test-container")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(Equal("Container with handle 'test-container' not found"))
			})
		})

		Context("when one or more pidfiles doesn't exist", func() {
			BeforeEach(func() {
				processPath := filepath.Join(fakeDepotDir, "test-container", "processes", "fake-process-id")

				Expect(os.MkdirAll(processPath, 0755)).To(Succeed())
			})

			It("returns a meaningful error", func() {
				_, err := linuxHost.ContainerPids("test-container")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(Equal("One or more pidfiles are missing"))
			})
		})
	})

	Describe("ContainerProcessName", func() {
		BeforeEach(func() {
			processProcPath := filepath.Join(fakeProc, "100")
			Expect(os.MkdirAll(processProcPath, 0755)).To(Succeed())
			processStatusfilePath := filepath.Join(processProcPath, "status")
			Expect(ioutil.WriteFile(processStatusfilePath, []byte("Name:\tfake-process\n"), 0644)).To(Succeed())
		})

		It("returns the name of the process identified by the pid", func() {
			processName, err := linuxHost.ContainerProcessName("100")
			Expect(err).NotTo(HaveOccurred())

			Expect(processName).To(Equal("fake-process"))
		})

		Context("when the process's status file doesn't exist", func() {
			It("returns a meaningful error", func() {
				_, err := linuxHost.ContainerProcessName("0")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(Equal(fmt.Sprintf("Unable to open %s/0/status", fakeProc)))
			})
		})

		Context("when the process Name isn't listed in the status file", func() {
			BeforeEach(func() {
				processProcPath := filepath.Join(fakeProc, "100")
				Expect(os.MkdirAll(processProcPath, 0755)).To(Succeed())
				processStatusfilePath := filepath.Join(processProcPath, "status")
				Expect(ioutil.WriteFile(processStatusfilePath, []byte("Cake: does-not-exist\n"), 0644)).To(Succeed())
			})

			It("returns a name of 'N/A'", func() {
				processName, err := linuxHost.ContainerProcessName("100")
				Expect(err).NotTo(HaveOccurred())

				Expect(processName).To(Equal("N/A"))
			})
		})
	})

	Describe("ContainerCreationTime", func() {
		var statefilePath string

		BeforeEach(func() {
			containerRunPath := filepath.Join(fakeRunDir, "test-container")
			Expect(os.MkdirAll(containerRunPath, 0755)).To(Succeed())
			statefilePath = filepath.Join(containerRunPath, "state.json")
			Expect(ioutil.WriteFile(statefilePath, []byte(`{"created":"2016-11-12T18:24:23.744239181Z"}`), 0644)).To(Succeed())
		})

		It("returns the time at which the container was created", func() {
			createdAt, err := (linuxHost.ContainerCreationTime("test-container"))
			Expect(err).NotTo(HaveOccurred())

			Expect(createdAt).To(Equal("2016-11-12T18:24:23.744239181Z"))
		})

		Context("when the statefile doesn't contain the created time", func() {
			BeforeEach(func() {
				Expect(ioutil.WriteFile(statefilePath, []byte(`{"notcreated":"2016-11-12T18:24:23.744239181Z"}`), 0644)).To(Succeed())
			})

			It("returns a time of 'N/A'", func() {
				createdAt, err := linuxHost.ContainerCreationTime("test-container")
				Expect(err).NotTo(HaveOccurred())

				Expect(createdAt).To(Equal("N/A"))
			})
		})

		Context("when there's an error decoding the JSON", func() {
			BeforeEach(func() {
				Expect(ioutil.WriteFile(statefilePath, []byte("this isn't JSON"), 0644)).To(Succeed())
			})

			It("returns the error", func() {
				_, err := linuxHost.ContainerCreationTime("test-container")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the statefile doesn't exist", func() {
			It("returns a meaningful error", func() {
				_, err := linuxHost.ContainerCreationTime("container-not-here")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(Equal(fmt.Sprintf("Unable to open %s/container-not-here/state.json", fakeRunDir)))
			})
		})
	})
})

func createFakeContainerProcess(depotDir, handle, processId, pid string) {
	processPath := filepath.Join(depotDir, handle, "processes", processId)
	pidfilePath := filepath.Join(processPath, "pidfile")
	Expect(os.MkdirAll(processPath, 0755)).To(Succeed())
	Expect(ioutil.WriteFile(pidfilePath, []byte(pid), 0644)).To(Succeed())
}
