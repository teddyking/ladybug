package ui

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type Table struct {
	Out io.Writer
}

type Row []string
type Rows []Row

func NewTable(out io.Writer) *Table {
	return &Table{
		Out: out,
	}
}

func (t *Table) Render(rows Rows) {
	w := new(tabwriter.Writer)

	w.Init(t.Out, 8, 8, 0, '\t', 0)

	for _, row := range rows {
		var formattedRow string
		for _, cell := range row {
			formattedRow = fmt.Sprintf("%s%s\t", formattedRow, cell)
		}
		fmt.Fprintln(w, formattedRow)
	}

	w.Flush()
}
