package config

import (
	"testing"
)

func TestValidationError(t *testing.T) {
	var e ValidationError
	e.Add("error 1")
	e.Add("error 2")

	if e.Valid() {
		t.Error("should not be valid")
	}
	if e.Error() != "error 1; error 2" {
		t.Errorf("unexpected error string: %s", e.Error())
	}
}

func TestValidateRequired(t *testing.T) {
	if err := ValidateRequired("x", "field"); err != nil {
		t.Errorf("non-empty should pass: %v", err)
	}
	if err := ValidateRequired("  ", "field"); err == nil {
		t.Error("whitespace only should fail")
	}
	if err := ValidateRequired("", "field"); err == nil {
		t.Error("empty should fail")
	}
}

func TestValidateOneOf(t *testing.T) {
	if err := ValidateOneOf("a", "x", []string{"a", "b"}); err != nil {
		t.Errorf("valid value should pass: %v", err)
	}
	if err := ValidateOneOf("c", "x", []string{"a", "b"}); err == nil {
		t.Error("invalid value should fail")
	}
}
