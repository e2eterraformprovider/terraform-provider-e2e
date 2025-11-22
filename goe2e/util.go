package goe2e

import "time"

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
