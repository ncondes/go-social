package handlers

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *Validator) validateStruct(data any) []string {
	if err := v.validate.Struct(data); err != nil {
		errs, ok := err.(validator.ValidationErrors)

		if ok {
			var messages []string

			for _, err := range errs {
				messages = append(messages, formatFieldError(err))
			}

			return messages
		}
	}

	return nil
}

func formatFieldError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()
	field = strings.ToLower(field)

	messages := map[string]string{
		"required": field + " is required",
		"max":      field + " must be at most " + param + " characters",
		"min":      field + " must be at least " + param + " characters",
	}

	message, exists := messages[tag]
	if exists {
		return message
	}

	return err.Error()
}
