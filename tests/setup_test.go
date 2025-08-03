package tests

// This file provides test utilities and setup functions that can be shared
// across all test files in the tests package.

import (
	"os"
	"testing"
)

// TestMain provides setup and teardown for all tests in this package
func TestMain(m *testing.M) {
	// Setup: Create any necessary test fixtures
	setupTestEnvironment()

	// Run all tests
	code := m.Run()

	// Teardown: Clean up any test artifacts
	teardownTestEnvironment()

	// Exit with the same code as the test run
	os.Exit(code)
}

func setupTestEnvironment() {
	// Create temporary directories for testing
	os.MkdirAll("/tmp/laravel-cli-tests", 0755)
}

func teardownTestEnvironment() {
	// Clean up temporary test files and directories
	os.RemoveAll("/tmp/laravel-cli-tests")
}

// Test helper functions that can be used across all test files

// WriteTestFile creates a temporary file with the given content
func WriteTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// ReadTestFile reads the content of a file
func ReadTestFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	return string(content), err
}

// RemoveTestFile removes a test file
func RemoveTestFile(path string) error {
	return os.Remove(path)
}

// CreateTempTestDir creates a temporary directory for testing
func CreateTempTestDir(name string) (string, error) {
	dir := "/tmp/laravel-cli-tests/" + name
	err := os.MkdirAll(dir, 0755)
	return dir, err
}

// RemoveTempTestDir removes a temporary test directory
func RemoveTempTestDir(path string) error {
	return os.RemoveAll(path)
}
