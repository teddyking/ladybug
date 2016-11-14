package exterminator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"
	"strings"
	"time"

	"code.cloudfoundry.org/garden"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("ladybug", func() {
	var (
		args    []string
		session *gexec.Session
		stdout  *gbytes.Buffer
	)

	BeforeEach(func() {
		// depot dir gets created in ci/scripts/test
		// and is set to /tmp/dir/depot
		args = []string{"-d", "/tmp/dir/depot"}
		stdout = gbytes.NewBuffer()
	})

	JustBeforeEach(func() {
		var err error

		command := exec.Command(pathToLadybug, args...)
		session, err = gexec.Start(command, stdout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when run without options or arguments", func() {
		It("aks the user to specify a command to run", func() {
			Eventually(stdout).Should(gbytes.Say("Please specify one command of"))
		})

		It("returns a 1 exit code", func() {
			Eventually(session).Should(gexec.Exit(1))
		})
	})

	Context("when run with -h/--help", func() {
		BeforeEach(func() {
			args = []string{"-h"}
		})

		It("prints usage info to stdout", func() {
			Eventually(stdout).Should(gbytes.Say("Usage:"))
		})

		It("returns a 0 exit code", func() {
			Eventually(session).Should(gexec.Exit(0))
		})
	})

	Context("when run with info", func() {
		BeforeEach(func() {
			args = []string{"info"}

			_, err := gardenClient.Create(garden.ContainerSpec{Handle: "info-container"})
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := gardenClient.Destroy("info-container")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a 0 exit code", func() {
			Eventually(session).Should(gexec.Exit(0))
		})

		It("prints the number of running containers to stdout", func() {
			Eventually(stdout).Should(gbytes.Say("Running containers: 1"))
		})
	})

	Context("when run with containers", func() {
		var (
			container   garden.Container
			currentDate string
		)

		BeforeEach(func() {
			// override the default depot dir path as the garden in the test
			// uses /tmp/dir/depot and not /var/vcap/data/garden/depot
			args = []string{"-d", "/tmp/dir/depot", "containers"}

			var err error
			currentDate = strings.Fields(time.Now().String())[0]
			container, err = gardenClient.Create(garden.ContainerSpec{Handle: "containers-container"})
			Expect(err).NotTo(HaveOccurred())

			_, err = container.Run(garden.ProcessSpec{
				Path: "sleep",
				Args: []string{"100"},
			}, garden.ProcessIO{})
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := gardenClient.Destroy("containers-container")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a 0 exit code", func() {
			Eventually(session).Should(gexec.Exit(0))
		})

		It("prints the containers' handles to stdout", func() {
			Eventually(stdout).Should(gbytes.Say("containers-container"))
		})

		It("prints the containers' creation times to stdout", func() {
			// actually prints the creation time but allow testing that
			Eventually(stdout).Should(gbytes.Say(currentDate))
		})

		Context("when one or more of the containers has one or more processes", func() {
			It("prints the process' names to stdout", func() {
				Eventually(stdout).Should(gbytes.Say("sleep"))
			})
		})

		Context("when one or more of the containers has one or more port mappings", func() {
			BeforeEach(func() {
				_, _, err := container.NetIn(1122, 1122)
				Expect(err).NotTo(HaveOccurred())
			})

			It("prints the port mappings to stdout", func() {
				Eventually(stdout).Should(gbytes.Say("1122->1122"))
			})
		})

		Context("when the container doesn't have any processes (other that init)", func() {
			It("doesn't error", func() {
				Eventually(session).Should(gexec.Exit(0))
			})
		})

		Context("when the depot dir does not exist", func() {
			BeforeEach(func() {
				args = []string{"-d", "/does/not/exist", "containers"}
			})

			It("prints a meaningful error message to stderr", func() {
				Eventually(stdout).Should(gbytes.Say("Depot directory at '/does/not/exist' not found"))
			})

			It("returns a 1 exit code", func() {
				Eventually(session).Should(gexec.Exit(1))
			})
		})
	})
})
