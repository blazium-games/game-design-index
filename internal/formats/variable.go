package formats

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/blazium-games/game-mechanics-index/internal"
)

func VariableMarkdown(v internal.GameVariable) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", v.Name)
	fmt.Fprintf(&b, "**Category:** %s · **Scope:** %s · **Value kind:** %s\n\n", v.Category, v.Scope, v.ValueKind)
	if v.Summary != "" {
		b.WriteString("## Summary\n\n")
		b.WriteString(v.Summary)
		b.WriteString("\n\n")
	}
	if v.SharedRationale != "" {
		b.WriteString("## Shared Rationale\n\n")
		b.WriteString(v.SharedRationale)
		b.WriteString("\n\n")
	}
	if v.PlayerFocus != "" {
		b.WriteString("## Player Focus\n\n")
		b.WriteString(v.PlayerFocus)
		b.WriteString("\n\n")
	}
	if v.TypicalRange != "" {
		b.WriteString("## Typical Range\n\n")
		b.WriteString(v.TypicalRange)
		b.WriteString("\n\n")
	}
	if len(v.RelatedMechanics) > 0 {
		b.WriteString("## Related Mechanics\n\n")
		for _, m := range v.RelatedMechanics {
			fmt.Fprintf(&b, "- %s\n", m)
		}
		b.WriteString("\n")
	}
	if len(v.FeaturedIn) > 0 {
		b.WriteString("## Featured In\n\n")
		for _, s := range v.FeaturedIn {
			fmt.Fprintf(&b, "- %s\n", s)
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func VariableText(v internal.GameVariable) string {
	return strings.ReplaceAll(VariableMarkdown(v), "## ", "")
}

type variableXML struct {
	XMLName   xml.Name `xml:"game_variable"`
	Slug      string   `xml:"slug"`
	Name      string   `xml:"name"`
	Category  string   `xml:"category"`
	Scope     string   `xml:"scope"`
	ValueKind string   `xml:"value_kind"`
	Summary   string   `xml:"summary,omitempty"`
}

func VariableXML(v internal.GameVariable) ([]byte, error) {
	root := variableXML{
		Slug:      v.Slug,
		Name:      v.Name,
		Category:  string(v.Category),
		Scope:     string(v.Scope),
		ValueKind: string(v.ValueKind),
		Summary:   v.Summary,
	}
	out, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), out...), nil
}
