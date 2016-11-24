package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/commands"

	"errors"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/gardenfakes"
	"github.com/teddyking/ladybug/output/outputfakes"
	"github.com/teddyking/ladybug/result"
	"github.com/teddyking/ladybug/sys/sysfakes"
)

var _ = Describe("Containers", func() {
	var (
		fakeGardenClient  gardenfakes.FakeClient
		fakeHost          sysfakes.FakeHost
		fakePrinter       outputfakes.FakePrinter
		containersCommand *Containers
	)

	BeforeEach(func() {
		fakeGardenClient = gardenfakes.FakeClient{}
		fakeHost = sysfakes.FakeHost{}
		fakePrinter = outputfakes.FakePrinter{}
	})

	JustBeforeEach(func() {
		containersCommand = &Containers{
			Client:  &fakeGardenClient,
			Host:    &fakeHost,
			Printer: &fakePrinter,
		}
	})

	Describe("Execute", func() {
		It("generates a result.Containers and sends it to be printed", func() {
			Expect(containersCommand.Execute(nil)).To(Succeed())

			Expect(fakePrinter.PrintContainersCallCount()).To(Equal(1))
			Expect(fakePrinter.PrintContainersArgsForCall(0)).To(Equal(result.Containers{}))
		})

		Context("when the garden Client returns an error", func() {
			BeforeEach(func() {
				fakeGardenClient.ContainersReturns(nil, errors.New("error-fetching-containers"))
			})

			It("returns the error", func() {
				err := containersCommand.Execute(nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error-fetching-containers"))
			})
		})

		Context("when one of the With* funcs returns an error", func() {
			BeforeEach(func() {
				fakeContainer := &gardenfakes.FakeContainer{}
				fakeGardenClient.ContainersReturns([]garden.Container{fakeContainer}, nil)
				fakeHost.ContainerCreationTimeReturns("", errors.New("error-fetching-created-at-time"))
			})

			It("returns the error", func() {
				err := containersCommand.Execute(nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error-fetching-created-at-time"))
			})
		})

		Context("when the Printer returns an error", func() {
			BeforeEach(func() {
				fakePrinter.PrintContainersReturns(errors.New("error-printing-result"))
			})

			It("returns the error", func() {
				err := containersCommand.Execute(nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error-printing-result"))
			})
		})
	})
})
