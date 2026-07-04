package internal

import (
	"sort"
	"strings"
)

type MechanicCooccurrence struct {
	MechanicA string
	MechanicB string
	Count     int
}

// MapsWithMechanic returns map slugs that bind the given mechanic.
func (b *Bundle) MapsWithMechanic(mechanicSlug string) []string {
	var out []string
	for slug, m := range b.Maps {
		for _, binding := range m.Mechanics {
			if binding.MechanicSlug == mechanicSlug {
				out = append(out, slug)
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

// MapsByQualityTier returns map slugs filtered by metadata quality tier.
func (b *Bundle) MapsByQualityTier(tier QualityTier) []string {
	var out []string
	for slug, m := range b.Maps {
		t := DetectQualityTier(m)
		if t == tier {
			out = append(out, slug)
		}
	}
	sort.Strings(out)
	return out
}

// MapsByGenre returns map slugs whose subject.genres contain the label (case-insensitive).
func (b *Bundle) MapsByGenre(genre string) []string {
	genre = stringsToLower(genre)
	var out []string
	for slug, m := range b.Maps {
		for _, g := range m.Subject.Genres {
			if stringsToLower(g) == genre {
				out = append(out, slug)
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

func stringsToLower(s string) string {
	return strings.ToLower(s)
}

// MechanicCooccurrenceMatrix counts how often mechanic pairs appear on the same map.
func (b *Bundle) MechanicCooccurrenceMatrix(minCount int) []MechanicCooccurrence {
	counts := map[[2]string]int{}
	for _, m := range b.Maps {
		slugs := make([]string, 0, len(m.Mechanics))
		seen := map[string]struct{}{}
		for _, bnd := range m.Mechanics {
			if _, ok := seen[bnd.MechanicSlug]; ok {
				continue
			}
			seen[bnd.MechanicSlug] = struct{}{}
			slugs = append(slugs, bnd.MechanicSlug)
		}
		sort.Strings(slugs)
		for i := 0; i < len(slugs); i++ {
			for j := i + 1; j < len(slugs); j++ {
				key := [2]string{slugs[i], slugs[j]}
				counts[key]++
			}
		}
	}
	var out []MechanicCooccurrence
	for key, n := range counts {
		if n < minCount {
			continue
		}
		out = append(out, MechanicCooccurrence{MechanicA: key[0], MechanicB: key[1], Count: n})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		if out[i].MechanicA != out[j].MechanicA {
			return out[i].MechanicA < out[j].MechanicA
		}
		return out[i].MechanicB < out[j].MechanicB
	})
	return out
}

// SimilarMaps finds maps sharing signature mechanics with the reference map.
func (b *Bundle) SimilarMaps(mapSlug string, minShared int) []string {
	ref, ok := b.Maps[mapSlug]
	if !ok {
		return nil
	}
	sigSet := make(map[string]struct{}, len(ref.SignatureGameplay))
	for _, s := range ref.SignatureGameplay {
		sigSet[s] = struct{}{}
	}
	type scored struct {
		slug  string
		score int
	}
	var matches []scored
	for slug, m := range b.Maps {
		if slug == mapSlug {
			continue
		}
		shared := 0
		for _, s := range m.SignatureGameplay {
			if _, ok := sigSet[s]; ok {
				shared++
			}
		}
		if shared >= minShared {
			matches = append(matches, scored{slug: slug, score: shared})
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].score != matches[j].score {
			return matches[i].score > matches[j].score
		}
		return matches[i].slug < matches[j].slug
	})
	out := make([]string, len(matches))
	for i, m := range matches {
		out[i] = m.slug
	}
	return out
}

// MapsWithVariable returns map slugs that bind the given variable.
func (b *Bundle) MapsWithVariable(variableSlug string) []string {
	var out []string
	for slug, m := range b.Maps {
		for _, binding := range m.Variables {
			if binding.VariableSlug == variableSlug {
				out = append(out, slug)
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

// MapsWithMenu returns map slugs that bind the given UI menu.
func (b *Bundle) MapsWithMenu(menuSlug string) []string {
	var out []string
	for slug, m := range b.Maps {
		for _, binding := range m.UIMenus {
			if binding.MenuSlug == menuSlug {
				out = append(out, slug)
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

// VariablesForMechanic returns variable slugs related to a mechanic.
func (b *Bundle) VariablesForMechanic(mechanicSlug string) []string {
	seen := map[string]struct{}{}
	for slug, v := range b.Variables {
		for _, m := range v.RelatedMechanics {
			if m == mechanicSlug {
				seen[slug] = struct{}{}
				break
			}
		}
	}
	for _, m := range b.Maps {
		for _, vb := range m.Variables {
			for _, mech := range vb.RelatedMechanics {
				if mech == mechanicSlug {
					seen[vb.VariableSlug] = struct{}{}
				}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// MenusForMechanic returns menu slugs related to a mechanic.
func (b *Bundle) MenusForMechanic(mechanicSlug string) []string {
	seen := map[string]struct{}{}
	for slug, menu := range b.UIMenus {
		for _, m := range menu.RelatedMechanics {
			if m == mechanicSlug {
				seen[slug] = struct{}{}
				break
			}
		}
	}
	for _, m := range b.Maps {
		for _, mb := range m.UIMenus {
			for _, mech := range mb.SupportsMechanics {
				if mech == mechanicSlug {
					seen[mb.MenuSlug] = struct{}{}
				}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// MenusForVariable returns menus that reference a variable.
func (b *Bundle) MenusForVariable(variableSlug string) []string {
	seen := map[string]struct{}{}
	for slug, menu := range b.UIMenus {
		for _, v := range menu.RelatedVariables {
			if v == variableSlug {
				seen[slug] = struct{}{}
				break
			}
		}
	}
	for _, m := range b.Maps {
		for _, mb := range m.UIMenus {
			for _, v := range mb.DisplaysVariables {
				if v == variableSlug {
					seen[mb.MenuSlug] = struct{}{}
				}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
