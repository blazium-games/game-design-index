package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type tagListFile struct {
	SchemaVersion string   `json:"schema_version"`
	Tags          []string `json:"tags"`
}

func LoadVariableTags(root string) (map[string]struct{}, error) {
	return loadTagVocab(root, "variable-tags.json")
}

func LoadMenuTags(root string) (map[string]struct{}, error) {
	return loadTagVocab(root, "menu-tags.json")
}

func loadTagVocab(root, filename string) (map[string]struct{}, error) {
	path := filepath.Join(root, "schema", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filename, err)
	}
	var doc tagListFile
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse %s: %w", filename, err)
	}
	out := make(map[string]struct{}, len(doc.Tags))
	for _, t := range doc.Tags {
		out[t] = struct{}{}
	}
	return out, nil
}
