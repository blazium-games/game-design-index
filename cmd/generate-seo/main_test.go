package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCollectPagesStaticRoutes(t *testing.T) {
	apiRoot := filepath.Join("..", "..", "data", "dist", "api", "v1")
	if _, err := os.Stat(apiRoot); err != nil {
		t.Skip("exported API not present:", err)
	}
	pages, _, err := collectPages(apiRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) < 100 {
		t.Fatalf("expected many pages, got %d", len(pages))
	}
	foundHome := false
	foundMechanic := false
	for _, p := range pages {
		if p.Path == "/" {
			foundHome = true
		}
		if p.Path == "/mechanics/melee-primary-combat" {
			foundMechanic = true
			if p.Title == "" {
				t.Fatal("mechanic page missing title")
			}
		}
	}
	if !foundHome {
		t.Fatal("missing home page")
	}
	if !foundMechanic {
		t.Fatal("missing melee-primary-combat page")
	}
}

func TestInjectMetaContainsOgTitle(t *testing.T) {
	template := `<!doctype html><html><head><!-- seo-defaults -->
    <title>Game Design Index</title>
    <meta name="description" content="default" />
    <!-- /seo-defaults --></head><body></body></html>`
	out := injectMeta(template, pageMeta{
		Path:        "/mechanics/test",
		Title:       "Test Mechanic · Game Design Index",
		Description: "A test mechanic summary.",
		OgType:      "article",
	}, "https://blazium-games.github.io/game-design-index", "https://blazium-games.github.io/game-design-index/og-default.png")
	if !strings.Contains(out, `property="og:title"`) {
		t.Fatal("missing og:title")
	}
	if !strings.Contains(out, "Test Mechanic · Game Design Index") {
		t.Fatal("missing title content")
	}
	if strings.Contains(out, "<!-- seo-defaults -->") {
		t.Fatal("seo-defaults marker should be replaced")
	}
}

func TestWriteSitemap(t *testing.T) {
	dir := t.TempDir()
	pages := []pageMeta{
		{Path: "/", Title: "Home", Description: "Home", OgType: "website"},
		{Path: "/games/foo", Title: "Foo", Description: "Foo game", OgType: "article"},
	}
	base := "https://example.com/game-design-index"
	if err := writeSitemap(dir, base, pages, "2026-07-04"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "sitemap.xml"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "<loc>https://example.com/game-design-index/</loc>") {
		t.Fatal("missing home URL in sitemap")
	}
	if !strings.Contains(content, "<loc>https://example.com/game-design-index/games/foo</loc>") {
		t.Fatal("missing game URL in sitemap")
	}
}

func TestWriteSPAFallback(t *testing.T) {
	dir := t.TempDir()
	indexHTML := `<!doctype html><html><head><title>Game Design Index</title></head><body></body></html>`
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(indexHTML), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeSPAFallback(dir); err != nil {
		t.Fatal(err)
	}
	notFound, err := os.ReadFile(filepath.Join(dir, "404.html"))
	if err != nil {
		t.Fatal("404.html missing:", err)
	}
	if string(notFound) != indexHTML {
		t.Fatal("404.html should match root index.html")
	}
	if _, err := os.Stat(filepath.Join(dir, ".nojekyll")); err != nil {
		t.Fatal(".nojekyll missing:", err)
	}
}
