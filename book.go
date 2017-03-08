package tablebook

//Book groups tables
type Book struct {
	tables []*Table
}

// NewBook creates a new Book
func NewBook() *Book {
	return &Book{}
}

// NewTable creates a Table, appends it to tables and returns it
// returns ErrTableExists if table already exist
func (b *Book) NewTable(name string, headers []string) (*Table, error) {
	if b.tableIndex(name) != -1 {
		return nil, ErrTableExists
	}

	table := NewTable(name, headers)
	b.tables = append(b.tables, table)

	return table, nil
}

// Table searches and returns table by given name
// returns tablebook.ErrNotFound if the table cannot be found
func (b *Book) Table(name string) (*Table, error) {
	index := b.tableIndex(name)

	if index == -1 {
		return nil, ErrNotFound
	}

	return b.tables[index], nil
}

// Tables returns all tables
func (b *Book) Tables() []*Table {
	return b.tables
}

// AddTable adds a existing table to tables
// returns tablebook.ErrTableExists if the table already exists
func (b *Book) AddTable(table *Table) error {
	if b.tableIndex(table.name) != -1 {
		return ErrTableExists
	}

	b.tables = append(b.tables, table)

	return nil
}

// index returns index of give table
// returns -1 if table is not found.
func (b *Book) tableIndex(name string) int {
	for i, table := range b.tables {
		if table.name == name {
			return i
		}
	}
	return -1
}
