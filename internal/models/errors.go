package models

import "errors"

// =============================================================================
// Custom Error Definitions
// =============================================================================

var (
	// ErrNoRecord is returned when a database query returns no rows
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail is returned when attempting to create a user with
	// an email address that already exists in the database
	ErrDuplicateEmail = errors.New("models: this email is already signed up")
)
