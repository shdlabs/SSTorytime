package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/markburgess/SSTorytime/internal/scraper"
)

func main() {
	// Flags
	output := flag.String("o", "", "output file (default: <package>.n4l)")
	chapter := flag.String("chapter", "", "chapter name (default: <package> package)")
	contexts := flag.String("context", "", "context tags (comma-separated)")
	verbose := flag.Bool("v", false, "verbose output")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(1)
	}

	url := flag.Arg(0)

	// Parse context tags
	var contextTags []string
	if *contexts != "" {
		for _, tag := range strings.Split(*contexts, ",") {
			if tag = strings.TrimSpace(tag); tag != "" {
				contextTags = append(contextTags, tag)
			}
		}
	}

	// Scrape documentation
	s := scraper.NewPkgGoDev(*verbose)
	pkg, err := s.Scrape(context.Background(), url)
	if err != nil {
		log.Fatalf("❌ Scraping failed: %v\n", err)
	}

	// Determine output file
	outputFile := *output
	if outputFile == "" {
		outputFile = pkg.Name + ".n4l"
	}

	// Create output file
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("❌ Creating output: %v\n", err)
	}
	defer f.Close()

	// Write N4L
	formatter := &scraper.N4LFormatter{
		ChapterName: *chapter,
		ContextTags: contextTags,
	}

	if err := formatter.WriteN4L(pkg, f); err != nil {
		log.Fatalf("❌ Writing N4L: %v\n", err)
	}

	// Stats
	fmt.Printf("\n✓ Successfully scraped %s\n", pkg.Name)
	fmt.Printf("  Import path: %s\n", pkg.ImportPath)
	fmt.Printf("  Functions: %d\n", len(pkg.Functions))
	fmt.Printf("  Types: %d\n", len(pkg.Types))
	fmt.Printf("  Constants: %d\n", len(pkg.Constants))
	fmt.Printf("  Examples: %d\n", len(pkg.Examples))
	fmt.Printf("  Output: %s\n", outputFile)

	fmt.Println("\nNext steps:")
	fmt.Printf("  📊 View: cat %s\n", outputFile)
	absPath, _ := filepath.Abs(outputFile)
	relPath, _ := filepath.Rel(filepath.Dir(filepath.Dir(absPath)), absPath)
	fmt.Printf("  📤 Upload: ../../src/N4L -u -force %s\n", relPath)
}

func usage() {
	fmt.Fprintf(os.Stderr, `godoc2n4l - Convert Go package documentation to N4L format

Usage:
  godoc2n4l [options] <package-url>

Options:
  -o <file>          Output file (default: <package>.n4l)
  -chapter <name>    Chapter name (default: <package> package)
  -context <tags>    Context tags, comma-separated (default: golang,package,<name>)
  -v, -verbose       Verbose output

Examples:
  # Scrape the flag package
  godoc2n4l https://pkg.go.dev/flag

  # Custom output and context
  godoc2n4l -o output.n4l -context "golang,stdlib,cli" https://pkg.go.dev/flag

  # Scrape and upload to N4L database
  godoc2n4l https://pkg.go.dev/fmt
  ../../src/N4L -u -force fmt.n4l

`)
}
