package util

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ToErrResponse(err error) []string {
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		errs := make([]string, len(fieldErrors))

		for i, err := range fieldErrors {

			lowercaseField := strings.ToLower(err.Field())

			switch err.Tag() {
			case "alphanum":
				errs[i] = fmt.Sprintf("%s can only contain alphabetic and numeric characters", lowercaseField)
			case "email":
				errs[i] = fmt.Sprintf("%s is not a valid email address", lowercaseField)
			case "required":
				errs[i] = fmt.Sprintf("%s is a required field", lowercaseField)
			case "max":
				errs[i] = fmt.Sprintf("%s must be a maximum of %s in length", lowercaseField, err.Param())
			default:
				errs[i] = fmt.Sprintf("something wrong on field %s; %s", lowercaseField, err.Tag())
			}
		}

		return errs
	}

	return nil
}
