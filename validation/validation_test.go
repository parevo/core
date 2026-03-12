package validation

import (
	"testing"
)

func TestValidate(t *testing.T) {
	type Req struct {
		Email string `json:"email" validate:"required,email"`
		Name  string `json:"name" validate:"required,min=2"`
	}

	err := Validate(&Req{Email: "a@b.com", Name: "Ab"})
	if err != nil {
		t.Errorf("valid: %v", err)
	}

	err = Validate(&Req{Email: "invalid", Name: "Ab"})
	if err == nil {
		t.Error("invalid email should fail")
	}

	err = Validate(&Req{Email: "a@b.com", Name: "A"})
	if err == nil {
		t.Error("name min=2 should fail")
	}
}

func TestValidateJSON(t *testing.T) {
	type Req struct {
		Email string `json:"email" validate:"required,email"`
	}

	body := []byte(`{"email":"a@b.com"}`)
	var req Req
	err := ValidateJSON(body, &req)
	if err != nil {
		t.Fatal(err)
	}
	if req.Email != "a@b.com" {
		t.Errorf("got %q", req.Email)
	}

	body = []byte(`{"email":"x"}`)
	err = ValidateJSON(body, &req)
	if err == nil {
		t.Error("invalid email should fail")
	}
}
