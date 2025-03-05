package errs

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrUnknownDBError = errors.New("database error")
)
