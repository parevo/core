// Package validation provides request and body validation.
//
// # Usage
//
//	// Struct tags (go-playground/validator)
//	type CreateUserRequest struct {
//	    Email string `json:"email" validate:"required,email"`
//	    Name  string `json:"name" validate:"required,min=2,max=100"`
//	}
//	err := validation.Validate(req)
//
//	// JSON body
//	err := validation.ValidateJSON(body, &req)
package validation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var defaultValidator = validator.New()

// Validate validates a struct using validator tags.
func Validate(v any) error {
	if v == nil {
		return fmt.Errorf("validation: nil value")
	}
	if err := defaultValidator.Struct(v); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			var msgs []string
			for _, e := range errs {
				msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
			}
			return fmt.Errorf("validation: %s", strings.Join(msgs, "; "))
		}
		return err
	}
	return nil
}

// ValidateJSON unmarshals JSON into v and validates.
func ValidateJSON(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("validation: invalid json: %w", err)
	}
	return Validate(v)
}
