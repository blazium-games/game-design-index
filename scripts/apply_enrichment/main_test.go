package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeCatalogSkills(t *testing.T) {
	root := t.TempDir()
	libDir := filepath.Join(root, "library")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatal(err)
	}

	initial := map[string]any{
		"schema_version": "1.0",
		"skills": []map[string]any{
			{
				"schema_version":    "1.0",
				"slug":              "exploration",
				"name":              "Exploration",
				"summary":           "Discovering routes and secrets.",
				"category":          "cognitive",
				"learning_outcome":  "Players read environmental cues.",
				"tags":              []string{"exploration"},
				"related_mechanics": []string{"open-world-exploration"},
			},
		},
	}
	raw, err := json.MarshalIndent(initial, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	raw = append(raw, '\n')
	skillsPath := filepath.Join(libDir, "skills.json")
	if err := os.WriteFile(skillsPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	patch := map[string]json.RawMessage{
		"exploration": json.RawMessage(`{
			"design_guidance": {
				"when_to_use": "When discovery drives engagement.",
				"where_to_use": "Open worlds and optional POIs.",
				"designer_notes": "Gate rewards with readable affordances."
			},
			"related_variables": ["stamina"]
		}`),
	}
	if err := mergeCatalog(skillsPath, "skills", patch); err != nil {
		t.Fatal(err)
	}

	updated, err := os.ReadFile(skillsPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Skills []map[string]any `json:"skills"`
	}
	if err := json.Unmarshal(updated, &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.Skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(doc.Skills))
	}
	skill := doc.Skills[0]
	dg, ok := skill["design_guidance"].(map[string]any)
	if !ok || dg["when_to_use"] != "When discovery drives engagement." {
		t.Fatalf("design_guidance not merged: %#v", skill["design_guidance"])
	}
	vars, ok := skill["related_variables"].([]any)
	if !ok || len(vars) != 1 || vars[0] != "stamina" {
		t.Fatalf("related_variables not merged: %#v", skill["related_variables"])
	}
	if skill["slug"] != "exploration" {
		t.Fatalf("slug changed: %v", skill["slug"])
	}
}

func TestAppendCatalogSkills(t *testing.T) {
	root := t.TempDir()
	libDir := filepath.Join(root, "library")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatal(err)
	}

	initial := map[string]any{
		"schema_version": "1.0",
		"skills": []map[string]any{
			{
				"schema_version":   "1.0",
				"slug":             "exploration",
				"name":             "Exploration",
				"summary":          "Discovering routes and secrets.",
				"category":         "cognitive",
				"learning_outcome": "Players read environmental cues.",
				"tags":             []string{"exploration"},
			},
		},
	}
	raw, err := json.MarshalIndent(initial, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	raw = append(raw, '\n')
	skillsPath := filepath.Join(libDir, "skills.json")
	if err := os.WriteFile(skillsPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	additions := []json.RawMessage{
		json.RawMessage(`{
			"schema_version": "1.0",
			"slug": "aim-precision",
			"name": "Aim Precision",
			"summary": "Landing shots on moving targets under pressure.",
			"category": "motor",
			"learning_outcome": "Players stabilize aim and lead targets.",
			"tags": ["combat", "timing"]
		}`),
	}
	if err := appendCatalog(skillsPath, "skills", additions); err != nil {
		t.Fatal(err)
	}

	updated, err := os.ReadFile(skillsPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Skills []map[string]any `json:"skills"`
	}
	if err := json.Unmarshal(updated, &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.Skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(doc.Skills))
	}
	if doc.Skills[1]["slug"] != "aim-precision" {
		t.Fatalf("unexpected append: %#v", doc.Skills[1])
	}
}

func TestWriteMapIfAbsent(t *testing.T) {
	root := t.TempDir()
	mapPath := filepath.Join(root, "maps", "puzzle-game.json")
	doc := json.RawMessage(`{"slug":"puzzle-game","title":"Puzzle Game","map_type":"genre"}`)
	if err := writeMapIfAbsent(mapPath, doc); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(mapPath); err != nil {
		t.Fatal(err)
	}
	if err := writeMapIfAbsent(mapPath, json.RawMessage(`{"slug":"changed"}`)); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(mapPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "puzzle-game") {
		t.Fatalf("map overwritten: %s", raw)
	}
}