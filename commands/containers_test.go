package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/commands"

	"errors"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/gardenfakes"
	"github.com/onsi/gomega/gbytes"
	"github.com/teddyking/ladybug/system/systemfakes"
)

var _ = Describe("Containers", func() {
	var (
		fakeGardenClient  gardenfakes.FakeClient
		fakeHost          systemfakes.FakeHost
		containersCommand *Containers
		stdout            *gbytes.Buffer
	)

	BeforeEach(func() {
		fakeGardenClient = gardenfakes.FakeClient{}
		fakeHost = systemfakes.FakeHost{}
		stdout = gbytes.NewBuffer()

		containersCommand = &Containers{
			Client: &fakeGardenClient,
			Host:   &fakeHost,
			Out:    stdout,
		}
	})

	Context("when garden reports 0 running containers", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{}, nil)
		})

		It("prints a message saying that there aren't any containers to stdout", func() {
			containersCommand.Execute(nil)

			Expect(stdout).To(gbytes.Say("0 running containers found on this host\n"))
		})

		It("doesn't return an error", func() {
			Expect(containersCommand.Execute(nil)).To(Succeed())
		})
	})

	Context("when garden reports 1 running container", func() {
		var (
			fakeContainer *gardenfakes.FakeContainer
			fakePids      []string
		)

		BeforeEach(func() {
			fakeContainer = &gardenfakes.FakeContainer{}
			fakePids = []string{"100"}
		})

		JustBeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{fakeContainer}, nil)

			fakeContainer.HandleReturns("test-container")
			fakeContainer.InfoReturns(
				garden.ContainerInfo{
					ContainerIP: "192.0.2.10",
					ProcessIDs:  fakePids,
				},
				nil,
			)

			fakeHost.ContainerPidsReturns(fakePids, nil)
			fakeHost.ContainerProcessNameReturns("test-process", nil)
		})

		It("prints detailed info about the container to stdout", func() {
			containersCommand.Execute(nil)

			Expect(fakeHost.ContainerProcessNameArgsForCall(0)).To(Equal("100"))

			Expect(stdout).To(gbytes.Say("test-container"))
			Expect(stdout).To(gbytes.Say("192.0.2.10"))
			Expect(stdout).To(gbytes.Say("test-process"))
		})

		It("doesn't return an error", func() {
			Expect(containersCommand.Execute(nil)).To(Succeed())
		})

		Context("when there is an error retrieving ContainerInfo", func() {
			JustBeforeEach(func() {
				fakeContainer.InfoReturns(
					garden.ContainerInfo{},
					errors.New("error-retrieving-container-info"),
				)
			})

			It("returns the error", func() {
				Expect(containersCommand.Execute(nil)).NotTo(Succeed())
			})
		})

		Context("when there is an error retrieving ContainerPids", func() {
			JustBeforeEach(func() {
				fakeHost.ContainerPidsReturns(
					nil,
					errors.New("error-retrieving-container-pids"),
				)
			})

			It("returns the error", func() {
				Expect(containersCommand.Execute(nil)).NotTo(Succeed())
			})
		})

		Context("when there is an error retrieving ContainerProcessName", func() {
			JustBeforeEach(func() {
				fakeHost.ContainerProcessNameReturns(
					"",
					errors.New("error-retrieving-container-process-name"),
				)
			})

			It("returns the error", func() {
				Expect(containersCommand.Execute(nil)).NotTo(Succeed())
			})
		})

		Context("when the container isn't running any processes (other than init)", func() {
			BeforeEach(func() {
				fakePids = []string{}
			})

			It("prints N/A to stdout", func() {
				containersCommand.Execute(nil)

				Expect(stdout).To(gbytes.Say("test-container"))
				Expect(stdout).To(gbytes.Say("192.0.2.10"))
				Expect(stdout).To(gbytes.Say("N/A"))
			})
		})
	})

	Context("there is an error retrieving containers", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns(nil, errors.New("error-getting-containers"))
		})

		It("returns the error", func() {
			err := containersCommand.Execute(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error-getting-containers"))
		})
	})
})
