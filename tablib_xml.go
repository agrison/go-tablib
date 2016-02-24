package tablib

import (
	"bytes"
	"github.com/clbanning/mxj"
)

// XML returns a XML representation of the Dataset as string.
func (d *Dataset) XML() string {
	return d.XMLWithTagNamePrefixIndent("row", "  ", "  ")
}

// XML returns a XML representation of the Databook as string.
func (d *Databook) XML() string {
	str := "<Databook>\n"
	for _, s := range d.sheets {
		str += "  <sheet>\n    <title>" + s.title + "</title>\n    "
		str += s.dataset.XMLWithTagNamePrefixIndent("row", "      ", "  ")
		str += "\n  </sheet>"
	}
	str += "\n</Databook>"
	return str
}

// XMLWithTagNamePrefixIndent returns a XML representation with custom tag, prefix and indent.
func (d *Dataset) XMLWithTagNamePrefixIndent(tagName, prefix, indent string) string {
	back := d.Dict()

	var b bytes.Buffer
	b.WriteString("<Dataset>\n")
	for _, r := range back {
		m := mxj.Map(r.(map[string]interface{}))
		m.XmlIndentWriter(&b, prefix, indent, tagName)
	}
	b.WriteString("\n" + prefix + "</Dataset>")

	return b.String()
}
