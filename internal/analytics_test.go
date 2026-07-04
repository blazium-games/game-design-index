package internal_test

import (
	"testing"

	"github.com/blazium-games/game-design-index/internal"
)

func TestBuildAnalytics(t *testing.T) {
	root := moduleRoot(t)
	b, err := internal.LoadDir(root)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	snap := internal.BuildAnalytics(b)

	if snap.SchemaVersion == "" {
		t.Error("expected schema_version")
	}
	if snap.Overview.GameCount < 1000 {
		t.Errorf("expected many games, got %d", snap.Overview.GameCount)
	}
	if len(snap.QualityTiers) == 0 {
		t.Error("expected quality_tiers")
	}
	if len(snap.TopGenres) == 0 {
		t.Error("expected top_genres")
	}
	if len(snap.MechanicDomains) == 0 {
		t.Error("expected mechanic_domains")
	}
	if len(snap.TopMechanics) == 0 {
		t.Error("expected top_mechanics")
	}
	if len(snap.GenreDomainHeatmap.Cells) == 0 {
		t.Error("expected genre_domain_heatmap cells")
	}
	if len(snap.Cooccurrence) == 0 {
		t.Error("expected cooccurrence")
	}
	for _, c := range snap.Cooccurrence {
		if c.Lift <= 0 {
			t.Errorf("expected positive lift for %s+%s", c.MechanicA, c.MechanicB)
		}
		break
	}
	if len(snap.Insights) < 3 {
		t.Errorf("expected at least 3 insights, got %d", len(snap.Insights))
	}
	if snap.Overview.VariableCount > 0 && snap.VariableStats.EnrichmentComplete+snap.VariableStats.EnrichmentNeedsInfo != snap.Overview.VariableCount {
		t.Error("variable enrichment counts should sum to total")
	}
}
