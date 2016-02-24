package tablib

import (
	"bytes"
	"github.com/tealeg/xlsx"
)

// XLSX exports the Dataset as a byte array representing the .xlsx format.
func (d *Dataset) XLSX() ([]byte, error) {
	file := xlsx.NewFile()
	if err := d.addXlsxSheetToFile(file, "Sheet 1"); err != nil {
		return nil, err
	}

	var b bytes.Buffer
	file.Write(&b)
	return b.Bytes(), nil
}

// XLSX returns a XLSX representation of the Databook as a byte array.
func (d *Databook) XLSX() ([]byte, error) {
	file := xlsx.NewFile()

	for _, s := range d.sheets {
		s.dataset.addXlsxSheetToFile(file, s.title)
	}

	var b bytes.Buffer
	file.Write(&b)
	return b.Bytes(), nil
}

func (d *Dataset) addXlsxSheetToFile(file *xlsx.File, sheetName string) error {
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		return nil
	}

	back := d.Records()
	for i, r := range back {
		row := sheet.AddRow()
		for _, c := range r {
			cell := row.AddCell()
			cell.Value = c
			if i == 0 {
				cell.GetStyle().Font.Bold = true
			}
		}
	}
	return nil
}
