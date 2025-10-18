package demo

import (
	"fmt"
	"log"
	"strings"

	"github.com/markburgess/SSTorytime/pkg/godoc2n4l"
)

func main() {
	fmt.Println("=== Testing godoc2n4l Scraper ===\n")

	// Test 1: Scrape the flag package (small and well-structured)
	fmt.Print("Test 1: Scraping flag package...")
	config1 := godoc2n4l.Config{
		PackageURL:  "https://pkg.go.dev/flag",
		OutputFile:  "flag.n4l",
		ChapterName: "flag package - Command-line flag parsing",
		ContextTags: []string{"golang", "stdlib", "cli", "flag"},
		Verbose:     true,
	}

	converter1 := godoc2n4l.NewConverter(config1)
	if err := converter1.Convert(); err != nil {
		log.Printf("Error converting flag package: %v\n", err)
	} else {
		stats := converter1.GetStats()
		fmt.Printf("\nStats for flag package:\n")
		for key, value := range stats {
			fmt.Printf("  %s: %d\n", key, value)
		}
	}

	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")

	// Test 2: Scrape the fmt package (has more examples)
	fmt.Println("Test 2: Scraping fmt package...")
	config2 := godoc2n4l.Config{
		PackageURL:  "https://pkg.go.dev/fmt",
		OutputFile:  "fmt.n4l",
		ChapterName: "fmt package - Formatted I/O",
		ContextTags: []string{"golang", "stdlib", "io", "formatting"},
		Verbose:     true,
	}

	converter2 := godoc2n4l.NewConverter(config2)
	if err := converter2.Convert(); err != nil {
		log.Printf("Error converting fmt package: %v\n", err)
	} else {
		stats := converter2.GetStats()
		fmt.Printf("\nStats for fmt package:\n")
		for key, value := range stats {
			fmt.Printf("  %s: %d\n", key, value)
		}
	}

	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")

	// Test 3: Scrape the http package (larger, more complex)
	fmt.Println("Test 3: Scraping net/http package...")
	config3 := godoc2n4l.Config{
		PackageURL:  "https://pkg.go.dev/net/http",
		OutputFile:  "http.n4l",
		ChapterName: "net/http package - HTTP client and server",
		ContextTags: []string{"golang", "stdlib", "http", "web"},
		Verbose:     true,
	}

	converter3 := godoc2n4l.NewConverter(config3)
	if err := converter3.Convert(); err != nil {
		log.Printf("Error converting http package: %v\n", err)
	} else {
		stats := converter3.GetStats()
		fmt.Printf("\nStats for net/http package:\n")
		for key, value := range stats {
			fmt.Printf("  %s: %d\n", key, value)
		}
	}

	fmt.Println("\n=== All tests complete! ===")
	fmt.Println("\nGenerated files:")
	fmt.Println("  - flag.n4l")
	fmt.Println("  - fmt.n4l")
	fmt.Println("  - http.n4l")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Examine the .n4l files to see what data was extracted")
	fmt.Println("  2. Identify what relationships exist in the real data")
	fmt.Println("  3. Design arrow types based on actual patterns")
	fmt.Println("  4. Upload to N4L database: ../../src/N4L -u -force <file>.n4l")
}
