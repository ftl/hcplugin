package plugin

import (
	"fmt"
	"strconv"
)

type GlobalSettings struct {
	Port int `json:"port"`
}

const DefaultPort = 8383

func parseGlobalSettings(settings map[string]any) (GlobalSettings, error) {
	result := GlobalSettings{Port: DefaultPort}
	if len(settings) == 0 {
		return result, nil
	}
	switch v := settings["port"].(type) {
	case nil:
		// keep default
	case float64:
		if v > 0 {
			result.Port = int(v)
		}
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return result, fmt.Errorf("global settings: port %q not an integer: %w", v, err)
		}
		if n > 0 {
			result.Port = n
		}
	default:
		return result, fmt.Errorf("global settings: unexpected port type %T", v)
	}
	return result, nil
}
