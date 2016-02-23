package tablib_test

import (
	"encoding/base64"
	tablib "github.com/agrison/go-tablib"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TablibSuite struct{}

var _ = Suite(&TablibSuite{})

func presidentDataset() *tablib.Dataset {
	ds := tablib.NewDataset([]string{"firstName", "lastName", "gpa"})
	ds.AppendValues("John", "Adams", 90)
	ds.AppendValues("George", "Washington", 67)
	ds.AppendValues("Thomas", "Jefferson", 50)
	return ds
}

func presidentDatasetWithTags() *tablib.Dataset {
	ds := tablib.NewDataset([]string{"firstName", "lastName", "gpa"})
	ds.AppendTagged([]interface{}{"John", "Adams", 90}, "Massachusetts")
	ds.AppendTagged([]interface{}{"George", "Washington", 67}, "Virginia")
	ds.AppendTagged([]interface{}{"Thomas", "Jefferson", 50}, "Virginia")
	return ds
}

func frenhPresidentDataset() *tablib.Dataset {
	ds := tablib.NewDataset([]string{"firstName", "lastName", "gpa"})
	ds.AppendValues("Jacques", "Chirac", 88)
	ds.AppendValues("Nicolas", "Sarkozy", 98)
	ds.AppendValues("François", "Hollande", 34)
	return ds
}

func frenhPresidentAdditionalDataset() *tablib.Dataset {
	ds := tablib.NewDataset([]string{"duration", "from"})
	ds.AppendValues(14, "Paris")
	ds.AppendValues(12, "Paris")
	ds.AppendValues(5, "Rouen")
	return ds
}

func rowAt(d *tablib.Dataset, index int) map[string]interface{} {
	return d.Dict()[index].(map[string]interface{})
}

func lastRow(d *tablib.Dataset) map[string]interface{} {
	return d.Dict()[d.Height()-1].(map[string]interface{})
}

func (s *TablibSuite) TestDimensions(c *C) {
	ds := presidentDataset()
	c.Assert(ds.Width(), Equals, 3)
	c.Assert(ds.Height(), Equals, 3)
	c.Assert(ds.Headers(), DeepEquals, []string{"firstName", "lastName", "gpa"})
}

func (s *TablibSuite) TestAppendRow(c *C) {
	ds := presidentDataset()
	// too much columns
	c.Assert(ds.AppendValues("a", "b", 50, "d"), Equals, tablib.ErrInvalidDimensions)
	// not enough columns
	c.Assert(ds.AppendValues("a", "b"), Equals, tablib.ErrInvalidDimensions)
	// ok
	c.Assert(ds.AppendValues("foo", "bar", 42), Equals, nil)
	// test values are there
	d := lastRow(ds)
	c.Assert(d["firstName"], Equals, "foo")
	c.Assert(d["lastName"], Equals, "bar")
	c.Assert(d["gpa"], Equals, 42)
}

func (s *TablibSuite) TestAppendColumn(c *C) {
	ds := presidentDataset()
	// too much rows
	c.Assert(ds.AppendColumnValues("foo", "a", "b", "c", "d"), Equals, tablib.ErrInvalidDimensions)
	// not enough columns
	c.Assert(ds.AppendColumnValues("foo", "a", "b"), Equals, tablib.ErrInvalidDimensions)
	// ok
	c.Assert(ds.AppendColumnValues("foo", "a", "b", "c"), Equals, nil)
	// test values are there
	d := ds.Column("foo")
	c.Assert(d[0], Equals, "a")
	c.Assert(d[1], Equals, "b")
	c.Assert(d[2], Equals, "c")
}

func (s *TablibSuite) TestInsert(c *C) {
	ds := presidentDataset()
	// invalid index
	c.Assert(ds.InsertValues(-1, "foo", "bar"), Equals, tablib.ErrInvalidRowIndex)
	c.Assert(ds.InsertValues(100, "foo", "bar"), Equals, tablib.ErrInvalidRowIndex)
	// too much columns
	c.Assert(ds.InsertValues(1, "foo", "bar", 42, "invalid"), Equals, tablib.ErrInvalidDimensions)
	// not enough columns
	c.Assert(ds.InsertValues(1, "foo", "bar"), Equals, tablib.ErrInvalidDimensions)
	// ok
	c.Assert(ds.InsertValues(1, "foo", "bar", 42), Equals, nil)
	// test values are there
	d := rowAt(ds, 1)
	c.Assert(d["firstName"], Equals, "foo")
	c.Assert(d["lastName"], Equals, "bar")
	c.Assert(d["gpa"], Equals, 42)
}

func (s *TablibSuite) TestInsertColumn(c *C) {
	ds := presidentDataset()
	// invalid index
	c.Assert(ds.InsertColumn(-1, "wut", []interface{}{"foo", "bar"}), Equals, tablib.ErrInvalidColumnIndex)
	c.Assert(ds.InsertColumn(100, "wut", []interface{}{"foo", "bar"}), Equals, tablib.ErrInvalidColumnIndex)
	// too much rows
	c.Assert(ds.InsertColumn(1, "wut", []interface{}{"foo", "bar", "baz", "kidding"}), Equals, tablib.ErrInvalidDimensions)
	// not enough rows
	c.Assert(ds.InsertColumn(1, "wut", []interface{}{"foo", "bar"}), Equals, tablib.ErrInvalidDimensions)
	// ok
	c.Assert(ds.InsertColumn(1, "wut", []interface{}{"foo", "bar", "baz"}), Equals, nil)
	// test values are there
	d := ds.Column("wut")
	c.Assert(d[0], Equals, "foo")
	c.Assert(d[1], Equals, "bar")
	c.Assert(d[2], Equals, "baz")
}

func firstNameB64(row []interface{}) interface{} {
	return base64.StdEncoding.EncodeToString([]byte(row[0].(string)))
}

func lastNameB64(row []interface{}) interface{} {
	return base64.StdEncoding.EncodeToString([]byte(row[1].(string)))
}

func (s *TablibSuite) TestDynamicColumn(c *C) {
	ds := presidentDataset()
	ds.AppendDynamicColumn("firstB64", firstNameB64)
	d := ds.Column("firstB64")
	c.Assert(d[0], Equals, "Sm9obg==") // John
	c.Assert(d[1], Equals, "R2Vvcmdl") // George
	c.Assert(d[2], Equals, "VGhvbWFz") // Thomas

	// invalid index
	c.Assert(ds.InsertDynamicColumn(-1, "foo", lastNameB64), Equals, tablib.ErrInvalidColumnIndex)
	c.Assert(ds.InsertDynamicColumn(100, "foo", lastNameB64), Equals, tablib.ErrInvalidColumnIndex)
	// ok
	c.Assert(ds.InsertDynamicColumn(2, "lastB64", lastNameB64), Equals, nil)
	// check values
	d = ds.Column("lastB64")
	c.Assert(d[0], Equals, "QWRhbXM=")         // Adams
	c.Assert(d[1], Equals, "V2FzaGluZ3Rvbg==") // Washington
	c.Assert(d[2], Equals, "SmVmZmVyc29u")     // Jefferson
}

func (s *TablibSuite) TestStack(c *C) {
	ds, _ := presidentDataset().Stack(frenhPresidentDataset())
	d := ds.Column("lastName")
	c.Assert(d[0], Equals, "Adams")
	c.Assert(d[1], Equals, "Washington")
	c.Assert(d[2], Equals, "Jefferson")
	c.Assert(d[3], Equals, "Chirac")
	c.Assert(d[4], Equals, "Sarkozy")
	c.Assert(d[5], Equals, "Hollande")

	// check invalid dimensions
	x := frenhPresidentDataset()
	x.DeleteColumn("lastName")
	ds, err := presidentDataset().Stack(x)
	c.Assert(err, Equals, tablib.ErrInvalidDimensions)
}

func (s *TablibSuite) TestStackColumn(c *C) {
	ds, _ := frenhPresidentDataset().StackColumn(frenhPresidentAdditionalDataset())
	d := lastRow(ds)
	c.Assert(d["firstName"], Equals, "François")
	c.Assert(d["lastName"], Equals, "Hollande")
	c.Assert(d["from"], Equals, "Rouen")
	c.Assert(d["duration"], Equals, 5)

	// check invalid dimensions
	x := frenhPresidentAdditionalDataset()
	x.DeleteRow(x.Height() - 1)
	ds, err := frenhPresidentDataset().StackColumn(x)
	c.Assert(err, Equals, tablib.ErrInvalidDimensions)
}

func (s *TablibSuite) TestFiltering(c *C) {
	ds := presidentDatasetWithTags()
	df := ds.Filter("Massachusetts")
	c.Assert(df.Height(), Equals, 1)
	r := lastRow(df)
	c.Assert(r["firstName"], Equals, "John")
	c.Assert(r["lastName"], Equals, "Adams")

	df = ds.Filter("Virginia")
	c.Assert(df.Height(), Equals, 2)
	r = rowAt(df, 0)
	c.Assert(r["firstName"], Equals, "George")
	c.Assert(r["lastName"], Equals, "Washington")
	r = rowAt(df, 1)
	c.Assert(r["firstName"], Equals, "Thomas")
	c.Assert(r["lastName"], Equals, "Jefferson")

	df = ds.Filter("Woot")
	c.Assert(df.Height(), Equals, 0)
	c.Assert(df.Width(), Equals, 3)
}

func (s *TablibSuite) TestSort(c *C) {
	ds := presidentDataset().Sort("gpa")
	c.Assert(ds.Height(), Equals, 3)

	r := rowAt(ds, 0)
	c.Assert(r["firstName"], Equals, "Thomas")
	c.Assert(r["lastName"], Equals, "Jefferson")
	c.Assert(r["gpa"], Equals, 50)

	r = rowAt(ds, 1)
	c.Assert(r["firstName"], Equals, "George")
	c.Assert(r["lastName"], Equals, "Washington")
	c.Assert(r["gpa"], Equals, 67)

	r = rowAt(ds, 2)
	c.Assert(r["firstName"], Equals, "John")
	c.Assert(r["lastName"], Equals, "Adams")
	c.Assert(r["gpa"], Equals, 90)

	ds = ds.SortReverse("lastName")
	c.Assert(ds.Height(), Equals, 3)

	r = rowAt(ds, 0)
	c.Assert(r["firstName"], Equals, "George")
	c.Assert(r["lastName"], Equals, "Washington")

	r = rowAt(ds, 1)
	c.Assert(r["firstName"], Equals, "Thomas")
	c.Assert(r["lastName"], Equals, "Jefferson")

	r = rowAt(ds, 2)
	c.Assert(r["firstName"], Equals, "John")
	c.Assert(r["lastName"], Equals, "Adams")
}

func (s *TablibSuite) TestJSON(c *C) {
	ds := frenhPresidentDataset()
	j, _ := ds.JSON()
	c.Assert(j, Equals, "[{\"firstName\":\"Jacques\",\"gpa\":88,\"lastName\":\"Chirac\"},{\"firstName\":\"Nicolas\",\"gpa\":98,\"lastName\":\"Sarkozy\"},{\"firstName\":\"François\",\"gpa\":34,\"lastName\":\"Hollande\"}]")
}

func (s *TablibSuite) TestYAML(c *C) {
	ds := frenhPresidentDataset()
	j, _ := ds.YAML()
	c.Assert(j, Equals, `- firstName: Jacques
  gpa: 88
  lastName: Chirac
- firstName: Nicolas
  gpa: 98
  lastName: Sarkozy
- firstName: François
  gpa: 34
  lastName: Hollande
`)
}

func (s *TablibSuite) TestCSV(c *C) {
	ds := frenhPresidentDataset()
	j, _ := ds.CSV()
	c.Assert(j, Equals, `firstName,lastName,gpa
Jacques,Chirac,88
Nicolas,Sarkozy,98
François,Hollande,34
`)
}

func (s *TablibSuite) TestTSV(c *C) {
	ds := frenhPresidentDataset()
	j, _ := ds.TSV()
	c.Assert(j, Equals, `firstName`+"\t"+`lastName`+"\t"+`gpa
Jacques`+"\t"+`Chirac`+"\t"+`88
Nicolas`+"\t"+`Sarkozy`+"\t"+`98
François`+"\t"+`Hollande`+"\t"+`34
`)
}

func (s *TablibSuite) TestHTML(c *C) {
	ds := frenhPresidentDataset()
	j := ds.HTML()
	c.Assert(j, Equals, `<table class="table table-striped">
	<thead>
		<tr>
			<th>firstName</th>
			<th>lastName</th>
			<th>gpa</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>Jacques</td>
			<td>Chirac</td>
			<td>88</td>
		</tr>
		<tr>
			<td>Nicolas</td>
			<td>Sarkozy</td>
			<td>98</td>
		</tr>
		<tr>
			<td>François</td>
			<td>Hollande</td>
			<td>34</td>
		</tr>
	</tbody>
</table>`)
}

func (s *TablibSuite) TestTabular(c *C) {
	ds := frenhPresidentDataset()
	j := ds.Tabular("grid")
	c.Assert(j, Equals, `+--------------+-------------+--------+
|    firstName |    lastName |    gpa |
+==============+=============+========+
|      Jacques |      Chirac |     88 |
+--------------+-------------+--------+
|      Nicolas |     Sarkozy |     98 |
+--------------+-------------+--------+
|     François |    Hollande |     34 |
+--------------+-------------+--------+
`)

	j = ds.Tabular("simple")
	c.Assert(j, Equals, `--------------  -------------  --------`+"\n"+
		`    firstName       lastName       gpa `+"\n"+
		`--------------  -------------  --------`+"\n"+
		`      Jacques         Chirac        88 `+"\n"+
		"\n"+
		`      Nicolas        Sarkozy        98 `+"\n"+
		"\n"+
		`     François       Hollande        34 `+"\n"+
		`--------------  -------------  --------`+
		"\n")
}
