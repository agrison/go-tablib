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
}

func TestDatabook(t *testing.T) {
	db := NewDatabook()
	ds := NewDataset([]string{"Maker", "Model"})
	ds.AppendTagged([]interface{}{"Porsche", "911"}, "fast", "luxury")
	ds.AppendTagged([]interface{}{"Skoda", "Octavia"}, "family")
	ds.AppendTagged([]interface{}{"Ferrari", "458"}, "fast", "luxury")
	ds.AppendValues("Citroen", "Picasso")
	ds.AppendTagged([]interface{}{"Bentley", "Continental"}, "luxury")
	db.AddSheet("Cars", ds)

	ds1 := NewDataset([]string{"firstName", "lastName"})
	ds1.AppendValues("George", "Washington")
	ds1.AppendValues("Henry", "Ford")
	ds1.AppendColumn("age", []interface{}{90, 67, 83})
	ds1.AppendDynamicColumn("Name length", lastNameLen)
	db.AddSheet("Presidents", ds1)
}

func TestTags(t *testing.T) {
	ds := NewDataset([]string{"Maker", "Model"})
	ds.AppendTagged([]interface{}{"Porsche", "911"}, "fast", "luxury")
	ds.AppendTagged([]interface{}{"Skoda", "Octavia"}, "family")
	ds.AppendTagged([]interface{}{"Ferrari", "458"}, "fast", "luxury")
	ds.AppendValues("Citroen", "Picasso")
	ds.AppendTagged([]interface{}{"Bentley", "Continental"}, "luxury")

	luxury := ds.Filter("luxury")
	if luxury.rows != 3 {
		t.Errorf("Should be 3 luxury cars")
	}

	fast := ds.Filter("fast")
	if fast.rows != 2 {
		t.Errorf("Should be 2 fast cars")
	}

	family := ds.Filter("family")
	if family.rows != 1 {
		t.Errorf("Should be 1 family car")
	}

	if ds.rows != 5 {
		t.Errorf("Should be 5 cars (original is untouched)")
	}
}

func TestDatabook2(t *testing.T) {
	db := NewDatabook()

	// a dataset of presidents
	presidents, _ := LoadJSON([]byte(`[
    {"Age":90,"First name":"John","Last name":"Adams"},
    {"Age":67,"First name":"George","Last name":"Washington"},
    {"Age":83,"First name":"Henry","Last name":"Ford"}
  ]`))

	// a dataset of cars
	cars := NewDataset([]string{"Maker", "Model", "Year"})
	cars.AppendValues("Porsche", "991", 2012)
	cars.AppendValues("Skoda", "Octavia", 2011)
	cars.AppendValues("Ferrari", "458", 2009)
	cars.AppendValues("Citroen", "Picasso II", 2013)
	cars.AppendValues("Bentley", "Continental GT", 2003)

	// add the sheets to the Databook
	db.AddSheet("Cars", cars.Sort("Year"))
	db.AddSheet("Presidents", presidents.SortReverse("Age"))
}

func TestSort(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.AppendValues("Foo", "Bar")
	ds.AppendColumn("age", []interface{}{90, 67, 83})

	/*nd1 := ds.Sort("lastName")
	x, _ := nd1.CSV()
	fmt.Printf("%s\n", x)
	nd2 := ds.SortReverse("age")
	x, _ = nd2.CSV()
	fmt.Printf("%s\n", x)*/
}

func TestInsert(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.AppendValues("Foo", "Bar")
	ds.AppendColumn("age", []interface{}{90, 67, 83})
	ds.Insert(1, []interface{}{"Alex", "Kid"})
}

func TestInsertColumn(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.AppendValues("Foo", "Bar")
	ds.InsertColumn(2, "age", []interface{}{90, 67, 83})
	ds.InsertDynamicColumn(10, "foo", lastNameLen)

	fmt.Println(ds.HTML())
}

func TestLoadYAML(t *testing.T) {
	ds, _ := LoadYAML([]byte(`- age: 90
  firstName: John
  lastName: Adams
- age: 67
  firstName: George
  lastName: Washington
- age: 83
  firstName: Henry
  lastName: Ford`))
	if ds.data[1][0] != "Washington" && ds.data[1][1] != "Washington" && ds.data[1][2] != "Washington" {
		t.Errorf("Error loadingYAML")
	}
	/*fmt.Printf("%+v\n", ds.headers)
	fmt.Printf("%+v\n", ds.data)
	j, _ := ds.JSON()
	fmt.Printf("JSONJSON: \n%s\n", j)*/
}

func TestLoadJSON(t *testing.T) {
	ds, _ := LoadJSON([]byte(`[
    {"age":90,"firstName":"John","lastName":"Adams"},
    {"age":67,"firstName":"George","lastName":"Washington"},
    {"age":83,"firstName":"Henry","lastName":"Ford"}
  ]`))
	if ds.data[1][0] != "Washington" && ds.data[1][1] != "Washington" && ds.data[1][2] != "Washington" {
		t.Errorf("Error LoadJSON")
	}
	/*fmt.Printf("%+v\n", ds.headers)
	fmt.Printf("%+v\n", ds.data)
	j, _ := ds.JSON()
	fmt.Printf("JSONJSON: \n%s\n", j)*/
}

func TestJSON(t *testing.T) {
	ds := NewDataset([]string{"firstName", "lastName"})
	ds.AppendValues("George", "Washington")
	ds.AppendValues("Henry", "Ford")
	ds.AppendColumn("age", []interface{}{90, 67})
	j, _ := ds.JSON()
	if j != `[{"age":90,"firstName":"George","lastName":"Washington"},{"age":67,"firstName":"Henry","lastName":"Ford"}]` {
		t.Errorf("error Json(): " + j)
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
