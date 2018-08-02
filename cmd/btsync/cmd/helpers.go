package cmd

import (
	"github.com/olekukonko/tablewriter"

	"io"
)

// NewTableWriter returns a `*tablewriter.Table` for outputting tables with the CLI.
func NewTableWriter(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Right: true})
	table.SetCenterSeparator("|")

	return table
}
