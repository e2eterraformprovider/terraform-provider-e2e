package goe2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newMalformedJSONServer creates a test server that returns intentionally malformed JSON
func newMalformedJSONServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{invalid json}`) // Intentionally malformed
	}))
}

// newNullFieldServer creates a test server that returns response with null fields
func newNullFieldServer(t *testing.T, nullFields map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonData, _ := json.Marshal(nullFields)
		writeJSON(w, string(jsonData))
	}))
}

// newMissingFieldServer creates a test server that returns response with missing required fields
func newMissingFieldServer(t *testing.T, incompleteData map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonData, _ := json.Marshal(incompleteData)
		writeJSON(w, string(jsonData))
	}))
}

// newInvalidFieldTypeServer creates a test server that returns fields with wrong data types
func newInvalidFieldTypeServer(t *testing.T, invalidFields map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonData, _ := json.Marshal(invalidFields)
		writeJSON(w, string(jsonData))
	}))
}
