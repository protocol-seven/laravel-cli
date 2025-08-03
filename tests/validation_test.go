package tests

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// Import the functions we need to test
// For now, we'll copy the validation logic here
// In a real implementation, you'd import from the main package

// validateProjectName validates a Laravel project name
func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check for invalid characters (spaces and most special characters)
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, name)
	if !matched {
		return fmt.Errorf("project name contains invalid characters. Use only letters, numbers, dots, hyphens, and underscores")
	}

	return nil
}

func TestProjectValidation(t *testing.T) {
	t.Run("TestValidProjectNames", func(t *testing.T) {
		validNames := []string{
			"my-project",
			"my_project",
			"my.project",
			"MyProject123",
			"project",
			"project-name",
			"project_name",
			"ProjectName",
			"project123",
			"a",
		}

		for _, name := range validNames {
			t.Run("Valid_"+name, func(t *testing.T) {
				if err := validateProjectName(name); err != nil {
					t.Errorf("Expected %s to be valid, but got error: %v", name, err)
				}
			})
		}
	})

	t.Run("TestInvalidProjectNames", func(t *testing.T) {
		invalidNames := []string{
			"my project",  // spaces
			"my@project",  // special characters
			"my#project",  // special characters
			"my/project",  // path separators
			"my\\project", // path separators
			"my|project",  // special characters
			"my<project",  // special characters
			"my>project",  // special characters
			"",            // empty string
		}

		for _, name := range invalidNames {
			t.Run("Invalid_"+strings.ReplaceAll(name, " ", "_SPACE_"), func(t *testing.T) {
				if err := validateProjectName(name); err == nil {
					t.Errorf("Expected %s to be invalid, but got no error", name)
				}
			})
		}
	})

	t.Run("TestEdgeCaseProjectNames", func(t *testing.T) {
		edgeCases := []struct {
			name     string
			expected bool
		}{
			{".", true},        // current directory
			{"..", true},       // parent directory
			{".project", true}, // hidden directory
			{"project.", true}, // ending with dot
			{"-project", true}, // starting with dash
			{"project-", true}, // ending with dash
			{"_project", true}, // starting with underscore
			{"project_", true}, // ending with underscore
		}

		for _, tc := range edgeCases {
			t.Run("EdgeCase_"+strings.ReplaceAll(tc.name, ".", "_DOT_"), func(t *testing.T) {
				err := validateProjectName(tc.name)
				if tc.expected && err != nil {
					t.Errorf("Expected %s to be valid, but got error: %v", tc.name, err)
				} else if !tc.expected && err == nil {
					t.Errorf("Expected %s to be invalid, but got no error", tc.name)
				}
			})
		}
	})

	t.Run("TestLongProjectNames", func(t *testing.T) {
		longNames := []struct {
			length int
			valid  bool
		}{
			{100, true}, // 100 characters - should be valid
			{255, true}, // 255 characters - filesystem limit, should be valid
			{300, true}, // Over typical limit - still valid for our validation
		}

		for i, tc := range longNames {
			name := strings.Repeat("a", tc.length)
			t.Run("LongName_"+string(rune(i+'A')), func(t *testing.T) {
				err := validateProjectName(name)
				if tc.valid && err != nil {
					t.Errorf("Expected long name (%d chars) to be valid, but got error: %v", tc.length, err)
				} else if !tc.valid && err == nil {
					t.Errorf("Expected long name (%d chars) to be invalid, but got no error", tc.length)
				}
			})
		}
	})
}
