package utils

import (
	"strconv"
	"strings"
)

// ToInt converts a value to an integer.
func ToInt[T interface{} | string | int | int64 | float32 | float64](value T) int {
	var x interface{} = value
	switch i := x.(type) {
	case int:
		return i
	case int64:
		return int(i)
	case float32:
		return int(i)
	case float64:
		return int(i)
	case string:
		intValue, err := strconv.Atoi(i)
		if err != nil {
			return 0
		}
		return intValue
	}

	return 0
}

// ToInt64 converts a value to an int64.
func ToInt64[T string | int | int64 | float32 | float64](value T) int64 {
	var x interface{} = value
	switch i := x.(type) {
	case int:
		return int64(i)
	case int64:
		return i
	case float32:
		return int64(i)
	case float64:
		return int64(i)
	case string:
		intValue, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			return 0
		}
		return intValue
	}
	return 0
}

// ToFloat converts a value to a float64.
func ToFloat[T string | int | int64 | float32 | float64](value T) float64 {
	var x interface{} = value
	switch i := x.(type) {
	case int:
		return float64(i)
	case int64:
		return float64(i)
	case float32:
		return float64(i)
	case float64:
		return i
	case string:
		f, _ := strconv.ParseFloat(strings.TrimSpace(i), 64)
		return f
	}
	return 0.0
}

// ToString converts a value to a string.
func ToString[T string | int | int64 | float32 | float64](value T) string {
	var x interface{} = value
	switch i := x.(type) {
	case int:
		return strconv.Itoa(i)
	case int64:
		return strconv.FormatInt(i, 10)
	case float32:
		return strconv.FormatFloat(float64(i), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(i, 'f', -1, 64)
	case string:
		return i
	}
	return ""
}
