package utils

import validator "github.com/asaskevich/govalidator"

// IsInteger to validate if all value is integer from string type values
func IsInteger(values ...string) bool {
	for _, value := range values {
		if !validator.IsInt(value) {
			return false
		}
	}
	return true
}

// IsNonNegative is to validate if all value is non-negative integer from integer type values
func IsNonNegative(values ...int) bool {
	for _, value := range values {
		if value < 0 {
			return false
		}
	}
	return true
}

// IsNonEmpty to check whether the values is set to empty/zero in their corresponding type
func IsNonEmpty(values ...interface{}) bool {
	for _, value := range values {
		switch value.(type) {
		case string:
			if value == "" {
				return false
			}
		case int:
			if value == 0 {
				return false
			}
		case int64:
			if value == 0 {
				return false
			}
		case uint:
			if value == 0 {
				return false
			}
		case uint32:
			if value == 0 {
				return false
			}
		case uint64:
			if value == 0 {
				return false
			}
		case float32:
			if value == 0 {
				return false
			}
		case float64:
			if value == 0 {
				return false
			}
		default:
			if value == nil {
				return false
			}
		}
	}

	return true
}
