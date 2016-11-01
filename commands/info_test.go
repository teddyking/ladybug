package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/commands"

	"errors"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/gardenfakes"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Info", func() {
	var (
		fakeGardenClient gardenfakes.FakeClient
		infoCommand      *Info
		stdout           *gbytes.Buffer
		stderr           *gbytes.Buffer
	)

	BeforeEach(func() {
		fakeGardenClient = gardenfakes.FakeClient{}
		stdout = gbytes.NewBuffer()
		stderr = gbytes.NewBuffer()

		infoCommand = &Info{
			Client: &fakeGardenClient,
			Out:    stdout,
			Err:    stderr,
		}
	})

	Context("when garden reports 0 running containers", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{}, nil)
		})

		It("prints 0 running containers to stdout", func() {
			infoCommand.Execute(nil)

			Expect(stdout).To(gbytes.Say("Running containers: 0\n"))
		})
	})

	Context("when garden reports 1 running container", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{nil}, nil)
		})

		It("prints 1 running containers to stdout", func() {
			infoCommand.Execute(nil)

			Expect(stdout).To(gbytes.Say("Running containers: 1\n"))
		})
	})

	Context("when garden reports > 1 running containers", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{nil, nil}, nil)
		})

		It("prints > 1 running containers to stdout", func() {
			infoCommand.Execute(nil)

			Expect(stdout).To(gbytes.Say("Running containers: 2\n"))
		})
	})

	Context("when garden returns an error", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns(nil, errors.New("error-getting-containers"))
		})

		It("prints a useful error message to stderr", func() {
			infoCommand.Execute(nil)

			Expect(fakeGardenClient.ContainersCallCount()).To(Equal(1))
			Expect(fakeGardenClient.ContainersArgsForCall(0)).To(Equal(garden.Properties{}))

			Expect(stderr).To(gbytes.Say("Garden returned an error - error-getting-containers\n"))
		})

		It("returns the error", func() {
			Expect(infoCommand.Execute(nil)).NotTo(Succeed())
		})
	})
})
