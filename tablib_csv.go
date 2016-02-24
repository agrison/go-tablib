package tablib

import (
	"bytes"
	"encoding/csv"
)

// CSV returns a CSV representation of the Dataset as string.
func (d *Dataset) CSV() (string, error) {
	records := d.Records()
	var b bytes.Buffer

	w := csv.NewWriter(&b)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		return "", err
	}

	return b.String(), nil
}

// TSV returns a TSV representation of the Dataset as string.
func (d *Dataset) TSV() (string, error) {
	records := d.Records()
	var b bytes.Buffer

	w := csv.NewWriter(&b)
	w.Comma = '\t'
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		return "", err
	}

	return b.String(), nil
}
