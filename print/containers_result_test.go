package print_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/print"

	"io"

	"code.cloudfoundry.org/garden"
	"github.com/onsi/gomega/gbytes"
	"github.com/teddyking/ladybug/result"
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
		func(containersResult result.Containers, expectedOutput string) {
			Expect(resultPrinter.PrintContainers(containersResult)).To(Succeed())

			Eventually(stdout).Should(gbytes.Say(expectedOutput))
		},

		Entry(
			"with 1 container",
			result.Containers{
				"test-handle": result.CInfo{
					Ip:           "192.0.2.10",
					ProcessName:  "ruby",
					CreatedAt:    "2016-11-15T06:48:15.137799416Z",
					PortMappings: []garden.PortMapping{},
				},
			},
			"test-handle  192.0.2.10  ruby  2016-11-15 06:48:15  N/A\n",
		),

		Entry(
			"with 1 container with 1 PortMapping",
			result.Containers{
				"test-handle": result.CInfo{
					Ip:           "192.0.2.10",
					ProcessName:  "ruby",
					CreatedAt:    "2016-11-15T06:48:15.137799416Z",
					PortMappings: []garden.PortMapping{{80, 8080}},
				},
			},
			"test-handle  192.0.2.10  ruby  2016-11-15 06:48:15  80->8080\n",
		),

		Entry(
			"with 1 container with 2 PortMappings",
			result.Containers{
				"test-handle": result.CInfo{
					Ip:           "192.0.2.10",
					ProcessName:  "ruby",
					CreatedAt:    "2016-11-15T06:48:15.137799416Z",
					PortMappings: []garden.PortMapping{{80, 8080}, {443, 4443}},
				},
			},
			"test-handle  192.0.2.10  ruby  2016-11-15 06:48:15  80->8080, 443->4443\n",
		),

		Entry(
			"with 2 containers",
			result.Containers{
				"test-handle": result.CInfo{
					Ip:           "192.0.2.10",
					ProcessName:  "ruby",
					CreatedAt:    "2016-11-15T06:48:15.137799416Z",
					PortMappings: []garden.PortMapping{},
				},
				"test-handle-2": result.CInfo{
					Ip:           "192.0.2.11",
					ProcessName:  "tree",
					CreatedAt:    "2016-11-15T06:48:15.137799416Z",
					PortMappings: []garden.PortMapping{},
				},
			},
			"test-handle    192.0.2.10  ruby  2016-11-15 06:48:15  N/A\ntest-handle-2  192.0.2.11  tree  2016-11-15 06:48:15  N/A\n",
		),
	)

	Context("when the table Render returns an error", func() {
		BeforeEach(func() {
			stdout = ErroringWriter{}
		})

		It("returns the error", func() {
			containersResult := result.Containers{
				"test-handle": result.CInfo{},
			}

			err := resultPrinter.PrintContainers(containersResult)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("error-writing-to-writer"))
		})
	})
})
