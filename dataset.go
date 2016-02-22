package tablib

import (
	"bytes"
	"container/list"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/clbanning/mxj"
	"gopkg.in/yaml.v2"
	"strconv"
	"time"
)

// Dataset represents a set of data, which is a list of data and header for each column.
type Dataset struct {
	headers *list.List
	data    *list.List
	rows    int
	cols    int
}

// NewDataset creates a new dataset.
func NewDataset(headers []string) *Dataset {
	return NewDatasetWithData(headers, nil)
}

// NewDataset creates a new dataset.
func NewDatasetWithData(headers []string, data [][]interface{}) *Dataset {
	h := list.New()
	for _, s := range headers {
		h.PushBack(s)
	}
	d := &Dataset{h, list.New(), 0, len(headers)}
	if data != nil {
		for _, r := range data {
			d.Append(r)
		}
	}
	return d
}

// Append appends a row of values to the dataset.
func (d *Dataset) Append(row []interface{}) *Dataset {
	d.data.PushBack(row)
	d.rows++
	return d
}

// AppendValues appends a row of values to the dataset.
func (d *Dataset) AppendValues(row ...interface{}) *Dataset {
	d.data.PushBack(row[:])
	d.rows++
	return d
}

// AppendColumn appends a new column with values to the dataset.
func (d *Dataset) AppendColumn(header string, cols []interface{}) *Dataset {
	d.headers.PushBack(header)
	d.cols++
	i := 0
	for e := d.data.Front(); e != nil; e = e.Next() {
		clone := make([]interface{}, d.cols)
		copy(clone, e.Value.([]interface{}))
		clone[d.cols-1] = cols[i]
		e.Value = clone
		i++
	}
	return d
}

// AppendColumnValues appends a new column with values to the dataset.
func (d *Dataset) AppendColumnValues(header string, cols ...interface{}) *Dataset {
	return d.AppendColumn(header, cols[:])
}

// Column returns all the values for a specific column
func (d *Dataset) Column(header string) []interface{} {
	i := indexOfColumn(header, d)
	if i == -1 {
		return nil
	}

	values := make([]interface{}, d.rows)
	idx := 0
	for e := d.data.Front(); e != nil; e = e.Next() {
		values[idx] = e.Value.([]interface{})[i]
		idx++
	}
	return values
}

// DeleteRow deletes a row at a specific index
func (d *Dataset) DeleteRow(row int) *Dataset {
	if row >= d.rows {
		return d
	}
	removed := false
	i := 0
	for e := d.data.Front(); e != nil && !removed; e = e.Next() {
		if i == row {
			d.data.Remove(e)
			removed = true
		}
		i++
	}
	d.rows--
	return d
}

// DeleteColumn deletes a column from the dataset.
func (d *Dataset) DeleteColumn(header string) *Dataset {
	i := indexOfColumn(header, d)
	if i == -1 {
		return d
	}
	// remove the column
	removed := false
	for e := d.headers.Front(); e != nil && !removed; e = e.Next() {
		if e.Value.(string) == header {
			d.headers.Remove(e)
			removed = true
		}
	}

	// remove the column values
	for e := d.data.Front(); e != nil; e = e.Next() {
		clone := make([]interface{}, d.cols)
		copy(clone[:i], e.Value.([]interface{})[:i])
		copy(clone[i:], e.Value.([]interface{})[i+1:])
		e.Value = clone
	}

	d.cols--
	return d
}

// Json returns a JSON representation of the dataset as string.
func (d *Dataset) Json() (string, error) {
	back := d.ArrayOfMap()

	b, err := json.Marshal(back)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Xml returns a XML representation of the dataset as string.
func (d *Dataset) Xml() string {
	return d.XmlWithTagNamePrefixIndent("row", "  ", "  ")
}

// XmlWithTagNamePrefixIndent returns a XML representation with custom tag, prefix and indent.
func (d *Dataset) XmlWithTagNamePrefixIndent(tagName, prefix, indent string) string {
	back := d.ArrayOfMap()

	var b bytes.Buffer
	b.WriteString("<dataset>\n")
	for _, r := range back {
		m := mxj.Map(r.(map[string]interface{}))
		m.XmlIndentWriter(&b, prefix, indent, tagName)
	}
	b.WriteString("\n" + prefix + "</dataset>")

	return b.String()
}

// Csv returns a CSV representation of the dataset as string.
func (d *Dataset) Csv() (string, error) {
	records := d.Records()
	var b bytes.Buffer

	w := csv.NewWriter(&b)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		return "", err
	}

	return b.String(), nil
}

// Tsv returns a TSV representation of the dataset as string.
func (d *Dataset) Tsv() (string, error) {
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

// Yaml returns a YAML representation of the dataset as string.
func (d *Dataset) Yaml() (string, error) {
	back := d.ArrayOfMap()

	b, err := yaml.Marshal(back)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func indexOfColumn(header string, d *Dataset) int {
	i := 0
	for e := d.headers.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == header {
			return i
		}
		i++
	}
	return -1
}

// ArrayOfMap returns the dataset as an array of map where each key is a column.
func (d *Dataset) ArrayOfMap() []interface{} {
	back := make([]interface{}, d.rows)
	i := 0
	for e := d.data.Front(); e != nil; e = e.Next() {
		m := make(map[string]interface{}, d.cols-1)
		j := 0
		for c := d.headers.Front(); c != nil; c = c.Next() {
			m[c.Value.(string)] = e.Value.([]interface{})[j]
			j++
		}
		back[i] = m
		i++
	}
	return back
}

// Records returns the dataset as an array of array where each entry is a string.
// The first row of the returned 2d array represents the columns of the dataset.
func (d *Dataset) Records() [][]string {
	records := make([][]string, d.rows+1 /* +1 for header */)
	i := 0
	j := 0
	records[i] = make([]string, d.cols)
	for e := d.headers.Front(); e != nil; e = e.Next() {
		records[i][j] = e.Value.(string)
		j++
	}
	i++
	for e := d.data.Front(); e != nil; e = e.Next() {
		records[i] = make([]string, d.cols)
		j = 0
		vals := e.Value.([]interface{})
		for _, v := range vals {
			switch v.(type) {
			case string:
				records[i][j] = v.(string)
			case int:
				records[i][j] = strconv.Itoa(v.(int))
			case int64:
				records[i][j] = strconv.FormatInt(v.(int64), 10)
			case uint64:
				records[i][j] = strconv.FormatUint(v.(uint64), 10)
			case bool:
				records[i][j] = strconv.FormatBool(v.(bool))
			case float64:
				records[i][j] = strconv.FormatFloat(v.(float64), 'G', -1, 32)
			case time.Time:
				records[i][j] = v.(time.Time).Format(time.RFC3339)
			default:
				fmt.Printf("Skipping value.")
			}
			j++
		}
		i++
	}
	return records
}
