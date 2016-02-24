package tablib

import "github.com/bndr/gotabulate"

// Tabular returns a tabular string representation of the Dataset.
// format is either grid or simple.
func (d *Dataset) Tabular(format string) string {
	back := d.Records()
	t := gotabulate.Create(back)

	return t.Render(format)
}
