package tablib

import "bytes"

// HTML returns the HTML representation of the Dataset as string.
func (d *Dataset) HTML() string {
	back := d.Records()
	var b bytes.Buffer

	b.WriteString("<table class=\"table table-striped\">\n\t<thead>")
	for i, r := range back {
		b.WriteString("\n\t\t<tr>")
		for _, c := range r {
			tag := "td"
			if i == 0 {
				tag = "th"
			}
			b.WriteString("\n\t\t\t<" + tag + ">")
			b.WriteString(c)
			b.WriteString("</" + tag + ">")
		}
		b.WriteString("\n\t\t</tr>")
		if i == 0 {
			b.WriteString("\n\t</thead>\n\t<tbody>")
		}
	}
	b.WriteString("\n\t</tbody>\n</table>")

	return b.String()
}

// HTML returns a HTML representation of the Databook as a byte array.
func (d *Databook) HTML() string {
	var b bytes.Buffer

	for _, s := range d.sheets {
		b.WriteString("<h1>" + s.title + "</h1>\n")
		b.WriteString(s.dataset.HTML())
		b.WriteString("\n\n")
	}

	return b.String()
}
