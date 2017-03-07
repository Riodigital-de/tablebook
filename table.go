package tablebook

// Table represents a single table with headers and rows
type Table struct {
	name    string
	headers []string
	rows    [][]interface{}
}

// NewTable creates a new Table.
func NewTable(name string, headers []string) *Table {
	return &Table{
		name:    name,
		headers: headers,
	}
}

// Width returns the number of columns in the Table.
func (t *Table) Width() int {
	return len(t.headers)
}

// Width returns the number of rows in the Table.
func (t *Table) Height() int {
	return len(t.rows)
}

// AppendRow appends a row of values to the Dataset.
// returns tablebook.ErrInvalidDimensions if the row is to long
func (t *Table) AppendRow(row []interface{}) error {
	if len(row) != t.Width() {
		return ErrInvalidDimensions
	}

	t.rows = append(t.rows, row)

	return nil
}

// Column returns all the values for a specific column
// returns tablebook.ErrNotFound if the column cannot be found
func (t *Table) Column(column string) ([]interface{}, error) {
	columnIndex := t.columnIndex(column)

	if columnIndex == -1 {
		return nil, ErrNotFound
	}

	columns := make([]interface{}, t.Height())

	for i, e := range t.rows {
		columns[i] = e[columnIndex]
	}

	return columns, nil
}

// RenameHeader renames header from given name to given name
// returns tablebook.ErrNotFound if the header cannot be found
func (t *Table) RenameHeader(from, to string) error {

	var found int

	for i, h := range t.headers {
		if h == from {
			t.headers[i] = to
			found++
		}
	}

	if found == 0 {
		return ErrNotFound
	}

	return nil
}

// Take joins given tables into current table on a row level.
func (t *Table) Take(tables []*Table) {

	for _, table := range tables {
		var row []interface{}
		for _, header := range t.headers {
			values, _ := table.Column(header)

			for _, col := range values {
				row = append(row, col)
			}
		}
		t.rows = append(t.rows, row)
	}
}

// columnIndex returns index of given column
// returns -1 if column is not found.
func (t *Table) columnIndex(column string) int {
	for i, e := range t.headers {
		if e == column {
			return i
		}
	}
	return -1
}
