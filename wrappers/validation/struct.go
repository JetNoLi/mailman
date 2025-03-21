package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type VSchema[T any] struct {
	fields map[string]Validator
}

func GetFieldValue(v any, fieldName string) (any, error) {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %T", v)
	}

	field := rv.FieldByName(fieldName)

	return field.Interface(), nil
}

func (s *VSchema[T]) Schema(fields map[string]Validator) {
	s.fields = fields
}

func (s *VSchema[T]) IsValid(v any) error {
	_, err := s.Validate(v)

	return err
}

func (s *VSchema[T]) Validate(v any) (T, error) {
	var errMap *ErrMap
	var val T

	for fieldName, validator := range s.fields {
		val, err := GetFieldValue(v, fieldName)

		if err != nil {
			if errMap == nil {
				errMap = &ErrMap{}
			}

			(*errMap)[fieldName] = err
			continue
		}

		err = validator.IsValid(val, "")

		if err != nil {
			if errMap == nil {
				errMap = &ErrMap{}
			}
			(*errMap)[fieldName] = err
		}
	}

	if errMap != nil {
		return val, *errMap
	}

	raw, err := json.Marshal(v)

	if err != nil {
		return val, fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := json.Unmarshal(raw, &val); err != nil {
		return val, fmt.Errorf("failed to unmarshal into struct: %w", err)
	}

	return val, nil
}

type SchemaMap map[string]Validator

func Schema[T any](v SchemaMap) *VSchema[T] {
	return &VSchema[T]{
		fields: v,
	}
}

func CreateSchema(v SchemaMap) {

}
