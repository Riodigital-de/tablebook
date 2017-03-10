package tablebook

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBook(t *testing.T) {
	Convey("Given a empty book", t, func() {
		book := NewBook()

		Convey("It can create tables", func() {
			_, err1 := book.NewTable("table", []string{"foo", "bar", "baz"})
			So(err1, ShouldBeNil)

			_, err2 := book.NewTable("table", []string{"foo", "bar", "baz"})
			So(err2, ShouldEqual, ErrTableExists)

			_, err3 := book.NewTable("table_duplicate_columns", []string{"bar", "bar", "baz"})
			So(err3, ShouldEqual, ErrColumnExists)
		})

		Convey("It can search tables", func() {
			table1, _ := book.NewTable("table1", []string{"foo", "bar", "baz"})
			table2, _ := book.NewTable("table2", []string{"foo", "bar", "baz"})

			t1, err1 := book.Table("table1")
			t2, err2 := book.Table("table2")
			t3, err3 := book.Table("tableunknown")

			// ok
			So(t1, ShouldResemble, table1)
			So(err1, ShouldBeNil)

			// ok
			So(t2, ShouldResemble, table2)
			So(err2, ShouldBeNil)

			// error
			So(t3, ShouldBeNil)
			So(err3, ShouldEqual, ErrNotFound)

		})

		Convey("It can return tables", func() {
			table1, _ := book.NewTable("table1", []string{"foo", "bar", "baz"})
			table2, _ := book.NewTable("table2", []string{"foo", "bar", "baz"})

			tables := []*Table{table1, table2}

			for i, t := range book.Tables() {
				So(t, ShouldResemble, tables[i])
			}

		})

		Convey("It can add tables", func() {
			table1, _ := NewTable("table1", []string{"foo", "bar", "baz"})

			ok1 := book.AddTable(table1)

			So(ok1, ShouldBeNil)

			t1, _ := book.Table("table1")

			So(t1, ShouldResemble, table1)

			ok2 := book.AddTable(table1)

			So(ok2, ShouldEqual, ErrTableExists)
		})
	})
}
