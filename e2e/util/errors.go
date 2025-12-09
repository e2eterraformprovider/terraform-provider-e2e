package util

import (
	"strings"
)

// IsNotFoundError checks if an error indicates a resource was not found
// This is useful for handling 404 errors and similar "not found" scenarios
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	// Check common error patterns for not found
	errStr := err.Error()
	return strings.Contains(strings.ToLower(errStr), "not found") ||
		strings.Contains(errStr, "404") ||
		strings.Contains(strings.ToLower(errStr), "does not exist")
}
