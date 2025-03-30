package regex

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	regex := `lock\.[a-z]+`
	regex2 := `\.lock$`
	files := []string{"lock.yaml", "lock.json", "package-lock.json", "pnpm-lock.yaml", "package-lock.yaml", "pnpm-lock.yaml", "main.rs"}

	results := []string{}
	for _, file := range files {
		isMatch1, _ := regexp.MatchString(regex, file)
		if isMatch1 {
			continue
		}
		isMatch2, _ := regexp.MatchString(regex2, file)
		if isMatch2 {
			continue
		}
		results = append(results, file)
	}

	fmt.Printf("Files that should not be filtered: %v\n", results)
}

func TestShouldIgnoreFile(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		patterns      []string
		wantIgnore    bool
		wantErr       bool
		errorContains string
	}{
		{
			name:       "should match simple pattern",
			filename:   "test.lock",
			patterns:   []string{`\.lock$`},
			wantIgnore: true,
			wantErr:    false,
		},
		{
			name:       "should not match when pattern doesn't match",
			filename:   "test.txt",
			patterns:   []string{`\.lock$`},
			wantIgnore: false,
			wantErr:    false,
		},
		{
			name:       "should match any pattern from multiple patterns",
			filename:   "package-lock.json",
			patterns:   []string{`\.txt$`, `lock\.json$`},
			wantIgnore: true,
			wantErr:    false,
		},
		{
			name:       "should handle empty patterns slice",
			filename:   "test.lock",
			patterns:   []string{},
			wantIgnore: false,
			wantErr:    false,
		},
		{
			name:       "should match case-sensitive patterns",
			filename:   "README.md",
			patterns:   []string{`^ReADME\.md$`},
			wantIgnore: false,
			wantErr:    false,
		},
		{
			name:       "should handle complex regex patterns",
			filename:   "test-123_file.lock",
			patterns:   []string{`^test-\d+_.*\.lock$`},
			wantIgnore: true,
			wantErr:    false,
		},
		{
			name:          "should handle invalid regex pattern",
			filename:      "test.txt",
			patterns:      []string{`[invalid(`},
			wantIgnore:    false,
			wantErr:       true,
			errorContains: "error parsing regexp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ShouldIgnoreFile(tt.filename, &tt.patterns)

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("ShouldIgnoreFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errorContains != "" {
				if !contains(err.Error(), tt.errorContains) {
					t.Errorf("ShouldIgnoreFile() error = %v, want error containing %v", err, tt.errorContains)
				}
				return
			}

			if got != tt.wantIgnore {
				t.Errorf("ShouldIgnoreFile() = %v, want %v", got, tt.wantIgnore)
			}
		})
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}
