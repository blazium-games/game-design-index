package formats

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/blazium-games/game-mechanics-index/internal"
)

func MenuMarkdown(m internal.UIMenu) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", m.Name)
	fmt.Fprintf(&b, "**Type:** %s · **Layer:** %s\n\n", m.MenuType, m.Layer)
	if m.Summary != "" {
		b.WriteString("## Summary\n\n")
		b.WriteString(m.Summary)
		b.WriteString("\n\n")
	}
	if m.SharedRationale != "" {
		b.WriteString("## Shared Rationale\n\n")
		b.WriteString(m.SharedRationale)
		b.WriteString("\n\n")
	}
	if len(m.TypicalActions) > 0 {
		b.WriteString("## Typical Actions\n\n")
		for _, a := range m.TypicalActions {
			fmt.Fprintf(&b, "- %s\n", a)
		}
		b.WriteString("\n")
	}
	if len(m.RelatedMechanics) > 0 {
		b.WriteString("## Related Mechanics\n\n")
		for _, mech := range m.RelatedMechanics {
			fmt.Fprintf(&b, "- %s\n", mech)
		}
		b.WriteString("\n")
	}
	if len(m.RelatedVariables) > 0 {
		b.WriteString("## Related Variables\n\n")
		for _, v := range m.RelatedVariables {
			fmt.Fprintf(&b, "- %s\n", v)
		}
		b.WriteString("\n")
	}
	if len(m.FeaturedIn) > 0 {
		b.WriteString("## Featured In\n\n")
		for _, s := range m.FeaturedIn {
			fmt.Fprintf(&b, "- %s\n", s)
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func MenuText(m internal.UIMenu) string {
	return strings.ReplaceAll(MenuMarkdown(m), "## ", "")
}

type menuXML struct {
	XMLName  xml.Name `xml:"ui_menu"`
	Slug     string   `xml:"slug"`
	Name     string   `xml:"name"`
	MenuType string   `xml:"menu_type"`
	Layer    string   `xml:"layer"`
	Summary  string   `xml:"summary,omitempty"`
}

func MenuXML(m internal.UIMenu) ([]byte, error) {
	root := menuXML{
		Slug:     m.Slug,
		Name:     m.Name,
		MenuType: string(m.MenuType),
		Layer:    string(m.Layer),
		Summary:  m.Summary,
	}
	out, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), out...), nil
}
