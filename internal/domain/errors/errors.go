package errors

import "errors"

// Genericos
var (
	ErrNotFound = errors.New("not found")
)

// Viper
var (
	ErrNotSetConfig = errors.New("not set config")
)
