package tablebook

// Table represents a single table with columnNames and rows
type Table struct {
	name        string
	columnNames []string
	rows        [][]interface{}
}

// EvaluatedColumn represents a function that can be evaluated dynamically
// when exporting to a predefined format.
type EvaluatedColumn func(table *Table, rowIndex int, cellIndex int) interface{}

// NewTable creates a new Table.
func NewTable(name string, columnNames []string) (*Table, error) {

	uniqueColumnNames := make([]string, len(columnNames))

	for i, cn := range columnNames {
		for _, ucn := range uniqueColumnNames {
			if ucn == cn {
				return nil, ErrColumnExists
			}
		}
		uniqueColumnNames[i] = cn
	}

	return &Table{
		name:        name,
		columnNames: columnNames,
	}, nil
}

// Width returns the number of columns in the Table.
func (t *Table) Width() int {
	return len(t.columnNames)
}

// Name returns the name
func (t *Table) Name() string {
	return t.name
}

// Height returns the number of rows in the Table.
func (t *Table) Height() int {
	return len(t.rows)
}

// Row returns row at given index
// returns tablebook.ErrNotFound if the given row cannot be found
func (t *Table) Row(index int) ([]interface{}, error) {
	if index < 0 || index > t.Height() {
		return nil, ErrInvalidDimensions
	}

	row := make([]interface{}, t.Width())

	for ci := range t.rows[index] {
		cv, _ := t.Cell(index, ci)
		row[ci] = cv

	}

	return row, nil
}

// Rows return rows
func (t *Table) Rows() [][]interface{} {
	rows := make([][]interface{}, t.Height())

	for ri := range t.rows {
		row, _ := t.Row(ri)

		rows[ri] = row
	}

	return rows
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

// ColumnNames return columnNames
func (t *Table) ColumnNames() []string {
	return t.columnNames
}

// Column returns a column by given columnName.
// returns tablebook.ErrNotFound if the given columnName cannot be found
func (t *Table) Column(columnName string) ([]interface{}, error) {
	index := t.ColumnIndex(columnName)

	if index == -1 {
		return nil, ErrNotFound
	}

	column := make([]interface{}, t.Height())

	for i, r := range t.Rows() {
		column[i] = r[index]
	}

	return column, nil
}

// RenameColumn renames columnName from given name to given name
// returns tablebook.ErrNotFound if the columnName cannot be found
// returns tablebook.ErrColumnExists if the columnName cannot be found
func (t *Table) RenameColumn(from, to string) error {
	if t.ColumnIndex(to) != -1 {
		return ErrColumnExists
	}

	index := t.ColumnIndex(from)

	if index == -1 {
		return ErrNotFound
	}

	t.columnNames[index] = to

	return nil
}

// AppendColumn appends a column of values and columnName to the Dataset.
// returns tablebook.ErrInvalidDimensions if the column is to large
// returns tablebook.ErrColumnExists if the column already exist
func (t *Table) AppendColumn(columnName string, column []interface{}) error {

	if len(column) > t.Height() {
		return ErrInvalidDimensions
	}

	if t.ColumnIndex(columnName) != -1 {
		return ErrColumnExists
	}

	t.columnNames = append(t.columnNames, columnName)

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

// AppendEvaluatedColumn appends a column of evaluated functions and columnName to the Dataset.
// returns tablebook.ErrColumnExists if the column already exist
func (t *Table) AppendEvaluatedColumn(columnName string, fn EvaluatedColumn) error {
	if t.ColumnIndex(columnName) != -1 {
		return ErrColumnExists
	}

	t.columnNames = append(t.columnNames, columnName)
	for i, r := range t.rows {
		t.rows[i] = append(r, fn)
	}

	return nil
}

// ColumnIndex returns index of given column
// returns -1 if column is not found.
func (t *Table) ColumnIndex(columnName string) int {
	for i, cn := range t.columnNames {
		if cn == columnName {
			return i
		}
	}
	return -1
}

// Cell return value for cell in given row and column
// returns tablebook.ErrInvalidDimensions if the given row or cell cannot be found
func (t *Table) Cell(rowIndex, columnIndex int) (interface{}, error) {

	if rowIndex < 0 || rowIndex > t.Height() || columnIndex < 0 || columnIndex > t.Width() {
		return nil, ErrInvalidDimensions
	}

	cell := t.rows[rowIndex][columnIndex]

	switch cell.(type) {
	case EvaluatedColumn:
		return cell.(EvaluatedColumn)(t, rowIndex, columnIndex), nil
	default:
		return cell, nil
	}
}

// Take joins given tables into current table on a row level.
func (t *Table) Take(tables []*Table) {

	for _, foreignTable := range tables {

		for foreignTableRowIndex := range foreignTable.rows {
			var row []interface{}

			for _, columnName := range t.columnNames {
				var value interface{}
				value = ""

				requiredIndex := foreignTable.ColumnIndex(columnName)

				if requiredIndex != -1 {
					cv, _ := foreignTable.Cell(foreignTableRowIndex, requiredIndex)
					value = cv
				}

				row = append(row, value)
			}

			t.AppendRow(row)

		}
	}
}
