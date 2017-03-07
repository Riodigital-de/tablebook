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

			//ok
			So(table.AppendRow([]interface{}{1, "foo", 2}), ShouldBeNil)
		})

		Convey("It has dimensions", func() {
			table.AppendRow([]interface{}{1, "foo", 2})
			table.AppendRow([]interface{}{1, "foo", 2})

			So(table.Width(), ShouldEqual, 3)
			So(table.Height(), ShouldEqual, 2)
		})

		Convey("It can rename headers", func() {
			//rename one
			So(table.RenameHeader("foo", "bar"), ShouldBeNil)
			So(table.headers, ShouldResemble, []string{"bar", "bar", "baz"})

			//rename multiple
			So(table.RenameHeader("bar", "foo"), ShouldBeNil)
			So(table.headers, ShouldResemble, []string{"foo", "foo", "baz"})

			//unknown header
			So(table.RenameHeader("unknown", "baz"), ShouldEqual, ErrNotFound)
		})

		Convey("It can get columns", func() {
			table.AppendRow([]interface{}{1, 2, 3})
			table.AppendRow([]interface{}{4, 5, 6})
			table.AppendRow([]interface{}{7, 8, 9})

			//ok
			ok, _ := table.Column("foo")
			So(ok, ShouldResemble, []interface{}{1, 4, 7})

			//unknown header
			_, err := table.Column("unknown")
			So(err, ShouldEqual, ErrNotFound)
		})

		Convey("It can take columns from other tables and merge it into own columns", func() {
			targetTable := NewTable("targetTable", []string{"foo", "baz"})
			targetTable.AppendRow([]interface{}{"foo targetTable", "baz targetTable"})

			foreignTable1 := NewTable("foreignTable1", []string{"foo", "bar", "baz"})
			foreignTable1.AppendRow([]interface{}{"foo foreignTable1", "bar foreignTable1", "baz foreignTable1"})

			foreignTable2 := NewTable("foreignTable2", []string{"baz", "bar", "foo"})
			foreignTable2.AppendRow([]interface{}{"baz foreignTable2", "bar foreignTable2", "foo foreignTable2"})

			targetTable.Take([]*Table{foreignTable1, foreignTable2})

			//it is the first row we added to the table... no taking
			So(targetTable.rows[0], ShouldResemble, []interface{}{"foo targetTable", "baz targetTable"})

			//it is the second row and it takes reflect foreignTable1 rows expect bar
			So(targetTable.rows[1], ShouldResemble, []interface{}{"foo foreignTable1", "baz foreignTable1"})

			//it is the second row and it takes reflect foreignTable2 rows expect bar
			So(targetTable.rows[2], ShouldResemble, []interface{}{"foo foreignTable2", "baz foreignTable2"})

			//it shows the first column (foo) and is composed of targetTable, foreignTable1, foreignTable2
			colFoo, _ := targetTable.Column("foo")
			So(colFoo, ShouldResemble, []interface{}{"foo targetTable", "foo foreignTable1", "foo foreignTable2"})

			//it shows the second and final column (baz) and is composed of targetTable, foreignTable1, foreignTable2
			colBaz, _ := targetTable.Column("baz")
			So(colBaz, ShouldResemble, []interface{}{"baz targetTable", "baz foreignTable1", "baz foreignTable2"})
		})
	})
}
