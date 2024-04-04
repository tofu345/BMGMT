package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
)

func init() {
	Validator = CustomValidator{validator: validator.New()}
	Validator.validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})
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
			if e.Param() != "" {
				err_map.Errors[e.Field()] = fmt.Sprintf("%v %v", e.Tag(), e.Param())
			} else {
				err_map.Errors[e.Field()] = e.Tag()
			}
		}
		return err_map
	}

	return ValidationErrs{Errors: nil}
}
