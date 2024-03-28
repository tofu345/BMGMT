package utils

import (
	"fmt"

	"github.com/go-playground/validator"
)

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
			err_map.Errors[e.Field()] = e.Tag()
		}
		return err_map
	}

	return ValidationErrs{Errors: nil}
}
