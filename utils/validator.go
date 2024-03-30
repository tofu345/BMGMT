package utils

import (
	"fmt"

	"github.com/go-playground/validator"
)

func init() {
	Validator = CustomValidator{validator: validator.New()}
}

var Validator CustomValidator

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return FmtValidationErrs(err)
	}
	return nil
}

type ValidationErrs struct {
	Errors map[string]any `json:"errors"`
}

func (e ValidationErrs) Error() string {
	return fmt.Sprint(e.Errors)
}

func FmtValidationErrs(err error) ValidationErrs {
	switch errs := err.(type) {
	case validator.ValidationErrors:
		err_map := ValidationErrs{Errors: make(map[string]any, len(errs))}
		for _, e := range errs {
			err_map.Errors[e.Field()] = fmt.Sprintf("%v %v", e.Tag(), e.Param())
		}
		return err_map
	}

	return ValidationErrs{Errors: nil}
}
