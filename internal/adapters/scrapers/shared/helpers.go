package shared

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ParseUint(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("formato de uint inválido %q: %w", value, err)
	}
	return uint(parsed), nil
}

func BuildURL(baseURL string, params map[string]string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}
	query := u.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func ParseFloat(value string) (float64, error) {
	if value == "" {
		return 0.0, fmt.Errorf("cadena vacía, no se puede convertir a float")
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, fmt.Errorf("formato de float inválido %q: %w", value, err)
	}
	return parsed, nil
}

func ExtractID(href string) (string, error) {
	parts := strings.Split(href, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("formato de href inválido: %s", href)
	}
	return parts[2], nil
}
