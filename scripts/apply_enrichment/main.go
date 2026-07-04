package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	root := filepath.Join("data", "source")
	cohortFlag := flag.String("cohort", "", "Apply a single cohort file from enrichment/cohorts/ (default: all)")
	flag.Parse()
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	cohortDir := filepath.Join(root, "enrichment", "cohorts")
	var files []string
	if *cohortFlag != "" {
		files = []string{filepath.Join(cohortDir, *cohortFlag)}
	} else {
		entries, err := os.ReadDir(cohortDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read cohorts: %v\n", err)
			os.Exit(1)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}
			files = append(files, filepath.Join(cohortDir, e.Name()))
		}
		sort.Strings(files)
	}
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "no cohort files found")
		os.Exit(1)
	}

	for _, path := range files {
		if err := applyCohort(root, path); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(path), err)
			os.Exit(1)
		}
		fmt.Printf("Applied %s\n", filepath.Base(path))
	}
}

func applyCohort(root, patchPath string) error {
	data, err := os.ReadFile(patchPath)
	if err != nil {
		return fmt.Errorf("read patch: %w", err)
	}
	var patch struct {
		Mechanics  map[string]json.RawMessage `json:"mechanics"`
		Variables  map[string]json.RawMessage `json:"variables"`
		UIMenus    map[string]json.RawMessage `json:"ui_menus"`
		Skills     map[string]json.RawMessage `json:"skills"`
		SkillsAdd  []json.RawMessage          `json:"skills_add"`
		MapPatches map[string]json.RawMessage `json:"maps"`
		MapsAdd    map[string]json.RawMessage `json:"maps_add"`
	}
	if err := json.Unmarshal(data, &patch); err != nil {
		return fmt.Errorf("parse patch: %w", err)
	}
	if err := mergeCatalog(filepath.Join(root, "library", "mechanics.json"), "mechanics", patch.Mechanics); err != nil {
		return err
	}
	if err := mergeCatalog(filepath.Join(root, "library", "variables.json"), "variables", patch.Variables); err != nil {
		return err
	}
	if err := mergeCatalog(filepath.Join(root, "library", "ui-menus.json"), "menus", patch.UIMenus); err != nil {
		return err
	}
	if err := mergeCatalog(filepath.Join(root, "library", "skills.json"), "skills", patch.Skills); err != nil {
		return err
	}
	if err := appendCatalog(filepath.Join(root, "library", "skills.json"), "skills", patch.SkillsAdd); err != nil {
		return err
	}
	for slug, raw := range patch.MapPatches {
		path := filepath.Join(root, "maps", slug+".json")
		if err := mergeMap(path, raw); err != nil {
			return err
		}
	}
	for slug, raw := range patch.MapsAdd {
		path := filepath.Join(root, "maps", slug+".json")
		if err := writeMapIfAbsent(path, raw); err != nil {
			return err
		}
	}
	return nil
}

func mergeCatalog(path, arrayKey string, patches map[string]json.RawMessage) error {
	if len(patches) == 0 {
		return nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	var schemaVer any
	_ = json.Unmarshal(doc["schema_version"], &schemaVer)
	var items []json.RawMessage
	if err := json.Unmarshal(doc[arrayKey], &items); err != nil {
		return err
	}
	for i, item := range items {
		var entry map[string]any
		if err := json.Unmarshal(item, &entry); err != nil {
			return err
		}
		slug, _ := entry["slug"].(string)
		patch, ok := patches[slug]
		if !ok {
			continue
		}
		var patchObj map[string]any
		if err := json.Unmarshal(patch, &patchObj); err != nil {
			return err
		}
		for k, v := range patchObj {
			entry[k] = v
		}
		merged, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		items[i] = merged
	}
	out, err := json.MarshalIndent(map[string]any{
		"schema_version": schemaVer,
		arrayKey:         items,
	}, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}

func appendCatalog(path, arrayKey string, additions []json.RawMessage) error {
	if len(additions) == 0 {
		return nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	var schemaVer any
	_ = json.Unmarshal(doc["schema_version"], &schemaVer)
	var items []json.RawMessage
	if err := json.Unmarshal(doc[arrayKey], &items); err != nil {
		return err
	}
	existing := make(map[string]struct{}, len(items))
	for _, item := range items {
		var entry map[string]any
		if err := json.Unmarshal(item, &entry); err != nil {
			return err
		}
		if slug, _ := entry["slug"].(string); slug != "" {
			existing[slug] = struct{}{}
		}
	}
	for _, add := range additions {
		var entry map[string]any
		if err := json.Unmarshal(add, &entry); err != nil {
			return err
		}
		slug, _ := entry["slug"].(string)
		if slug == "" {
			return fmt.Errorf("append %s: entry missing slug", arrayKey)
		}
		if _, ok := existing[slug]; ok {
			continue
		}
		items = append(items, add)
		existing[slug] = struct{}{}
	}
	out, err := json.MarshalIndent(map[string]any{
		"schema_version": schemaVer,
		arrayKey:         items,
	}, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}

func writeMapIfAbsent(path string, doc json.RawMessage) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var pretty any
	if err := json.Unmarshal(doc, &pretty); err != nil {
		return err
	}
	out, err := json.MarshalIndent(pretty, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}

func mergeMap(path string, patch json.RawMessage) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	var patchObj map[string]any
	if err := json.Unmarshal(patch, &patchObj); err != nil {
		return err
	}
	for k, v := range patchObj {
		doc[k] = v
	}
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}
