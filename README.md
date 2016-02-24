# go-tablib: format-agnostic tabular dataset library

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![Go Report Card](https://goreportcard.com/badge/github.com/agrison/go-commons-lang)][goreportcard]
[![Build Status](https://travis-ci.org/agrison/go-tablib.svg?branch=master)](https://travis-ci.org/agrison/go-tablib)

[license]: https://github.com/agrison/go-tablib/blob/master/LICENSE
[godocs]: https://godoc.org/github.com/agrison/go-tablib
[goreportcard]: https://goreportcard.com/report/github.com/agrison/go-tablib

Go-Tablib is a format-agnostic tabular dataset library, written in Go.
This is a port of the famous Python's [tablib](https://github.com/kennethreitz/tablib) by Kenneth Reitz.

Export formats supported:

* JSON (Sets + Books)
* YAML (Sets + Books)
* XLSX (Sets + Books)
* XML (Sets + Books)
* TSV (Sets)
* CSV (Sets)
* ASCII (Sets)

Loading formats supported:

* JSON (Sets)
* YAML (Sets)


## Overview

### tablib.Dataset
A Dataset is a table of tabular data. It may or may not have a header row. They can be build and manipulated as raw Python datatypes (Lists of tuples|dictionaries). Datasets can be exported to JSON, YAML, CSV, TSV, and XML.

### tablib.Databook
A Databook is a set of Datasets. The most common form of a Databook is an Excel file with multiple spreadsheets. Databooks can be exported to JSON, YAML and XML.

## Usage

Creates a dataset and populate it:

```go
ds := NewDataset([]string{"firstName", "lastName"})
```

Add new rows:
```go
ds.Append([]interface{}{"John", "Adams"})
ds.AppendValues("George", "Washington")
```

Add new columns:
```go
ds.AppendColumn("age", []interface{}{90, 67})
ds.AppendColumnValues("sex", "male", "male")
```

Add a dynamic column, by passing a function which has access to the current row, and must
return a value:
```go
func lastNameLen(row []interface{}) interface{} {
	return len(row[1].(string))
}
ds.AppendDynamicColumn("lastName length", lastNameLen)
ds.CSV()
// >>
// firstName, lastName, age, sex, lastName length
// John, Adams, 90, male, 5
// George, Washington, 67, male, 10
```

Delete rows:
```go
ds.DeleteRow(1) // starts at 0
```

Delete columns:
```go
ds.DeleteColumn("sex")
```

Get a row or multiple rows:
```go
row, _ := ds.Row(0)
fmt.Println(row["firstName"]) // George

rows, _ := ds.Rows(0, 1)
fmt.Println(rows[0]["firstName"]) // George
fmt.Println(rows[1]["firstName"]) // Thomas
```

Slice a Dataset:
```go
newDs, _ := ds.Slice(1, 5) // returns a fresh Dataset with rows [1..5[
```


## Filtering

You can add **tags** to rows by using a specific `Dataset` method. This allows you to filter your `Dataset` later. This can be useful to separate rows of data based on arbitrary criteria (e.g. origin) that you don’t want to include in your `Dataset`.
```go
ds := NewDataset([]string{"Maker", "Model"})
ds.AppendTagged([]interface{}{"Porsche", "911"}, "fast", "luxury")
ds.AppendTagged([]interface{}{"Skoda", "Octavia"}, "family")
ds.AppendTagged([]interface{}{"Ferrari", "458"}, "fast", "luxury")
ds.AppendValues("Citroen", "Picasso")
ds.AppendTagged([]interface{}{"Bentley", "Continental"}, "luxury")
```

Filtering the `Dataset` is possible by calling `Filter(column)`:
```go
luxuryCars := ds.Filter("luxury").CSV()
// >>>
// Maker,Model
// Porsche,911
// Ferrari,458
// Bentley,Continental
```

```go
fastCars := ds.Filter("fast").CSV()
// >>>
// Maker,Model
// Porsche,911
// Ferrari,458
```

## Sorting

Datasets can be sorted by a specific column.
```go
ds := NewDataset([]string{"Maker", "Model", "Year"})
ds.AppendValues("Porsche", "991", 2012)
ds.AppendValues("Skoda", "Octavia", 2011)
ds.AppendValues("Ferrari", "458", 2009)
ds.AppendValues("Citroen", "Picasso II", 2013)
ds.AppendValues("Bentley", "Continental GT", 2003)

ds.Sort("Year").CSV()
// >>
// Maker, Model, Year
// Bentley, Continental GT, 2003
// Ferrari, 458, 2009
// Skoda, Octavia, 2011
// Porsche, 991, 2012
// Citroen, Picasso II, 2013
```

## Loading

### JSON
```go
ds, _ := LoadJSON([]byte(`[
  {"age":90,"firstName":"John","lastName":"Adams"},
  {"age":67,"firstName":"George","lastName":"Washington"},
  {"age":83,"firstName":"Henry","lastName":"Ford"}
]`))
```

### YAML
```go
ds, _ := LoadYAML([]byte(`- age: 90
  firstName: John
  lastName: Adams
- age: 67
  firstName: George
  lastName: Washington
- age: 83
  firstName: Henry
  lastName: Ford`))
```

## Exports

### JSON
```go
json, _ := ds.JSON()
fmt.Printf("%s\n", json)
```

Will output:
```json
[{"age":90,"firstName":"John","lastName":"Adams"},{"age":67,"firstName":"George","lastName":"Washington"},{"age":83,"firstName":"Henry","lastName":"Ford"}]
```

### XML
```go
xml := ds.XML()
fmt.Printf("%s\n", xml)
```

Will ouput:
```xml
<dataset>
 <row>
   <age>90</age>
   <firstName>John</firstName>
   <lastName>Adams</lastName>
 </row>  <row>
   <age>67</age>
   <firstName>George</firstName>
   <lastName>Washington</lastName>
 </row>  <row>
   <age>83</age>
   <firstName>Henry</firstName>
   <lastName>Ford</lastName>
 </row>
</dataset>
```

### CSV
```go
csv, _ := ds.CSV()
fmt.Printf("%s\n", csv)
```

Will ouput:
```csv
firstName,lastName,age
John,Adams,90
George,Washington,67
Henry,Ford,83
```

### TSV
```go
tsv, _ := ds.TSV()
fmt.Printf("%s\n", tsv)
```

Will ouput:
```tsv
firstName lastName  age
John  Adams  90
George  Washington  67
Henry Ford 83
```

### YAML
```go
yaml, _ := ds.YAML()
fmt.Printf("%s\n", yaml)
```

Will ouput:
```yaml
- age: 90
  firstName: John
  lastName: Adams
- age: 67
  firstName: George
  lastName: Washington
- age: 83
  firstName: Henry
  lastName: Ford
```

### HTML
```go
html := ds.HTML()
fmt.Printf("%s\n", html)
```

Will output:
```html
<table class="table table-striped">
	<thead>
		<tr>
			<th>firstName</th>
			<th>lastName</th>
			<th>age</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>George</td>
			<td>Washington</td>
			<td>90</td>
		</tr>
		<tr>
			<td>Henry</td>
			<td>Ford</td>
			<td>67</td>
		</tr>
		<tr>
			<td>Foo</td>
			<td>Bar</td>
			<td>83</td>
		</tr>
	</tbody>
</table>
```

### XLSX
```go
xlsx, _ := ds.XLSX()
// >>>
// binary content
```

### ASCII

#### Grid format
```go
ascii := ds.Tabular("grid")
```

Will output:
```
+--------------+---------------+--------+
|    firstName |      lastName |    age |
+==============+===============+========+
|       George |    Washington |     90 |
+--------------+---------------+--------+
|        Henry |          Ford |     67 |
+--------------+---------------+--------+
|          Foo |           Bar |     83 |
+--------------+---------------+--------+
```

#### Simple format
```go
ascii := ds.Tabular("simple")
```

Will output:
```
--------------  ---------------  --------
    firstName         lastName       age
--------------  ---------------  --------
       George       Washington        90

        Henry             Ford        67

          Foo              Bar        83
--------------  ---------------  --------
```

## Databooks

This is an example of how to use Databooks.

```go
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

fmt.Println(db.JSON())
```

Will output the following JSON representation of the Databook:
```json
[
  {
    "title": "Cars",
    "data": [
      {"Maker":"Bentley","Model":"Continental GT","Year":2003},
      {"Maker":"Ferrari","Model":"458","Year":2009},
      {"Maker":"Skoda","Model":"Octavia","Year":2011},
      {"Maker":"Porsche","Model":"991","Year":2012},
      {"Maker":"Citroen","Model":"Picasso II","Year":2013}
    ]
  },
  {
    "title": "Presidents",
    "data": [
      {"Age":90,"First name":"John","Last name":"Adams"},
      {"Age":83,"First name":"Henry","Last name":"Ford"},
      {"Age":67,"First name":"George","Last name":"Washington"}
    ]
  }
]
```

## Installation

```bash
go get github.com/agrison/go-tablib
```

## TODO

* Loading in more formats
* Support more formats: DBF, XLS, LATEX, ...

## Contribute
PRs more than welcomed, come and join :)

## Acknowledgement

Thanks to kennethreitz for the first implementation in Python, [`github.com/bndr/gotabulate`](https://github.com/bndr/gotabulate), [`github.com/clbanning/mxj`](https://github.com/clbanning/mxj), [`github.com/tealeg/xlsx`](https://github.com/tealeg/xlsx), [`gopkg.in/yaml.v2`](https://gopkg.in/yaml.v2)
