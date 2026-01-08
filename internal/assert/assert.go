package assert

import (
	"strings"
	"testing"
)

// =============================================================================
// Test Assertion Helpers
// =============================================================================

// StringContains checks if a string contains a substring
func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

// Equal checks if two values are equal using Go's == operator
//
// Uses generics to work with any comparable type
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

// NilError checks if an error is nil
func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
