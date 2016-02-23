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
