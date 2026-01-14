package helpers

import "testing"

func TestGetContent_DoesNotPanic(t *testing.T) {
	// This test will fail with "panic: assignment to entry in nil map"
	// if the nested map is not properly initialized
	result := GetContent()

	if result == nil {
		t.Error("GetContent() returned nil")
	}
}
