package ui_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/ladybug/ui"

	"fmt"

	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Table", func() {
	var (
		table  *Table
		stdout *gbytes.Buffer
	)

	BeforeEach(func() {
		stdout = gbytes.NewBuffer()
		table = NewTable(stdout)
	})

	DescribeTable("Render",
		func(rows Rows, expectedTable string) {
			table.Render(rows)

			Eventually(stdout).Should(gbytes.Say(expectedTable))
		},
		Entry(
			"with 0 Cols, 0 Rows",
			Rows{},
			"",
		),
		Entry(
			"with 1 Col, 1 Row",
			Rows{
				{"Header"},
			},
			fmt.Sprintf("Header\t\n"),
		),
		Entry(
			"with 2 Cols, 2 Rows",
			Rows{
				{"Header", "Header"},
				{"Cell", "Cell"},
			},
			fmt.Sprintf("Header\tHeader\t\nCell\tCell\t\n"),
		),
		Entry(
			"with 2 Cols, 2 Rows and an empty string in one of the Cells",
			Rows{
				{"Header", "Header"},
				{"", "Cell"},
			},
			fmt.Sprintf("Header\tHeader\t\n\tCell\t\n"),
		),
		Entry(
			"with 2 Cols, 2 Rows and a really long string in one of the Cells",
			Rows{
				{"Header", "Header"},
				{"This is a really long cell, longer than 8 chars", "Cell"},
			},
			fmt.Sprintf("Header\t\t\t\t\t\tHeader\t\nThis is a really long cell, longer than 8 chars\tCell\t\n"),
		),
	)
})
