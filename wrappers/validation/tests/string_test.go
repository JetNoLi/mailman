package tests

import (
	"ewails/wrappers/validation"
	"strings"
	"testing"
)

const (
	ValidEmail            = "test@gmail.com"
	InvalidEmailNoAt      = "testgmail.com"
	InvalidEmailNoDomain1 = "test@gmail."
	InvalidEmailNoDomain2 = "test@gmail"
	InvalidEmailNoStart   = "@gmail.com"

	FieldErrMsg = "issues with field %s :"
	EmailErrMsg = "invalid email provided %s"
)

func TestVStringTypeFailure(t *testing.T) {
	s := validation.String()

	str, err := s.Validate(100, "")

	if err == nil {
		t.Fatal("expected validation of int to fail, but passed for String", str, err)
	}
}

func TestVStringTypeEmail(t *testing.T) {
	s := validation.String().Email()

	_, err := s.Validate(InvalidEmailNoAt, "")

	if err == nil {
		t.Fail()
		t.Log("expected validation error, no @ symbol in email", InvalidEmailNoAt)
	} else if !strings.Contains(err.Error(), "invalid email provided") {
		t.Fail()
		t.Log("expected error message to contain 'invalid email provided', got", err.Error())
	}

	_, err = s.Validate(InvalidEmailNoDomain1, "")

	if err == nil {
		t.Log("expected validation error, no domain in email", InvalidEmailNoDomain1)
		t.Fail()
	} else if !strings.Contains(err.Error(), "invalid email provided") {
		t.Log("")
		t.Fail()
	}

}
