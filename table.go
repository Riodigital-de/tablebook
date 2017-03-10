package tablebook

// Table represents a single table with headers and rows
type Table struct {
	name    string
	headers []string
	rows    [][]interface{}
}

// DynamicColumn represents a function that can be evaluated dynamically
// when exporting to a predefined format.
type DynamicColumn func(*Table, int, int, []interface{}) interface{}

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

// Row return row at given index
// returns tablebook.ErrNotFound if the given row cannot be found
func (t *Table) Row(index int) ([]interface{}, error) {
	if index < 0 || index > t.Height() {
		return nil, ErrNotFound
	}

	row := make([]interface{}, t.Width())

	for ci, c := range t.rows[index] {
		switch c.(type) {
		case DynamicColumn:
			row[ci] = c.(DynamicColumn)(t, index, ci, t.rows[index])
		default:
			row[ci] = c
		}

	}

	return row, nil
}

// Rows return rows
func (t *Table) Rows() [][]interface{} {
	rows := make([][]interface{}, t.Height())

	for ri := range t.rows {
		row, err := t.Row(ri)

		if err != nil {
			continue
		}

		rows[ri] = row
	}

	return rows
}

// Headers return headers
func (t *Table) Headers() []string {
	return t.headers
}

// AppendRow appends a row of values to the Dataset.
// returns tablebook.ErrInvalidDimensions if the row is to long
func (t *Table) AppendRow(cells []interface{}) error {
	if len(cells) > t.Width() {
		return ErrInvalidDimensions
	}

	row := make([]interface{}, t.Width())

	for i := range row {
		var value interface{}
		value = ""

		if i < len(cells) {
			value = cells[i]
		}

		row[i] = value
	}

	t.rows = append(t.rows, row)

	return nil
}

// AppendColumn appends a column of values and header to the Dataset.
// returns tablebook.ErrInvalidDimensions if the column is to large
// returns tablebook.ErrColumnExists if the column already exist
func (t *Table) AppendColumn(header string, column []interface{}) error {

	if len(column) > t.Height() {
		return ErrInvalidDimensions
	}

	if t.columnIndex(header) != -1 {
		return ErrColumnExists
	}

	t.headers = append(t.headers, header)

	for i, r := range t.rows {
		var value interface{}
		value = ""

		if i < len(column) {
			value = column[i]
		}

		t.rows[i] = append(r, value)
	}

	return nil
}

// AppendDynamicColumn appends a column of evaluated functions and header to the Dataset.
// returns tablebook.ErrColumnExists if the column already exist
func (t *Table) AppendDynamicColumn(header string, fn DynamicColumn) error {
	if t.columnIndex(header) != -1 {
		return ErrColumnExists
	}

	t.headers = append(t.headers, header)
	for i, r := range t.rows {
		t.rows[i] = append(r, fn)
	}

	return nil
}

// Column returns a column by given header.
// returns tablebook.ErrNotFound if the given header cannot be found
func (t *Table) Column(header string) ([]interface{}, error) {
	index := t.columnIndex(header)

	if index == -1 {
		return nil, ErrNotFound
	}

	column := make([]interface{}, t.Height())

	for i, r := range t.Rows() {
		column[i] = r[index]
	}

	return column, nil
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
