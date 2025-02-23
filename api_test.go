package main

import (
	"testing"
)

func TestSampleFunction(t *testing.T) {
	expectedValue := nil // Expected output

	APIServer := NewAPIServer(":8000", store)

	result := APIServer.Serve()

	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}
