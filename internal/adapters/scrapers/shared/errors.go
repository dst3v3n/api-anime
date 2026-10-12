package shared

import (
	"context"
	"errors"
	"net/http"

	domainErrors "github.com/dst3v3n/api-anime/internal/domain/errors"
)

func kindFromError(err error) domainErrors.Kind {
	switch {
	case errors.Is(err, context.Canceled):
		return domainErrors.KindCancel
	case errors.Is(err, context.DeadlineExceeded):
		return domainErrors.KindTimeout
	default:
		return domainErrors.KindExternal
	}
}

func kindFromStatus(code int) domainErrors.Kind {
	switch code {
	case http.StatusNotFound:
		return domainErrors.KindNotFound
	case http.StatusTooManyRequests:
		return domainErrors.KindRateLimited
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return domainErrors.KindTimeout
	default:
		return domainErrors.KindExternal
	}
}
