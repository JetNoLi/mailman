package common

import "fmt"

func Validate[T any](v any) (T, error) {
	value, ok := v.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("invalid value, expected type %T", zero)
	}
	return value, nil
}
