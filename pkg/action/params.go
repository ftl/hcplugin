package action

import (
	"fmt"
	"strings"
)

// parseParams reads a comma separated list of name=value pairs, for example
// "number=5" or "amount=100,unit=hz".
func parseParams(s string) (map[string]string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	result := make(map[string]string)
	for _, pair := range strings.Split(s, ",") {
		name, value, found := strings.Cut(pair, "=")
		name = strings.TrimSpace(name)
		if !found || name == "" {
			return nil, fmt.Errorf("invalid parameter %q, use name=value", strings.TrimSpace(pair))
		}
		result[name] = strings.TrimSpace(value)
	}
	return result, nil
}
