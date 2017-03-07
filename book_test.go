package tablebook

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBook(t *testing.T) {
	Convey("Given a empty book", t, func() {
		book := NewBook()

		Convey("It can create tables", func() {
			_, ok1 := book.NewTable("table", []string{"foo", "bar", "baz"})
			So(ok1, ShouldBeNil)

			_, ok2 := book.NewTable("table", []string{"foo", "bar", "baz"})
			So(ok2, ShouldEqual, ErrTableExists)
		})

		Convey("It can search tables", func() {
			table1, _ := book.NewTable("table1", []string{"foo", "bar", "baz"})
			table2, _ := book.NewTable("table2", []string{"foo", "bar", "baz"})

			So(book.Table("table1"), ShouldResemble, table1)
			So(book.Table("table2"), ShouldResemble, table2)
			So(book.Table("tableunknown"), ShouldBeNil)
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
			table1 := NewTable("table1", []string{"foo", "bar", "baz"})

			ok1 := book.AddTable(table1)

			So(ok1, ShouldBeNil)
			So(book.Table("table1"), ShouldResemble, table1)

			ok2 := book.AddTable(table1)

			So(ok2, ShouldEqual, ErrTableExists)
		})
	})
}
