# Text2N4L Package

This package provides functionality for extracting high-intentionality sentences from text documents and converting them to N4L format.

## Overview

The package analyzes text documents and identifies sentences that are measured to be high in "intentionality" or potential knowledge significance using two complementary methods:

1. **Dynamic Running Assessment** - Analyzes patterns as the document unfolds
2. **Static Post-hoc Assessment** - Performs statistical analysis after processing the entire document

## Usage

### As a Library

```go
import "github.com/shdlabs/SSTorytime/services/text2n4l"

// Process a file with 50% extraction rate
err := text2n4l.ProcessFile("input.txt", 50.0)
if err != nil {
    log.Fatal(err)
}
```

### As a Command-Line Tool

```bash
# Extract 50% of high-intentionality sentences (default)
./text2N4L input.txt

# Extract 25% of sentences
./text2N4L -% 25 input.txt
```

## Output

The tool generates an N4L file with the following structure:

- **Selected sentences** with intentionality rankings
- **Contextual information** for each partition
- **Ambient phrases** for semantic context
- **Statistical summary** of the extraction process

## Functions

### Core Functions

- `ProcessFile(filename, percentage)` - Main processing function
- `SelectByRunningIntent()` - Dynamic analysis method
- `SelectByStaticIntent()` - Static analysis method
- `MergeSelections()` - Combines results from both methods

### Utility Functions

- `Sanitize()` - Cleans text for N4L format
- `SpliceSet()` - Joins string slices with commas
- `PartName()` - Generates descriptive partition names
- `OrderAndRank()` - Sorts and filters by significance

## Dependencies

- `github.com/shdlabs/SSTorytime/services/sstorytime` - Core SST analysis functions

## Testing

### Golden Test Strategy

We use a golden test approach to ensure that refactoring doesn't break the core functionality. The golden test:

1. **Input**: Uses `testdata/promisetheory1.dat` as a stable test document
2. **Core Validation**: Checks that the same sentences are selected (by sentence number)
3. **Structure Validation**: Ensures the N4L output format remains consistent
4. **Non-deterministic Handling**: The test focuses on deterministic aspects (sentence selection) while being flexible about non-deterministic contextual phrases

### Why This Approach?

The text analysis algorithm has some inherent non-determinism in contextual phrase extraction, which is actually beneficial for the end-user experience. Our golden test captures the essential functionality while allowing for this beneficial variability.

### Running Tests

```bash
# Run all tests
go test -v

# Run just the golden test
go test -run TestProcessFileGolden -v

# Run with coverage
go test -cover
```

### Test Data

- `testdata/promisetheory1.dat`: Input document from Promise Theory documentation
- Expected output: 58 sentences selected from 10% selection threshold
- Key sentences tested: First 10 selected sentences to verify core algorithm stability

### Refactoring Guidelines

When refactoring this package:

1. **Run the golden test first** to establish baseline
2. **Make incremental changes** and re-run tests frequently  
3. **If golden test fails**, investigate whether it's due to:
   - Bug in refactored code (fix the code)
   - Intentional algorithm improvement (update test expectations)
   - Infrastructure change (update test setup)

The golden test will catch any unintended changes to the core sentence selection algorithm while allowing beneficial improvements to be validated and incorporated.