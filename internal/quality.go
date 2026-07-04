package internal

import (
	"regexp"
)

var (
	templateCoreLoopRe   = regexp.MustCompile(`(?i)explore interconnected areas, gain tools`)
	templateDescRe       = regexp.MustCompile(`(?i)acclaimed action-adventure title`)
	templateSigNotesRe   = regexp.MustCompile(`(?i)^Core .+ loop defines .+'s identity\.$`)
	templateSupportNotes = regexp.MustCompile(`(?i) supports the main loop in `)
	upliftNotesRe        = regexp.MustCompile(`(?i)^Uplifted from .+ recipe`)
)

func isUpliftedMap(m GameplayMap) bool {
	for _, b := range m.Mechanics {
		if upliftNotesRe.MatchString(b.MapNotes) {
			return true
		}
	}
	return false
}

// IsTemplateMap returns true when narrative or map_notes match known batch boilerplate.
func IsTemplateMap(m GameplayMap) bool {
	if isUpliftedMap(m) {
		return true
	}
	if m.Narrative.CoreLoop != "" && templateCoreLoopRe.MatchString(m.Narrative.CoreLoop) {
		return true
	}
	if m.Narrative.Description != "" && templateDescRe.MatchString(m.Narrative.Description) {
		return true
	}
	for _, b := range m.Mechanics {
		if b.MapNotes == "" {
			continue
		}
		if templateSigNotesRe.MatchString(b.MapNotes) {
			return true
		}
		if templateSupportNotes.MatchString(b.MapNotes) {
			return true
		}
	}
	return false
}

// InferQualityTier classifies a map from content only, ignoring stored metadata.
func InferQualityTier(m GameplayMap) QualityTier {
	if m.MapType == MapTypeGenre {
		return QualityCurated
	}
	if IsTemplateMap(m) {
		return QualityTemplate
	}
	if hasCuratedSignals(m) {
		return QualityCurated
	}
	return QualityCatalog
}

// DetectQualityTier classifies a map for programmatic trust filtering.
func DetectQualityTier(m GameplayMap) QualityTier {
	if m.Metadata != nil && m.Metadata.QualityTier != "" {
		return m.Metadata.QualityTier
	}
	if m.MapType == MapTypeGenre {
		return QualityCurated
	}
	if IsTemplateMap(m) {
		return QualityTemplate
	}
	if hasCuratedSignals(m) {
		return QualityCurated
	}
	return QualityCatalog
}

func hasCuratedSignals(m GameplayMap) bool {
	if len(m.Subject.Influences) > 0 {
		return true
	}
	specificNotes := 0
	for _, b := range m.Mechanics {
		if b.MapNotes == "" {
			continue
		}
		if templateSigNotesRe.MatchString(b.MapNotes) || templateSupportNotes.MatchString(b.MapNotes) {
			continue
		}
		if upliftNotesRe.MatchString(b.MapNotes) {
			continue
		}
		if len(b.MapNotes) > 40 {
			specificNotes++
		}
	}
	return specificNotes >= 3
}

// ActionAdventureTemplatePack lists supporting mechanics wrongly assigned by genre templates.
var ActionAdventureTemplatePack = map[string]struct{}{
	"boss-gated-progression":            {},
	"melee-primary-combat":              {},
	"environmental-hazard-navigation":     {},
	"depletable-health-pool":            {},
	"mid-level-checkpoint":              {},
	"patterned-boss-encounter":           {},
	"telegraphed-threat-loop":            {},
	"progress-gated-dungeon-completion":   {},
	"health-capacity-upgrade":           {},
	"environmental-trigger-switch":      {},
	"spin-attack-combat":                {},
	"ability-gated-exploration":           {},
	"hub-world-stage-gates":               {},
	"exploration-gated-upgrades":          {},
}

// NonCombatSignatures indicate games that should not carry action-adventure combat packs.
var NonCombatSignatures = map[string]struct{}{
	"document-inspection-stamping": {},
	"dialogue-choice-consequences": {},
	"city-building-simulation":     {},
	"city-building-economy":        {},
	"rhythm-battle-system":         {},
	"rhythm-tap-timing":            {},
}

// ShouldStripTemplatePack reports whether action-adventure supporting mechanics are inappropriate.
func ShouldStripTemplatePack(m GameplayMap) bool {
	for _, sig := range m.SignatureGameplay {
		if _, ok := NonCombatSignatures[sig]; ok {
			return true
		}
	}
	return false
}

// StripInappropriateTemplateMechanics removes genre-template supporting mechanics when signatures disagree.
func StripInappropriateTemplateMechanics(m GameplayMap) GameplayMap {
	if !ShouldStripTemplatePack(m) {
		return m
	}
	sigSet := make(map[string]struct{}, len(m.SignatureGameplay))
	for _, s := range m.SignatureGameplay {
		sigSet[s] = struct{}{}
	}
	filtered := make([]MapMechanicBinding, 0, len(m.Mechanics))
	for _, b := range m.Mechanics {
		if b.Role == RoleSignature {
			filtered = append(filtered, b)
			continue
		}
		if _, drop := ActionAdventureTemplatePack[b.MechanicSlug]; drop {
			continue
		}
		filtered = append(filtered, b)
	}
	m.Mechanics = filtered
	return m
}
