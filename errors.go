package tablebook

import "errors"

// ErrInvalidDimensions is returned when trying to append/insert too much
// or not enough values to a row or column
var ErrInvalidDimensions = errors.New("tablebook: Invalid dimension")

// ErrNotFound is returned when trying to acces a uknown ro or column.
var ErrNotFound = errors.New("tablebook: Not found")

// ErrTableExists is returned when trying to acces a uknown ro or column.
var ErrTableExists = errors.New("tablebook: given table already exists")

// ErrColumnExists is returned when trying to acces a uknown ro or column.
var ErrColumnExists = errors.New("tablebook: given column already exists")
