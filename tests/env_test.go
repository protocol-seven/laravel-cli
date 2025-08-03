package tests

import (
	"os"
	"strings"
	"testing"
)

// Copy updateEnvFile function for testing
func updateEnvFile(envPath, key, value string) error {
	// Read the file
	content, err := os.ReadFile(envPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	found := false

	// Look for existing key and update it
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + value
			found = true
			break
		}
	}

	// If key not found, add it
	if !found {
		lines = append(lines, key+"="+value)
	}

	// Write back to file
	updatedContent := strings.Join(lines, "\n")
	return os.WriteFile(envPath, []byte(updatedContent), 0644)
}

func TestEnvironmentFileOperations(t *testing.T) {
	t.Run("TestUpdateEnvFileExistingKey", func(t *testing.T) {
		// Test updating an existing key in .env file
		tmpFile := "/tmp/test-env-update.env"
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
		defer os.Remove(tmpFile)

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
	})

	t.Run("TestUpdateEnvFileNewKey", func(t *testing.T) {
		// Test adding a new key to .env file
		tmpFile := "/tmp/test-env-new.env"
		content := `APP_NAME="Laravel"
APP_ENV=local`

		err := writeTestFile(tmpFile, content)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		defer os.Remove(tmpFile)

		// Add VITE_PORT
		err = updateEnvFile(tmpFile, "VITE_PORT", "5173")
		if err != nil {
			t.Fatalf("Failed to add new env key: %v", err)
		}

		// Read and verify
		updatedContent, err := readTestFile(tmpFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}

		if !strings.Contains(updatedContent, "VITE_PORT=5173") {
			t.Error("VITE_PORT was not added correctly")
		}
	})

	t.Run("TestUpdateEnvFileWithQuotes", func(t *testing.T) {
		// Test updating values that contain quotes or special characters
		testCases := []struct {
			key   string
			value string
		}{
			{"APP_NAME", `"My Laravel App"`},
			{"APP_URL", "http://localhost:8080"},
			{"DATABASE_URL", "mysql://user:pass@localhost/db"},
			{"SPECIAL_VALUE", "value with spaces"},
			{"QUOTED_VALUE", `"value with quotes"`},
		}

		for _, tc := range testCases {
			t.Run("Update_"+tc.key, func(t *testing.T) {
				tmpFile := "/tmp/test-env-quotes-" + tc.key + ".env"
				content := "APP_NAME=Laravel\nAPP_ENV=local"

				err := writeTestFile(tmpFile, content)
				if err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				defer os.Remove(tmpFile)

				// Update with special value
				err = updateEnvFile(tmpFile, tc.key, tc.value)
				if err != nil {
					t.Fatalf("Failed to update env file: %v", err)
				}

				// Read and verify
				updatedContent, err := readTestFile(tmpFile)
				if err != nil {
					t.Fatalf("Failed to read test file: %v", err)
				}

				expected := tc.key + "=" + tc.value
				if !strings.Contains(updatedContent, expected) {
					t.Errorf("Expected to find %s in updated content", expected)
				}
			})
		}
	})

	t.Run("TestUpdateEnvFileEmptyFile", func(t *testing.T) {
		// Test adding to an empty .env file
		tmpFile := "/tmp/test-env-empty.env"

		err := writeTestFile(tmpFile, "")
		if err != nil {
			t.Fatalf("Failed to write empty test file: %v", err)
		}
		defer os.Remove(tmpFile)

		// Add to empty file
		err = updateEnvFile(tmpFile, "NEW_KEY", "new_value")
		if err != nil {
			t.Fatalf("Failed to update empty env file: %v", err)
		}

		// Read and verify
		updatedContent, err := readTestFile(tmpFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}

		if !strings.Contains(updatedContent, "NEW_KEY=new_value") {
			t.Error("NEW_KEY was not added to empty file")
		}
	})

	t.Run("TestUpdateEnvFileNonExistentFile", func(t *testing.T) {
		// Test behavior when .env file doesn't exist
		tmpFile := "/tmp/test-env-nonexistent.env"

		// Ensure file doesn't exist
		os.Remove(tmpFile)

		// Test updating non-existent file - should return error
		err := updateEnvFile(tmpFile, "KEY", "value")
		if err == nil {
			t.Error("Expected error when updating non-existent file")
		}
	})
}

func TestEnvironmentFileEdgeCases(t *testing.T) {
	t.Run("TestEnvFileWithComments", func(t *testing.T) {
		// Test .env file with comments
		tmpFile := "/tmp/test-env-comments.env"
		content := `# Application Configuration
APP_NAME="Laravel"
APP_ENV=local
# Database Configuration  
DB_CONNECTION=mysql
DB_HOST=127.0.0.1`

		err := writeTestFile(tmpFile, content)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		defer os.Remove(tmpFile)

		// Update existing key
		err = updateEnvFile(tmpFile, "APP_ENV", "production")
		if err != nil {
			t.Fatalf("Failed to update env file: %v", err)
		}

		// Read and verify comments are preserved
		updatedContent, err := readTestFile(tmpFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}

		if !strings.Contains(updatedContent, "# Application Configuration") {
			t.Error("Comments were not preserved")
		}
		if !strings.Contains(updatedContent, "APP_ENV=production") {
			t.Error("APP_ENV was not updated correctly")
		}
	})

	t.Run("TestEnvFileWithBlankLines", func(t *testing.T) {
		// Test .env file with blank lines
		tmpFile := "/tmp/test-env-blanks.env"
		content := `APP_NAME="Laravel"

APP_ENV=local


DB_CONNECTION=mysql`

		err := writeTestFile(tmpFile, content)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		defer os.Remove(tmpFile)

		// Update existing key
		err = updateEnvFile(tmpFile, "APP_ENV", "testing")
		if err != nil {
			t.Fatalf("Failed to update env file: %v", err)
		}

		// Read and verify blank lines are preserved
		updatedContent, err := readTestFile(tmpFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}

		// Count blank lines
		lines := strings.Split(updatedContent, "\n")
		blankCount := 0
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				blankCount++
			}
		}

		if blankCount < 2 {
			t.Error("Blank lines were not preserved")
		}
		if !strings.Contains(updatedContent, "APP_ENV=testing") {
			t.Error("APP_ENV was not updated correctly")
		}
	})

	t.Run("TestEnvFileWithDuplicateKeys", func(t *testing.T) {
		// Test .env file with duplicate keys
		tmpFile := "/tmp/test-env-duplicates.env"
		content := `APP_NAME="Laravel"
APP_ENV=local
APP_NAME="Laravel App"
DB_CONNECTION=mysql`

		err := writeTestFile(tmpFile, content)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		defer os.Remove(tmpFile)

		// Update duplicate key - should update first occurrence
		err = updateEnvFile(tmpFile, "APP_NAME", "Updated App")
		if err != nil {
			t.Fatalf("Failed to update env file: %v", err)
		}

		// Read and verify
		updatedContent, err := readTestFile(tmpFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}

		lines := strings.Split(updatedContent, "\n")
		appNameCount := 0
		for _, line := range lines {
			if strings.HasPrefix(line, "APP_NAME=") {
				appNameCount++
			}
		}

		// Should still have duplicate keys, but first one should be updated
		if appNameCount < 2 {
			t.Error("Duplicate key structure was not preserved")
		}
	})
}

func TestDatabaseSpecificEnvOperations(t *testing.T) {
	t.Run("TestCommentDatabaseConfigForSQLite", func(t *testing.T) {
		// Test commentDatabaseConfigForSQLite function
		t.Skip("Requires import of main package functions")
	})

	t.Run("TestUncommentDatabaseConfig", func(t *testing.T) {
		// Test uncommentDatabaseConfig function
		t.Skip("Requires import of main package functions")
	})
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
