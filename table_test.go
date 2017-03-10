package tablebook

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTable(t *testing.T) {
	Convey("Given a empty table", t, func() {
		table := NewTable("table", []string{"foo", "bar", "baz"})

		Convey("It has a Name", func() {
			// ok
			So(table.Name(), ShouldEqual, "table")
		})

		Convey("It has dimensions", func() {
			table.AppendRow([]interface{}{1, "foo", 2})
			table.AppendRow([]interface{}{1, "foo", 2})

			So(table.Width(), ShouldEqual, 3)
			So(table.Height(), ShouldEqual, 2)
		})

		Convey("It can get headers", func() {
			So(table.Headers(), ShouldResemble, []string{"foo", "bar", "baz"})
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

			//ok
			So(table.Rows(), ShouldResemble,
				[][]interface{}{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			)

			tr1, err1 := table.Row(1000)

			So(err1, ShouldEqual, ErrNotFound)
			So(tr1, ShouldBeNil)

			tr2, _ := table.Row(1)
			So(tr2, ShouldResemble, []interface{}{4, 5, 6})
		})

		Convey("It can append rows", func() {
			// too much columns
			So(table.AppendRow([]interface{}{1, "foo", "bar", "too_much"}), ShouldEqual, ErrInvalidDimensions)

			// ok, not enough columns
			table.AppendRow([]interface{}{1})
			tr1, _ := table.Row(0)
			So(tr1, ShouldResemble, []interface{}{1, "", ""})

			// ok
			tbl2 := table.AppendRow([]interface{}{1, "foo", 2})
			tr2, _ := table.Row(1)
			So(tr2, ShouldResemble, []interface{}{1, "foo", 2})
			So(tbl2, ShouldBeNil)
		})

		Convey("It can get columns", func() {
			table.AppendRow([]interface{}{1, 2, 3})
			table.AppendRow([]interface{}{4, 5, 6})
			table.AppendRow([]interface{}{7, 8, 9})

			col1, err1 := table.Column("bar")
			col2, err2 := table.Column("unknown")

			So(col1, ShouldResemble, []interface{}{2, 5, 8})
			So(err1, ShouldBeNil)

			So(col2, ShouldBeNil)
			So(err2, ShouldEqual, ErrNotFound)
		})

		Convey("It can append columns", func() {
			table.AppendRow([]interface{}{1, 2, 3})
			table.AppendRow([]interface{}{5, 6, 7})
			table.AppendRow([]interface{}{9, 10, 11})

			// too much rows
			So(table.AppendColumn("new_column", []interface{}{1, "foo", "bar", "too_much"}), ShouldEqual, ErrInvalidDimensions)

			// column already exists
			So(table.AppendColumn("foo", nil), ShouldEqual, ErrColumnExists)

			// ok
			table.AppendColumn("new_column", []interface{}{4, 8})
			So(table.Headers()[3], ShouldEqual, "new_column")

			// ok
			So(
				table.Rows(),
				ShouldResemble,
				[][]interface{}{
					{1, 2, 3, 4},
					{5, 6, 7, 8},
					{9, 10, 11, ""},
				},
			)

			// ok
			So(table.AppendRow([]interface{}{1, "foo", 2}), ShouldBeNil)
		})

		Convey("It can append dynamic columns", func() {
			table.AppendRow([]interface{}{1, 2, 3})
			table.AppendRow([]interface{}{5, 6, 7})
			table.AppendRow([]interface{}{9, 10, 11})

			dynamic_sum := func(table *Table, ri int, ci int, row []interface{}) interface{} {
				return row[0].(int) + row[1].(int) + row[2].(int)
			}

			// column already exists
			So(table.AppendDynamicColumn("foo", dynamic_sum), ShouldEqual, ErrColumnExists)

			// ok
			table.AppendDynamicColumn("new_dynamic_column", dynamic_sum)
			So(table.Headers()[3], ShouldEqual, "new_dynamic_column")
			So(
				table.Rows(),
				ShouldResemble,
				[][]interface{}{
					{1, 2, 3, 6},
					{5, 6, 7, 18},
					{9, 10, 11, 30},
				},
			)
		})

		Convey("It can take columns from other tables and merge it into own columns", func() {
			targetTable := NewTable("targetTable", []string{"foo", "baz"})
			targetTable.AppendRow([]interface{}{"foo 1 targetTable", "baz 1 targetTable"})

			foreignTable1 := NewTable("foreignTable1", []string{"foo", "bar", "unknown"})
			foreignTable1.AppendRow([]interface{}{"foo 1 foreignTable1", "bar 1 foreignTable1", "baz 1 foreignTable1"})
			foreignTable1.AppendRow([]interface{}{"foo 2 foreignTable1", "bar 2 foreignTable1", "baz 2 foreignTable1"})
			foreignTable1.RenameHeader("unknown", "baz")

			foreignTable2 := NewTable("foreignTable2", []string{"bar", "foo"})
			foreignTable2.AppendRow([]interface{}{"bar 1 foreignTable2", "foo 1 foreignTable2"})
			foreignTable2.AppendRow([]interface{}{"bar 2 foreignTable2", "foo 2 foreignTable2"})

			targetTable.Take([]*Table{foreignTable1, foreignTable2})

			// ok
			So(
				targetTable.Rows(),
				ShouldResemble,
				[][]interface{}{
					{"foo 1 targetTable", "baz 1 targetTable"},
					{"foo 1 foreignTable1", "baz 1 foreignTable1"},
					{"foo 2 foreignTable1", "baz 2 foreignTable1"},
					{"foo 1 foreignTable2", ""},
					{"foo 2 foreignTable2", ""},
				},
			)

		})
	})
}
