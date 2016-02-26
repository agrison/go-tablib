package tablib

import (
	"github.com/bndr/gotabulate"
	"regexp"
)

var (
	// TabularGrid is the value to be passed to gotabulate to render the table
	// as ASCII table with grid format
	TabularGrid = "grid"
	// TabularGrid is the value to be passed to gotabulate to render the table
	// as ASCII table with simple format
	TabularSimple = "simple"
	// TabularGrid is the value to be passed to gotabulate to render the table
	// as ASCII table with condensed format
	TabularCondensed = "condensed"
)

// Tabular returns a tabular string representation of the Dataset.
// format is either grid, simple or condensed.
func (d *Dataset) Tabular(format string) string {
	back := d.Records()
	t := gotabulate.Create(back)

	if format == TabularCondensed {
		return regexp.MustCompile("\n\n\\s").ReplaceAllString(t.Render("simple"), "\n ")
	}
	return t.Render(format)
}
