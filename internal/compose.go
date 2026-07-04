package internal

import (
	"fmt"
	"sort"
	"strings"
)

// ComposeMapOptions configures stub map generation from genre recipes and references.
type ComposeMapOptions struct {
	Slug           string
	Title          string
	GameName       string
	Genres         []string
	GenreMaps      []string
	ReferenceGames []string
	Notes          string
	CuratedOnly    bool
}

// ComposeMap merges genre recipe maps and optional reference game maps into a stub gameplay map.
func (b *Bundle) ComposeMap(opts ComposeMapOptions) (*GameplayMap, error) {
	if opts.Slug == "" {
		return nil, fmt.Errorf("slug is required")
	}
	if len(opts.GenreMaps) == 0 && len(opts.ReferenceGames) == 0 {
		return nil, fmt.Errorf("at least one genre map or reference game is required")
	}

	sigOrder := []string{}
	sigSet := map[string]struct{}{}
	supportSet := map[string]MapMechanicBinding{}
	refsUsed := []string{}
	genresUsed := []string{}
	loopPhases := []LoopPhase{}
	conflicts := map[string]struct{}{}

	addSig := func(slug string) {
		if _, ok := sigSet[slug]; ok {
			return
		}
		sigSet[slug] = struct{}{}
		sigOrder = append(sigOrder, slug)
	}

	mergeMap := func(m GameplayMap) {
		if opts.CuratedOnly && DetectQualityTier(m) != QualityCurated && DetectQualityTier(m) != QualityStub {
			return
		}
		for _, s := range m.SignatureGameplay {
			addSig(s)
		}
		for _, binding := range m.Mechanics {
			if binding.Role == RoleSignature {
				if _, ok := supportSet[binding.MechanicSlug]; !ok {
					nb := binding
					nb.MapNotes = ""
					supportSet[binding.MechanicSlug] = nb
				}
				continue
			}
			if _, ok := sigSet[binding.MechanicSlug]; ok {
				continue
			}
			if existing, ok := supportSet[binding.MechanicSlug]; ok {
				if binding.Weight > existing.Weight {
					existing.Weight = binding.Weight
				}
				if binding.Phase != "" && existing.Phase == "" {
					existing.Phase = binding.Phase
				}
				if binding.Domain != "" && existing.Domain == "" {
					existing.Domain = binding.Domain
				}
				supportSet[binding.MechanicSlug] = existing
				continue
			}
			nb := binding
			nb.Role = RoleSupporting
			nb.MapNotes = ""
			supportSet[binding.MechanicSlug] = nb
		}
		if m.Systems != nil && len(m.Systems.PrimaryLoopPhases) > 0 && len(loopPhases) == 0 {
			loopPhases = append(loopPhases, m.Systems.PrimaryLoopPhases...)
		}
	}

	for _, gslug := range opts.GenreMaps {
		m, ok := b.Maps[gslug]
		if !ok {
			return nil, fmt.Errorf("unknown genre map %q", gslug)
		}
		if m.MapType != MapTypeGenre {
			return nil, fmt.Errorf("%q is not a genre map", gslug)
		}
		genresUsed = append(genresUsed, gslug)
		mergeMap(m)
	}

	for _, rslug := range opts.ReferenceGames {
		m, ok := b.Maps[rslug]
		if !ok {
			return nil, fmt.Errorf("unknown reference map %q", rslug)
		}
		if m.MapType != MapTypeGame {
			return nil, fmt.Errorf("%q is not a game map", rslug)
		}
		if opts.CuratedOnly && DetectQualityTier(m) == QualityTemplate {
			continue
		}
		refsUsed = append(refsUsed, rslug)
		mergeMap(m)
	}

	if len(sigOrder) == 0 {
		return nil, fmt.Errorf("no signatures merged; check quality tier filters or inputs")
	}

	for _, sig := range sigOrder {
		entry, ok := b.Mechanics[sig]
		if !ok {
			continue
		}
		for _, c := range entry.Conflicts {
			conflicts[strings.ToLower(c)] = struct{}{}
		}
	}

	filteredSupport := map[string]MapMechanicBinding{}
	for slug, binding := range supportSet {
		if _, isSig := sigSet[slug]; isSig {
			continue
		}
		entry, ok := b.Mechanics[slug]
		if ok {
			for _, c := range entry.Conflicts {
				conflicts[strings.ToLower(c)] = struct{}{}
			}
			nameLower := strings.ToLower(entry.Name)
			if _, bad := conflicts[nameLower]; bad {
				continue
			}
		}
		filteredSupport[slug] = binding
	}

	bindings := make([]MapMechanicBinding, 0, len(sigOrder)+len(filteredSupport))
	for _, sig := range sigOrder {
		entry, ok := b.Mechanics[sig]
		note := fmt.Sprintf("Describe how %s implements %s.", opts.GameName, sig)
		if ok {
			note = fmt.Sprintf("Describe how %s implements %s (%s).", opts.GameName, entry.Name, entry.Summary)
		}
		bindings = append(bindings, MapMechanicBinding{
			MechanicSlug: sig,
			Role:         RoleSignature,
			MapNotes:     note,
		})
	}
	supportSlugs := make([]string, 0, len(filteredSupport))
	for s := range filteredSupport {
		supportSlugs = append(supportSlugs, s)
	}
	sort.Strings(supportSlugs)
	for _, slug := range supportSlugs {
		binding := filteredSupport[slug]
		binding.Role = RoleSupporting
		entry, ok := b.Mechanics[slug]
		if binding.MapNotes == "" && ok {
			binding.MapNotes = fmt.Sprintf("Supporting role for %s: %s", opts.GameName, entry.Summary)
		}
		bindings = append(bindings, binding)
	}

	title := opts.Title
	if title == "" {
		title = opts.GameName + " Gameplay Map"
	}
	gameName := opts.GameName
	if gameName == "" {
		gameName = strings.TrimSuffix(title, " Gameplay Map")
	}

	meta := false
	out := &GameplayMap{
		SchemaVersion:     SchemaVersion,
		Slug:              opts.Slug,
		Title:             title,
		MapType:           MapTypeGame,
		SourceURL:         "https://github.com/blazium-games/game-design-index",
		SignatureGameplay: append([]string(nil), sigOrder...),
		Mechanics:         bindings,
		Views:             StandardMapViews(),
		Subject: Subject{
			Kind:   "game",
			Name:   gameName,
			Genres: append([]string(nil), opts.Genres...),
		},
		Narrative: Narrative{
			Description: fmt.Sprintf("%s — composed stub map. Fill narrative after playtesting or research.", gameName),
			CoreLoop:  "Define the minute-to-minute loop for this game.",
			SkillsTested: []string{
				"timing", "positioning", "resource management",
			},
		},
		Metadata: &MapMetadata{QualityTier: QualityStub},
		Composition: &MapComposition{
			BaseGenreMaps:      append([]string(nil), genresUsed...),
			ReferenceGames:     append([]string(nil), refsUsed...),
			CustomizationNotes: opts.Notes,
		},
		Relationships: &MapRelationships{
			GenreMaps:      append([]string(nil), genresUsed...),
			InfluenceSlugs: append([]string(nil), refsUsed...),
		},
	}
	if len(loopPhases) > 0 {
		out.Systems = &MapSystems{PrimaryLoopPhases: loopPhases}
		meta = true
	}
	_ = meta

	if len(conflicts) > 0 {
		// surfaced via composition notes for human review
		var names []string
		for c := range conflicts {
			names = append(names, c)
		}
		sort.Strings(names)
		if out.Composition.CustomizationNotes != "" {
			out.Composition.CustomizationNotes += " "
		}
		out.Composition.CustomizationNotes += "Review conflicts: " + strings.Join(names, ", ")
	}

	return out, nil
}

// BuildDesignBriefFromInputs merges genre maps and/or game maps into a design brief.
func (b *Bundle) BuildDesignBriefFromInputs(genreSlugs, gameSlugs []string, minCooccurrence int, curatedOnly bool) (*DesignBrief, error) {
	var refs []string
	for _, g := range genreSlugs {
		m, ok := b.Maps[g]
		if !ok {
			return nil, fmt.Errorf("unknown genre map %q", g)
		}
		if m.MapType != MapTypeGenre {
			return nil, fmt.Errorf("%q is not a genre map", g)
		}
		refs = append(refs, g)
	}
	for _, s := range gameSlugs {
		m, ok := b.Maps[s]
		if !ok {
			return nil, fmt.Errorf("unknown map %q", s)
		}
		if curatedOnly && DetectQualityTier(m) == QualityTemplate {
			continue
		}
		refs = append(refs, s)
	}
	if len(refs) == 0 {
		return nil, fmt.Errorf("no reference maps after filtering")
	}
	return b.BuildDesignBrief(refs, minCooccurrence)
}
