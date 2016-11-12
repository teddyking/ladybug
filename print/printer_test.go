package print_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/print"

	"errors"
	"io"

	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Resultprinter", func() {
	Describe("PrintContainers", func() {
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

		DescribeTable("formatted output",
			func(result ContainersResult, expectedOutput string) {
				Expect(resultPrinter.PrintContainers(result)).To(Succeed())

				Eventually(stdout).Should(gbytes.Say(expectedOutput))
			},

			Entry(
				"with 1 ContainerInfos",
				ContainersResult{
					ContainerInfos: []ContainerInfo{
						ContainerInfo{
							Handle:      "test-handle",
							Ip:          "192.0.2.10",
							ProcessName: "ruby",
						},
					},
				},
				"test-handle  192.0.2.10  ruby\n",
			),

			Entry(
				"with 2 ContainerInfos",
				ContainersResult{
					ContainerInfos: []ContainerInfo{
						ContainerInfo{
							Handle:      "test-handle",
							Ip:          "192.0.2.10",
							ProcessName: "ruby",
						},
						ContainerInfo{
							Handle:      "test-handle-2",
							Ip:          "192.0.2.11",
							ProcessName: "tree",
						},
					},
				},
				"test-handle    192.0.2.10  ruby\ntest-handle-2  192.0.2.11  tree\n",
			),
		)

		Context("when the table Render returns an error", func() {
			BeforeEach(func() {
				stdout = erroringWriter{}
			})

			It("returns the error", func() {
				result := ContainersResult{
					ContainerInfos: []ContainerInfo{
						ContainerInfo{
							Handle:      "test-handle",
							Ip:          "192.0.2.10",
							ProcessName: "ruby",
						},
					},
				}

				err := resultPrinter.PrintContainers(result)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error-writing-to-writer"))
			})
		})
	})
})

type erroringWriter struct{}

func (e erroringWriter) Write(p []byte) (int, error) {
	return 0, errors.New("error-writing-to-writer")
}
