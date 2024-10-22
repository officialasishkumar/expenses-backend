// utils/utils.go
package utils

import "math"

// AlmostEqual checks if two floats are approximately equal
func AlmostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// ConvertToFloat64 converts various numeric types to float64
func ConvertToFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}
