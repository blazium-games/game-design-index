package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadMechanicTags reads the controlled tag vocabulary from schema/mechanic-tags.json.
func LoadMechanicTags(root string) (map[string]struct{}, error) {
	path := filepath.Join(root, "schema", "mechanic-tags.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read mechanic-tags.json: %w", err)
	}
	var doc MechanicTagsCatalog
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse mechanic-tags.json: %w", err)
	}
	out := make(map[string]struct{}, len(doc.Tags))
	for _, t := range doc.Tags {
		out[t] = struct{}{}
	}
	return out, nil
}

// ValidateMechanicsCorpus returns non-fatal quality warnings for the mechanics library.
func ValidateMechanicsCorpus(b *Bundle, tagVocab map[string]struct{}) []string {
	var warnings []string
	domainSet := map[MechanicDomain]struct{}{
		DomainLocomotion: {}, DomainCombat: {}, DomainProgression: {},
		DomainEconomy: {}, DomainLevel: {}, DomainSession: {},
	}
	for slug, entry := range b.Mechanics {
		if entry.Domain == "" {
			warnings = append(warnings, fmt.Sprintf("mechanic %q missing domain", slug))
		}
		if len(entry.Tags) == 0 {
			warnings = append(warnings, fmt.Sprintf("mechanic %q missing tags", slug))
		}
		for _, t := range entry.Tags {
			if _, ok := tagVocab[t]; !ok {
				warnings = append(warnings, fmt.Sprintf("mechanic %q has unknown tag %q", slug, t))
			}
		}
		for _, req := range entry.Requirements {
			switch req.Kind {
			case RequirementMechanic:
				if _, ok := b.Mechanics[req.Ref]; !ok {
					warnings = append(warnings, fmt.Sprintf("mechanic %q requirement references unknown mechanic %q", slug, req.Ref))
				}
			case RequirementDomain:
				if _, ok := domainSet[MechanicDomain(req.Ref)]; !ok {
					warnings = append(warnings, fmt.Sprintf("mechanic %q requirement references unknown domain %q", slug, req.Ref))
				}
			case RequirementTag:
				if _, ok := tagVocab[req.Ref]; !ok {
					warnings = append(warnings, fmt.Sprintf("mechanic %q requirement references unknown tag %q", slug, req.Ref))
				}
			}
		}
		for _, syn := range entry.Synergies {
			if _, ok := b.Mechanics[syn]; !ok {
				warnings = append(warnings, fmt.Sprintf("mechanic %q synergy references unknown mechanic %q", slug, syn))
			}
		}
		conflictNames := make(map[string]struct{}, len(entry.Conflicts))
		for _, c := range entry.Conflicts {
			conflictNames[c] = struct{}{}
		}
		for _, syn := range entry.Synergies {
			if other, ok := b.Mechanics[syn]; ok {
				for _, c := range other.Conflicts {
					if _, hit := conflictNames[c]; hit {
						warnings = append(warnings, fmt.Sprintf("mechanic %q synergizes with %q but lists conflict %q", slug, syn, c))
					}
				}
			}
		}
	}
	return warnings
}
