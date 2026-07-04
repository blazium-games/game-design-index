package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blazium-games/game-design-index/internal"
)

func main() {
	source := flag.String("source", "data/source", "corpus root (maps, library, schema)")
	out := flag.String("out", "data/dist", "output directory")
	version := flag.String("version", "0.1.0", "release version stamp")
	flag.Parse()

	root, err := filepath.Abs(*source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	outDir, err := filepath.Abs(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	b, err := internal.LoadDir(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load corpus: %v\n", err)
		os.Exit(1)
	}

	if err := exportAPI(b, outDir, *version); err != nil {
		fmt.Fprintf(os.Stderr, "export: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Exported %d maps, %d mechanics -> %s/api/v1 and %s/formats/v1\n", len(b.Maps), len(b.Mechanics), outDir, outDir)
}
