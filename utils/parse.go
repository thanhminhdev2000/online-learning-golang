package utils

import (
	"strconv"
)

func ParseIntWithDefault(str string, defaultValue int) int {
	val, err := strconv.Atoi(str)
	if err != nil || val < 1 {
		return defaultValue
	}
	return val
}

func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
} 