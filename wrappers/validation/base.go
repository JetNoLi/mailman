package validation

import "fmt"

type Validator interface {
	IsValid(v any, fieldName string) error
}

type ValidatorFn[T any] func(v T) error

type ErrMap map[string]error

func (errs ErrMap) Error() string {
	message := ""

	if errs == nil {
		return ""
	}

	for fieldName, err := range errs {
		if err == nil {
			fmt.Println("DEBUG: invalid error of nil provided for", fieldName)
			continue
		}

		baseMessage := ""

		if fieldName != "" {
			baseMessage = fmt.Sprintf("issues with field %s: ", fieldName)
		}

		if message == "" {
			message = fmt.Sprintf("%s%s", baseMessage, err.Error())
			continue
		}

		message += fmt.Sprintf("\n%s%s", baseMessage, err.Error())
	}

	return message
}
