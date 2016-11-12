package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/commands"

	"errors"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/gardenfakes"
	"github.com/teddyking/ladybug/print/printfakes"
	"github.com/teddyking/ladybug/system/systemfakes"
)

var _ = Describe("Containers", func() {
	var (
		fakeGardenClient  gardenfakes.FakeClient
		fakeHost          systemfakes.FakeHost
		fakePrinter       printfakes.FakePrinter
		fakeContainer     *gardenfakes.FakeContainer
		containersCommand *Containers
	)

	BeforeEach(func() {
		fakeGardenClient = gardenfakes.FakeClient{}
		fakeHost = systemfakes.FakeHost{}
		fakePrinter = printfakes.FakePrinter{}
		fakeContainer = &gardenfakes.FakeContainer{}

		containersCommand = &Containers{
			Client:  &fakeGardenClient,
			Host:    &fakeHost,
			Printer: &fakePrinter,
		}
	})

	Context("when garden reports 0 running containers", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{}, nil)
		})

		It("generates the correct ContainersResult and prints it", func() {
			containersCommand.Execute(nil)

			generatedResult := fakePrinter.PrintContainersArgsForCall(0)

			Expect(len(generatedResult.ContainerInfos)).To(Equal(0))
			Expect(fakePrinter.PrintContainersCallCount()).To(Equal(1))
		})

		It("doesn't return an error", func() {
			Expect(containersCommand.Execute(nil)).To(Succeed())
		})
	})

	Context("when garden reports 1 running container", func() {
		var (
			fakePids []string
		)

		BeforeEach(func() {
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

		It("generates the correct ContainersResult and prints it", func() {
			containersCommand.Execute(nil)

			generatedResult := fakePrinter.PrintContainersArgsForCall(0)

			Expect(len(generatedResult.ContainerInfos)).To(Equal(1))
			Expect(generatedResult.ContainerInfos[0].Handle).To(Equal("test-container"))
			Expect(generatedResult.ContainerInfos[0].Ip).To(Equal("192.0.2.10"))
			Expect(generatedResult.ContainerInfos[0].ProcessName).To(Equal("test-process"))
			Expect(fakePrinter.PrintContainersCallCount()).To(Equal(1))
		})

		It("doesn't return an error", func() {
			Expect(containersCommand.Execute(nil)).To(Succeed())
		})
	})

	Context("when garden reports 1 running container", func() {
		var (
			fakeContainer2 *gardenfakes.FakeContainer
		)

		BeforeEach(func() {
			fakeContainer2 = &gardenfakes.FakeContainer{}
		})

		JustBeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{fakeContainer, fakeContainer2}, nil)

			fakeContainer.HandleReturns("test-container")
			fakeContainer2.HandleReturns("test-container-2")
			fakeContainer.InfoReturns(
				garden.ContainerInfo{
					ContainerIP: "192.0.2.10",
				},
				nil,
			)
			fakeContainer2.InfoReturns(
				garden.ContainerInfo{
					ContainerIP: "192.0.2.11",
				},
				nil,
			)
		})

		It("generates the correct ContainersResult and prints it", func() {
			containersCommand.Execute(nil)

			generatedResult := fakePrinter.PrintContainersArgsForCall(0)

			Expect(len(generatedResult.ContainerInfos)).To(Equal(2))
			Expect(generatedResult.ContainerInfos[0].Handle).To(Equal("test-container"))
			Expect(generatedResult.ContainerInfos[0].Ip).To(Equal("192.0.2.10"))
			Expect(generatedResult.ContainerInfos[0].ProcessName).To(Equal("N/A"))
			Expect(generatedResult.ContainerInfos[1].Handle).To(Equal("test-container-2"))
			Expect(generatedResult.ContainerInfos[1].Ip).To(Equal("192.0.2.11"))
			Expect(generatedResult.ContainerInfos[1].ProcessName).To(Equal("N/A"))
			Expect(fakePrinter.PrintContainersCallCount()).To(Equal(1))
		})

		It("doesn't return an error", func() {
			Expect(containersCommand.Execute(nil)).To(Succeed())
		})
	})

	Context("when there is an error retrieving containers", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns(nil, errors.New("error-getting-containers"))
		})

		It("returns the error", func() {
			err := containersCommand.Execute(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error-getting-containers"))
		})
	})

	Context("when there is an error retrieving ContainerInfo", func() {
		JustBeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{fakeContainer}, nil)
			fakeContainer.InfoReturns(garden.ContainerInfo{}, errors.New("error-retrieving-container-info"))
		})

		It("returns the error", func() {
			err := containersCommand.Execute(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error-retrieving-container-info"))
		})
	})

	Context("when there is an error retrieving ContainerPids", func() {
		JustBeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{fakeContainer}, nil)
			fakeContainer.InfoReturns(garden.ContainerInfo{ProcessIDs: []string{"100"}}, nil)
			fakeHost.ContainerPidsReturns(nil, errors.New("error-retrieving-container-pids"))
		})

		It("returns the error", func() {
			err := containersCommand.Execute(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error-retrieving-container-pids"))
		})
	})

	Context("when there is an error retrieving ContainerProcessName", func() {
		JustBeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{fakeContainer}, nil)
			fakeContainer.InfoReturns(garden.ContainerInfo{ProcessIDs: []string{"100"}}, nil)
			fakeHost.ContainerPidsReturns([]string{"100"}, nil)
			fakeHost.ContainerProcessNameReturns("", errors.New("error-retrieving-container-process-name"))
		})

		It("returns the error", func() {
			err := containersCommand.Execute(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error-retrieving-container-process-name"))
		})
	})

	Context("when there is an error printing the result", func() {
		BeforeEach(func() {
			fakePrinter.PrintContainersReturns(errors.New("error-printing-result"))
		})

		It("returns the error", func() {
			err := containersCommand.Execute(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error-printing-result"))
		})
	})
})
