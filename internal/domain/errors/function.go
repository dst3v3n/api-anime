package errors

import "fmt"

func NewNotFoundErr(message string, id string) error {
	return fmt.Errorf("%s (Resource: %s): %w", message, id, ErrNotFound)
}

func NewNotSetConfig(message string, id string) error {
	return fmt.Errorf("%s (Resource: %s): %w", message, id, ErrNotSetConfig)
}
