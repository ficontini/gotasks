package db

import "errors"

var (
	ErrorNotFound           = errors.New("resource not found")
	ErrInvalidID            = errors.New("invalid ID")
	ErrInvalidOperationType = errors.New("invalid operation")
	ErrInvalidBatchSize     = errors.New("invalid batch size")
)
