package tablebook

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTable(t *testing.T) {
	Convey("Given a empty table", t, func() {
		table := NewTable("table", []string{"foo", "bar", "baz"})

		Convey("It can append rows", func() {
			// too much columns
			So(table.AppendRow([]interface{}{1, "foo", "bar", "too_much"}), ShouldEqual, ErrInvalidDimensions)

			// not enough columns
			So(table.AppendRow([]interface{}{1}), ShouldEqual, ErrInvalidDimensions)

			// ok
			So(table.AppendRow([]interface{}{1, "foo", 2}), ShouldBeNil)
		})

		Convey("It has dimensions", func() {
			table.AppendRow([]interface{}{1, "foo", 2})
			table.AppendRow([]interface{}{1, "foo", 2})

			So(table.Width(), ShouldEqual, 3)
			So(table.Height(), ShouldEqual, 2)
		})

		Convey("It can rename headers", func() {
			// rename one
			So(table.RenameHeader("foo", "bar"), ShouldBeNil)
			So(table.Headers(), ShouldResemble, []string{"bar", "bar", "baz"})

			// rename multiple
			So(table.RenameHeader("bar", "foo"), ShouldBeNil)
			So(table.Headers(), ShouldResemble, []string{"foo", "foo", "baz"})

			// unknown header
			So(table.RenameHeader("unknown", "baz"), ShouldEqual, ErrNotFound)
		})

		Convey("It can get rows", func() {
			table.AppendRow([]interface{}{1, 2, 3})
			table.AppendRow([]interface{}{4, 5, 6})
			table.AppendRow([]interface{}{7, 8, 9})

			So(len(table.Rows()), ShouldEqual, 3)
		})

		Convey("It can get headers", func() {
			So(table.Headers(), ShouldResemble, []string{"foo", "bar", "baz"})
		})

		Convey("It can take columns from other tables and merge it into own columns", func() {
			targetTable := NewTable("targetTable", []string{"foo", "baz"})
			targetTable.AppendRow([]interface{}{"foo 1 targetTable", "baz 1 targetTable"})

			foreignTable1 := NewTable("foreignTable1", []string{"foo", "bar", "baz"})
			foreignTable1.AppendRow([]interface{}{"foo 1 foreignTable1", "bar 1 foreignTable1", "baz 1 foreignTable1"})
			foreignTable1.AppendRow([]interface{}{"foo 2 foreignTable1", "bar 2 foreignTable1", "baz 2 foreignTable1"})

			foreignTable2 := NewTable("foreignTable2", []string{"bar", "foo"})
			foreignTable2.AppendRow([]interface{}{"bar 1 foreignTable2", "foo 1 foreignTable2"})
			foreignTable2.AppendRow([]interface{}{"bar 2 foreignTable2", "foo 2 foreignTable2"})

			targetTable.Take([]*Table{foreignTable1, foreignTable2})

			So(len(targetTable.Rows()), ShouldEqual, 5)

			// no foreign merge
			So(targetTable.Rows()[0], ShouldResemble, []interface{}{"foo 1 targetTable", "baz 1 targetTable"})

			// merge foo, baz from foreignTable1
			So(targetTable.Rows()[1], ShouldResemble, []interface{}{"foo 1 foreignTable1", "baz 1 foreignTable1"})

			// merge foo, baz from foreignTable2.
			// foreignTable2 has no baz... it has a empty string as replacement
			So(targetTable.Rows()[3], ShouldResemble, []interface{}{"foo 1 foreignTable2", ""})

			// shows the first column (foo) and is composed of targetTable, foreignTable1, foreignTable2
			So(
				[]interface{}{targetTable.Rows()[0][0], targetTable.Rows()[1][0], targetTable.Rows()[2][0], targetTable.Rows()[3][0], targetTable.Rows()[4][0]},
				ShouldResemble,
				[]interface{}{"foo 1 targetTable", "foo 1 foreignTable1", "foo 2 foreignTable1", "foo 1 foreignTable2", "foo 2 foreignTable2"},
			)

			// shows the second column (baz) and is composed of targetTable, foreignTable1, foreignTable2.
			// foreignTable2 has no baz... it has a empty string as replacement
			So(
				[]interface{}{targetTable.Rows()[0][1], targetTable.Rows()[1][1], targetTable.Rows()[2][1], targetTable.Rows()[3][1], targetTable.Rows()[4][1]},
				ShouldResemble,
				[]interface{}{"baz 1 targetTable", "baz 1 foreignTable1", "baz 2 foreignTable1", "", ""},
			)
		})
	})
}
