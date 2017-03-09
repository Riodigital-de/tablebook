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

// Name returns the name
func (t *Table) Name() string {
	return t.name
}

// Height returns the number of rows in the Table.
func (t *Table) Height() int {
	return len(t.rows)
}

// Rows return rows
func (t *Table) Rows() [][]interface{} {
	return t.rows
}

// Headers return headers
func (t *Table) Headers() []string {
	return t.headers
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

	for _, foreignTable := range tables {

		for _, foreignTableRow := range foreignTable.rows {
			var row []interface{}

			for _, header := range t.headers {
				var value interface{}

				value = ""

				requiredIndex := foreignTable.columnIndex(header)

				if requiredIndex != -1 {
					value = foreignTableRow[requiredIndex]
				}

				row = append(row, value)
			}

			t.AppendRow(row)

		}
	}
}

// columnIndex returns index of given column
// returns -1 if column is not found.
func (t *Table) columnIndex(column string) int {
	for i, header := range t.headers {
		if header == column {
			return i
		}
	}
	return -1
}
