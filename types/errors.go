package types

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid")
	ErrQuota    = errors.New("quota")
	ErrCanceled = errors.New("canceled")
	ErrConflict = errors.New("conflict")
)
