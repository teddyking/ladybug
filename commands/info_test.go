package commands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/commands"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/gardenfakes"
	"github.com/teddyking/ladybug/print/printfakes"
)

var _ = Describe("Info", func() {
	var (
		fakeGardenClient gardenfakes.FakeClient
		infoCommand      *Info
		fakePrinter      printfakes.FakePrinter
	)

	BeforeEach(func() {
		fakeGardenClient = gardenfakes.FakeClient{}
		fakePrinter = printfakes.FakePrinter{}

		infoCommand = &Info{
			Client:  &fakeGardenClient,
			Printer: &fakePrinter,
		}
	})

	Context("when garden reports 0 running containers", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{}, nil)
		})

		It("generates the correct InfoResult and prints it", func() {
			infoCommand.Execute(nil)

			Expect(fakePrinter.PrintInfoCallCount()).To(Equal(1))
			generatedResult := fakePrinter.PrintInfoArgsForCall(0)

			Expect(generatedResult.ContainersCount).To(Equal(0))
		})
	})

	Context("when garden reports 1 running container", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{nil}, nil)
		})

		It("generates the correct InfoResult and prints it", func() {
			infoCommand.Execute(nil)

			Expect(fakePrinter.PrintInfoCallCount()).To(Equal(1))
			generatedResult := fakePrinter.PrintInfoArgsForCall(0)

			Expect(generatedResult.ContainersCount).To(Equal(1))
		})
	})

	Context("when garden reports > 1 running container", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns([]garden.Container{nil, nil}, nil)
		})

		It("generates the correct InfoResult and prints it", func() {
			infoCommand.Execute(nil)

			Expect(fakePrinter.PrintInfoCallCount()).To(Equal(1))
			generatedResult := fakePrinter.PrintInfoArgsForCall(0)

			Expect(generatedResult.ContainersCount).To(Equal(2))
		})
	})

	Context("there is an error retrieving containers", func() {
		BeforeEach(func() {
			fakeGardenClient.ContainersReturns(nil, errors.New("error-getting-containers"))
		})

		It("returns the error", func() {
			err := infoCommand.Execute(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error-getting-containers"))
		})
	})

	Context("when there is an error printing the result", func() {
		BeforeEach(func() {
			fakePrinter.PrintInfoReturns(errors.New("error-printing-result"))
		})

		It("returns the error", func() {
			err := infoCommand.Execute(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error-printing-result"))
		})
	})
})
