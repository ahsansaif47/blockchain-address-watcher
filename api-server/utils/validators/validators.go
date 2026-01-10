package validators

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// NewValidator creates a new validator instance and registers custom validators
func NewValidator() *validator.Validate {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("strong_password", validateStrongPassword)

	return v
}

// validatePhone field validator
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true
	}

	re := regexp.MustCompile(`^\+?\d{10,15}$`)
	return re.MatchString(phone)

}

// validateStrongPassword validates password strength
// Minimum 8 characters, at least one uppercase, one lowercase, one digit, one special character
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// GetValidationErrors extracts validation error messages
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			switch tag {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errors[field] = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
			case "phone":
				errors[field] = fmt.Sprintf("%s must be a valid phone number", field)
			case "strong_password":
				errors[field] = fmt.Sprintf("%s must be at least 8 characters with uppercase, lowercase, digit, and special character", field)
			default:
				errors[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
	}

	return errors
}
