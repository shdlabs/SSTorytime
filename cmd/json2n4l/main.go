//******************************************************************
//
// json2n4l - Convert JSON files to N4L format
//
// Usage:
//   json2n4l input.json [output.n4l] [options]
//
// Options:
//   -chapter string      Chapter name (default: "JSON Import")
//   -context string      Context tags (comma-separated)
//   -types               Include type annotations
//   -pretty              Pretty print with comments
//   -maxdepth int        Maximum nesting depth (0 = unlimited)
//   -root string         Root node name (default: filename)
//   -arrows string       Arrow style: simple, semantic, bidirectional (default: simple)
//   -aliases             Generate @aliases for bidirectional mode
//   -autocontext         Automatically add context markers
//
//******************************************************************

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/markburgess/SSTorytime/pkg/json2n4l"
)

func main() {
	// Define flags
	chapter := flag.String("chapter", "JSON Import", "Chapter name for the N4L document")
	contextStr := flag.String("context", "", "Context tags (comma-separated)")
	includeTypes := flag.Bool("types", false, "Include type annotations")
	prettyPrint := flag.Bool("pretty", false, "Pretty print with comments")
	maxDepth := flag.Int("maxdepth", 0, "Maximum nesting depth (0 = unlimited)")
	rootName := flag.String("root", "", "Root node name (default: filename)")
	arrowStyle := flag.String("arrows", "simple", "Arrow style: simple, semantic, bidirectional")
	generateAliases := flag.Bool("aliases", false, "Generate @aliases (for bidirectional mode)")
	autoContext := flag.Bool("autocontext", false, "Automatically add context markers")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] input.json [output.n4l]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Convert JSON files to N4L (Narrative for Learning) format.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nArrow Styles:\n")
		fmt.Fprintf(os.Stderr, "  simple        - Use only (contain) arrows (backward compatible)\n")
		fmt.Fprintf(os.Stderr, "  semantic      - Use semantic arrows: (propt) for values, (in-set)/(setof) for arrays\n")
		fmt.Fprintf(os.Stderr, "  bidirectional - Add inverse arrows: (belong), (in-set), and @aliases\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s data.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -pretty -types data.json output.n4l\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -chapter \"API Data\" -context \"api,json\" data.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -arrows semantic -autocontext data.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -arrows bidirectional -aliases data.json output.n4l\n", os.Args[0])
	}

	flag.Parse()

	// Check arguments
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: input file required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	inputFile := args[0]
	outputFile := ""
	if len(args) >= 2 {
		outputFile = args[1]
	}

	// Parse context tags
	var contextTags []string
	if *contextStr != "" {
		contextTags = strings.Split(*contextStr, ",")
		for i := range contextTags {
			contextTags[i] = strings.TrimSpace(contextTags[i])
		}
	}

	// Parse arrow style
	var arrowStyleEnum json2n4l.ArrowStyle
	switch strings.ToLower(*arrowStyle) {
	case "simple":
		arrowStyleEnum = json2n4l.ArrowStyleSimple
	case "semantic":
		arrowStyleEnum = json2n4l.ArrowStyleSemantic
	case "bidirectional":
		arrowStyleEnum = json2n4l.ArrowStyleBidirectional
	default:
		fmt.Fprintf(os.Stderr, "Error: invalid arrow style '%s'. Valid options: simple, semantic, bidirectional\n", *arrowStyle)
		os.Exit(1)
	}

	// Create config
	config := json2n4l.Config{
		InputFile:       inputFile,
		OutputFile:      outputFile,
		RootName:        *rootName,
		ChapterName:     *chapter,
		ContextTags:     contextTags,
		IncludeTypes:    *includeTypes,
		MaxDepth:        *maxDepth,
		PrettyPrint:     *prettyPrint,
		ArrowStyle:      arrowStyleEnum,
		GenerateAliases: *generateAliases,
		AutoContext:     *autoContext,
	}

	// Create parser (using proper CS terminology)
	parser := json2n4l.NewParser(config)

	// Parse JSON and generate N4L
	fmt.Printf("Converting %s to N4L format...\n", inputFile)

	if err := parser.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get stats and config
	stats := parser.GetStats()
	finalConfig := parser.GetConfig()

	fmt.Printf("✓ Successfully converted to %s\n", finalConfig.OutputFile)
	fmt.Printf("  Nodes created: %d\n", stats["nodeCount"])
	fmt.Printf("  Output size: %d bytes\n", stats["outputSize"])

	if config.MaxDepth > 0 {
		fmt.Printf("  Max depth limit: %d\n", config.MaxDepth)
	}

	// Show arrow style info
	arrowStyleName := "simple"
	switch finalConfig.ArrowStyle {
	case json2n4l.ArrowStyleSemantic:
		arrowStyleName = "semantic"
	case json2n4l.ArrowStyleBidirectional:
		arrowStyleName = "bidirectional"
	}
	fmt.Printf("  Arrow style: %s\n", arrowStyleName)

	if finalConfig.GenerateAliases {
		fmt.Printf("  Generated @aliases\n")
	}
	if finalConfig.AutoContext {
		fmt.Printf("  Auto-context enabled\n")
	}
}
