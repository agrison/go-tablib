package tablib

import (
	"bytes"
	"github.com/tealeg/xlsx"
	"gopkg.in/yaml.v2"
)

// Sheet represents a sheet in a Databook, holding a title (if any) and a dataset.
type Sheet struct {
	title   string
	dataset *Dataset
}

// Title return the title of the sheet.
func (s Sheet) Title() string {
	return s.title
}

// Dataset returns the dataset of the sheet.
func (s Sheet) Dataset() *Dataset {
	return s.dataset
}

// Databook represents a Databook which is an array of sheets.
type Databook struct {
	sheets map[string]Sheet
}

// NewDatabook constructs a new Databook.
func NewDatabook() *Databook {
	return &Databook{make(map[string]Sheet)}
}

// Sheets returns the sheets in the Databook.
func (d *Databook) Sheets() map[string]Sheet {
	return d.sheets
}

// Sheet returns the sheet with a specific title.
func (d *Databook) Sheet(title string) Sheet {
	return d.sheets[title]
}

// AddSheet adds a sheet to the Databook.
func (d *Databook) AddSheet(title string, dataset *Dataset) *Databook {
	d.sheets[title] = Sheet{title, dataset}
	return d
}

// Size returns the number of sheets in the databook.
func (d *Databook) Size() int {
	return len(d.sheets)
}

// JSON returns a JSON representation of the databook as string.
func (d *Databook) JSON() (string, error) {
	str := "["
	for _, s := range d.sheets {
		str += "{\"title\": \"" + s.title + "\", \"data\": "
		js, err := s.dataset.JSON()
		if err != nil {
			return "", err
		}
		str += js + "},"
	}
	str = str[:len(str)-1] + "]"
	return str, nil
}

// XML returns a XML representation of the Databook as string.
func (d *Databook) XML() string {
	str := "<databook>\n"
	for _, s := range d.sheets {
		str += "  <sheet>\n    <title>" + s.title + "</title>\n    "
		str += s.dataset.XMLWithTagNamePrefixIndent("row", "      ", "  ")
		str += "\n  </sheet>"
	}
	str += "\n</databook>"
	return str
}

// YAML returns a YAML representation of the databook as string.
func (d *Databook) YAML() (string, error) {
	y := make([]map[string]interface{}, len(d.sheets))
	i := 0
	for _, s := range d.sheets {
		y[i] = make(map[string]interface{})
		y[i]["title"] = s.title
		y[i]["dataset"] = s.dataset.Dict()
	}
	b, err := yaml.Marshal(y)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// XLSX returns a XLSX representation of the databook as a byte array.
func (d *Databook) XLSX() ([]byte, error) {
	file := xlsx.NewFile()

	for _, s := range d.sheets {
		s.dataset.addXlsxSheetToFile(file, s.title)
	}

	var b bytes.Buffer
	file.Write(&b)
	return b.Bytes(), nil
}

// HTML returns a HTML representation of the databook as a byte array.
func (d *Databook) HTML() string {
	var b bytes.Buffer

	for _, s := range d.sheets {
		b.WriteString("<h1>" + s.title + "</h1>\n")
		b.WriteString(s.dataset.HTML())
		b.WriteString("\n\n")
	}

	return b.String()
}
