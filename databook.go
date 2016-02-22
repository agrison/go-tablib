package tablib

import (
	"gopkg.in/yaml.v2"
)

// Sheet represents a sheet in a Databook, holding a title (if any) and a dataset.
type Sheet struct {
	title   string
	dataset *Dataset
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
	str += str[:len(str)-1] + "]"
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
		y[i]["dataset"] = s.dataset.ArrayOfMap()
	}
	b, err := yaml.Marshal(y)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
