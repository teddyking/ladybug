package exterminator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"

	"code.cloudfoundry.org/garden"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("the ladybug info command", func() {
	var (
		session *gexec.Session
		stdout  *gbytes.Buffer
	)

	BeforeEach(func() {
		var err error

		args := []string{"-d", "/tmp/dir/depot", "info"}
		stdout = gbytes.NewBuffer()

		_, err = gardenClient.Create(garden.ContainerSpec{Handle: "info-container"})
		Expect(err).NotTo(HaveOccurred())

		command := exec.Command(pathToLadybug, args...)
		session, err = gexec.Start(command, stdout, GinkgoWriter)
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
