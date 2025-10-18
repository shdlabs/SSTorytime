package main

import (
	"fmt"
	"log"

	"github.com/markburgess/SSTorytime/pkg/json2n4l"
)

func main() {
	inputFile := "example_semantic.json"

	// Example 1: Simple style (default - backward compatible)
	fmt.Println("=== Example 1: Simple Style (default) ===")
	parser1 := json2n4l.NewParser(json2n4l.Config{
		InputFile:   inputFile,
		OutputFile:  "example_simple.n4l",
		ChapterName: "Simple Style Example",
		PrettyPrint: true,
	})
	if err := parser1.Parse(); err != nil {
		log.Fatalf("Simple style failed: %v", err)
	}
	fmt.Println("Generated: example_simple.n4l")

	// Example 2: Semantic style with auto-context
	fmt.Println("\n=== Example 2: Semantic Style with Auto-Context ===")
	parser2 := json2n4l.NewParser(json2n4l.Config{
		InputFile:   inputFile,
		OutputFile:  "example_semantic.n4l",
		ChapterName: "Semantic Style Example",
		PrettyPrint: true,
		ArrowStyle:  json2n4l.ArrowStyleSemantic,
		AutoContext: true,
	})
	if err := parser2.Parse(); err != nil {
		log.Fatalf("Semantic style failed: %v", err)
	}
	fmt.Println("Generated: example_semantic.n4l")

	// Example 3: Bidirectional style with aliases
	fmt.Println("\n=== Example 3: Bidirectional Style with Aliases ===")
	parser3 := json2n4l.NewParser(json2n4l.Config{
		InputFile:         inputFile,
		OutputFile:        "example_bidirectional.n4l",
		ChapterName:       "Bidirectional Style Example",
		PrettyPrint:       true,
		ArrowStyle:        json2n4l.ArrowStyleBidirectional,
		AutoContext:       true,
		GenerateAliases:   true,
		PreserveStructure: true,
	})
	if err := parser3.Parse(); err != nil {
		log.Fatalf("Bidirectional style failed: %v", err)
	}
	fmt.Println("Generated: example_bidirectional.n4l")

	// Example 4: Using convenience functions
	fmt.Println("\n=== Example 4: Using Option Functions ===")
	if err := json2n4l.ParseFile(
		inputFile,
		"example_options.n4l",
		json2n4l.WithChapter("Option Functions Example"),
		json2n4l.WithSemantic(),
		json2n4l.WithAutoContext(),
		json2n4l.WithAliases(),
		json2n4l.WithPrettyPrint(),
		json2n4l.WithTypes(),
	); err != nil {
		log.Fatalf("Option functions failed: %v", err)
	}
	fmt.Println("Generated: example_options.n4l")

	fmt.Println("\n=== All examples generated successfully! ===")
	fmt.Println("Compare the outputs to see the different semantic arrow styles.")
}
