package internal

import (
	"regexp"
	"strings"
	"unicode"
)

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

// DeriveSlug converts a display name to a kebab-case map slug.
func DeriveSlug(name string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && b.Len() > 0 {
			b.WriteByte('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	out = nonSlugChars.ReplaceAllString(out, "-")
	out = strings.Trim(out, "-")
	for strings.Contains(out, "--") {
		out = strings.ReplaceAll(out, "--", "-")
	}
	return out
}

// NewGameMapScaffold returns a prefilled gameplay map skeleton for analysis.
func NewGameMapScaffold(gameName, slug string, genres []string) GameplayMap {
	if slug == "" {
		slug = DeriveSlug(gameName)
	}
	title := gameName + " Gameplay Map"
	return GameplayMap{
		SchemaVersion:     SchemaVersion,
		Slug:              slug,
		Title:             title,
		MapType:           MapTypeGame,
		SourceURL:         "",
		SignatureGameplay: []string{},
		Mechanics:         []MapMechanicBinding{},
		Subject: Subject{
			Kind:       "game",
			Name:       gameName,
			Genres:     append([]string(nil), genres...),
			Influences: []string{},
		},
		Narrative: Narrative{
			Description:  "",
			CoreLoop:     "",
			SkillsTested: []string{},
		},
		Views: StandardMapViews(),
	}
}

// StandardMapViews returns Notion-style filter tabs used across seed maps.
func StandardMapViews() []MapView {
	return []MapView{
		{ID: "all", Label: "All", Filter: ViewFilter{}},
		{ID: "signature", Label: "Signature Gameplay", Filter: ViewFilter{Role: ptrRole(RoleSignature)}},
		{ID: "action", Label: "Action", Filter: ViewFilter{Flavor: ptrFlavor(FlavorAction)}},
		{ID: "adventure", Label: "Adventure", Filter: ViewFilter{Flavor: ptrFlavor(FlavorAdventure)}},
		{ID: "strategy", Label: "Strategy", Filter: ViewFilter{Flavor: ptrFlavor(FlavorStrategy)}},
	}
}

func ptrRole(r MechanicRole) *MechanicRole {
	return &r
}

func ptrFlavor(f Flavor) *Flavor {
	return &f
}
