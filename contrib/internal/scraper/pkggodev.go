// Package scraper provides documentation scrapers for various sources
package scraper

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Package represents extracted Go package documentation
type Package struct {
	Name       string     `json:"name"`
	ImportPath string     `json:"import_path"`
	Synopsis   string     `json:"synopsis"`
	Doc        string     `json:"doc"`
	Functions  []Function `json:"functions,omitempty"`
	Types      []Type     `json:"types,omitempty"`
	Constants  []Constant `json:"constants,omitempty"`
	Variables  []Variable `json:"variables,omitempty"`
	Examples   []Example  `json:"examples,omitempty"`
	Imports    []string   `json:"imports,omitempty"`
}

// Function represents a function or method
type Function struct {
	Name      string `json:"name"`
	Signature string `json:"signature"`
	Doc       string `json:"doc"`
	Receiver  string `json:"receiver,omitempty"`
}

// Type represents a type definition
type Type struct {
	Name    string     `json:"name"`
	Kind    string     `json:"kind"` // struct, interface, alias, etc.
	Doc     string     `json:"doc"`
	Methods []Function `json:"methods,omitempty"`
}

// Constant represents a constant declaration
type Constant struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
	Doc   string `json:"doc"`
}

// Variable represents a variable declaration
type Variable struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
	Doc  string `json:"doc"`
}

// Example represents a code example
type Example struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Doc  string `json:"doc"`
}

// Scraper defines the interface for documentation scrapers
type Scraper interface {
	Scrape(ctx context.Context, url string) (*Package, error)
}

// PkgGoDev scrapes documentation from pkg.go.dev
type PkgGoDev struct {
	client  *http.Client
	verbose bool
}

// NewPkgGoDev creates a new pkg.go.dev scraper
func NewPkgGoDev(verbose bool) *PkgGoDev {
	return &PkgGoDev{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		verbose: verbose,
	}
}

// Scrape extracts package documentation from pkg.go.dev
func (s *PkgGoDev) Scrape(ctx context.Context, url string) (*Package, error) {
	if s.verbose {
		fmt.Printf("📡 Scraping %s\n", url)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	pkg := &Package{
		ImportPath: s.extractImportPath(doc, url),
	}
	pkg.Name = filepath.Base(pkg.ImportPath)
	pkg.Synopsis = s.extractSynopsis(doc)
	pkg.Doc = s.extractPackageDoc(doc)
	pkg.Functions = s.extractFunctions(doc)
	pkg.Types = s.extractTypes(doc)
	pkg.Constants = s.extractConstants(doc)
	pkg.Variables = s.extractVariables(doc)
	pkg.Examples = s.extractExamples(doc)
	pkg.Imports = s.extractImports(doc)

	if s.verbose {
		fmt.Printf("✓ %s: %d functions, %d types, %d examples\n",
			pkg.Name, len(pkg.Functions), len(pkg.Types), len(pkg.Examples))
	}

	return pkg, nil
}

func (s *PkgGoDev) extractImportPath(doc *goquery.Document, url string) string {
	// Best approach: extract from URL directly
	// This avoids picking up breadcrumb navigation text
	parts := strings.Split(url, "/")

	// Find where the package path starts
	for i, part := range parts {
		if part == "std" || (i > 0 && parts[i-1] == "localhost:6060") || (i > 0 && parts[i-1] == "pkg.go.dev") {
			if i+1 < len(parts) {
				return strings.Join(parts[i+1:], "/")
			}
		}
	}

	// Simpler fallback: take everything after domain
	if len(parts) > 3 {
		path := strings.Join(parts[3:], "/")
		// Remove any query parameters
		if idx := strings.Index(path, "?"); idx >= 0 {
			path = path[:idx]
		}
		return path
	}

	return "unknown"
}

func (s *PkgGoDev) extractSynopsis(doc *goquery.Document) string {
	selectors := []string{
		".Documentation-synopsis",
		"[data-test-id='UnitHeader-synopsis']",
		".UnitMeta-synopsis",
	}

	for _, sel := range selectors {
		if text := strings.TrimSpace(doc.Find(sel).Text()); text != "" {
			return text
		}
	}

	return ""
}

func (s *PkgGoDev) extractPackageDoc(doc *goquery.Document) string {
	var paragraphs []string

	doc.Find(".Documentation-overview p, [data-test-id='Documentation-overview'] p").Each(func(i int, sel *goquery.Selection) {
		if text := strings.TrimSpace(sel.Text()); text != "" && i < 3 {
			paragraphs = append(paragraphs, text)
		}
	})

	return strings.Join(paragraphs, "\n\n")
}

func (s *PkgGoDev) extractFunctions(doc *goquery.Document) []Function {
	var functions []Function

	doc.Find(".Documentation-function, [data-test-id='function']").Each(func(i int, sel *goquery.Selection) {
		name := s.extractName(sel, "func ")
		if name == "" {
			return
		}

		functions = append(functions, Function{
			Name:      name,
			Signature: strings.TrimSpace(sel.Find("pre").First().Text()),
			Doc:       s.extractFirstParagraph(sel),
		})
	})

	return functions
}

func (s *PkgGoDev) extractTypes(doc *goquery.Document) []Type {
	var types []Type

	doc.Find(".Documentation-type, [data-test-id='type']").Each(func(i int, sel *goquery.Selection) {
		name := s.extractName(sel, "type ")
		if name == "" {
			return
		}

		sig := strings.TrimSpace(sel.Find("pre").First().Text())
		kind := "other"
		switch {
		case strings.Contains(sig, "struct"):
			kind = "struct"
		case strings.Contains(sig, "interface"):
			kind = "interface"
		case strings.Contains(sig, "="):
			kind = "alias"
		}

		types = append(types, Type{
			Name: name,
			Kind: kind,
			Doc:  s.extractFirstParagraph(sel),
		})
	})

	return types
}

func (s *PkgGoDev) extractConstants(doc *goquery.Document) []Constant {
	var constants []Constant

	doc.Find(".Documentation-constant, [data-test-id='constant']").Each(func(i int, sel *goquery.Selection) {
		name := s.extractName(sel, "const ")
		if name != "" {
			constants = append(constants, Constant{
				Name: name,
				Doc:  s.extractFirstParagraph(sel),
			})
		}
	})

	return constants
}

func (s *PkgGoDev) extractVariables(doc *goquery.Document) []Variable {
	var variables []Variable

	doc.Find(".Documentation-variable, [data-test-id='variable']").Each(func(i int, sel *goquery.Selection) {
		name := s.extractName(sel, "var ")
		if name != "" {
			variables = append(variables, Variable{
				Name: name,
				Doc:  s.extractFirstParagraph(sel),
			})
		}
	})

	return variables
}

func (s *PkgGoDev) extractExamples(doc *goquery.Document) []Example {
	var examples []Example

	doc.Find(".Documentation-example, [data-test-id='example']").Each(func(i int, sel *goquery.Selection) {
		name := strings.TrimSpace(sel.Find("h4").Text())
		if name == "" {
			return
		}

		code := strings.TrimSpace(sel.Find("pre").Text())
		if len(code) > 300 {
			code = code[:300] + "\n// ..."
		}

		examples = append(examples, Example{
			Name: name,
			Code: code,
			Doc:  s.extractFirstParagraph(sel),
		})
	})

	return examples
}

func (s *PkgGoDev) extractImports(doc *goquery.Document) []string {
	seen := make(map[string]bool)
	var imports []string

	doc.Find("pre").Each(func(i int, sel *goquery.Selection) {
		code := sel.Text()
		for _, line := range strings.Split(code, "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "import ") {
				continue
			}

			line = strings.TrimPrefix(line, "import ")
			line = strings.Trim(line, "\"()")
			if line != "" && !seen[line] {
				imports = append(imports, line)
				seen[line] = true
			}
		}
	})

	return imports
}

func (s *PkgGoDev) extractName(sel *goquery.Selection, prefix string) string {
	name := strings.TrimSpace(sel.Find("h4").Text())
	if !strings.HasPrefix(name, prefix) {
		return strings.TrimSpace(name)
	}

	name = strings.TrimPrefix(name, prefix)
	if idx := strings.IndexAny(name, "( "); idx > 0 {
		name = name[:idx]
	}

	return strings.TrimSpace(name)
}

func (s *PkgGoDev) extractFirstParagraph(sel *goquery.Selection) string {
	text := strings.TrimSpace(sel.Find("p").First().Text())
	if len(text) > 500 {
		text = text[:500] + "..."
	}
	return text
}
