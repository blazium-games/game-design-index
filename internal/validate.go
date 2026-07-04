package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func validateSlug(label, slug string) error {
	if slug == "" {
		return fmt.Errorf("%s: slug is required", label)
	}
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("%s: invalid slug %q", label, slug)
	}
	return nil
}

func validateFlavor(label string, f Flavor) error {
	switch f {
	case FlavorAction, FlavorAdventure, FlavorStrategy:
		return nil
	default:
		return fmt.Errorf("%s: invalid flavor %q", label, f)
	}
}

func validateQualityTier(label string, t QualityTier) error {
	switch t {
	case "", QualityCurated, QualityCatalog, QualityTemplate, QualityStub:
		return nil
	default:
		return fmt.Errorf("%s: invalid quality_tier %q", label, t)
	}
}

func validateMechanicDomain(label string, d MechanicDomain) error {
	switch d {
	case "", DomainLocomotion, DomainCombat, DomainProgression, DomainEconomy, DomainLevel, DomainSession:
		return nil
	default:
		return fmt.Errorf("%s: invalid domain %q", label, d)
	}
}

func validateMechanicPhase(label string, p MechanicPhase) error {
	switch p {
	case "", PhaseEarly, PhaseMid, PhaseLate, PhaseOptional:
		return nil
	default:
		return fmt.Errorf("%s: invalid phase %q", label, p)
	}
}

func validateImplementationComplexity(label string, c ImplementationComplexity) error {
	switch c {
	case "", ComplexitySmall, ComplexityMedium, ComplexityLarge:
		return nil
	default:
		return fmt.Errorf("%s: invalid implementation_complexity %q", label, c)
	}
}

func validateRequirementKind(label string, k RequirementKind) error {
	switch k {
	case RequirementMechanic, RequirementDomain, RequirementTag:
		return nil
	default:
		return fmt.Errorf("%s: invalid requirement kind %q", label, k)
	}
}

func validateMechanicEntry(entry MechanicEntry, tagVocab map[string]struct{}, mechanicSlugs map[string]struct{}) error {
	label := entry.Slug
	if label == "" {
		label = entry.Name
	}
	if entry.SchemaVersion != SchemaVersion {
		return fmt.Errorf("%s: schema_version must be %q", label, SchemaVersion)
	}
	if err := validateSlug(label, entry.Slug); err != nil {
		return err
	}
	if entry.Name == "" {
		return fmt.Errorf("%s: name is required", label)
	}
	if err := validateFlavor(label, entry.Flavor); err != nil {
		return err
	}
	if entry.Summary == "" {
		return fmt.Errorf("%s: summary is required", label)
	}
	if entry.Domain == "" {
		return fmt.Errorf("%s: domain is required", label)
	}
	if err := validateMechanicDomain(label, entry.Domain); err != nil {
		return err
	}
	if len(entry.Tags) == 0 {
		return fmt.Errorf("%s: at least one tag is required", label)
	}
	for _, t := range entry.Tags {
		if tagVocab != nil {
			if _, ok := tagVocab[t]; !ok {
				return fmt.Errorf("%s: unknown tag %q", label, t)
			}
		}
	}
	for _, req := range entry.Requirements {
		if err := validateRequirementKind(label+".requirements", req.Kind); err != nil {
			return err
		}
		switch req.Kind {
		case RequirementMechanic:
			if mechanicSlugs != nil {
				if _, ok := mechanicSlugs[req.Ref]; !ok {
					return fmt.Errorf("%s: requirement references unknown mechanic %q", label, req.Ref)
				}
			}
		case RequirementDomain:
			if err := validateMechanicDomain(label+".requirements", MechanicDomain(req.Ref)); err != nil {
				return err
			}
		case RequirementTag:
			if tagVocab != nil {
				if _, ok := tagVocab[req.Ref]; !ok {
					return fmt.Errorf("%s: requirement references unknown tag %q", label, req.Ref)
				}
			}
		}
	}
	for _, s := range entry.Synergies {
		if err := validateSlug(label+".synergies", s); err != nil {
			return err
		}
		if mechanicSlugs != nil {
			if _, ok := mechanicSlugs[s]; !ok {
				return fmt.Errorf("%s: synergy references unknown mechanic %q", label, s)
			}
		}
	}
	for _, sn := range entry.SynergyNotes {
		if err := validateSlug(label+".synergy_notes", sn.Slug); err != nil {
			return err
		}
		if mechanicSlugs != nil {
			if _, ok := mechanicSlugs[sn.Slug]; !ok {
				return fmt.Errorf("%s: synergy_notes references unknown mechanic %q", label, sn.Slug)
			}
		}
	}
	if entry.Complexity != "" {
		if err := validateImplementationComplexity(label+".complexity", entry.Complexity); err != nil {
			return err
		}
	}
	if entry.RelationshipModel != nil {
		switch entry.RelationshipModel.Type {
		case "directed_cycle", "graph", "":
		default:
			return fmt.Errorf("%s: invalid relationship_model.type %q", label, entry.RelationshipModel.Type)
		}
	}
	for _, s := range entry.Prerequisites {
		if err := validateSlug(label+".prerequisites", s); err != nil {
			return err
		}
	}
	for _, s := range entry.FeaturedIn {
		if err := validateSlug(label+".featured_in", s); err != nil {
			return err
		}
		if strings.Contains(s, "https") {
			return fmt.Errorf("%s: corrupted featured_in slug %q", label, s)
		}
	}
	if entry.Implementation != nil && entry.Implementation.ImplementationComplexity != nil {
		if err := validateImplementationComplexity(label, *entry.Implementation.ImplementationComplexity); err != nil {
			return err
		}
	}
	return nil
}

func validateMapType(label string, t MapType) error {
	switch t {
	case MapTypeGame, MapTypeGenre:
		return nil
	default:
		return fmt.Errorf("%s: invalid map_type %q", label, t)
	}
}

func validateRole(label string, r MechanicRole) error {
	switch r {
	case RoleSignature, RoleSupporting, RoleCommon:
		return nil
	default:
		return fmt.Errorf("%s: invalid role %q", label, r)
	}
}

func validateGameplayMap(m GameplayMap) error {
	label := m.Slug
	if label == "" {
		label = m.Title
	}
	if m.SchemaVersion != MapSchemaVersion && m.SchemaVersion != SchemaVersion {
		return fmt.Errorf("%s: schema_version must be %q or %q", label, MapSchemaVersion, SchemaVersion)
	}
	if err := validateSlug(label, m.Slug); err != nil {
		return err
	}
	if m.Title == "" {
		return fmt.Errorf("%s: title is required", label)
	}
	if err := validateMapType(label, m.MapType); err != nil {
		return err
	}
	if m.SourceURL == "" {
		return fmt.Errorf("%s: source_url is required", label)
	}
	if m.Subject.Name == "" {
		return fmt.Errorf("%s: subject.name is required", label)
	}
	if m.Subject.Kind != "game" && m.Subject.Kind != "genre" {
		return fmt.Errorf("%s: invalid subject.kind %q", label, m.Subject.Kind)
	}
	if m.Narrative.Description == "" {
		return fmt.Errorf("%s: narrative.description is required", label)
	}
	if len(m.Mechanics) == 0 {
		return fmt.Errorf("%s: at least one mechanic binding is required", label)
	}
	if len(m.Views) == 0 {
		return fmt.Errorf("%s: at least one view is required", label)
	}
	if m.Metadata != nil {
		if err := validateQualityTier(label+".metadata", m.Metadata.QualityTier); err != nil {
			return err
		}
	}
	for _, s := range m.SignatureGameplay {
		if err := validateSlug(label+".signature_gameplay", s); err != nil {
			return err
		}
	}
	for _, binding := range m.Mechanics {
		if err := validateSlug(label+".mechanics", binding.MechanicSlug); err != nil {
			return err
		}
		if err := validateRole(label+".mechanics", binding.Role); err != nil {
			return err
		}
		if binding.Domain == "" {
			return fmt.Errorf("%s: mechanic binding %q missing domain", label, binding.MechanicSlug)
		}
		if err := validateMechanicDomain(label+".mechanics", binding.Domain); err != nil {
			return err
		}
		if err := validateMechanicPhase(label+".mechanics", binding.Phase); err != nil {
			return err
		}
		if binding.Weight != 0 && (binding.Weight < 1 || binding.Weight > 5) {
			return fmt.Errorf("%s: mechanic weight must be 1-5, got %d", label, binding.Weight)
		}
		for _, dep := range binding.DependsOn {
			if err := validateSlug(label+".mechanics.depends_on", dep); err != nil {
				return err
			}
		}
	}
	for _, rel := range m.MechanicRelationships {
		if err := validateSlug(label+".mechanic_relationships.from", rel.FromMechanic); err != nil {
			return err
		}
		if err := validateSlug(label+".mechanic_relationships.to", rel.ToMechanic); err != nil {
			return err
		}
	}
	for _, view := range m.Views {
		if err := validateSlug(label+".views", view.ID); err != nil {
			return err
		}
		if view.Label == "" {
			return fmt.Errorf("%s: view %q label is required", label, view.ID)
		}
		if view.Filter.Role != nil {
			if err := validateRole(label+".views", *view.Filter.Role); err != nil {
				return err
			}
		}
		if view.Filter.Flavor != nil {
			if err := validateFlavor(label+".views", *view.Filter.Flavor); err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateMechanic(_ string, entry MechanicEntry) error {
	return validateMechanicEntry(entry, nil, nil)
}

func ValidateMap(_ string, m GameplayMap) error {
	return validateGameplayMap(m)
}

func ValidateCatalog(root string, catalog MechanicsCatalog) error {
	if catalog.SchemaVersion != "" && catalog.SchemaVersion != SchemaVersion {
		return fmt.Errorf("catalog schema_version must be %q", SchemaVersion)
	}
	tagVocab, err := LoadMechanicTags(root)
	if err != nil {
		return err
	}
	slugs := make(map[string]struct{}, len(catalog.Mechanics))
	for _, entry := range catalog.Mechanics {
		slugs[entry.Slug] = struct{}{}
	}
	for _, entry := range catalog.Mechanics {
		if err := validateMechanicEntry(entry, tagVocab, slugs); err != nil {
			return err
		}
	}
	return nil
}

func crossValidate(_ string, b *Bundle) error {
	for slug, m := range b.Maps {
		sigSet := make(map[string]struct{}, len(m.SignatureGameplay))
		for _, s := range m.SignatureGameplay {
			sigSet[s] = struct{}{}
		}

		boundSlugs := make(map[string]struct{}, len(m.Mechanics))
		for _, bnd := range m.Mechanics {
			boundSlugs[bnd.MechanicSlug] = struct{}{}
		}

		for _, binding := range m.Mechanics {
			if _, ok := b.Mechanics[binding.MechanicSlug]; !ok {
				return fmt.Errorf("map %q references unknown mechanic %q", slug, binding.MechanicSlug)
			}
			if lib, ok := b.Mechanics[binding.MechanicSlug]; ok && binding.Domain != "" && lib.Domain != "" {
				if binding.Domain != lib.Domain {
					return fmt.Errorf("map %q: binding %q domain %q != library domain %q", slug, binding.MechanicSlug, binding.Domain, lib.Domain)
				}
			}
			for _, dep := range binding.DependsOn {
				if _, ok := boundSlugs[dep]; !ok {
					return fmt.Errorf("map %q: binding %q depends_on %q not on same map", slug, binding.MechanicSlug, dep)
				}
			}
			if binding.Role == RoleSignature {
				if _, ok := sigSet[binding.MechanicSlug]; !ok {
					return fmt.Errorf("map %q: signature binding %q missing from signature_gameplay", slug, binding.MechanicSlug)
				}
			}
		}

		for _, s := range m.SignatureGameplay {
			found := false
			for _, binding := range m.Mechanics {
				if binding.MechanicSlug == s {
					found = true
					if binding.Role != RoleSignature {
						return fmt.Errorf("map %q: signature_gameplay entry %q must have role signature", slug, s)
					}
					break
				}
			}
			if !found {
				return fmt.Errorf("map %q: signature_gameplay entry %q has no mechanics binding", slug, s)
			}
		}

		if m.Relationships != nil {
			for _, ref := range m.Relationships.SimilarTo {
				if err := validateSlug(slug+".relationships.similar_to", ref); err != nil {
					return err
				}
				if _, ok := b.Maps[ref]; !ok {
					return fmt.Errorf("map %q: similar_to references unknown map %q", slug, ref)
				}
			}
			for _, ref := range m.Relationships.GenreMaps {
				if err := validateSlug(slug+".relationships.genre_maps", ref); err != nil {
					return err
				}
				if gm, ok := b.Maps[ref]; !ok {
					return fmt.Errorf("map %q: genre_maps references unknown map %q", slug, ref)
				} else if gm.MapType != MapTypeGenre {
					return fmt.Errorf("map %q: genre_maps entry %q is not a genre map", slug, ref)
				}
			}
			for _, ref := range m.Relationships.InfluenceSlugs {
				if err := validateSlug(slug+".relationships.influence_slugs", ref); err != nil {
					return err
				}
				if _, ok := b.Maps[ref]; !ok {
					return fmt.Errorf("map %q: influence_slugs references unknown map %q", slug, ref)
				}
			}
		}

		if err := validateGenreMechanicSanity(slug, m); err != nil {
			return err
		}
	}
	return nil
}

func validateGenreMechanicSanity(slug string, m GameplayMap) error {
	if !ShouldStripTemplatePack(m) {
		return nil
	}
	for _, binding := range m.Mechanics {
		if binding.Role != RoleSupporting {
			continue
		}
		if _, bad := ActionAdventureTemplatePack[binding.MechanicSlug]; bad {
			return fmt.Errorf("map %q: inappropriate supporting mechanic %q for non-combat signatures", slug, binding.MechanicSlug)
		}
	}
	return nil
}

// CollectTemplateWarnings returns non-fatal quality warnings for template-tier maps.
func CollectTemplateWarnings(b *Bundle) []string {
	var warnings []string
	for slug, m := range b.Maps {
		if DetectQualityTier(m) == QualityTemplate {
			warnings = append(warnings, fmt.Sprintf("map %q is template-tier (batch boilerplate)", slug))
		}
	}
	return warnings
}

func ValidateBundle(root string, b *Bundle) error {
	catalog := MechanicsCatalog{
		SchemaVersion: SchemaVersion,
		Mechanics:     make([]MechanicEntry, 0, len(b.Mechanics)),
	}
	for _, m := range b.Mechanics {
		catalog.Mechanics = append(catalog.Mechanics, m)
	}
	if err := ValidateCatalog(root, catalog); err != nil {
		return err
	}
	for _, m := range b.Maps {
		if err := ValidateMap(root, m); err != nil {
			return err
		}
	}
	return crossValidate(root, b)
}

// ValidateSchemaFiles ensures schema JSON files are present and parseable.
func ValidateSchemaFiles(root string) error {
	for _, name := range []string{"mechanic-entry.schema.json", "gameplay-map.schema.json", "mechanic-tags.json"} {
		path := filepath.Join(root, "schema", name)
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read schema %s: %w", name, err)
		}
		var doc map[string]any
		if err := json.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("parse schema %s: %w", name, err)
		}
	}
	return nil
}
