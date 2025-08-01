package main

import (
	"os"
	"strings"
	"testing"
)

func TestIsPortAvailable(t *testing.T) {
	// Test with a port that should be available (high port number)
	if !isPortAvailable(45678) {
		t.Error("Expected port 45678 to be available")
	}
}

func TestUpdateEnvFile(t *testing.T) {
	// Create a temporary .env file for testing
	tmpFile := "/tmp/test.env"
	content := `APP_NAME="Laravel"
APP_ENV=local
APP_KEY=
APP_DEBUG=true
APP_PORT=80
APP_URL=http://localhost`

	// Write test content
	err := writeTestFile(tmpFile, content)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Update APP_PORT
	err = updateEnvFile(tmpFile, "APP_PORT", "8080")
	if err != nil {
		t.Fatalf("Failed to update env file: %v", err)
	}

	// Read and verify
	updatedContent, err := readTestFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if !contains(updatedContent, "APP_PORT=8080") {
		t.Error("APP_PORT was not updated correctly")
	}

	// Add a new key
	err = updateEnvFile(tmpFile, "VITE_PORT", "5173")
	if err != nil {
		t.Fatalf("Failed to add new env key: %v", err)
	}

	// Read and verify
	updatedContent, err = readTestFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if !contains(updatedContent, "VITE_PORT=5173") {
		t.Error("VITE_PORT was not added correctly")
	}

	// Clean up
	removeTestFile(tmpFile)
}

// Helper functions for testing
func writeTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func readTestFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	return string(content), err
}

func removeTestFile(path string) error {
	return os.Remove(path)
}

func contains(text, substr string) bool {
	return strings.Contains(text, substr)
}
