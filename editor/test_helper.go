package editor

import (
	"os"
	"testing"
)

// setupTestFile creates a temporary file with given contents for testing.
// It returns a path to the file.
func setupTestFile(t *testing.T, contents string) string {
	t.Helper()
	f, err := os.CreateTemp("", "test-*.hcl")
	if err != nil {
		t.Fatalf("failed to create test file: %s", err)
	}

	path := f.Name()
	t.Cleanup(func() { os.Remove(path) })

	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatalf("failed to write test file: %s", err)
	}

	return path
}

// readTestFile is a test helper for reading file with error handling.
func readTestFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test file: %s", err)
	}

	return string(b)
}
