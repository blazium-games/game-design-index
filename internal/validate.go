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
	if entry.SchemaVersion != SchemaVersion && entry.SchemaVersion != "1.3" {
		return fmt.Errorf("%s: schema_version must be %q or %q", label, SchemaVersion, "1.3")
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
	if m.SchemaVersion != MapSchemaVersion && m.SchemaVersion != SchemaVersion && m.SchemaVersion != "1.2" {
		return fmt.Errorf("%s: schema_version must be %q, %q, or %q", label, MapSchemaVersion, SchemaVersion, "1.2")
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
	if err := validateMapVariableBindings(label, m.Variables); err != nil {
		return err
	}
	if err := validateMapUIMenuBindings(label, m.UIMenus); err != nil {
		return err
	}
	for _, rel := range m.VariableRelationships {
		if err := validateSlug(label+".variable_relationships.from", rel.FromVariable); err != nil {
			return err
		}
		if err := validateSlug(label+".variable_relationships.to", rel.ToVariable); err != nil {
			return err
		}
	}
	for _, edge := range m.MenuFlow {
		if err := validateSlug(label+".menu_flow.from", edge.FromMenu); err != nil {
			return err
		}
		if err := validateSlug(label+".menu_flow.to", edge.ToMenu); err != nil {
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

		boundVars := make(map[string]struct{}, len(m.Variables))
		for _, vb := range m.Variables {
			if _, ok := b.Variables[vb.VariableSlug]; !ok {
				return fmt.Errorf("map %q references unknown variable %q", slug, vb.VariableSlug)
			}
			boundVars[vb.VariableSlug] = struct{}{}
			for _, mech := range vb.RelatedMechanics {
				if _, ok := b.Mechanics[mech]; !ok {
					return fmt.Errorf("map %q: variable %q related_mechanics references unknown mechanic %q", slug, vb.VariableSlug, mech)
				}
			}
		}

		boundMenus := make(map[string]struct{}, len(m.UIMenus))
		for _, mb := range m.UIMenus {
			if _, ok := b.UIMenus[mb.MenuSlug]; !ok {
				return fmt.Errorf("map %q references unknown menu %q", slug, mb.MenuSlug)
			}
			boundMenus[mb.MenuSlug] = struct{}{}
			for _, from := range mb.OpensFrom {
				if _, ok := boundMenus[from]; !ok {
					if _, libOk := b.UIMenus[from]; !libOk {
						return fmt.Errorf("map %q: menu %q opens_from unknown menu %q", slug, mb.MenuSlug, from)
					}
				}
			}
			for _, v := range mb.DisplaysVariables {
				if _, ok := b.Variables[v]; !ok {
					return fmt.Errorf("map %q: menu %q displays unknown variable %q", slug, mb.MenuSlug, v)
				}
			}
			for _, mech := range mb.SupportsMechanics {
				if _, ok := b.Mechanics[mech]; !ok {
					return fmt.Errorf("map %q: menu %q supports unknown mechanic %q", slug, mb.MenuSlug, mech)
				}
			}
		}

		for _, rel := range m.VariableRelationships {
			if _, ok := boundVars[rel.FromVariable]; !ok {
				if _, libOk := b.Variables[rel.FromVariable]; !libOk {
					return fmt.Errorf("map %q: variable_relationship from unknown variable %q", slug, rel.FromVariable)
				}
			}
			if _, ok := boundVars[rel.ToVariable]; !ok {
				if _, libOk := b.Variables[rel.ToVariable]; !libOk {
					return fmt.Errorf("map %q: variable_relationship to unknown variable %q", slug, rel.ToVariable)
				}
			}
		}

		for _, edge := range m.MenuFlow {
			if _, ok := boundMenus[edge.FromMenu]; !ok {
				if _, libOk := b.UIMenus[edge.FromMenu]; !libOk {
					return fmt.Errorf("map %q: menu_flow from unknown menu %q", slug, edge.FromMenu)
				}
			}
			if _, ok := boundMenus[edge.ToMenu]; !ok {
				if _, libOk := b.UIMenus[edge.ToMenu]; !libOk {
					return fmt.Errorf("map %q: menu_flow to unknown menu %q", slug, edge.ToMenu)
				}
			}
		}

		for _, sk := range m.SkillSlugs {
			if err := validateSlug(slug+".skill_slugs", sk); err != nil {
				return err
			}
			if len(b.Skills) > 0 {
				if _, ok := b.Skills[sk]; !ok {
					return fmt.Errorf("map %q: skill_slugs references unknown skill %q", slug, sk)
				}
			}
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
	mechanicSlugs := make(map[string]struct{}, len(b.Mechanics))
	for s := range b.Mechanics {
		mechanicSlugs[s] = struct{}{}
	}
	variableSlugs := make(map[string]struct{}, len(b.Variables))
	for s := range b.Variables {
		variableSlugs[s] = struct{}{}
	}

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

	if len(b.Variables) > 0 {
		varCatalog := VariablesCatalog{
			SchemaVersion: VariableSchemaVersion,
			Variables:     make([]GameVariable, 0, len(b.Variables)),
		}
		for _, v := range b.Variables {
			varCatalog.Variables = append(varCatalog.Variables, v)
		}
		if err := validateVariablesCatalog(root, varCatalog, mechanicSlugs); err != nil {
			return err
		}
	}

	if len(b.UIMenus) > 0 {
		menuCatalog := UIMenusCatalog{
			SchemaVersion: UIMenuSchemaVersion,
			Menus:         make([]UIMenu, 0, len(b.UIMenus)),
		}
		for _, m := range b.UIMenus {
			menuCatalog.Menus = append(menuCatalog.Menus, m)
		}
		if err := validateUIMenusCatalog(root, menuCatalog, mechanicSlugs, variableSlugs); err != nil {
			return err
		}
	}

	if len(b.Skills) > 0 {
		skillsCatalog := SkillsCatalog{
			SchemaVersion: SkillSchemaVersion,
			Skills:        make([]DesignSkill, 0, len(b.Skills)),
		}
		for _, s := range b.Skills {
			skillsCatalog.Skills = append(skillsCatalog.Skills, s)
		}
		if err := validateSkillsCatalog(root, skillsCatalog, mechanicSlugs, variableSlugs); err != nil {
			return err
		}
	}

	skillSlugs := make(map[string]struct{}, len(b.Skills))
	for s := range b.Skills {
		skillSlugs[s] = struct{}{}
	}
	for slug, entry := range b.Mechanics {
		for _, sk := range entry.SkillsDeveloped {
			if len(skillSlugs) == 0 {
				continue
			}
			if _, ok := skillSlugs[sk]; !ok {
				return fmt.Errorf("%s: skills_developed references unknown skill %q", slug, sk)
			}
		}
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
	for _, name := range []string{
		"mechanic-entry.schema.json",
		"gameplay-map.schema.json",
		"mechanic-tags.json",
		"game-variable.schema.json",
		"ui-menu.schema.json",
		"variable-tags.json",
		"menu-tags.json",
		"skill.schema.json",
		"skill-tags.json",
	} {
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
