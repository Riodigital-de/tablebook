package tablebook

import (
	"testing"
	"fmt"
)

func BenchmarkBookWithTables(b *testing.B) {
	b.SetBytes(2)

	book := NewBook()
	for i := 0; i < b.N; i++ {

		table, _ := book.NewTable(fmt.Sprintf("table-%s", i), []string{"foo", "bar", "baz"})

		table.AppendRow([]interface{}{
			fmt.Sprintf("baz %s table %s", i),
			fmt.Sprintf("bar %s table %s", i),
			fmt.Sprintf("baz %s table %s", i),
		})

	}
}
