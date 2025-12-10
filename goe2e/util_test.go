package goe2e

import (
	"testing"
	"time"
)

func TestPtrTo(t *testing.T) {
	// Test int
	i := 42
	ptr := PtrTo(i)
	if ptr == nil {
		t.Fatal("Expected non-nil pointer")
	}
	//nolint:SA5011 // t.Fatal stops execution, ptr cannot be nil here
	//nolint:SA5011 // t.Fatal stops execution
	if *ptr != 42 {
		t.Errorf("Expected *ptr = 42, got %d", *ptr)
	}

	// Test string
	s := "test"
	ptrStr := PtrTo(s)
	if ptrStr == nil {
		t.Fatal("Expected non-nil pointer")
	}
	//nolint:SA5011 // t.Fatal stops execution
	//nolint:SA5011 // t.Fatal stops execution
	if *ptrStr != "test" {
		t.Errorf("Expected *ptrStr = test, got %s", *ptrStr)
	}
}

func TestString(t *testing.T) {
	s := "hello"
	ptr := String(s)
	if ptr == nil {
		t.Fatal("Expected non-nil pointer")
	}
	//nolint:SA5011 // t.Fatal stops execution
	if *ptr != "hello" {
		t.Errorf("Expected *ptr = hello, got %s", *ptr)
	}

	// Test empty string
	empty := String("")
	if empty == nil {
		t.Fatal("Expected non-nil pointer for empty string")
	}
	//nolint:SA5011 // t.Fatal stops execution
	if *empty != "" {
		t.Errorf("Expected *empty = '', got %s", *empty)
	}
}

func TestInt(t *testing.T) {
	tests := []int{0, 1, -1, 42, -42, 999}
	for _, val := range tests {
		ptr := Int(val)
		if ptr == nil {
			t.Fatalf("Expected non-nil pointer for %d", val)
		}
		//nolint:SA5011 // t.Fatal stops execution
		if *ptr != val {
			t.Errorf("Expected *ptr = %d, got %d", val, *ptr)
		}
	}
}

func TestBool(t *testing.T) {
	// Test true
	truePtr := Bool(true)
	if truePtr == nil {
		t.Fatal("Expected non-nil pointer")
	}
	//nolint:SA5011 // t.Fatal stops execution
	if *truePtr != true {
		t.Error("Expected *truePtr = true")
	}

	// Test false
	falsePtr := Bool(false)
	if falsePtr == nil {
		t.Fatal("Expected non-nil pointer")
	}
	//nolint:SA5011 // t.Fatal stops execution
	if *falsePtr != false {
		t.Error("Expected *falsePtr = false")
	}
}

func TestTime(t *testing.T) {
	now := time.Now()
	ptr := Time(now)
	if ptr == nil {
		t.Fatal("Expected non-nil pointer")
	}
	//nolint:SA5011 // t.Fatal stops execution
	if !ptr.Equal(now) {
		t.Errorf("Expected *ptr = %v, got %v", now, *ptr)
	}

	// Test zero time
	zero := time.Time{}
	zeroPtr := Time(zero)
	if zeroPtr == nil {
		t.Fatal("Expected non-nil pointer for zero time")
	}
	//nolint:SA5011 // t.Fatal stops execution
	if !zeroPtr.Equal(zero) {
		t.Errorf("Expected *zeroPtr = zero time, got %v", *zeroPtr)
	}
}
