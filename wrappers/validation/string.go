package validation

import (
	"fmt"
	"regexp"
)

const EmailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

type StringOptions struct {
	MinLen int
	MaxLen int
	Email  bool
}

type StringValidators = map[string]Validator

func isString(v any) (string, error) {
	s, ok := v.(string)

	if !ok {
		return "", fmt.Errorf("invalid string provided %v", v)
	}

	return s, nil
}

func email(v string) error {
	ok, err := regexp.Match(EmailRegex, []byte(v))

	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("invalid email provided %s", v)
	}

	return nil
}

func minLen(val string, min int) error {
	if len(val) <= min {
		return fmt.Errorf("invalid string, minimum length %d exceeded", min)
	}

	return nil
}

func maxLen(val string, max int) error {
	if len(val) <= max {
		return fmt.Errorf("invalid string, maximum length %d exceeded", max)
	}

	return nil
}

type VString struct {
	validators []ValidatorFn[string]
}

func (s *VString) Min(len int) *VString {
	s.validators = append(s.validators, func(v string) error {
		return minLen(v, len)
	})
	return s
}

func (s *VString) Max(len int) *VString {
	s.validators = append(s.validators, func(v string) error {
		return maxLen(v, len)
	})
	return s
}

func (s *VString) Email() *VString {
	s.validators = append(s.validators, func(v string) error {
		return email(v)
	})

	return s
}

func (s *VString) Validate(v any, fieldName string) (string, error) {
	strVal, err := isString(v)

	if err != nil {
		return "", ErrMap{
			fieldName: err,
		}
	}

	var issues *ErrMap = nil

	for _, validator := range s.validators {
		err := validator(strVal)

		if err != nil {
			if issues == nil {
				issues = &ErrMap{}
			}

			(*issues)[fieldName] = err
		}
	}

	if issues != nil {
		return strVal, *issues
	}

	return strVal, nil
}

func (s *VString) IsValid(v any, fieldName string) error {
	_, err := s.Validate(v, fieldName)

	return err
}

// Create Method
func String() *VString {
	return &VString{
		validators: make([]ValidatorFn[string], 0),
	}
}
