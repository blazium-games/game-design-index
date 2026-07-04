package main_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/blazium-games/game-design-index/internal"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func TestLoadSourceCorpus(t *testing.T) {
	root := repoRoot(t)
	src := filepath.Join(root, "data", "source")
	b, err := internal.LoadDir(src)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	if len(b.Maps) != 1401 {
		t.Errorf("expected 1401 maps, got %d", len(b.Maps))
	}
	if len(b.Mechanics) < 248 {
		t.Errorf("expected at least 248 mechanics, got %d", len(b.Mechanics))
	}
}

func TestExportProducesCatalog(t *testing.T) {
	root := repoRoot(t)
	catalog := filepath.Join(root, "data", "dist", "api", "v1", "catalog.json")
	if _, err := os.Stat(catalog); os.IsNotExist(err) {
		t.Skip("run go run ./cmd/export first")
	}
	data, err := os.ReadFile(catalog)
	if err != nil {
		t.Fatal(err)
	}
	lower := string(data)
	for _, term := range []string{"godot", "ludonaut", "source_url"} {
		if contains(lower, term) {
			t.Errorf("catalog.json contains banned term %q", term)
		}
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexFold(s, sub) >= 0)
}

func indexFold(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if equalFold(s[i:i+len(sub)], sub) {
			return i
		}
	}
	return -1
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

func TestExportProducesAnalyticsAndChangelog(t *testing.T) {
	root := repoRoot(t)
	src := filepath.Join(root, "data", "source")
	b, err := internal.LoadDir(src)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	analytics := internal.BuildAnalytics(b)
	if len(analytics.TopMechanics) == 0 {
		t.Error("analytics missing top_mechanics")
	}
	if len(analytics.Insights) == 0 {
		t.Error("analytics missing insights")
	}
	changelogPath := filepath.Join(src, "changelog.json")
	if _, err := os.Stat(changelogPath); err != nil {
		t.Fatalf("changelog.json: %v", err)
	}
	analyticsPath := filepath.Join(root, "data", "dist", "api", "v1", "analytics.json")
	if _, err := os.Stat(analyticsPath); os.IsNotExist(err) {
		t.Skip("run go run ./cmd/export first")
	}
	changelogOut := filepath.Join(root, "data", "dist", "api", "v1", "changelog.json")
	if _, err := os.Stat(changelogOut); os.IsNotExist(err) {
		t.Skip("run go run ./cmd/export first")
	}
}
