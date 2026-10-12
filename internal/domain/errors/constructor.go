package errors

import "fmt"

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Kind, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Kind)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func E(service string, op string, kind Kind, err error) *Error {
	return &Error{Service: service, Op: op, Kind: kind, Err: err}
}
