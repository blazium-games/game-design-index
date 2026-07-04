package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func LoadDir(root string) (*Bundle, error) {
	b := &Bundle{
		Root:      root,
		Mechanics: make(map[string]MechanicEntry),
		Variables: make(map[string]GameVariable),
		UIMenus:   make(map[string]UIMenu),
		Skills:    make(map[string]DesignSkill),
		Maps:      make(map[string]GameplayMap),
	}

	catalogPath := filepath.Join(root, "library", "mechanics.json")
	catalogData, err := os.ReadFile(catalogPath)
	if err != nil {
		return nil, fmt.Errorf("read mechanics catalog: %w", err)
	}
	var catalog MechanicsCatalog
	if err := json.Unmarshal(catalogData, &catalog); err != nil {
		return nil, fmt.Errorf("parse mechanics catalog: %w", err)
	}
	for _, entry := range catalog.Mechanics {
		if entry.Slug == "" {
			return nil, fmt.Errorf("mechanic entry missing slug")
		}
		if _, exists := b.Mechanics[entry.Slug]; exists {
			return nil, fmt.Errorf("duplicate mechanic slug %q", entry.Slug)
		}
		b.Mechanics[entry.Slug] = entry
	}

	varPath := filepath.Join(root, "library", "variables.json")
	if varData, err := os.ReadFile(varPath); err == nil {
		var varCatalog VariablesCatalog
		if err := json.Unmarshal(varData, &varCatalog); err != nil {
			return nil, fmt.Errorf("parse variables catalog: %w", err)
		}
		for _, entry := range varCatalog.Variables {
			if entry.Slug == "" {
				return nil, fmt.Errorf("variable entry missing slug")
			}
			if _, exists := b.Variables[entry.Slug]; exists {
				return nil, fmt.Errorf("duplicate variable slug %q", entry.Slug)
			}
			b.Variables[entry.Slug] = entry
		}
	}

	menuPath := filepath.Join(root, "library", "ui-menus.json")
	if menuData, err := os.ReadFile(menuPath); err == nil {
		var menuCatalog UIMenusCatalog
		if err := json.Unmarshal(menuData, &menuCatalog); err != nil {
			return nil, fmt.Errorf("parse ui-menus catalog: %w", err)
		}
		for _, entry := range menuCatalog.Menus {
			if entry.Slug == "" {
				return nil, fmt.Errorf("ui menu entry missing slug")
			}
			if _, exists := b.UIMenus[entry.Slug]; exists {
				return nil, fmt.Errorf("duplicate ui menu slug %q", entry.Slug)
			}
			b.UIMenus[entry.Slug] = entry
		}
	}

	skillsPath := filepath.Join(root, "library", "skills.json")
	if skillsData, err := os.ReadFile(skillsPath); err == nil {
		var skillsCatalog SkillsCatalog
		if err := json.Unmarshal(skillsData, &skillsCatalog); err != nil {
			return nil, fmt.Errorf("parse skills catalog: %w", err)
		}
		for _, entry := range skillsCatalog.Skills {
			if entry.Slug == "" {
				return nil, fmt.Errorf("skill entry missing slug")
			}
			if _, exists := b.Skills[entry.Slug]; exists {
				return nil, fmt.Errorf("duplicate skill slug %q", entry.Slug)
			}
			b.Skills[entry.Slug] = entry
		}
	}

	mapsDir := filepath.Join(root, "maps")
	entries, err := os.ReadDir(mapsDir)
	if err != nil {
		return nil, fmt.Errorf("read maps directory: %w", err)
	}
	for _, ent := range entries {
		if ent.IsDir() || filepath.Ext(ent.Name()) != ".json" {
			continue
		}
		if strings.HasPrefix(ent.Name(), "_") {
			continue
		}
		path := filepath.Join(mapsDir, ent.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read map %s: %w", ent.Name(), err)
		}
		var m GameplayMap
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("parse map %s: %w", ent.Name(), err)
		}
		if m.Slug == "" {
			return nil, fmt.Errorf("map %s missing slug", ent.Name())
		}
		if _, exists := b.Maps[m.Slug]; exists {
			return nil, fmt.Errorf("duplicate map slug %q", m.Slug)
		}
		b.Maps[m.Slug] = m
	}

	if err := ValidateBundle(root, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Bundle) Mechanic(slug string) (MechanicEntry, bool) {
	m, ok := b.Mechanics[slug]
	return m, ok
}

func (b *Bundle) Map(slug string) (GameplayMap, bool) {
	m, ok := b.Maps[slug]
	return m, ok
}

func (b *Bundle) MapMechanics(mapSlug string) []MapMechanicBinding {
	m, ok := b.Maps[mapSlug]
	if !ok {
		return nil
	}
	out := make([]MapMechanicBinding, len(m.Mechanics))
	copy(out, m.Mechanics)
	return out
}

func (b *Bundle) SignatureMechanics(mapSlug string) []MechanicEntry {
	m, ok := b.Maps[mapSlug]
	if !ok {
		return nil
	}
	var out []MechanicEntry
	for _, binding := range m.Mechanics {
		if binding.Role != RoleSignature {
			continue
		}
		if entry, ok := b.Mechanics[binding.MechanicSlug]; ok {
			out = append(out, entry)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Slug < out[j].Slug
	})
	return out
}

func (b *Bundle) MechanicsByFlavor(mapSlug string, flavor Flavor) []MechanicEntry {
	m, ok := b.Maps[mapSlug]
	if !ok {
		return nil
	}
	var out []MechanicEntry
	for _, binding := range m.Mechanics {
		entry, ok := b.Mechanics[binding.MechanicSlug]
		if !ok || entry.Flavor != flavor {
			continue
		}
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Slug < out[j].Slug
	})
	return out
}

func (b *Bundle) MapSlugs() []string {
	slugs := make([]string, 0, len(b.Maps))
	for s := range b.Maps {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)
	return slugs
}

func (b *Bundle) MechanicSlugs() []string {
	slugs := make([]string, 0, len(b.Mechanics))
	for s := range b.Mechanics {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)
	return slugs
}

func (b *Bundle) Variable(slug string) (GameVariable, bool) {
	v, ok := b.Variables[slug]
	return v, ok
}

func (b *Bundle) UIMenu(slug string) (UIMenu, bool) {
	m, ok := b.UIMenus[slug]
	return m, ok
}

func (b *Bundle) VariableSlugs() []string {
	slugs := make([]string, 0, len(b.Variables))
	for s := range b.Variables {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)
	return slugs
}

func (b *Bundle) UIMenuSlugs() []string {
	slugs := make([]string, 0, len(b.UIMenus))
	for s := range b.UIMenus {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)
	return slugs
}

// VariableEnrichmentStatus returns "complete" or "needs_info".
func VariableEnrichmentStatus(v GameVariable) string {
	if v.SharedRationale != "" && v.PlayerFocus != "" {
		return "complete"
	}
	return "needs_info"
}

// MenuEnrichmentStatus returns "complete" or "needs_info".
func MenuEnrichmentStatus(m UIMenu) string {
	if m.SharedRationale != "" {
		return "complete"
	}
	return "needs_info"
}

func (b *Bundle) Skill(slug string) (DesignSkill, bool) {
	s, ok := b.Skills[slug]
	return s, ok
}

func (b *Bundle) SkillSlugs() []string {
	slugs := make([]string, 0, len(b.Skills))
	for s := range b.Skills {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)
	return slugs
}

// SkillEnrichmentStatus returns "complete" or "needs_info".
func SkillEnrichmentStatus(s DesignSkill) string {
	if s.LearningOutcome != "" && len(s.RelatedMechanics) > 0 {
		return "complete"
	}
	return "needs_info"
}

// MechanicEnrichmentStatus returns "complete" or "needs_info".
func MechanicEnrichmentStatus(m MechanicEntry) string {
	if m.DesignGuidance != nil && m.DesignGuidance.WhenToUse != "" && m.AgentContext != nil && m.AgentContext.SummaryForAgents != "" {
		return "complete"
	}
	return "needs_info"
}
