# Validation Module

Request ve body validation (go-playground/validator).

## Usage

```go
import "github.com/parevo/core/validation"

type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2,max=100"`
}

// Struct validation
err := validation.Validate(req)

// JSON body
err := validation.ValidateJSON(bodyBytes, &req)
```

## Validator Tags

- `required` — boş olamaz
- `email` — geçerli email
- `min=N`, `max=N` — uzunluk
- `oneof=a b c` — değerlerden biri
- `uuid` — UUID formatı
- `url` — geçerli URL

Tam liste: https://pkg.go.dev/github.com/go-playground/validator/v10
