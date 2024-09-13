package validation_test

import (
	"testing"

	"github.com/konocha/qr-generate/internal/app/validation"
	"github.com/stretchr/testify/assert"
)

func Test_validateEmail(t *testing.T){
	badEmail := "hellgrtgrgro@"
	badEmailLength := "h@a.ru"

	goodEmail := "hello@gmail.com"

	
	assert.NoError(t, validation.ValidateEmail(badEmail))
	assert.Error(t, validation.ValidateEmail(badEmailLength))

	assert.NoError(t, validation.ValidateEmail(goodEmail))
}

func Test_validatePassword(t *testing.T){
	badPassword := "hr"
	goodPassword := "Hello123World"

	assert.Error(t, validation.ValidatePassword(badPassword))

	assert.NoError(t, validation.ValidatePassword(goodPassword))
}