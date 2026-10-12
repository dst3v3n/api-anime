package shared

import "strings"

func NormalizeID(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}
