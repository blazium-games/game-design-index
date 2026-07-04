package internal

import (
	"encoding/json"
	"fmt"
	"sort"
)

type DesignBriefMechanic struct {
	Slug      string       `json:"slug"`
	Name      string       `json:"name"`
	Flavor    Flavor       `json:"flavor"`
	Domain    MechanicDomain `json:"domain,omitempty"`
	Summary   string       `json:"summary"`
	Conflicts []string     `json:"conflicts,omitempty"`
	Synergies []string     `json:"synergies,omitempty"`
	Role      MechanicRole `json:"role"`
	MapNotes  []string     `json:"map_notes_from_refs"`
}

type DesignBriefReference struct {
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	Genres      []string    `json:"genres,omitempty"`
	CoreLoop    string      `json:"core_loop,omitempty"`
	Signatures  []string    `json:"signatures"`
	QualityTier QualityTier `json:"quality_tier"`
}

type DesignBrief struct {
	SchemaVersion string                 `json:"schema_version"`
	References    []DesignBriefReference `json:"references"`
	Signatures    []DesignBriefMechanic  `json:"signatures"`
	Supporting    []DesignBriefMechanic  `json:"supporting"`
	SynergyPairs  []MechanicCooccurrence `json:"synergy_pairs,omitempty"`
}

// BuildDesignBrief merges signature and supporting mechanics from reference maps into a GDD seed.
func (b *Bundle) BuildDesignBrief(refSlugs []string, minCooccurrence int) (*DesignBrief, error) {
	if len(refSlugs) == 0 {
		return nil, fmt.Errorf("at least one reference map slug is required")
	}

	brief := &DesignBrief{
		SchemaVersion: SchemaVersion,
		SynergyPairs:  b.MechanicCooccurrenceMatrix(minCooccurrence),
	}

	sigSeen := map[string]*DesignBriefMechanic{}
	supportSeen := map[string]*DesignBriefMechanic{}

	for _, slug := range refSlugs {
		m, ok := b.Maps[slug]
		if !ok {
			return nil, fmt.Errorf("unknown map slug %q", slug)
		}
		brief.References = append(brief.References, DesignBriefReference{
			Slug:        m.Slug,
			Title:       m.Title,
			Genres:      append([]string(nil), m.Subject.Genres...),
			CoreLoop:    m.Narrative.CoreLoop,
			Signatures:  append([]string(nil), m.SignatureGameplay...),
			QualityTier: DetectQualityTier(m),
		})

		for _, binding := range m.Mechanics {
			entry, ok := b.Mechanics[binding.MechanicSlug]
			if !ok {
				continue
			}
			target := supportSeen
			if binding.Role == RoleSignature {
				target = sigSeen
			}
			dm, exists := target[binding.MechanicSlug]
			if !exists {
				dm = &DesignBriefMechanic{
					Slug:      entry.Slug,
					Name:      entry.Name,
					Flavor:    entry.Flavor,
					Domain:    entry.Domain,
					Summary:   entry.Summary,
					Conflicts: append([]string(nil), entry.Conflicts...),
					Synergies: append([]string(nil), entry.Synergies...),
					Role:      binding.Role,
				}
				target[binding.MechanicSlug] = dm
			}
			if binding.MapNotes != "" {
				dm.MapNotes = appendUniqueString(dm.MapNotes, binding.MapNotes)
			}
		}
	}

	brief.Signatures = mapToSortedMechanics(sigSeen)
	brief.Supporting = mapToSortedMechanics(supportSeen)
	return brief, nil
}

func mapToSortedMechanics(m map[string]*DesignBriefMechanic) []DesignBriefMechanic {
	out := make([]DesignBriefMechanic, 0, len(m))
	for _, v := range m {
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Slug < out[j].Slug })
	return out
}

func appendUniqueString(list []string, s string) []string {
	for _, existing := range list {
		if existing == s {
			return list
		}
	}
	return append(list, s)
}

// DesignBriefJSON marshals a design brief with indentation.
func DesignBriefJSON(brief *DesignBrief) ([]byte, error) {
	return json.MarshalIndent(brief, "", "  ")
}
