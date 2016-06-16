package packet

import (
	"strings"
)

func Clamp(value int, min int, max int) int {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}

func StringClamp(value string, max int) string {
	if len(value) > max {
		value = value[:max]
	}
	return value
}

func NullClean(value string) string {
	value = strings.Replace(value, "\x00", "", -1)
	value = strings.TrimSpace(value)
	return value
}
