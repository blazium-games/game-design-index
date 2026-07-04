package main

import (
	"encoding/json"
	"strings"

	"github.com/blazium-games/game-mechanics-index/internal"
)

var bannedSubstrings = []string{
	"godot",
	"ludonaut",
	"ludens",
	"retrododo",
	"source_url",
}

func sanitizeMechanic(e internal.MechanicEntry) map[string]any {
	raw, _ := json.Marshal(e)
	var out map[string]any
	_ = json.Unmarshal(raw, &out)
	delete(out, "implementation")
	delete(out, "media")
	return out
}

func sanitizeMap(m internal.GameplayMap) map[string]any {
	raw, _ := json.Marshal(m)
	var out map[string]any
	_ = json.Unmarshal(raw, &out)
	delete(out, "source_url")
	return out
}

func scanBannedJSON(data []byte, path string) error {
	lower := strings.ToLower(string(data))
	for _, term := range bannedSubstrings {
		if strings.Contains(lower, term) {
			return &bannedTermError{path: path, term: term}
		}
	}
	return nil
}

type bannedTermError struct {
	path string
	term string
}

func (e *bannedTermError) Error() string {
	return "banned term " + e.term + " in " + e.path
}

func sanitizeBundle(b *internal.Bundle) (maps map[string]map[string]any, mechanics map[string]map[string]any) {
	maps = make(map[string]map[string]any, len(b.Maps))
	for slug, m := range b.Maps {
		maps[slug] = sanitizeMap(m)
	}
	mechanics = make(map[string]map[string]any, len(b.Mechanics))
	for slug, e := range b.Mechanics {
		mechanics[slug] = sanitizeMechanic(e)
	}
	return maps, mechanics
}
