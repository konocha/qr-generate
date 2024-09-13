package validation

import (

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func ValidateEmail(email string) error {
	return validation.Validate(email, validation.Required, validation.Length(7, 100), is.Email)
}

func ValidatePassword(password string) error {
	return validation.Validate(password, validation.Required, validation.Length(5, 100))
}
