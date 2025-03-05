package validator

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
// and it is concurrent-safe
var validate *validator.Validate

func Validate(i interface{}) error {
	return validate.Struct(i)
}

func GetValidationErrors(err error) (*validator.ValidationErrors, bool) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return &validationErrors, true
	}

	return nil, false
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// TODO: support custom tag names other than json
		return fld.Tag.Get("json")
	})
}
