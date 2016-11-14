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

var _ = Describe("the ladybug containers command", func() {
	var (
		args        []string
		session     *gexec.Session
		stdout      *gbytes.Buffer
		container   garden.Container
		currentDate string
	)

	BeforeEach(func() {
		var err error

		args = []string{"-d", "/tmp/dir/depot", "containers"}
		stdout = gbytes.NewBuffer()

		currentDate = strings.Fields(time.Now().String())[0]

		container, err = gardenClient.Create(garden.ContainerSpec{Handle: "containers-container"})
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		var err error

		command := exec.Command(pathToLadybug, args...)
		session, err = gexec.Start(command, stdout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := gardenClient.Destroy("containers-container")
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns a 0 exit code", func() {
		Eventually(session).Should(gexec.Exit(0))
	})

	It("prints container handles to stdout", func() {
		Eventually(stdout).Should(gbytes.Say("containers-container"))
	})

	It("prints container creation times to stdout", func() {
		// actually prints the creation time but allow testing that
		Eventually(stdout).Should(gbytes.Say(currentDate))
	})

	Context("when containers are running a processes (other than /proc/self/exe)", func() {
		BeforeEach(func() {
			_, err := container.Run(garden.ProcessSpec{
				Path: "sleep",
				Args: []string{"100"},
			}, garden.ProcessIO{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("prints process names to stdout", func() {
			Eventually(stdout).Should(gbytes.Say("sleep"))
		})
	})

	Context("when containers are not running a process (other that /proc/self/exe)", func() {
		It("prints N/A to stdout", func() {
			Eventually(stdout).Should(gbytes.Say("N/A"))
		})
	})

	Context("when containers have port mappings", func() {
		BeforeEach(func() {
			_, _, err := container.NetIn(1122, 1122)
			Expect(err).NotTo(HaveOccurred())
		})

		It("prints port mappings to stdout", func() {
			Eventually(stdout).Should(gbytes.Say("1122->1122"))
		})
	})

	Context("when containers don't have port mappings", func() {
		It("prints N/A to stdout", func() {
			Eventually(stdout).Should(gbytes.Say("N/A"))
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
