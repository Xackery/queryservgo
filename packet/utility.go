package packet

import ()

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
