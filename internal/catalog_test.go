package internal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/blazium-games/game-design-index/internal"
)

func moduleRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "..", "data", "source"))
}

func TestLoadDirVariablesAndMenus(t *testing.T) {
	root := moduleRoot(t)
	b, err := internal.LoadDir(root)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	if got := len(b.VariableSlugs()); got < 40 {
		t.Errorf("expected at least 40 variables, got %d", got)
	}
	if got := len(b.UIMenuSlugs()); got < 15 {
		t.Errorf("expected at least 15 ui menus, got %d", got)
	}
}

func TestPilotMapVariableBindings(t *testing.T) {
	root := moduleRoot(t)
	b, err := internal.LoadDir(root)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	m, ok := b.Map("hollow-knight")
	if !ok {
		t.Fatal("hollow-knight not found")
	}
	if len(m.Variables) == 0 {
		t.Error("hollow-knight expected variable bindings")
	}
	if len(m.UIMenus) == 0 {
		t.Error("hollow-knight expected ui menu bindings")
	}
}

func TestVariablesForMechanic(t *testing.T) {
	root := moduleRoot(t)
	b, err := internal.LoadDir(root)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	vars := b.VariablesForMechanic("depletable-health-pool")
	if len(vars) == 0 {
		t.Error("expected variables for depletable-health-pool")
	}
}
