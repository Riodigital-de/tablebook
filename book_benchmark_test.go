package tablebook

import (
	"fmt"
	"testing"
)

func BenchmarkBookWithTables(b *testing.B) {
	b.SetBytes(2)

	book := NewBook()
	for i := 0; i < b.N; i++ {

		table, _ := book.NewTable(fmt.Sprintf("table-%d", i), []string{"foo", "bar", "baz"})

		table.AppendRow([]interface{}{
			fmt.Sprintf("foo %d table", i),
			fmt.Sprintf("bar %d table", i),
			fmt.Sprintf("baz %d table", i),
		})

	}
}
