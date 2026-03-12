package config

import (
	"fmt"
	"strings"
)

// StrictValidator validates configuration with strict rules.
type StrictValidator interface {
	ValidateStrict() error
}

// ValidationError collects multiple validation errors.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Errors, "; ")
}

// Add adds an error.
func (e *ValidationError) Add(msg string) {
	e.Errors = append(e.Errors, msg)
}

// Valid returns true if no errors.
func (e *ValidationError) Valid() bool {
	return len(e.Errors) == 0
}

// ValidateRequired checks that s is non-empty after trim.
func ValidateRequired(s, name string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

// ValidateMinLength checks minimum length.
func ValidateMinLength(s, name string, min int) error {
	if len(s) < min {
		return fmt.Errorf("%s must be at least %d characters", name, min)
	}
	return nil
}

// ValidateOneOf checks s is in allowed values.
func ValidateOneOf(s, name string, allowed []string) error {
	for _, a := range allowed {
		if s == a {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of: %s", name, strings.Join(allowed, ", "))
}
