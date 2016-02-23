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
	ds.AppendDynamicColumn("lastB64", lastNameB64)
	d := ds.Column("lastB64")
	c.Assert(d[0], Equals, "Sm9obg==")
	c.Assert(d[1], Equals, "R2Vvcmdl")
	c.Assert(d[2], Equals, "VGhvbWFz")
}
