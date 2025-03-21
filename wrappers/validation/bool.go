package validation

import "fmt"

func IsBool(v any) (bool, error) {
	b, ok := v.(bool)

	if !ok {
		return false, fmt.Errorf("invalid bool provided %v", v)
	}

	return b, nil
}

type VBool struct{}

func (b *VBool) IsValid(v any, fieldName string) (bool, error) {
	vb, err := IsBool(v)

	return vb, &ErrMap{fieldName: err}
}

func Bool() *VBool {
	return &VBool{}
}
