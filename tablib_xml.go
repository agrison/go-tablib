package tablib

import (
	"bytes"
	"github.com/clbanning/mxj"
)

// XML returns a XML representation of the Dataset as string.
func (d *Dataset) XML() (string, error) {
	return d.XMLWithTagNamePrefixIndent("row", "  ", "  ")
}

// XML returns a XML representation of the Databook as string.
func (d *Databook) XML() (string, error) {
	str := "<databook>\n"
	for _, s := range d.sheets {
		str += "  <sheet>\n    <title>" + s.title + "</title>\n    "
        row, err := s.dataset.XMLWithTagNamePrefixIndent("row", "      ", "  ")
        if err != nil {
            return "", nil
        }
		str += row
		str += "\n  </sheet>"
	}
	str += "\n</databook>"
	return str, nil
}

// XMLWithTagNamePrefixIndent returns a XML representation with custom tag, prefix and indent.
func (d *Dataset) XMLWithTagNamePrefixIndent(tagName, prefix, indent string) (string, error) {
	back := d.Dict()

	var b bytes.Buffer
	b.WriteString("<dataset>\n")
	for _, r := range back {
		m := mxj.Map(r.(map[string]interface{}))
		if err :=m.XmlIndentWriter(&b, prefix, indent, tagName); err != nil {
            return "", err
        }
	}
	b.WriteString("\n" + prefix + "</dataset>")

	return b.String(), nil
}

// LoadXML loads a Dataset from an XML source.
func LoadXML(input []byte) (*Dataset, error) {
	m, _, err := mxj.NewMapXmlReaderRaw(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	// this seems quite a bit hacky
	datasetNode, _ := m.ValueForPath("dataset")
	rowNode := datasetNode.(map[string]interface{})["row"].([]interface{})

	back := make([]map[string]interface{}, 0, len(rowNode))
	for _, r := range rowNode {
		back = append(back, r.(map[string]interface{}))
	}

	return internalLoadFromDict(back)
}
