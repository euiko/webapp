package validator

import "github.com/go-playground/validator/v10"

// use a single instance of Validate, it caches struct info
// and it is concurrent-safe
var validate *validator.Validate

func Validate(i interface{}) error {
	return validate.Struct(i)
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
