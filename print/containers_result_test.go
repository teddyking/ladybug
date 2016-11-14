package print_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/print"

	"errors"
	"io"

	"code.cloudfoundry.org/garden"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("PrintContainers", func() {
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
			"with 1 ContainerInfos, 0 PortMappings",
			ContainersResult{
				ContainerInfos: []ContainerInfo{
					ContainerInfo{
						Handle:       "test-handle",
						Ip:           "192.0.2.10",
						ProcessName:  "ruby",
						CreatedAt:    "12-11-2016 23:00:00",
						PortMappings: []garden.PortMapping{},
					},
				},
			},
			"test-handle  192.0.2.10  ruby  12-11-2016 23:00:00  N/A\n",
		),

		Entry(
			"with 1 ContainerInfos, 1 PortMapping",
			ContainersResult{
				ContainerInfos: []ContainerInfo{
					ContainerInfo{
						Handle:       "test-handle",
						Ip:           "192.0.2.10",
						ProcessName:  "ruby",
						CreatedAt:    "12-11-2016 23:00:00",
						PortMappings: []garden.PortMapping{{80, 8080}},
					},
				},
			},
			"test-handle  192.0.2.10  ruby  12-11-2016 23:00:00  80->8080\n",
		),

		Entry(
			"with 1 ContainerInfos and 2 PortMappings",
			ContainersResult{
				ContainerInfos: []ContainerInfo{
					ContainerInfo{
						Handle:       "test-handle",
						Ip:           "192.0.2.10",
						ProcessName:  "ruby",
						CreatedAt:    "12-11-2016 23:00:00",
						PortMappings: []garden.PortMapping{{80, 8080}, {443, 4443}},
					},
				},
			},
			"test-handle  192.0.2.10  ruby  12-11-2016 23:00:00  80->8080, 443->4443\n",
		),

		Entry(
			"with 2 ContainerInfos, each with 0 PortMappings",
			ContainersResult{
				ContainerInfos: []ContainerInfo{
					ContainerInfo{
						Handle:      "test-handle",
						Ip:          "192.0.2.10",
						ProcessName: "ruby",
						CreatedAt:   "12-11-2016 23:00:00",
					},
					ContainerInfo{
						Handle:      "test-handle-2",
						Ip:          "192.0.2.11",
						ProcessName: "tree",
						CreatedAt:   "12-11-2016 23:00:01",
					},
				},
			},
			"test-handle    192.0.2.10  ruby  12-11-2016 23:00:00  N/A\ntest-handle-2  192.0.2.11  tree  12-11-2016 23:00:01  N/A\n",
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

type erroringWriter struct{}

func (e erroringWriter) Write(p []byte) (int, error) {
	return 0, errors.New("error-writing-to-writer")
}
