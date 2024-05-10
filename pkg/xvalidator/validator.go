package xvalidator

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

// Validator is a struct that contains a pointer to a validator.Validate instance.
type Validator struct {
	validate *validator.Validate
}

// NewValidator is a function that initializes a new Validator instance.
// It registers a tag name function that returns the "name" tag of a struct field.
// It logs that the validator has been initialized and returns the new Validator instance.
func NewValidator() *Validator {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("name")
	})

	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phoneNumber := fl.Field().String()

		phoneString, err := phonenumbers.Parse(phoneNumber, "MY")
		if err != nil {
			return false
		}
		return phonenumbers.IsValidNumber(phoneString)
	})

	validate.RegisterValidation("email_if_type", func(fl validator.FieldLevel) bool {
		typeField := fl.Parent().FieldByName("Type").String()
		emailField := fl.Field().String()
		if typeField == "email" {
			if err := validator.New().Var(emailField, "email"); err != nil {
				return false
			}
		}
		return true
	})

	validate.RegisterValidation("phone_if_type", func(fl validator.FieldLevel) bool {
		typeField := fl.Parent().FieldByName("Type").String()
		phoneNumber := fl.Field().String()

		if typeField == "phone_number" {
			phoneString, err := phonenumbers.Parse(phoneNumber, "MY")
			if err != nil {
				return false
			}
			return phonenumbers.IsValidNumber(phoneString)
		}
		return true
	})

	slog.Info("validator initialized")
	return &Validator{validate: validate}
}

// Struct is a method of the Validator struct that validates a struct.
// It returns a slice of strings containing the validation errors.
// If there are no validation errors, it returns nil.
func (v *Validator) Struct(s interface{}) map[string]string {
	err := v.validate.Struct(s)
	if err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

// Var is a method of the Validator struct that validates a single variable.
// It returns a slice of strings containing the validation errors.
// If there are no validation errors, it returns nil.
func (v *Validator) Var(field interface{}, tag string) map[string]string {
	err := v.validate.Var(field, tag)
	if err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

// formatValidationError is a method of the Validator struct that formats validation errors.
// It returns a slice of strings containing the formatted validation errors.
func (v *Validator) formatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errors[err.Field()] = fmt.Sprintf("%s is required", err.Field())
		case "email":
			errors[err.Field()] = fmt.Sprintf("%s is not a valid email", err.Field())
		case "min":
			errors[err.Field()] = fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
		case "max":
			errors[err.Field()] = fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
		case "len":
			errors[err.Field()] = fmt.Sprintf("%s must be %s characters long", err.Field(), err.Param())
		case "gte":
			errors[err.Field()] = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
		case "gt":
			errors[err.Field()] = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
		case "lte":
			errors[err.Field()] = fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
		case "lt":
			errors[err.Field()] = fmt.Sprintf("%s must be less than %s", err.Field(), err.Param())
		case "numeric":
			errors[err.Field()] = fmt.Sprintf("%s must be numeric", err.Field())
		case "number":
			errors[err.Field()] = fmt.Sprintf("%s must be a number", err.Field())
		case "phone":
			errors[err.Field()] = fmt.Sprintf("%s invalid phone number", err.Field())
		default:
			errors[err.Field()] = fmt.Sprintf("%s is not valid", err.Field())
		}
	}
	return errors
}
