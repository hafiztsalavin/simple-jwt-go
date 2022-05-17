package utils

import validator "github.com/asaskevich/govalidator"

func IsInteger(values ...string) bool {
	for _, value := range values {
		if !validator.IsInt(value) {
			return false
		}
	}
	return true
}

func IsNonNegative(values ...int) bool {
	for _, value := range values {
		if value < 0 {
			return false
		}
	}
	return true
}

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
