package cache

import (
	"encoding/json"
	"fmt"

	apperrors "github.com/dst3v3n/api-anime/internal/domain/errors"
)

func invalid(err error) error {
	return apperrors.E(serviceName, "encodeJSON", apperrors.KindInvalidInput, err)
}

func encodeJSON(value any) ([]byte, error) {
	if value == nil {
		return nil, invalid(fmt.Errorf("%w: valor nil", ErrInvalidCacheValue))
	}

	data, err := json.Marshal(value)
	if err != nil {
		return nil, invalid(fmt.Errorf("%w: %w", ErrInvalidCacheValue, err))
	}

	if string(data) == "null" {
		return nil, invalid(fmt.Errorf("%w: valor nil", ErrInvalidCacheValue))
	}

	if !json.Valid(data) {
		return nil, invalid(fmt.Errorf("%w: el resultado no es JSON válido", ErrInvalidCacheValue))
	}

	if data[0] != '{' && data[0] != '[' {
		return nil, invalid(fmt.Errorf("%w: la raíz debe ser un objeto o un arreglo, no un valor escalar", ErrInvalidCacheValue))
	}

	return data, nil
}
