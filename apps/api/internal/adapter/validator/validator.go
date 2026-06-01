package validator

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the go-playground validator with custom rules.
type CustomValidator struct {
	validate *validator.Validate
}

// New creates a new custom validator with registered rules.
func New() *CustomValidator {
	v := validator.New()

	// Register custom validations
	v.RegisterValidation("date", validateDate)
	v.RegisterValidation("future_date", validateFutureDate)
	v.RegisterValidation("phone_id", validatePhoneID)
	v.RegisterValidation("uuid", validateUUID)

	return &CustomValidator{validate: v}
}

// Validate validates a struct using registered rules.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

// Engine returns the underlying validator for Gin binding.
func (cv *CustomValidator) Engine() *validator.Validate {
	return cv.validate
}

// validateDate checks if value is a valid date (YYYY-MM-DD).
func validateDate(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

// validateFutureDate checks if value is a date in the future.
func validateFutureDate(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	return date.After(time.Now().Truncate(24 * time.Hour))
}

// validatePhoneID checks if value is a valid Indonesian phone number.
func validatePhoneID(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^(\+62|62|0)8[1-9][0-9]{6,11}$`, phone)
	return matched
}

// validateUUID checks if value is a valid UUID v4.
func validateUUID(fl validator.FieldLevel) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(fl.Field().String())
}
