package errors

import (
	"context"
	stderrors "errors"
)

var (
	ErrEmptyAnimeID   = stderrors.New("anime id is empty")
	ErrInvalidEpisode = stderrors.New("episode must be greater than zero")
)

func Classify(err error) Kind {
	var e *Error
	switch {
	case err == nil:
		return KindUnknown
	case stderrors.As(err, &e):
		return e.Kind
	case stderrors.Is(err, context.DeadlineExceeded):
		return KindTimeout
	case stderrors.Is(err, context.Canceled):
		return KindCancel
	default:
		return KindExternal
	}
}

func Wrap(service, op string, err error) error {
	if err == nil {
		return nil
	}
	var e *Error
	if stderrors.As(err, &e) {
		return err
	}
	return E(service, op, Classify(err), err)
}
