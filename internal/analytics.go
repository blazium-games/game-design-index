package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const AnalyticsSchemaVersion = "1.0"

type CountRow struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

type MechanicAdoptionRow struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type CooccurrenceAnalyticsRow struct {
	MechanicA string  `json:"mechanic_a"`
	MechanicB string  `json:"mechanic_b"`
	Count     int     `json:"count"`
	Lift      float64 `json:"lift"`
}

type HeatmapCell struct {
	Genre  string `json:"genre"`
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

type GenreDomainHeatmap struct {
	Genres  []string      `json:"genres"`
	Domains []string      `json:"domains"`
	Cells   []HeatmapCell `json:"cells"`
}

type EntityEnrichmentStats struct {
	ByCategory          []CountRow `json:"by_category,omitempty"`
	ByScope             []CountRow `json:"by_scope,omitempty"`
	ByMenuType          []CountRow `json:"by_menu_type,omitempty"`
	ByLayer             []CountRow `json:"by_layer,omitempty"`
	EnrichmentComplete  int        `json:"enrichment_complete"`
	EnrichmentNeedsInfo int        `json:"enrichment_needs_info"`
	WithMapBindings     int        `json:"with_map_bindings"`
}

type MenuHubRow struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	InDegree    int    `json:"in_degree"`
	OutDegree   int    `json:"out_degree"`
	TotalDegree int    `json:"total_degree"`
}

type MenuFlowStats struct {
	EdgeCountsByRelationship []CountRow   `json:"edge_counts_by_relationship"`
	TopHubs                  []MenuHubRow `json:"top_hubs"`
}

type VariableMechanicPair struct {
	Variable string `json:"variable"`
	Mechanic string `json:"mechanic"`
}

type AnalyticsOverview struct {
	MapCount               int     `json:"map_count"`
	GameCount              int     `json:"game_count"`
	GenreRecipeCount       int     `json:"genre_recipe_count"`
	MechanicCount          int     `json:"mechanic_count"`
	VariableCount          int     `json:"variable_count"`
	MenuCount              int     `json:"menu_count"`
	SkillCount             int     `json:"skill_count"`
	MechanicEnrichmentPct  float64 `json:"mechanic_enrichment_pct"`
	SkillEnrichmentPct     float64 `json:"skill_enrichment_pct"`
	VariableEnrichmentPct  float64 `json:"variable_enrichment_pct"`
	MenuEnrichmentPct      float64 `json:"menu_enrichment_pct"`
	AvgSignatureCount      float64 `json:"avg_signature_count"`
}

type AnalyticsSnapshot struct {
	SchemaVersion         string                   `json:"schema_version"`
	Overview              AnalyticsOverview        `json:"overview"`
	QualityTiers          []CountRow               `json:"quality_tiers"`
	GamesByDecade         []CountRow               `json:"games_by_decade"`
	TopGenres             []CountRow               `json:"top_genres"`
	MechanicDomains       []CountRow               `json:"mechanic_domains"`
	MechanicFlavors       []CountRow               `json:"mechanic_flavors"`
	MechanicComplexity    []CountRow               `json:"mechanic_complexity"`
	TopMechanics          []MechanicAdoptionRow    `json:"top_mechanics"`
	GenreDomainHeatmap    GenreDomainHeatmap       `json:"genre_domain_heatmap"`
	Cooccurrence          []CooccurrenceAnalyticsRow `json:"cooccurrence"`
	MechanicStats         EntityEnrichmentStats    `json:"mechanic_stats"`
	SkillStats            EntityEnrichmentStats    `json:"skill_stats"`
	VariableStats         EntityEnrichmentStats    `json:"variable_stats"`
	MenuStats             EntityEnrichmentStats    `json:"menu_stats"`
	VariableMechanicPairs []VariableMechanicPair   `json:"variable_mechanic_pairs"`
	MenuFlow              MenuFlowStats            `json:"menu_flow"`
	SignatureDistribution []CountRow               `json:"signature_distribution"`
	Insights              []string                 `json:"insights"`
}

// BuildAnalytics computes corpus-wide statistics for the public analytics page.
func BuildAnalytics(b *Bundle) AnalyticsSnapshot {
	gameCount := 0
	genreCount := 0
	for _, m := range b.Maps {
		if m.MapType == MapTypeGenre {
			genreCount++
		} else {
			gameCount++
		}
	}

	mechComplete, mechNeeds := 0, 0
	for _, m := range b.Mechanics {
		if MechanicEnrichmentStatus(m) == "complete" {
			mechComplete++
		} else {
			mechNeeds++
		}
	}
	skillComplete, skillNeeds := 0, 0
	for _, s := range b.Skills {
		if SkillEnrichmentStatus(s) == "complete" {
			skillComplete++
		} else {
			skillNeeds++
		}
	}
	varComplete, varNeeds := 0, 0
	for _, v := range b.Variables {
		if VariableEnrichmentStatus(v) == "complete" {
			varComplete++
		} else {
			varNeeds++
		}
	}
	menuComplete, menuNeeds := 0, 0
	for _, m := range b.UIMenus {
		if MenuEnrichmentStatus(m) == "complete" {
			menuComplete++
		} else {
			menuNeeds++
		}
	}

	sigTotal, sigGames := 0, 0
	for _, m := range b.Maps {
		if m.MapType != MapTypeGame {
			continue
		}
		n := len(m.SignatureGameplay)
		if n > 0 {
			sigTotal += n
			sigGames++
		}
	}
	avgSig := 0.0
	if sigGames > 0 {
		avgSig = float64(sigTotal) / float64(sigGames)
	}

	mechEnrichPct := 0.0
	if len(b.Mechanics) > 0 {
		mechEnrichPct = float64(mechComplete) / float64(len(b.Mechanics)) * 100
	}
	skillEnrichPct := 0.0
	if len(b.Skills) > 0 {
		skillEnrichPct = float64(skillComplete) / float64(len(b.Skills)) * 100
	}
	varEnrichPct := 0.0
	if len(b.Variables) > 0 {
		varEnrichPct = float64(varComplete) / float64(len(b.Variables)) * 100
	}
	menuEnrichPct := 0.0
	if len(b.UIMenus) > 0 {
		menuEnrichPct = float64(menuComplete) / float64(len(b.UIMenus)) * 100
	}

	snap := AnalyticsSnapshot{
		SchemaVersion: AnalyticsSchemaVersion,
		Overview: AnalyticsOverview{
			MapCount:              len(b.Maps),
			GameCount:             gameCount,
			GenreRecipeCount:      genreCount,
			MechanicCount:         len(b.Mechanics),
			VariableCount:         len(b.Variables),
			MenuCount:             len(b.UIMenus),
			SkillCount:            len(b.Skills),
			MechanicEnrichmentPct: round2(mechEnrichPct),
			SkillEnrichmentPct:    round2(skillEnrichPct),
			VariableEnrichmentPct: round2(varEnrichPct),
			MenuEnrichmentPct:     round2(menuEnrichPct),
			AvgSignatureCount:     round2(avgSig),
		},
		QualityTiers:          countQualityTiers(b),
		GamesByDecade:         countGamesByDecade(b),
		TopGenres:             topGenres(b, 20),
		MechanicDomains:       countMechanicDomains(b),
		MechanicFlavors:       countMechanicFlavors(b),
		MechanicComplexity:    countMechanicComplexity(b),
		TopMechanics:          topMechanicsByAdoption(b, 25),
		GenreDomainHeatmap:    buildGenreDomainHeatmap(b, 12),
		Cooccurrence:          buildCooccurrenceWithLift(b, 50),
		MechanicStats:         buildMechanicStats(b),
		SkillStats:            buildSkillStats(b),
		VariableStats:         buildVariableStats(b),
		MenuStats:             buildMenuStats(b),
		VariableMechanicPairs: loadVariableMechanicPairs(b, 20),
		MenuFlow:              buildMenuFlowStats(b),
		SignatureDistribution: countSignatureDistribution(b),
	}
	snap.Insights = generateInsights(snap, b)
	return snap
}

func countQualityTiers(b *Bundle) []CountRow {
	counts := map[string]int{}
	for _, m := range b.Maps {
		if m.MapType != MapTypeGame {
			continue
		}
		tier := string(DetectQualityTier(m))
		counts[tier]++
	}
	order := []string{"curated", "catalog", "template", "stub"}
	return rowsFromMap(counts, order)
}

func mapYear(m GameplayMap) int {
	if m.Context != nil && m.Context.Year > 0 {
		return m.Context.Year
	}
	return extractYearFromDescription(m.Narrative.Description)
}

func extractYearFromDescription(desc string) int {
	if desc == "" {
		return 0
	}
	idx := strings.LastIndex(desc, "(")
	if idx < 0 {
		return 0
	}
	end := strings.Index(desc[idx:], ")")
	if end <= 1 {
		return 0
	}
	inner := strings.TrimSpace(desc[idx+1 : idx+end])
	if len(inner) == 4 {
		var y int
		if _, err := fmt.Sscanf(inner, "%d", &y); err == nil {
			return y
		}
	}
	return 0
}

func countGamesByDecade(b *Bundle) []CountRow {
	counts := map[string]int{}
	for _, m := range b.Maps {
		if m.MapType != MapTypeGame {
			continue
		}
		y := mapYear(m)
		if y < 1970 {
			counts["unknown"]++
			continue
		}
		decade := fmt.Sprintf("%ds", (y/10)*10)
		counts[decade]++
	}
	labels := make([]string, 0, len(counts))
	for k := range counts {
		labels = append(labels, k)
	}
	sort.Slice(labels, func(i, j int) bool {
		if labels[i] == "unknown" {
			return false
		}
		if labels[j] == "unknown" {
			return true
		}
		return labels[i] < labels[j]
	})
	out := make([]CountRow, len(labels))
	for i, l := range labels {
		out[i] = CountRow{Label: l, Count: counts[l]}
	}
	return out
}

func topGenres(b *Bundle, limit int) []CountRow {
	counts := map[string]int{}
	for _, m := range b.Maps {
		if m.MapType != MapTypeGame {
			continue
		}
		for _, g := range m.Subject.Genres {
			counts[g]++
		}
	}
	return topRows(counts, limit)
}

func countMechanicDomains(b *Bundle) []CountRow {
	counts := map[string]int{}
	for _, e := range b.Mechanics {
		counts[string(e.Domain)]++
	}
	order := []string{"locomotion", "combat", "progression", "economy", "level", "session"}
	return rowsFromMap(counts, order)
}

func countMechanicFlavors(b *Bundle) []CountRow {
	counts := map[string]int{}
	for _, e := range b.Mechanics {
		counts[string(e.Flavor)]++
	}
	order := []string{"action", "adventure", "strategy"}
	return rowsFromMap(counts, order)
}

func countMechanicComplexity(b *Bundle) []CountRow {
	counts := map[string]int{"S": 0, "M": 0, "L": 0, "unset": 0}
	for _, e := range b.Mechanics {
		if e.Complexity == "" {
			counts["unset"]++
		} else {
			counts[string(e.Complexity)]++
		}
	}
	order := []string{"S", "M", "L", "unset"}
	return rowsFromMap(counts, order)
}

func topMechanicsByAdoption(b *Bundle, limit int) []MechanicAdoptionRow {
	mechToMaps := map[string][]string{}
	for slug, m := range b.Maps {
		seen := map[string]struct{}{}
		for _, bind := range m.Mechanics {
			if _, ok := seen[bind.MechanicSlug]; ok {
				continue
			}
			seen[bind.MechanicSlug] = struct{}{}
			mechToMaps[bind.MechanicSlug] = append(mechToMaps[bind.MechanicSlug], slug)
		}
	}
	type row struct {
		slug  string
		count int
	}
	var rows []row
	for slug, maps := range mechToMaps {
		rows = append(rows, row{slug: slug, count: len(maps)})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].count != rows[j].count {
			return rows[i].count > rows[j].count
		}
		return rows[i].slug < rows[j].slug
	})
	if len(rows) > limit {
		rows = rows[:limit]
	}
	out := make([]MechanicAdoptionRow, len(rows))
	for i, r := range rows {
		name := r.slug
		if e, ok := b.Mechanics[r.slug]; ok {
			name = e.Name
		}
		out[i] = MechanicAdoptionRow{Slug: r.slug, Name: name, Count: r.count}
	}
	return out
}

func buildGenreDomainHeatmap(b *Bundle, genreLimit int) GenreDomainHeatmap {
	genreCounts := map[string]int{}
	for _, m := range b.Maps {
		if m.MapType != MapTypeGame {
			continue
		}
		for _, g := range m.Subject.Genres {
			genreCounts[g]++
		}
	}
	topGenreLabels := topLabels(genreCounts, genreLimit)
	genreSet := map[string]struct{}{}
	for _, g := range topGenreLabels {
		genreSet[g] = struct{}{}
	}

	domains := []string{"locomotion", "combat", "progression", "economy", "level", "session"}
	cellCounts := map[string]int{}

	for _, m := range b.Maps {
		if m.MapType != MapTypeGame {
			continue
		}
		mapGenres := make([]string, 0)
		for _, g := range m.Subject.Genres {
			if _, ok := genreSet[g]; ok {
				mapGenres = append(mapGenres, g)
			}
		}
		if len(mapGenres) == 0 {
			continue
		}
		domainSeen := map[string]map[string]struct{}{}
		for _, bind := range m.Mechanics {
			domain := string(bind.Domain)
			if domain == "" {
				if e, ok := b.Mechanics[bind.MechanicSlug]; ok {
					domain = string(e.Domain)
				}
			}
			if domain == "" {
				continue
			}
			if domainSeen[domain] == nil {
				domainSeen[domain] = map[string]struct{}{}
			}
			for _, g := range mapGenres {
				key := g + "\x00" + domain
				if _, ok := domainSeen[domain][g]; ok {
					continue
				}
				domainSeen[domain][g] = struct{}{}
				cellCounts[key]++
			}
		}
	}

	var cells []HeatmapCell
	for _, g := range topGenreLabels {
		for _, d := range domains {
			key := g + "\x00" + d
			cells = append(cells, HeatmapCell{Genre: g, Domain: d, Count: cellCounts[key]})
		}
	}
	return GenreDomainHeatmap{Genres: topGenreLabels, Domains: domains, Cells: cells}
}

func buildCooccurrenceWithLift(b *Bundle, limit int) []CooccurrenceAnalyticsRow {
	totalMaps := len(b.Maps)
	if totalMaps == 0 {
		return nil
	}

	mechMapCount := map[string]int{}
	for _, m := range b.Maps {
		seen := map[string]struct{}{}
		for _, bind := range m.Mechanics {
			if _, ok := seen[bind.MechanicSlug]; ok {
				continue
			}
			seen[bind.MechanicSlug] = struct{}{}
			mechMapCount[bind.MechanicSlug]++
		}
	}

	pairs := b.MechanicCooccurrenceMatrix(1)
	out := make([]CooccurrenceAnalyticsRow, 0, limit)
	for _, p := range pairs {
		if len(out) >= limit {
			break
		}
		countA := mechMapCount[p.MechanicA]
		countB := mechMapCount[p.MechanicB]
		expected := float64(totalMaps) * (float64(countA) / float64(totalMaps)) * (float64(countB) / float64(totalMaps))
		lift := 0.0
		if expected > 0 {
			lift = float64(p.Count) / expected
		}
		out = append(out, CooccurrenceAnalyticsRow{
			MechanicA: p.MechanicA,
			MechanicB: p.MechanicB,
			Count:     p.Count,
			Lift:      round2(lift),
		})
	}
	return out
}

func buildMechanicStats(b *Bundle) EntityEnrichmentStats {
	domainCounts := map[string]int{}
	flavorCounts := map[string]int{}
	complete, needs := 0, 0
	for _, m := range b.Mechanics {
		domainCounts[string(m.Domain)]++
		flavorCounts[string(m.Flavor)]++
		if MechanicEnrichmentStatus(m) == "complete" {
			complete++
		} else {
			needs++
		}
	}
	return EntityEnrichmentStats{
		ByCategory:          topRows(domainCounts, 20),
		ByScope:             topRows(flavorCounts, 20),
		EnrichmentComplete:  complete,
		EnrichmentNeedsInfo: needs,
	}
}

func buildSkillStats(b *Bundle) EntityEnrichmentStats {
	catCounts := map[string]int{}
	complete, needs := 0, 0
	for _, s := range b.Skills {
		catCounts[string(s.Category)]++
		if SkillEnrichmentStatus(s) == "complete" {
			complete++
		} else {
			needs++
		}
	}
	return EntityEnrichmentStats{
		ByCategory:          topRows(catCounts, 20),
		EnrichmentComplete:  complete,
		EnrichmentNeedsInfo: needs,
	}
}

func buildVariableStats(b *Bundle) EntityEnrichmentStats {
	catCounts := map[string]int{}
	scopeCounts := map[string]int{}
	complete, needs := 0, 0
	withBindings := 0
	for _, v := range b.Variables {
		catCounts[string(v.Category)]++
		scopeCounts[string(v.Scope)]++
		if VariableEnrichmentStatus(v) == "complete" {
			complete++
		} else {
			needs++
		}
		if len(b.MapsWithVariable(v.Slug)) > 0 {
			withBindings++
		}
	}
	return EntityEnrichmentStats{
		ByCategory:          topRows(catCounts, 20),
		ByScope:             topRows(scopeCounts, 20),
		EnrichmentComplete:  complete,
		EnrichmentNeedsInfo: needs,
		WithMapBindings:     withBindings,
	}
}

func buildMenuStats(b *Bundle) EntityEnrichmentStats {
	typeCounts := map[string]int{}
	layerCounts := map[string]int{}
	complete, needs := 0, 0
	withBindings := 0
	for _, m := range b.UIMenus {
		typeCounts[string(m.MenuType)]++
		layerCounts[string(m.Layer)]++
		if MenuEnrichmentStatus(m) == "complete" {
			complete++
		} else {
			needs++
		}
		if len(b.MapsWithMenu(m.Slug)) > 0 {
			withBindings++
		}
	}
	return EntityEnrichmentStats{
		ByMenuType:          topRows(typeCounts, 20),
		ByLayer:             topRows(layerCounts, 20),
		EnrichmentComplete:  complete,
		EnrichmentNeedsInfo: needs,
		WithMapBindings:     withBindings,
	}
}

type variableToMechanicsFile struct {
	Variables map[string][]string `json:"variables"`
}

func loadVariableMechanicPairs(b *Bundle, limit int) []VariableMechanicPair {
	path := filepath.Join(b.Root, "indexes", "variable-to-mechanics.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var file variableToMechanicsFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil
	}
	type pair struct {
		variable string
		mechanic string
		weight   int
	}
	var pairs []pair
	for varSlug, mechs := range file.Variables {
		weight := len(b.MapsWithVariable(varSlug))
		for _, mech := range mechs {
			pairs = append(pairs, pair{variable: varSlug, mechanic: mech, weight: weight})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].weight != pairs[j].weight {
			return pairs[i].weight > pairs[j].weight
		}
		if pairs[i].variable != pairs[j].variable {
			return pairs[i].variable < pairs[j].variable
		}
		return pairs[i].mechanic < pairs[j].mechanic
	})
	if len(pairs) > limit {
		pairs = pairs[:limit]
	}
	out := make([]VariableMechanicPair, len(pairs))
	for i, p := range pairs {
		out[i] = VariableMechanicPair{Variable: p.variable, Mechanic: p.mechanic}
	}
	return out
}

type menuFlowEdgesFile struct {
	Edges []MenuFlowEdge `json:"edges"`
}

func buildMenuFlowStats(b *Bundle) MenuFlowStats {
	path := filepath.Join(b.Root, "indexes", "menu-flow-edges.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return MenuFlowStats{}
	}
	var file menuFlowEdgesFile
	if err := json.Unmarshal(data, &file); err != nil {
		return MenuFlowStats{}
	}

	relCounts := map[string]int{}
	inDeg := map[string]int{}
	outDeg := map[string]int{}
	for _, e := range file.Edges {
		rel := e.Relationship
		if rel == "" {
			rel = "unspecified"
		}
		relCounts[rel]++
		outDeg[e.FromMenu]++
		inDeg[e.ToMenu]++
	}

	type hub struct {
		slug  string
		in    int
		out   int
		total int
	}
	var hubs []hub
	seen := map[string]struct{}{}
	for slug := range inDeg {
		seen[slug] = struct{}{}
	}
	for slug := range outDeg {
		seen[slug] = struct{}{}
	}
	for slug := range seen {
		in, out := inDeg[slug], outDeg[slug]
		hubs = append(hubs, hub{slug: slug, in: in, out: out, total: in + out})
	}
	sort.Slice(hubs, func(i, j int) bool {
		if hubs[i].total != hubs[j].total {
			return hubs[i].total > hubs[j].total
		}
		return hubs[i].slug < hubs[j].slug
	})
	if len(hubs) > 10 {
		hubs = hubs[:10]
	}

	topHubs := make([]MenuHubRow, len(hubs))
	for i, h := range hubs {
		name := h.slug
		if menu, ok := b.UIMenus[h.slug]; ok {
			name = menu.Name
		}
		topHubs[i] = MenuHubRow{
			Slug:        h.slug,
			Name:        name,
			InDegree:    h.in,
			OutDegree:   h.out,
			TotalDegree: h.total,
		}
	}

	return MenuFlowStats{
		EdgeCountsByRelationship: topRows(relCounts, 20),
		TopHubs:                  topHubs,
	}
}

func countSignatureDistribution(b *Bundle) []CountRow {
	buckets := map[string]int{}
	for _, m := range b.Maps {
		if m.MapType != MapTypeGame {
			continue
		}
		n := len(m.SignatureGameplay)
		label := fmt.Sprintf("%d", n)
		if n >= 7 {
			label = "7+"
		}
		buckets[label]++
	}
	order := []string{"3", "4", "5", "6", "7+"}
	return rowsFromMap(buckets, order)
}

func generateInsights(snap AnalyticsSnapshot, b *Bundle) []string {
	var insights []string

	if len(snap.TopGenres) > 0 {
		top := snap.TopGenres[0]
		insights = append(insights, fmt.Sprintf(
			"%q is the most common genre label (%d games), making it the strongest axis for genre×domain heatmap comparisons.",
			top.Label, top.Count,
		))
	}

	if len(snap.QualityTiers) > 0 {
		var curated, template int
		for _, q := range snap.QualityTiers {
			if q.Label == "curated" {
				curated = q.Count
			}
			if q.Label == "template" {
				template = q.Count
			}
		}
		if curated > 0 && template > 0 {
			pct := float64(curated) / float64(curated+template) * 100
			insights = append(insights, fmt.Sprintf(
				"Among quality-tracked games, %.0f%% of curated+template maps are curated (%d) vs template stubs (%d).",
				pct, curated, template,
			))
		}
	}

	if len(snap.MechanicDomains) > 0 {
		topDomain := snap.MechanicDomains[0]
		maxCount := topDomain.Count
		for _, d := range snap.MechanicDomains {
			if d.Count > maxCount {
				maxCount = d.Count
				topDomain = d
			}
		}
		insights = append(insights, fmt.Sprintf(
			"The mechanic library is densest in the %q domain (%d entries), which often drives cross-genre co-occurrence patterns.",
			topDomain.Label, topDomain.Count,
		))
	}

	if len(snap.Cooccurrence) > 0 {
		top := snap.Cooccurrence[0]
		highLift := top
		for _, c := range snap.Cooccurrence {
			if c.Lift > highLift.Lift {
				highLift = c
			}
		}
		insights = append(insights, fmt.Sprintf(
			"Highest raw co-occurrence: %s + %s (%d maps). Strongest lift (association beyond marginal frequency): %s + %s (lift %.2f).",
			top.MechanicA, top.MechanicB, top.Count,
			highLift.MechanicA, highLift.MechanicB, highLift.Lift,
		))
	}

	if len(snap.TopMechanics) > 0 {
		tm := snap.TopMechanics[0]
		insights = append(insights, fmt.Sprintf(
			"%q (%s) appears on the most maps (%d), indicating broad genre adoption rather than a niche signature.",
			tm.Name, tm.Slug, tm.Count,
		))
	}

	if snap.Overview.MechanicCount > 0 {
		insights = append(insights, fmt.Sprintf(
			"Mechanics: %.0f%% design pedagogy enrichment complete (%d/%d); %d still need design_guidance and agent_context.",
			snap.Overview.MechanicEnrichmentPct,
			snap.MechanicStats.EnrichmentComplete, snap.Overview.MechanicCount,
			snap.MechanicStats.EnrichmentNeedsInfo,
		))
	}

	if snap.Overview.SkillCount > 0 {
		insights = append(insights, fmt.Sprintf(
			"Skills catalog: %.0f%% enrichment complete (%d/%d).",
			snap.Overview.SkillEnrichmentPct,
			snap.SkillStats.EnrichmentComplete, snap.Overview.SkillCount,
		))
	}

	if snap.Overview.VariableCount > 0 {
		insights = append(insights, fmt.Sprintf(
			"Variables: %.0f%% enrichment complete (%d/%d); %d variables have map bindings in pilot maps.",
			snap.Overview.VariableEnrichmentPct,
			snap.VariableStats.EnrichmentComplete, snap.Overview.VariableCount,
			snap.VariableStats.WithMapBindings,
		))
	}

	if len(snap.MenuFlow.TopHubs) > 0 {
		hub := snap.MenuFlow.TopHubs[0]
		insights = append(insights, fmt.Sprintf(
			"Menu navigation hub: %q (%s) has the highest total degree (%d in + %d out edges).",
			hub.Name, hub.Slug, hub.InDegree, hub.OutDegree,
		))
	}

	if len(snap.GamesByDecade) > 1 {
		var peak CountRow
		for _, d := range snap.GamesByDecade {
			if d.Label != "unknown" && d.Count > peak.Count {
				peak = d
			}
		}
		if peak.Count > 0 {
			insights = append(insights, fmt.Sprintf(
				"Release-year coverage peaks in the %s decade with %d indexed games.",
				peak.Label, peak.Count,
			))
		}
	}

	if len(snap.GenreDomainHeatmap.Cells) > 0 {
		maxCell := snap.GenreDomainHeatmap.Cells[0]
		for _, c := range snap.GenreDomainHeatmap.Cells {
			if c.Count > maxCell.Count {
				maxCell = c
			}
		}
		if maxCell.Count > 0 {
			insights = append(insights, fmt.Sprintf(
				"Strongest genre×domain signal: %q games most often bind %q-domain mechanics (%d maps).",
				maxCell.Genre, maxCell.Domain, maxCell.Count,
			))
		}
	}

	if len(insights) > 10 {
		insights = insights[:10]
	}
	return insights
}

func rowsFromMap(counts map[string]int, order []string) []CountRow {
	out := make([]CountRow, 0, len(order))
	for _, label := range order {
		if c, ok := counts[label]; ok && c > 0 {
			out = append(out, CountRow{Label: label, Count: c})
		}
	}
	return out
}

func topRows(counts map[string]int, limit int) []CountRow {
	labels := topLabels(counts, limit)
	out := make([]CountRow, len(labels))
	for i, l := range labels {
		out[i] = CountRow{Label: l, Count: counts[l]}
	}
	return out
}

func topLabels(counts map[string]int, limit int) []string {
	type kv struct {
		k string
		v int
	}
	var items []kv
	for k, v := range counts {
		items = append(items, kv{k, v})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].v != items[j].v {
			return items[i].v > items[j].v
		}
		return items[i].k < items[j].k
	})
	if len(items) > limit {
		items = items[:limit]
	}
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.k
	}
	return out
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
