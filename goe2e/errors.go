package goe2e

import (
	"fmt"
	"net/http"
)

// ErrorResponse represents an error response from the E2E Networks API.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response `json:"-"`

	// Error message
	Message string `json:"message"`

	// HTTP status code
	Code int `json:"code"`

	// List of error details
	Errors []string `json:"errors"`
}

func (r *ErrorResponse) Error() string {
	if r.Response != nil {
		return fmt.Sprintf("%v %v: %d %v %v",
			r.Response.Request.Method,
			r.Response.Request.URL,
			r.Response.StatusCode,
			r.Message,
			r.Errors,
		)
	}
	return fmt.Sprintf("API error: %d %v %v", r.Code, r.Message, r.Errors)
}

// ArgError identifies an invalid argument passed to a method.
type ArgError struct {
	arg    string
	reason string
}

// NewArgError creates a new argument error
func NewArgError(arg, reason string) *ArgError {
	return &ArgError{
		arg:    arg,
		reason: reason,
	}
}

func (e *ArgError) Error() string {
	return fmt.Sprintf("%s is invalid because %s", e.arg, e.reason)
}
