package taskfile

import (
	"os"
	"testing"
)

func TestHasTask_ReturnsFalseWhenFileDoesNotExist(t *testing.T) {
	// Ensure file doesn't exist
	os.Remove("Taskfile.yml")

	hasTask, err := HasTask("test-task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasTask {
		t.Error("expected hasTask to be false when file doesn't exist")
	}
}
