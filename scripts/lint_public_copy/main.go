package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var scanTerms = []string{
	"godot",
	"ludonaut",
	"ludens",
	"retrododo",
	"source_url",
	"metacritic.com",
	"gamefaqs.gamespot.com",
	"imdb.com/list",
}

var skipDirs = map[string]bool{
	"node_modules": true,
	".git":           true,
	"dist":           true,
}

func main() {
	root := flag.String("root", ".", "repo root to scan")
	flag.Parse()

	abs, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	targets := []string{
		filepath.Join(abs, "data", "dist"),
		filepath.Join(abs, "docs"),
		filepath.Join(abs, "site", "src"),
	}

	var violations []string
	for _, base := range targets {
		if _, err := os.Stat(base); os.IsNotExist(err) {
			continue
		}
		_ = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				if info != nil && info.IsDir() && skipDirs[info.Name()] {
					return filepath.SkipDir
				}
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".json", ".md", ".ts", ".tsx", ".css", ".html":
			default:
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			lower := strings.ToLower(string(data))
			rel, _ := filepath.Rel(abs, path)
			for _, term := range scanTerms {
				if strings.Contains(lower, term) {
					violations = append(violations, fmt.Sprintf("%s: contains %q", rel, term))
				}
			}
			return nil
		})
	}

	if len(violations) > 0 {
		fmt.Println("Public content policy violations:")
		for _, v := range violations {
			fmt.Println("  -", v)
		}
		os.Exit(1)
	}
	fmt.Println("OK: no banned terms in public paths")
}
