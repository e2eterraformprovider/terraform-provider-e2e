package acceptance_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e"
)

func TestProviderCredentialSchema_NoHardcodedDefaults(t *testing.T) {
	p := e2e.Provider()

	apiKey := p.Schema["api_key"]
	if apiKey == nil {
		t.Fatalf("expected provider schema to include api_key")
	}
	if apiKey.DefaultFunc == nil {
		t.Fatalf("expected api_key to be sourced via DefaultFunc (env/config), got nil")
	}
	if apiKey.Default != nil && apiKey.Default != "" {
		t.Fatalf("expected api_key Default to be empty/nil, got: %#v", apiKey.Default)
	}

	authToken := p.Schema["auth_token"]
	if authToken == nil {
		t.Fatalf("expected provider schema to include auth_token")
	}
	if authToken.DefaultFunc == nil {
		t.Fatalf("expected auth_token to be sourced via DefaultFunc (env/config), got nil")
	}
	if authToken.Default != nil && authToken.Default != "" {
		t.Fatalf("expected auth_token Default to be empty/nil, got: %#v", authToken.Default)
	}
}

func TestNoHardcodedCredentialEnvSetInGoFiles(t *testing.T) {
	root := repoRoot(t)
	var offenders []string

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case "vendor", ".git":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		b, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		s := string(b)

		if strings.Contains(s, `os.Setenv("E2E_API_KEY"`) ||
			strings.Contains(s, `os.Setenv("E2E_AUTH_TOKEN"`) {
			offenders = append(offenders, path)
		}

		return nil
	})

	if len(offenders) > 0 {
		t.Fatalf("found os.Setenv for E2E_API_KEY/E2E_AUTH_TOKEN in non-test Go files: %v", offenders)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("unable to determine caller path")
	}

	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("unable to find repo root (go.mod) from %s", thisFile)
		}
		dir = parent
	}
}
