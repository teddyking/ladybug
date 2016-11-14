package exterminator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("ladybug usage", func() {
	var (
		args    []string
		session *gexec.Session
		stdout  *gbytes.Buffer
	)

	BeforeEach(func() {
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
})
