package goe2e

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// PtrTo returns a pointer to the value passed in.
func PtrTo[T any](v T) *T {
	return &v
}

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// Int returns a pointer to the int value passed in.
func Int(v int) *int {
	return &v
}

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// Time returns a pointer to the time.Time value passed in.
func Time(v time.Time) *time.Time {
	return &v
}

// SanitizePathParam sanitizes a path parameter for safe use in URLs.
// It trims whitespace and applies URL path escaping to prevent path traversal attacks.
func SanitizePathParam(param string) string {
	return url.PathEscape(strings.TrimSpace(param))
}

// ExtractStringID attempts to extract a string ID from a map[string]any.
// It handles multiple types: string, float64, int, and int64.
// This is useful for parsing API responses that may return IDs in different formats.
func ExtractStringID(data map[string]any, key string) (string, bool) {
	if data == nil {
		return "", false
	}

	val, exists := data[key]
	if !exists {
		return "", false
	}

	switch v := val.(type) {
	case string:
		return v, true
	case float64:
		return fmt.Sprintf("%.0f", v), true
	case int:
		return fmt.Sprintf("%d", v), true
	case int64:
		return fmt.Sprintf("%d", v), true
	default:
		return "", false
	}
}

// IsNotFoundResponse checks if the response is a 404 Not Found response.
// This provides a consistent way to check for resource-not-found scenarios
// across all service methods.
func IsNotFoundResponse(resp *Response) bool {
	return resp != nil && resp.StatusCode == 404
}
