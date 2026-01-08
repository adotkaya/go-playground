package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// =============================================================================
// Validator Type
// =============================================================================

// Validator struct contains validation errors for forms
type Validator struct {
	NonFieldErrors []string          // Errors not specific to any field
	FieldErrors    map[string]string // Field-specific validation errors
}

// =============================================================================
// Email Regular Expression
// =============================================================================

// EmailRX is a regular expression for validating email addresses
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// =============================================================================
// Validator Methods
// =============================================================================

// Valid returns true if there are no validation errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddNonFieldError adds a non-field-specific error message
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// AddFieldError adds an error message for a specific form field
//
// If an error already exists for the field, it will not be overwritten
func (v *Validator) AddFieldError(key, message string) {
	// Initialize the map if it doesn't exist
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	// Only add if an error doesn't already exist for this field
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds a field error if the validation check fails
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// =============================================================================
// Validation Functions
// =============================================================================

// NotBlank returns true if a value is not an empty string or whitespace
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinChars returns true if a value contains at least n characters
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// MaxChars returns true if a value contains no more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue returns true if a value matches one of the permitted values
//
// Uses Go generics to work with any comparable type (strings, ints, etc.)
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// Matches returns true if a value matches a provided regular expression pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
