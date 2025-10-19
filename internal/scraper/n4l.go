// Package scraper provides N4L output formatters
package scraper

import (
	"fmt"
	"io"
	"strings"
)

// N4LFormatter writes Package data in N4L format
type N4LFormatter struct {
	ChapterName string
	ContextTags []string
}

// WriteN4L outputs package data in N4L format using bulletproof structure
// Goal: ALWAYS uploadable, user can refine later
// Uses tabs for indentation, arrows for relationships, @sepN for unique section IDs
func (f *N4LFormatter) WriteN4L(pkg *Package, w io.Writer) error {
	// Chapter header
	chapterName := f.ChapterName
	if chapterName == "" {
		chapterName = pkg.Name + " package"
	}
	fmt.Fprintf(w, "- %s\n\n", chapterName)

	// Context tags
	tags := f.ContextTags
	if len(tags) == 0 {
		tags = []string{"golang", "package", pkg.Name}
	}
	for _, tag := range tags {
		fmt.Fprintf(w, "+ %s\n", tag)
	}
	fmt.Fprintln(w)

	sepCounter := 0 // Track @sep numbers for unique IDs

	// Package root node
	fmt.Fprintf(w, " %s\n", escapeN4L(pkg.Name))

	// Import path with arrow - use a labeled node to avoid self-reference
	if pkg.ImportPath != "" {
		fmt.Fprintf(w, "      \" (source) %s\n", escapeN4L("import-path: "+pkg.ImportPath))
	}

	// Description with arrow - ditto before arrow
	if pkg.Synopsis != "" {
		fmt.Fprintf(w, "      \" (description) %s\n", escapeN4L(pkg.Synopsis))
	} else if pkg.Doc != "" {
		paras := strings.Split(pkg.Doc, "\n\n")
		if len(paras) > 0 && len(paras[0]) > 0 {
			fmt.Fprintf(w, "      \" (description) %s\n", escapeN4L(paras[0]))
		}
	}

	// Functions section with unique @sep ID (on same line!)
	if len(pkg.Functions) > 0 {
		sepCounter++
		fmt.Fprintf(w, "\n@sep%d %s\n\n", sepCounter, escapeN4L("Functions in "+pkg.Name))
		for i := range pkg.Functions {
			f.writeFunction(w, &pkg.Functions[i])
		}
	}

	// Types section with unique @sep ID (on same line!)
	if len(pkg.Types) > 0 {
		sepCounter++
		fmt.Fprintf(w, "\n@sep%d %s\n\n", sepCounter, escapeN4L("Types in "+pkg.Name))
		for i := range pkg.Types {
			f.writeType(w, &pkg.Types[i])
		}
	}

	// Constants section with unique @sep ID (on same line!)
	if len(pkg.Constants) > 0 {
		sepCounter++
		fmt.Fprintf(w, "\n@sep%d %s\n\n", sepCounter, escapeN4L("Constants in "+pkg.Name))
		for i := range pkg.Constants {
			c := &pkg.Constants[i]
			fmt.Fprintf(w, "      %s\n", escapeN4L(c.Name))
			if c.Doc != "" && len(c.Doc) < 150 {
				// Arrow with ditto (child of child = 11 spaces)
				fmt.Fprintf(w, "           \" (description) %s\n", escapeN4L(c.Doc))
			}
		}
	}

	fmt.Fprintln(w)
	return nil
}

func (f *N4LFormatter) writeFunction(w io.Writer, fn *Function) {
	// Function name as child node (6 spaces)
	fmt.Fprintf(w, "      %s\n", escapeN4L(fn.Name))

	// Signature with arrow: ditto + (def) + signature (11 spaces for grandchild)
	if fn.Signature != "" && len(fn.Signature) < 150 {
		sig := escapeN4L(fn.Signature)
		fmt.Fprintf(w, "           \" (def) %s\n", sig)
	}

	// Documentation with description arrow (11 spaces)
	if fn.Doc != "" && len(fn.Doc) < 150 {
		doc := escapeN4L(fn.Doc)
		fmt.Fprintf(w, "           \" (description) %s\n", doc)
	}
}

func (f *N4LFormatter) writeType(w io.Writer, t *Type) {
	// Type name as child node (6 spaces)
	fmt.Fprintf(w, "      %s\n", escapeN4L(t.Name))

	// Kind with description arrow (11 spaces)
	if t.Kind != "" {
		kind := escapeN4L(t.Kind)
		fmt.Fprintf(w, "           \" (description) %s\n", kind)
	}

	// Documentation with description arrow (11 spaces)
	if t.Doc != "" && len(t.Doc) < 150 {
		doc := escapeN4L(t.Doc)
		fmt.Fprintf(w, "           \" (description) %s\n", doc)
	}
}

// escapeN4L prepares text for N4L format
// Only quotes when necessary (special chars, annotations, etc)
func escapeN4L(s string) string {
	if len(s) > 500 {
		s = s[:500] + "..."
	}

	// Normalize whitespace
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", " ")

	// Collapse multiple spaces
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}

	s = strings.TrimSpace(s)

	// Never return empty string
	if s == "" {
		return "no-description"
	}

	// Check if we need quotes
	needsQuotes := false

	// Need quotes if contains special annotation symbols
	if strings.ContainsAny(s, "%=<>*\"") {
		needsQuotes = true
	}

	// Need quotes if contains parentheses (might conflict with arrows)
	if strings.Contains(s, "(") || strings.Contains(s, ")") {
		needsQuotes = true
	}

	// Need quotes if starts with special chars
	if len(s) > 0 && (s[0] == '@' || s[0] == '+' || s[0] == '-' || s[0] == '#') {
		needsQuotes = true
	}

	// If we need quotes, clean up the content
	if needsQuotes {
		// Replace internal quotes with single quotes
		s = strings.ReplaceAll(s, `"`, "'")

		// Remove annotation symbols - order matters!
		s = strings.ReplaceAll(s, "**", " ")
		s = strings.ReplaceAll(s, ">>", " RSHIFT ")
		s = strings.ReplaceAll(s, "==", " equals ")
		s = strings.ReplaceAll(s, "!=", " not-equals ")
		s = strings.ReplaceAll(s, ">=", " GTE ")
		s = strings.ReplaceAll(s, "<=", " LTE ")
		s = strings.ReplaceAll(s, "=>", " implies ")
		s = strings.ReplaceAll(s, "=", " equals ")
		s = strings.ReplaceAll(s, ">", " GT ")
		s = strings.ReplaceAll(s, "<", " LT ")
		s = strings.ReplaceAll(s, "%", " PERCENT ")

		return `"` + s + `"`
	}

	// Return plain text without quotes
	return s
}
