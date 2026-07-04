package formats

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func ToYAML(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func CatalogMarkdown(schemaVersion, releaseVersion string, mapCount, mechCount int) string {
	var b strings.Builder
	b.WriteString("# Game Design Index Catalog\n\n")
	fmt.Fprintf(&b, "- Schema: %s\n", schemaVersion)
	fmt.Fprintf(&b, "- Release: %s\n", releaseVersion)
	fmt.Fprintf(&b, "- Maps: %d\n", mapCount)
	fmt.Fprintf(&b, "- Mechanics: %d\n", mechCount)
	return b.String()
}

func IndexMarkdown(title string, rows []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", title)
	for _, r := range rows {
		fmt.Fprintf(&b, "- %s\n", r)
	}
	return b.String()
}
