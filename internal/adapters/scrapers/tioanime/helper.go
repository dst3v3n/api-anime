package tioanime

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (v *video) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("video: se esperaba un array: %w", err)
	}

	if len(raw) < 4 {
		return errors.New("video: el array tiene menos de 4 elementos")
	}

	if err := json.Unmarshal(raw[0], &v.Server); err != nil {
		return fmt.Errorf("video: servidor inválido: %w", err)
	}
	if err := json.Unmarshal(raw[1], &v.URL); err != nil {
		return fmt.Errorf("video: url inválida: %w", err)
	}
	if err := json.Unmarshal(raw[2], &v.Number); err != nil {
		return fmt.Errorf("video: campo 3 inválido: %w", err)
	}
	if err := json.Unmarshal(raw[3], &v.Number1); err != nil {
		return fmt.Errorf("video: campo 4 inválido: %w", err)
	}

	return nil
}

func removeLast(s, substr string) string {
	i := strings.LastIndex(s, substr)
	if i == -1 {
		return s
	}
	return s[:i] + s[i+len(substr):]
}
