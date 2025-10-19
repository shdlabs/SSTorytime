package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/markburgess/SSTorytime/contrib/internal/n4l"
)

func main() {
	// Parse command line flags
	watch := flag.Bool("w", false, "Watch mode - re-validate on file changes")
	watchLong := flag.Bool("watch", false, "Watch mode - re-validate on file changes")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file.n4l>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nValidates N4L files without uploading to database\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	watchMode := *watch || *watchLong

	// Get absolute path for watching
	absPath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	if watchMode {
		fmt.Printf("👁️  Watch mode enabled - monitoring %s\n", absPath)
		fmt.Printf("Press Ctrl+C to stop\n\n")
		watchFile(absPath, *verbose)
	} else {
		// Single validation
		validateFile(filename, *verbose)
	}
}

func validateFile(filename string, verbose bool) {
	fmt.Printf("🔍 Validating %s\n\n", filename)

	// Compile the file
	compiler := n4l.NewCompiler(filename)
	ast, err := compiler.Compile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Success
	fmt.Printf("✓ Validation passed!\n\n")
	if verbose {
		fmt.Printf("  Chapter: %s\n", ast.Chapter)
		fmt.Printf("  Contexts: %v\n", ast.Contexts)
		fmt.Printf("  Nodes: %d\n", len(ast.Nodes))
		fmt.Printf("  Core arrows available: %d\n", len(n4l.CoreArrows))
	}
}

func watchFile(filename string, verbose bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating watcher: %v\n", err)
		os.Exit(1)
	}
	defer watcher.Close()

	// Initial validation
	validateAndReport(filename, verbose)

	// Add file to watcher
	err = watcher.Add(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error watching file: %v\n", err)
		os.Exit(1)
	}

	// Debounce timer - wait for editor to finish writing
	var debounceTimer *time.Timer
	const debounceDelay = 300 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Only care about write events
			if event.Op&fsnotify.Write == fsnotify.Write {
				// Reset debounce timer
				if debounceTimer != nil {
					debounceTimer.Stop()
				}

				debounceTimer = time.AfterFunc(debounceDelay, func() {
					fmt.Printf("\n� File changed - re-validating...\n\n")
					validateAndReport(filename, verbose)
				})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
		}
	}
}

func validateAndReport(filename string, verbose bool) {
	startTime := time.Now()

	compiler := n4l.NewCompiler(filename)
	ast, err := compiler.Compile()

	duration := time.Since(startTime)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fmt.Printf("❌ Validation failed (%v)\n", duration)
		return
	}

	// Success
	fmt.Printf("✅ Validation passed! (%v)\n", duration)
	if verbose {
		fmt.Printf("  Chapter: %s\n", ast.Chapter)
		fmt.Printf("  Contexts: %v\n", ast.Contexts)
		fmt.Printf("  Nodes: %d\n", len(ast.Nodes))
		fmt.Printf("  Core arrows: %d\n", len(n4l.CoreArrows))
	}
	fmt.Printf("\n👁️  Watching for changes...\n")
}
