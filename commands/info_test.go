package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/commands"

	"errors"

	"code.cloudfoundry.org/garden/gardenfakes"
	"github.com/teddyking/ladybug/print/printfakes"
	"github.com/teddyking/ladybug/result"
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

	Describe("Execute", func() {
		It("generates a result.Info and sends it to be printed", func() {
			Expect(infoCommand.Execute(nil)).To(Succeed())

			Expect(fakePrinter.PrintInfoCallCount()).To(Equal(1))
			Expect(fakePrinter.PrintInfoArgsForCall(0)).To(Equal(result.Info{}))
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
})
