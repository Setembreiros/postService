package database

import "fmt"

type NotFoundError struct {
	table string
	key   any
}

func (e *NotFoundError) Error() string {
	errorMessage := fmt.Sprintf("Data in table %s not found for key %v", e.table, e.key)
	return errorMessage
}

func NewNotFoundError(table string, key any) *NotFoundError {
	return &NotFoundError{
		table: table,
		key:   key,
	}
}
