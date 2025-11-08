package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationError(validationErrors)
		}
		return err
	}
	return nil
}

func formatValidationError(errs validator.ValidationErrors) error {
	var messages []string
	for _, err := range errs {
		var message string
		field := strings.ToLower(err.Field())

		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
		case "max":
			message = fmt.Sprintf("%s must not exceed %s characters", field, err.Param())
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}
		messages = append(messages, message)
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}
