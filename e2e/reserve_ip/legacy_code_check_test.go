package reserve_ip

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ============================================================================
// Legacy Code Verification Tests
// ============================================================================

func TestLegacyCodeVerification_NoLegacyClientInMainCode(t *testing.T) {
	// Get the directory of the current test file
	testDir := "."
	resourceFile := filepath.Join(testDir, "resource_reserve_ip.go")

	// Read the file content
	content, err := os.ReadFile(resourceFile)
	if err != nil {
		t.Fatalf("Failed to read resource file: %v", err)
	}

	contentStr := string(content)

	// Check for legacy client patterns
	legacyPatterns := []struct {
		name    string
		pattern string
		allowed bool // Whether this pattern is allowed (e.g., in comments)
	}{
		{
			name:    "No client.Client references",
			pattern: "client.Client",
			allowed: false,
		},
		{
			name:    "No apiClient references",
			pattern: "apiClient.",
			allowed: false,
		},
	}

	for _, check := range legacyPatterns {
		t.Run(check.name, func(t *testing.T) {
			// Count occurrences
			count := strings.Count(contentStr, check.pattern)
			if count > 0 {
				// Check if it's only in comments
				lines := strings.Split(contentStr, "\n")
				nonCommentCount := 0
				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					// Skip comment lines
					if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
						continue
					}
					if strings.Contains(line, check.pattern) {
						nonCommentCount++
					}
				}
				if nonCommentCount > 0 && !check.allowed {
					t.Errorf("Found %d occurrences of legacy pattern %q in non-comment code", nonCommentCount, check.pattern)
				}
			}
		})
	}
}

func TestLegacyCodeVerification_NoLegacyClientInFloatingIPAttachment(t *testing.T) {
	// Check floating_ip_attachment directory
	attachmentDir := "../floating_ip_attachment"
	resourceFile := filepath.Join(attachmentDir, "resource_floating_ip_attachment.go")

	// Check if file exists
	if _, err := os.Stat(resourceFile); os.IsNotExist(err) {
		t.Skip("floating_ip_attachment resource file not found")
		return
	}

	content, err := os.ReadFile(resourceFile)
	if err != nil {
		t.Fatalf("Failed to read floating_ip_attachment resource file: %v", err)
	}

	contentStr := string(content)

	// Check for legacy client patterns
	legacyPatterns := []string{
		"client.Client",
		"apiClient.",
	}

	for _, pattern := range legacyPatterns {
		t.Run("No "+pattern+" in floating_ip_attachment", func(t *testing.T) {
			// Count occurrences in non-comment code
			lines := strings.Split(contentStr, "\n")
			nonCommentCount := 0
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				// Skip comment lines
				if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
					continue
				}
				if strings.Contains(line, pattern) {
					nonCommentCount++
				}
			}
			if nonCommentCount > 0 {
				t.Errorf("Found %d occurrences of legacy pattern %q in non-comment code", nonCommentCount, pattern)
			}
		})
	}
}

func TestLegacyCodeVerification_UsesGoe2eClient(t *testing.T) {
	// Verify that the code uses goe2e client
	testDir := "."
	resourceFile := filepath.Join(testDir, "resource_reserve_ip.go")

	content, err := os.ReadFile(resourceFile)
	if err != nil {
		t.Fatalf("Failed to read resource file: %v", err)
	}

	contentStr := string(content)

	// Verify goe2e client patterns exist
	goe2ePatterns := []string{
		"goe2eClient",
		"Goe2eClient",
		"goe2e.Client",
	}

	foundPattern := false
	for _, pattern := range goe2ePatterns {
		if strings.Contains(contentStr, pattern) {
			foundPattern = true
			break
		}
	}

	if !foundPattern {
		t.Error("Resource should use goe2e client patterns")
	}
}

// ============================================================================
// Static Analysis Helper
// ============================================================================

func runGrepCheck(t *testing.T, pattern, path string) []string {
	cmd := exec.Command("grep", "-r", pattern, path)
	output, err := cmd.Output()
	if err != nil {
		// grep returns non-zero exit code when no matches found, which is OK
		return []string{}
	}
	lines := strings.Split(string(output), "\n")
	// Filter out empty lines
	result := make([]string, 0)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func TestLegacyCodeVerification_GrepCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping grep checks in short mode")
	}

	tests := []struct {
		name    string
		pattern string
		path    string
		expect  int // Expected number of matches (0 = no matches allowed)
	}{
		{
			name:    "No client.Client in reserve_ip",
			pattern: "client\\.Client",
			path:    ".",
			expect:  0,
		},
		{
			name:    "No apiClient in reserve_ip",
			pattern: "apiClient\\.",
			path:    ".",
			expect:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := runGrepCheck(t, tt.pattern, tt.path)
			// Filter out test files and comments
			nonTestMatches := make([]string, 0)
			for _, match := range matches {
				// Skip test files
				if strings.Contains(match, "_test.go") {
					continue
				}
				// Skip if it's in a comment
				if strings.Contains(match, "//") || strings.Contains(match, "/*") {
					continue
				}
				nonTestMatches = append(nonTestMatches, match)
			}
			if len(nonTestMatches) != tt.expect {
				t.Errorf("Expected %d matches, got %d. Matches: %v", tt.expect, len(nonTestMatches), nonTestMatches)
			}
		})
	}
}
