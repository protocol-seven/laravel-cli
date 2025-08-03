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

	if !strings.Contains(updatedContent, "APP_PORT=8080") {
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

	if !strings.Contains(updatedContent, "VITE_PORT=5173") {
		t.Error("VITE_PORT was not added correctly")
	}

	// Clean up
	removeTestFile(tmpFile)
}

func TestValidateProjectName(t *testing.T) {
	validNames := []string{"my-project", "my_project", "my.project", "MyProject123"}
	invalidNames := []string{"my project", "my@project", "my#project"}

	for _, name := range validNames {
		if err := validateProjectName(name); err != nil {
			t.Errorf("Expected %s to be valid, but got error: %v", name, err)
		}
	}

	for _, name := range invalidNames {
		if err := validateProjectName(name); err == nil {
			t.Errorf("Expected %s to be invalid, but got no error", name)
		}
	}
}

func TestGetStarterKit(t *testing.T) {
	// Reset all flags
	react = false
	vue = false
	livewire = false
	using = ""

	// Test React starter kit
	react = true
	if kit := getStarterKit(); kit != "laravel/react-starter-kit" {
		t.Errorf("Expected React starter kit, got %s", kit)
	}
	react = false

	// Test Vue starter kit
	vue = true
	if kit := getStarterKit(); kit != "laravel/vue-starter-kit" {
		t.Errorf("Expected Vue starter kit, got %s", kit)
	}
	vue = false

	// Test Livewire starter kit
	livewire = true
	if kit := getStarterKit(); kit != "laravel/livewire-starter-kit" {
		t.Errorf("Expected Livewire starter kit, got %s", kit)
	}
	livewire = false

	// Test custom starter kit
	using = "custom/starter-kit"
	if kit := getStarterKit(); kit != "custom/starter-kit" {
		t.Errorf("Expected custom starter kit, got %s", kit)
	}
	using = ""

	// Test no starter kit
	if kit := getStarterKit(); kit != "" {
		t.Errorf("Expected no starter kit, got %s", kit)
	}
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
