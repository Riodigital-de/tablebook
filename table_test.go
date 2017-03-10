package tablebook

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTable(t *testing.T) {
	Convey("Given a empty table", t, func() {
		table, _ := NewTable("table", []string{"foo", "bar", "baz"})

		Convey("It can create a table", func() {
			// duplicated columns
			_, err := NewTable("table", []string{"foo", "foo"})
			So(err, ShouldEqual, ErrColumnExists)
		})

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

		Convey("It can get columnNames", func() {
			So(table.ColumnNames(), ShouldResemble, []string{"foo", "bar", "baz"})
		})

		Convey("It can rename columns", func() {
			// rename one
			So(table.RenameColumn("unknown", "newColumnName"), ShouldEqual, ErrNotFound)

			// rename existing column
			So(table.RenameColumn("bar", "foo"), ShouldEqual, ErrColumnExists)

			table.RenameColumn("bar", "new")
			So(table.ColumnNames(), ShouldResemble, []string{"foo", "new", "baz"})
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

			So(err1, ShouldEqual, ErrInvalidDimensions)
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
			So(table.ColumnNames()[3], ShouldEqual, "new_column")

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

			table.AppendEvaluatedColumn(
				"new_dynamic_column_1",
				func(table *Table, rowIndex int, cellIndex int) interface{} {
					return "from_new_dynamic_column_1"
				},
			)

			table.AppendEvaluatedColumn(
				"new_dynamic_column_2",
				func(table *Table, rowIndex int, cellIndex int) interface{} {
					v, _ := table.Cell(rowIndex, 3)
					return v
				},
			)

			table.AppendEvaluatedColumn(
				"new_dynamic_column_3",
				func(table *Table, rowIndex int, cellIndex int) interface{} {
					v, _ := table.Cell(rowIndex, 4)
					return v
				},
			)

			targetTable, _ := NewTable("targetTable", []string{"foo", "baz", "new_dynamic_column_3"})
			targetTable.Take([]*Table{table})

			// column already exists
			So(table.AppendEvaluatedColumn(
				"foo",
				func(table *Table, rowIndex int, cellIndex int) interface{} {
					return nil
				},
			), ShouldEqual, ErrColumnExists)

			// ok
			So(table.ColumnNames(), ShouldResemble, []string{"foo", "bar", "baz", "new_dynamic_column_1", "new_dynamic_column_2", "new_dynamic_column_3"})
			So(targetTable.ColumnNames(), ShouldResemble, []string{"foo", "baz", "new_dynamic_column_3"})
			So(
				targetTable.Rows(),
				ShouldResemble,
				[][]interface{}{
					{1, 3, "from_new_dynamic_column_1"},
					{5, 7, "from_new_dynamic_column_1"},
					{9, 11, "from_new_dynamic_column_1"},
				},
			)
		})

		Convey("It can take columns from other tables and merge it into own columns", func() {
			targetTable, _ := NewTable("targetTable", []string{"foo", "baz"})
			targetTable.AppendRow([]interface{}{"foo 1 targetTable", "baz 1 targetTable"})

			foreignTable1, _ := NewTable("foreignTable1", []string{"foo", "bar", "unknown"})
			foreignTable1.AppendRow([]interface{}{"foo 1 foreignTable1", "bar 1 foreignTable1", "baz 1 foreignTable1"})
			foreignTable1.AppendRow([]interface{}{"foo 2 foreignTable1", "bar 2 foreignTable1", "baz 2 foreignTable1"})
			foreignTable1.RenameColumn("unknown", "baz")

			foreignTable2, _ := NewTable("foreignTable2", []string{"bar", "foo"})
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
