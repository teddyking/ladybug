package output_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/output"

	"io"

	"github.com/onsi/gomega/gbytes"
	"github.com/teddyking/ladybug/result"
)

var _ = Describe("PrintInfo", func() {
	var (
		stdout        io.Writer
		resultPrinter *ResultPrinter
	)

	BeforeEach(func() {
		stdout = gbytes.NewBuffer()
	})

	JustBeforeEach(func() {
		resultPrinter = NewResultPrinter(stdout)
	})

	It("prints the number of running containers to stdout", func() {
		infoResult := result.Info{
			ContainersCount: 3,
		}
		resultPrinter.PrintInfo(infoResult)

		Eventually(stdout).Should(gbytes.Say("Running containers: 3"))
	})

	Context("when there is an error writing to Out", func() {
		BeforeEach(func() {
			stdout = ErroringWriter{}
		})

		It("returns the error", func() {
			infoResult := result.Info{}
			err := resultPrinter.PrintInfo(infoResult)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("error-writing-to-writer"))
		})
	})
})
