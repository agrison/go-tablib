package tablib

import (
	"fmt"
	"strings"
	"testing"
)

func TestAppend(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.Append([]interface{}{"John", "Adams"})
	if ds.rows != 1 {
		t.Errorf("Number of rows should be 1")
	}
}

func TestAppendValues(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("John", "Adams")
	if ds.rows != 1 {
		t.Errorf("Number of rows should be 1")
	}
}

func TestAppendColumn(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendColumn("age", nil)
	if ds.cols != 3 {
		t.Errorf("Number of rows should be 1")
	}
}

func TestColumn(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	x := ds.Column("lastName")
	if !stringInSlice("Washington", x) && !stringInSlice("Ford", x) {
		t.Errorf("Washington and Ford should be in the column result")
	}
}

func TestDeleteRow(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.DeleteRow(1)
	x := ds.Column("lastName")
	if ds.rows != 1 && !stringInSlice("Washington", x) && stringInSlice("Ford", x) {
		t.Errorf("Ford should be in the column result")
	}
}

func TestDeleteColumn(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.DeleteColumn("lastName")
	if ds.cols != 1 && ds.Column("lastName") != nil {
		t.Errorf("lastName should not be part of the dataset")
	}
}

func lastNameLen(row []interface{}) interface{} {
	return len(row[1].(string))
}

func TestDynamicColumn(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.AppendColumn("age", []interface{}{90, 67, 83})
	ds.AppendDynamicColumn("Name length", lastNameLen)
	x, _ := ds.CSV()
	fmt.Printf("%s\n", x)
}

func TestJSON(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.AppendColumn("age", []interface{}{90, 67, 83})
	j, _ := ds.JSON()
	if j != `[{"age":90,"firstName":"George","lastName":"Washington"},{"age":67,"firstName":"Henry","lastName":"Ford"}]` {
		t.Errorf("error Json()")
	}
}

func TestYAML(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.AppendColumn("age", []interface{}{90, 67})
	y, _ := ds.YAML()
	if !strings.Contains(y, "- age:") && !strings.Contains(y, "firstName:") {
		t.Errorf("error Yaml()")
	}
}

func TestXML(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.AppendColumn("age", []interface{}{90, 67})
	x := ds.XML()
	if !strings.Contains(x, "<age>") && !strings.Contains(x, "<firstName>") {
		t.Errorf("error XML()")
	}
}

func stringInSlice(a string, list []interface{}) bool {
	for _, b := range list {
		if b.(string) == a {
			return true
		}
	}
	return false
}
