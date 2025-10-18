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

// WriteN4L outputs package data in N4L format
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

	// Package node
	fmt.Fprintf(w, " %s\n", escapeN4L(pkg.Name))

	// Import path
	if pkg.ImportPath != "" {
		fmt.Fprintln(w, "      \" (contain) \"import path\"")
		fmt.Fprintf(w, "           \" (contain) \"path: %s\"\n", pkg.ImportPath)
	}

	// Synopsis
	if pkg.Synopsis != "" {
		fmt.Fprintln(w, "      \" (contain) synopsis")
		fmt.Fprintf(w, "           \" (contain) %s\n", escapeN4L(pkg.Synopsis))
	}

	// Documentation
	if pkg.Doc != "" {
		fmt.Fprintln(w, "      \" (contain) documentation")
		for i, para := range strings.Split(pkg.Doc, "\n\n") {
			if i >= 2 {
				break
			}
			fmt.Fprintf(w, "           \" (contain) %s\n", escapeN4L(para))
		}
	}

	// Functions
	if len(pkg.Functions) > 0 {
		fmt.Fprintln(w, "      \" (contain) functions")
		for i := range pkg.Functions {
			f.writeFunction(w, &pkg.Functions[i])
		}
	}

	// Types
	if len(pkg.Types) > 0 {
		fmt.Fprintln(w, "      \" (contain) types")
		for i := range pkg.Types {
			f.writeType(w, &pkg.Types[i])
		}
	}

	// Constants
	if len(pkg.Constants) > 0 {
		fmt.Fprintln(w, "      \" (contain) constants")
		for i := range pkg.Constants {
			fmt.Fprintf(w, "           \" (contain) %s\n", escapeN4L(pkg.Constants[i].Name))
			if pkg.Constants[i].Doc != "" {
				fmt.Fprintf(w, "                \" (contain) %s\n", escapeN4L(pkg.Constants[i].Doc))
			}
		}
	}

	// Examples
	if len(pkg.Examples) > 0 {
		fmt.Fprintln(w, "      \" (contain) examples")
		for i := range pkg.Examples {
			f.writeExample(w, &pkg.Examples[i])
		}
	}

	// Imports
	if len(pkg.Imports) > 0 {
		fmt.Fprintln(w, "      \" (contain) \"imported packages\"")
		for _, imp := range pkg.Imports {
			fmt.Fprintf(w, "           \" (contain) %s\n", escapeN4L(imp))
		}
	}

	return nil
}

func (f *N4LFormatter) writeFunction(w io.Writer, fn *Function) {
	fmt.Fprintf(w, "           \" (contain) %s\n", escapeN4L(fn.Name))

	if fn.Signature != "" {
		fmt.Fprintln(w, "                \" (contain) signature")
		fmt.Fprintf(w, "                     \" (contain) %s\n", escapeN4L(fn.Signature))
	}

	if fn.Doc != "" {
		fmt.Fprintln(w, "                \" (contain) description")
		fmt.Fprintf(w, "                     \" (contain) %s\n", escapeN4L(fn.Doc))
	}
}

func (f *N4LFormatter) writeType(w io.Writer, t *Type) {
	fmt.Fprintf(w, "           \" (contain) %s\n", escapeN4L(t.Name))

	if t.Kind != "" {
		fmt.Fprintln(w, "                \" (contain) kind")
		fmt.Fprintf(w, "                     \" (contain) %s\n", escapeN4L(t.Kind))
	}

	if t.Doc != "" {
		fmt.Fprintln(w, "                \" (contain) description")
		fmt.Fprintf(w, "                     \" (contain) %s\n", escapeN4L(t.Doc))
	}
}

func (f *N4LFormatter) writeExample(w io.Writer, ex *Example) {
	fmt.Fprintf(w, "           \" (contain) %s\n", escapeN4L(ex.Name))

	if ex.Code != "" {
		fmt.Fprintln(w, "                \" (contain) code")
		fmt.Fprintf(w, "                     \" (contain) %s\n", escapeN4L(ex.Code))
	}
}

// escapeN4L escapes text for N4L format
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

	// Never return empty string - use placeholder
	if s == "" {
		return "\"(no description)\""
	}

	// Replace == with = to avoid N4L arrow parsing issues
	s = strings.ReplaceAll(s, "==", "=")

	// Quote if needed
	if strings.ContainsAny(s, " \t()[]{}") {
		s = `"` + strings.ReplaceAll(s, `"`, "'") + `"`
	}

	return s
}
