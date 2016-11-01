package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"

	"code.cloudfoundry.org/garden"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration", func() {
	var (
		args    []string
		session *gexec.Session
		stdout  *gbytes.Buffer
	)

	BeforeEach(func() {
		args = []string{}
		stdout = gbytes.NewBuffer()
	})

	JustBeforeEach(func() {
		var err error

		command := exec.Command(pathToLadybug, args...)
		session, err = gexec.Start(command, stdout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when run without options or arguments", func() {
		It("returns a 0 exit code", func() {
			Eventually(session).Should(gexec.Exit(0))
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
})
