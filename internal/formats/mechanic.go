package formats

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/blazium-games/game-mechanics-index/internal"
)

func MechanicMarkdown(e internal.MechanicEntry) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", e.Name)
	fmt.Fprintf(&b, "**Flavor:** %s · **Domain:** %s", e.Flavor, e.Domain)
	if e.Complexity != "" {
		fmt.Fprintf(&b, " · **Complexity:** %s", e.Complexity)
	}
	b.WriteString("\n\n")
	if e.Summary != "" {
		b.WriteString("## Description\n\n")
		b.WriteString(e.Summary)
		b.WriteString("\n\n")
	}
	if e.FlavorRationale != "" {
		b.WriteString("## Category Insights\n\n")
		b.WriteString(e.FlavorRationale)
		b.WriteString("\n\n")
	}
	if e.PlayerExperience != "" {
		b.WriteString("## Player Experience\n\n")
		b.WriteString(e.PlayerExperience)
		b.WriteString("\n\n")
	}
	if len(e.FeaturedIn) > 0 {
		b.WriteString("## Featured In\n\n")
		for _, s := range e.FeaturedIn {
			fmt.Fprintf(&b, "- %s\n", s)
		}
		b.WriteString("\n")
	}
	if len(e.CommonIn) > 0 {
		b.WriteString("## Common In\n\n")
		for _, s := range e.CommonIn {
			fmt.Fprintf(&b, "- %s\n", s)
		}
		b.WriteString("\n")
	}
	if len(e.SynergyNotes) > 0 {
		b.WriteString("## Synergies\n\n")
		for _, sn := range e.SynergyNotes {
			if sn.Note != "" {
				fmt.Fprintf(&b, "- **%s**: %s\n", sn.Slug, sn.Note)
			} else {
				fmt.Fprintf(&b, "- %s\n", sn.Slug)
			}
		}
		b.WriteString("\n")
	} else if len(e.Synergies) > 0 {
		b.WriteString("## Synergies\n\n")
		for _, s := range e.Synergies {
			fmt.Fprintf(&b, "- %s\n", s)
		}
		b.WriteString("\n")
	}
	if e.DesignGuidance != nil {
		if e.DesignGuidance.WhenToUse != "" {
			b.WriteString("## When to Use\n\n")
			b.WriteString(e.DesignGuidance.WhenToUse)
			b.WriteString("\n\n")
		}
		if e.DesignGuidance.WhereToUse != "" {
			b.WriteString("## Where to Use\n\n")
			b.WriteString(e.DesignGuidance.WhereToUse)
			b.WriteString("\n\n")
		}
		if len(e.DesignGuidance.WhenToAvoid) > 0 {
			b.WriteString("## When to Avoid\n\n")
			for _, a := range e.DesignGuidance.WhenToAvoid {
				fmt.Fprintf(&b, "- %s\n", a)
			}
			b.WriteString("\n")
		}
	}
	if e.SignatureOf != nil && len(e.SignatureOf.Games) > 0 {
		b.WriteString("## Signature of\n\n")
		for _, g := range e.SignatureOf.Games {
			fmt.Fprintf(&b, "- %s\n", g)
		}
		b.WriteString("\n")
	}
	if len(e.Examples) > 0 {
		b.WriteString("## Examples\n\n")
		for _, ex := range e.Examples {
			label := ex.Label
			if label == "" {
				label = ex.MapSlug
			}
			if ex.Description != "" {
				fmt.Fprintf(&b, "- **%s**: %s\n", label, ex.Description)
			} else {
				fmt.Fprintf(&b, "- %s\n", label)
			}
		}
		b.WriteString("\n")
	}
	if e.AgentContext != nil && e.AgentContext.SummaryForAgents != "" {
		b.WriteString("## Agent Summary\n\n")
		b.WriteString(e.AgentContext.SummaryForAgents)
		b.WriteString("\n\n")
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func MechanicText(e internal.MechanicEntry) string {
	return strings.ReplaceAll(MechanicMarkdown(e), "## ", "")
}

type mechanicXML struct {
	XMLName          xml.Name `xml:"mechanic"`
	Slug             string   `xml:"slug"`
	Name             string   `xml:"name"`
	Flavor           string   `xml:"flavor"`
	Domain           string   `xml:"domain"`
	Summary          string   `xml:"summary,omitempty"`
	FlavorRationale  string   `xml:"flavor_rationale,omitempty"`
	PlayerExperience string   `xml:"player_experience,omitempty"`
	Complexity       string   `xml:"complexity,omitempty"`
}

func MechanicXML(e internal.MechanicEntry) ([]byte, error) {
	root := mechanicXML{
		Slug:             e.Slug,
		Name:             e.Name,
		Flavor:           string(e.Flavor),
		Domain:           string(e.Domain),
		Summary:          e.Summary,
		FlavorRationale:  e.FlavorRationale,
		PlayerExperience: e.PlayerExperience,
		Complexity:       string(e.Complexity),
	}
	out, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), out...), nil
}

func MapMarkdown(m internal.GameplayMap) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", m.Subject.Name)
	if m.Narrative.Description != "" {
		b.WriteString("## Overview\n\n")
		b.WriteString(m.Narrative.Description)
		b.WriteString("\n\n")
	}
	if m.Narrative.CoreLoop != "" {
		b.WriteString("## Core Loop\n\n")
		b.WriteString(m.Narrative.CoreLoop)
		b.WriteString("\n\n")
	}
	if m.GDDOutline != nil {
		if len(m.GDDOutline.PlayerGoals) > 0 {
			b.WriteString("## Player Goals\n\n")
			for _, g := range m.GDDOutline.PlayerGoals {
				fmt.Fprintf(&b, "- %s\n", g)
			}
			b.WriteString("\n")
		}
		if m.GDDOutline.CombatNotes != "" {
			b.WriteString("## Combat Notes\n\n")
			b.WriteString(m.GDDOutline.CombatNotes)
			b.WriteString("\n\n")
		}
	}
	if len(m.SignatureGameplay) > 0 {
		b.WriteString("## Signature Gameplay\n\n")
		for _, s := range m.SignatureGameplay {
			fmt.Fprintf(&b, "- %s\n", s)
		}
		b.WriteString("\n")
	}
	if len(m.Variables) > 0 {
		b.WriteString("## Variable Bindings\n\n")
		for _, vb := range m.Variables {
			fmt.Fprintf(&b, "- **%s** (%s)", vb.VariableSlug, vb.Role)
			if vb.Expression != "" {
				fmt.Fprintf(&b, ": %s", vb.Expression)
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	if len(m.UIMenus) > 0 {
		b.WriteString("## UI Menu Bindings\n\n")
		for _, mb := range m.UIMenus {
			fmt.Fprintf(&b, "- **%s** (%s)", mb.MenuSlug, mb.Role)
			if mb.MapNotes != "" {
				fmt.Fprintf(&b, ": %s", mb.MapNotes)
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func MapText(m internal.GameplayMap) string {
	return strings.ReplaceAll(MapMarkdown(m), "## ", "")
}

type mapXML struct {
	XMLName xml.Name `xml:"gameplay_map"`
	Slug    string   `xml:"slug"`
	Title   string   `xml:"title"`
	Name    string   `xml:"name"`
	MapType string   `xml:"map_type"`
}

func MapXML(m internal.GameplayMap) ([]byte, error) {
	root := mapXML{
		Slug:    m.Slug,
		Title:   m.Title,
		Name:    m.Subject.Name,
		MapType: string(m.MapType),
	}
	out, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), out...), nil
}
