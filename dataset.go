// Package tablib is a format-agnostic tabular dataset library, written in Go.
// It allows you to import, export, and manipulate tabular data sets.
// Advanced features include, dynamic columns, tags & filtering, and seamless format import & export.
package tablib

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bndr/gotabulate"
	"github.com/clbanning/mxj"
	"github.com/tealeg/xlsx"
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

var (
	// ErrInvalidDimensions is returned when trying to append/insert too much
	// or not enough values to a row or column
	ErrInvalidDimensions = errors.New("tablib: Invalid dimension")
	// ErrInvalidColumnIndex is returned when trying to insert a column at an
	// invalid index
	ErrInvalidColumnIndex = errors.New("tablib: Invalid column index")
	// ErrInvalidRowIndex is returned when trying to insert a row at an
	// invalid index
	ErrInvalidRowIndex = errors.New("tablib: Invalid row index")
)

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

// Width returns the number of columns in the dataset.
func (d *Dataset) Width() int {
	return d.cols
}

// Height returns the number of rows in the dataset.
func (d *Dataset) Height() int {
	return d.rows
}

// Append appends a row of values to the dataset.
func (d *Dataset) Append(row []interface{}) error {
	if len(row) != d.cols {
		return ErrInvalidDimensions
	}
	d.data = append(d.data, row)
	d.tags = append(d.tags, make([]string, 0))
	d.rows++
	return nil
}

// AppendTagged appends a row of values to the dataset with one or multiple tags
// for filtering purposes.
func (d *Dataset) AppendTagged(row []interface{}, tags ...string) error {
	if err := d.Append(row); err != nil {
		return err
	}
	d.tags[d.rows-1] = tags[:]
	return nil
}

// AppendValues appends a row of values to the dataset.
func (d *Dataset) AppendValues(row ...interface{}) error {
	return d.Append(row[:])
}

// AppendValuesTagged appends a row of values to the dataset with one or multiple tags
// for filtering purposes.
func (d *Dataset) AppendValuesTagged(row ...interface{}) error {
	return d.AppendTagged(row[:])
}

// Insert inserts a row at a given index.
func (d *Dataset) Insert(index int, row []interface{}) error {
	if index < 0 || index >= d.rows {
		return ErrInvalidRowIndex
	}

	if len(row) != d.cols {
		return ErrInvalidDimensions
	}

	ndata := make([][]interface{}, 0, d.rows+1)
	ndata = append(ndata, d.data[:index]...)
	ndata = append(ndata, row)
	ndata = append(ndata, d.data[index:]...)
	d.data = ndata
	d.rows++

	ntags := make([][]string, 0, d.rows+1)
	ntags = append(ntags, d.tags[:index]...)
	ntags = append(ntags, make([]string, 0))
	ntags = append(ntags, d.tags[index:]...)
	d.tags = ntags

	return nil
}

// InsertValues inserts a row of values at a given index.
func (d *Dataset) InsertValues(index int, values ...interface{}) error {
	return d.Insert(index, values[:])
}

// InsertTagged inserts a row at a given index with specific tags.
func (d *Dataset) InsertTagged(index int, row []interface{}, tags ...string) error {
	if err := d.Insert(index, row); err != nil {
		return err
	}
	d.Insert(index, row)
	d.tags[index] = tags[:]

	return nil
}

// AppendColumn appends a new column with values to the dataset.
func (d *Dataset) AppendColumn(header string, cols []interface{}) error {
	if len(cols) != d.rows {
		return ErrInvalidDimensions
	}
	d.headers = append(d.headers, header)
	d.cols++
	for i, e := range d.data {
		d.data[i] = append(e, cols[i])
	}
	return nil
}

// AppendColumnValues appends a new column with values to the dataset.
func (d *Dataset) AppendColumnValues(header string, cols ...interface{}) error {
	return d.AppendColumn(header, cols[:])
}

// AppendDynamicColumn appends a dynamic column to the dataset.
func (d *Dataset) AppendDynamicColumn(header string, fn DynamicColumn) {
	d.headers = append(d.headers, header)
	d.cols++
	for i, e := range d.data {
		d.data[i] = append(e, fn)
	}
}

// InsertColumn insert a new column at a given index.
func (d *Dataset) InsertColumn(index int, header string, cols []interface{}) error {
	if index < 0 || index >= d.cols {
		return ErrInvalidColumnIndex
	}

	if len(cols) != d.rows {
		return ErrInvalidDimensions
	}

	d.insertHeader(index, header)

	// for each row, insert the column
	for i, r := range d.data {
		row := make([]interface{}, 0, d.cols)
		row = append(row, r[:index]...)
		row = append(row, cols[i])
		row = append(row, r[index:]...)
		d.data[i] = row
	}

	return nil
}

// InsertDynamicColumn insert a new dynamic column at a given index.
func (d *Dataset) InsertDynamicColumn(index int, header string, fn DynamicColumn) error {
	if index < 0 || index >= d.cols {
		return ErrInvalidColumnIndex
	}

	d.insertHeader(index, header)

	// for each row, insert the column
	for i, r := range d.data {
		row := make([]interface{}, 0, d.cols)
		row = append(row, r[:index]...)
		row = append(row, fn)
		row = append(row, r[index:]...)
		d.data[i] = row
	}

	return nil
}

func (d *Dataset) insertHeader(index int, header string) {
	headers := make([]string, 0, d.cols+1)
	headers = append(headers, d.headers[:index]...)
	headers = append(headers, header)
	headers = append(headers, d.headers[index:]...)
	d.headers = headers
	d.cols++
}

// Stack stacks two Dataset by joining at the row level, and return new combined Dataset.
func (d *Dataset) Stack(other *Dataset) (*Dataset, error) {
	if d.Width() != other.Width() {
		return nil, ErrInvalidDimensions
	}

	nd := NewDataset(d.headers)
	nd.cols = d.cols
	nd.rows = d.rows + other.rows

	nd.tags = make([][]string, 0, nd.rows)
	nd.tags = append(nd.tags, d.tags...)
	nd.tags = append(nd.tags, other.tags...)

	nd.data = make([][]interface{}, 0, nd.rows)
	nd.data = append(nd.data, d.data...)
	nd.data = append(nd.data, other.data...)

	return nd, nil
}

// StackColumn stacks two Dataset by joining them at the column level, and return new combined Dataset.
func (d *Dataset) StackColumn(other *Dataset) (*Dataset, error) {
	if d.Height() != other.Height() {
		return nil, ErrInvalidDimensions
	}

	nheaders := d.headers
	nheaders = append(nheaders, other.headers...)

	nd := NewDataset(nheaders)
	nd.cols = d.cols + nd.cols
	nd.rows = d.rows
	nd.data = make([][]interface{}, nd.rows, nd.rows)
	nd.tags = make([][]string, nd.rows, nd.rows)

	for i := range d.data {
		nd.data[i] = make([]interface{}, 0, nd.cols)
		nd.data[i] = append(nd.data[i], d.data[i]...)
		nd.data[i] = append(nd.data[i], other.data[i]...)

		nd.tags[i] = make([]string, 0, nd.cols)
		nd.tags[i] = append(nd.tags[i], d.tags[i]...)
		nd.tags[i] = append(nd.tags[i], other.tags[i]...)
	}

	return nd, nil
}

// Column returns all the values for a specific column
// returns nil if column is not found.
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

// Transpose transposes a Dataset, turning rows into columns and vice versa,
// returning a new Dataset instance. The first row of the original instance
// becomes the new header row.
// TODO
func (d *Dataset) Transpose() *Dataset {
	panic("Transpose() not yet implemented")
}

// DeleteRow deletes a row at a specific index
func (d *Dataset) DeleteRow(row int) error {
	if row < 0 || row >= d.rows {
		return ErrInvalidRowIndex
	}
	d.data = append(d.data[:row], d.data[row+1:]...)
	d.rows--
	return nil
}

// DeleteColumn deletes a column from the dataset.
func (d *Dataset) DeleteColumn(header string) error {
	colIndex := indexOfColumn(header, d)
	if colIndex == -1 {
		return ErrInvalidColumnIndex
	}
	d.cols--
	d.headers = append(d.headers[:colIndex], d.headers[colIndex+1:]...)
	// remove the column
	for i := range d.data {
		d.data[i] = append(d.data[i][:colIndex], d.data[i][colIndex+1:]...)
	}
	return nil
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

// HTML returns the HTML representation of the dataset as string.
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

// Tabular returns a tabular string representation of the dataset.
// format is either grid or simple.
func (d *Dataset) Tabular(format string) string {
	back := d.Records()
	t := gotabulate.Create(back)

	return t.Render(format)
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
