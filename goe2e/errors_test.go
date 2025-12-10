package goe2e

import (
	"testing"
)

func TestArgError_Error(t *testing.T) {
	err := NewArgError("testArg", "cannot be empty")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expected := "testArg is invalid because cannot be empty"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}

	// Test that it satisfies error interface
	var _ error = err
}

func TestErrorResponse_Error(t *testing.T) {
	errResp := &ErrorResponse{
		Code:    400,
		Message: "API error occurred",
		Errors:  []string{},
	}

	expected := "API error: 400 API error occurred []"
	if errResp.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, errResp.Error())
	}

	// Test with errors list
	errWithDetails := &ErrorResponse{
		Code:    422,
		Message: "Validation failed",
		Errors:  []string{"field1 is required", "field2 is invalid"},
	}
	result := errWithDetails.Error()
	if result == "" {
		t.Error("Expected non-empty error message")
	}

	// Test that it satisfies error interface
	var _ error = errResp
}

func TestNewArgError(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		message string
		want    string
	}{
		{
			name:    "simple error",
			arg:     "id",
			message: "cannot be empty",
			want:    "id is invalid because cannot be empty",
		},
		{
			name:    "long error",
			arg:     "functionName",
			message: "must be between 3 and 64 characters",
			want:    "functionName is invalid because must be between 3 and 64 characters",
		},
		{
			name:    "empty arg",
			arg:     "",
			message: "test message",
			want:    " is invalid because test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewArgError(tt.arg, tt.message)
			if err.Error() != tt.want {
				t.Errorf("NewArgError() = %v, want %v", err.Error(), tt.want)
			}
			// Note: ArgError fields are private, we can only test Error() output
			if err.Error() != tt.want {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.want)
			}
		})
	}
}
