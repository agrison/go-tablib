package tablib

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/clbanning/mxj"
	"gopkg.in/yaml.v2"
	"sort"
	"strconv"
	"time"
)

// Dataset represents a set of data, which is a list of data and header for each column.
type Dataset struct {
	headers []string
	data    [][]interface{}
	tags    [][]string
	rows    int
	cols    int
}

// DynamicColumn represents a function that can be evaluated dynamically
// when exporting to a predefined format.
type DynamicColumn func([]interface{}) interface{}

// NewDataset creates a new dataset.
func NewDataset(headers []string) *Dataset {
	return NewDatasetWithData(headers, nil)
}

// NewDatasetWithData creates a new dataset.
func NewDatasetWithData(headers []string, data [][]interface{}) *Dataset {
	d := &Dataset{headers, data, make([][]string, 0), len(data), len(headers)}
	return d
}

// Headers return the headers of the dataset.
func (d *Dataset) Headers() []string {
	return d.headers
}

// Append appends a row of values to the dataset.
func (d *Dataset) Append(row []interface{}) *Dataset {
	d.data = append(d.data, row)
	d.tags = append(d.tags, make([]string, 0))
	d.rows++
	return d
}

// AppendTagged appends a row of values to the dataset with one or multiple tags
// for filtering purposes.
func (d *Dataset) AppendTagged(row []interface{}, tags ...string) *Dataset {
	d.Append(row)
	d.tags[d.rows-1] = tags[:]
	return d
}

// AppendValues appends a row of values to the dataset.
func (d *Dataset) AppendValues(row ...interface{}) *Dataset {
	return d.Append(row[:])
}

// AppendValuesTagged appends a row of values to the dataset with one or multiple tags
// for filtering purposes.
func (d *Dataset) AppendValuesTagged(row ...interface{}) *Dataset {
	return d.AppendTagged(row[:])
}

// AppendColumn appends a new column with values to the dataset.
func (d *Dataset) AppendColumn(header string, cols []interface{}) *Dataset {
	d.headers = append(d.headers, header)
	d.cols++
	for i, e := range d.data {
		d.data[i] = append(e, cols[i])
	}
	return d
}

// AppendColumnValues appends a new column with values to the dataset.
func (d *Dataset) AppendColumnValues(header string, cols ...interface{}) *Dataset {
	return d.AppendColumn(header, cols[:])
}

// AppendDynamicColumn appends a dynamic column to the dataset.
func (d *Dataset) AppendDynamicColumn(header string, fn DynamicColumn) *Dataset {
	d.headers = append(d.headers, header)
	d.cols++
	for i, e := range d.data {
		d.data[i] = append(e, fn)
	}
	return d
}

// Column returns all the values for a specific column
func (d *Dataset) Column(header string) []interface{} {
	colIndex := indexOfColumn(header, d)
	if colIndex == -1 {
		return nil
	}

	values := make([]interface{}, d.rows)
	for i, e := range d.data {
		switch e[colIndex].(type) {
		case DynamicColumn:
			values[i] = e[colIndex].(DynamicColumn)(e)
		default:
			values[i] = e[colIndex]
		}
	}
	return values
}

// Filter filters a dataset, returning a fresh dataset including only the rows
// previously tagged with one of the given tags. Returns a new Dataset.
func (d *Dataset) Filter(tags ...string) *Dataset {
	nd := NewDataset(d.headers)
	for rowIndex, rowValue := range d.data {
		for _, filterTag := range tags {
			if isTagged(filterTag, d.tags[rowIndex]) {
				nd.AppendTagged(rowValue, d.tags[rowIndex]...) // copy tags
			}
		}
	}
	return nd
}

// Sort sorts the Dataset by a specific column. Returns a new Dataset.
func (d *Dataset) Sort(column string) *Dataset {
	return d.internalSort(column, false)
}

// SortReverse sorts the Dataset by a specific column in reverse order. Returns a new Dataset.
func (d *Dataset) SortReverse(column string) *Dataset {
	return d.internalSort(column, true)
}

func (d *Dataset) internalSort(column string, reverse bool) *Dataset {
	nd := NewDataset(d.headers)
	pairs := make([]entryPair, 0, nd.rows)
	for i, v := range d.Column(column) {
		pairs = append(pairs, entryPair{i, v})
	}

	var how sort.Interface
	// sort by column
	switch pairs[0].value.(type) {
	case string:
		how = byStringValue(pairs)
	case int:
		how = byIntValue(pairs)
	case int64:
		how = byInt64Value(pairs)
	case uint64:
		how = byUint64Value(pairs)
	case float64:
		how = byFloatValue(pairs)
	case time.Time:
		how = byTimeValue(pairs)
	default:
		// nothing
	}

	if !reverse {
		sort.Sort(how)
	} else {
		sort.Sort(sort.Reverse(how))
	}

	// now iterate on the pairs and add the data sorted to the new dataset
	for _, p := range pairs {
		nd.AppendTagged(d.data[p.index], d.tags[p.index]...)
	}

	return nd
}

// DeleteRow deletes a row at a specific index
func (d *Dataset) DeleteRow(row int) *Dataset {
	if row >= d.rows {
		return d
	}
	d.data = append(d.data[:row], d.data[row+1:]...)
	d.rows--
	return d
}

// DeleteColumn deletes a column from the dataset.
func (d *Dataset) DeleteColumn(header string) *Dataset {
	colIndex := indexOfColumn(header, d)
	if colIndex == -1 {
		return d
	}
	d.cols--
	d.headers = append(d.headers[:colIndex], d.headers[colIndex+1:]...)
	// remove the column
	for i := range d.data {
		d.data[i] = append(d.data[i][:colIndex], d.data[i][colIndex+1:]...)
	}
	return d
}

// JSON returns a JSON representation of the dataset as string.
func (d *Dataset) JSON() (string, error) {
	back := d.Dict()

	b, err := json.Marshal(back)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// XML returns a XML representation of the dataset as string.
func (d *Dataset) XML() string {
	return d.XMLWithTagNamePrefixIndent("row", "  ", "  ")
}

// XMLWithTagNamePrefixIndent returns a XML representation with custom tag, prefix and indent.
func (d *Dataset) XMLWithTagNamePrefixIndent(tagName, prefix, indent string) string {
	back := d.Dict()

	var b bytes.Buffer
	b.WriteString("<dataset>\n")
	for _, r := range back {
		m := mxj.Map(r.(map[string]interface{}))
		m.XmlIndentWriter(&b, prefix, indent, tagName)
	}
	b.WriteString("\n" + prefix + "</dataset>")

	return b.String()
}

// CSV returns a CSV representation of the dataset as string.
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

// TSV returns a TSV representation of the dataset as string.
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

// YAML returns a YAML representation of the dataset as string.
func (d *Dataset) YAML() (string, error) {
	back := d.Dict()

	b, err := yaml.Marshal(back)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func indexOfColumn(header string, d *Dataset) int {
	for i, e := range d.headers {
		if e == header {
			return i
		}
	}
	return -1
}

// Dict returns the dataset as an array of map where each key is a column.
func (d *Dataset) Dict() []interface{} {
	back := make([]interface{}, d.rows)
	for i, e := range d.data {
		m := make(map[string]interface{}, d.cols-1)
		for j, c := range d.headers {
			switch e[j].(type) {
			case DynamicColumn:
				m[c] = e[j].(DynamicColumn)(e)
			default:
				m[c] = e[j]
			}
		}
		back[i] = m
	}
	return back
}

// Records returns the dataset as an array of array where each entry is a string.
// The first row of the returned 2d array represents the columns of the dataset.
func (d *Dataset) Records() [][]string {
	records := make([][]string, d.rows+1 /* +1 for header */)
	records[0] = make([]string, d.cols)
	for j, e := range d.headers {
		records[0][j] = e
	}
	for i, e := range d.data {
		rowIndex := i + 1
		j := 0
		records[rowIndex] = make([]string, d.cols)
		for _, v := range e {
			vv := v
			switch v.(type) {
			case DynamicColumn:
				vv = v.(DynamicColumn)(e)
			default:
				// nothing
			}
			switch vv.(type) {
			case string:
				records[rowIndex][j] = vv.(string)
			case int:
				records[rowIndex][j] = strconv.Itoa(vv.(int))
			case int64:
				records[rowIndex][j] = strconv.FormatInt(vv.(int64), 10)
			case uint64:
				records[rowIndex][j] = strconv.FormatUint(vv.(uint64), 10)
			case bool:
				records[rowIndex][j] = strconv.FormatBool(vv.(bool))
			case float64:
				records[rowIndex][j] = strconv.FormatFloat(vv.(float64), 'G', -1, 32)
			case time.Time:
				records[rowIndex][j] = vv.(time.Time).Format(time.RFC3339)
			default:
				fmt.Printf("Skipping value.")
			}
			j++
		}
	}

	return records
}

func isTagged(tag string, tags []string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
