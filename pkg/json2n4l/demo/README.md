# JSON to N4L Semantic Arrow Styles Demo

This demo showcases the enhanced semantic arrow selection capabilities added to the json2n4l package based on the guidelines in `N4L_STRUCTURED_DATA_GUIDE.md`.

## Overview

The json2n4l package now supports three arrow styles:

1. **Simple** (default) - Uses only `(contain)` arrows for backward compatibility
2. **Semantic** - Uses appropriate semantic arrows based on context:
   - `(propt)` for property values (instead of contain)
   - `(in-set)`/`(setof)` for array elements
3. **Bidirectional** - Adds inverse relationships:
   - `(belong)` inverse of `(contain)`
   - Plus all Semantic features

## New Features

### 1. Auto-Context Detection

Automatically detects semantic context from JSON keys:

```go
// Detects patterns like:
"user" → ["user", "profile", "account"]
"api" → ["api", "endpoint", "response"]
"config" → ["configuration", "settings"]
```

### 2. Alias Generation

Creates `@alias` references for complex objects:

```n4l
@user_2 user
@profile_3 profile
```

### 3. Type Annotations

Optionally includes type information:

```n4l
" (propt) 12345 {number}
" (propt) "Alice Johnson" {string}
" (propt) true {bool}
```

### 4. Semantic Arrows

Uses appropriate arrows from SSTconfig:

- CN-2 Category: `contain/belong`, `setof/in-set`
- EP-3 Category: `propt/propt-of`

## Running the Demo

```bash
cd /home/alex/SSTorytime/pkg/json2n4l/demo
go run demo_styles.go
```

This generates four N4L files:

- `example_simple.n4l` - Default style (backward compatible)
- `example_semantic.n4l` - Semantic arrows with auto-context
- `example_bidirectional.n4l` - Bidirectional with aliases
- `example_options.n4l` - Using option functions with types

## Comparison

### Simple Style (Before)

```n4l
 user
      " (contain) id
           " (contain) 12345
      " (contain) name
           " (contain) "Alice Johnson"
```

### Semantic Style (After)

```n4l
 user
      " (contain) id
           " (propt) 12345
      " (contain) name
           " (propt) "Alice Johnson"
```

### Bidirectional Style (Advanced)

```n4l
@user_2 user
      " (contain) id
           " (propt) 12345
      " (belong) user
      " (contain) name
           " (propt) "Alice Johnson"
      " (belong) user
```

## API Usage

### Using Config directly:

```go
parser := json2n4l.NewParser(json2n4l.Config{
    InputFile:   "input.json",
    OutputFile:  "output.n4l",
    ArrowStyle:  json2n4l.ArrowStyleSemantic,
    AutoContext: true,
})
parser.Parse()
```

### Using option functions:

```go
json2n4l.ParseFile(
    "input.json",
    "output.n4l",
    json2n4l.WithSemantic(),
    json2n4l.WithAutoContext(),
    json2n4l.WithAliases(),
    json2n4l.WithPrettyPrint(),
    json2n4l.WithTypes(),
)
```

## Configuration Options

| Option               | Description                        | Default |
| -------------------- | ---------------------------------- | ------- |
| `ArrowStyle`         | Simple, Semantic, or Bidirectional | Simple  |
| `AutoContext`        | Auto-detect context tags from JSON | false   |
| `GenerateAliases`    | Create @alias references           | false   |
| `IncludeTypes`       | Add {type} annotations             | false   |
| `PreserveStructure`  | Use (pt-of) for hierarchy          | false   |
| `UseSequenceForList` | Use _sequence_ for arrays          | false   |
| `PrettyPrint`        | Add readable comments              | false   |

## Convenience Functions

```go
WithArrowStyle(style)    // Set arrow style
WithSemantic()           // Enable semantic arrows
WithBidirectional()      // Enable bidirectional arrows
WithAutoContext()        // Enable auto-context detection
WithAliases()            // Enable alias generation
WithTypes()              // Enable type annotations
WithPreserveStructure()  // Enable structure preservation
WithSequenceForList()    // Enable sequence mode for arrays
WithPrettyPrint()        // Enable pretty printing
WithChapter(name)        // Set chapter name
WithContext(tags...)     // Set context tags
WithMaxDepth(depth)      // Set max nesting depth
```

## Benefits

1. **Better Semantics** - Uses appropriate arrows that reflect actual relationships
2. **Backward Compatible** - Default behavior unchanged
3. **Flexible** - Choose the level of semantic detail you need
4. **Standards Compliant** - Follows N4L_STRUCTURED_DATA_GUIDE.md best practices
5. **Type Safe** - Go's type system ensures correct usage
6. **Well Documented** - Clear examples and API

## Implementation Details

The enhancements follow Computer Science principles:

- **Recursive Descent Parser** - Clean separation of lexing and parsing
- **Buffered I/O** - Efficient memory usage with `bufio`
- **Semantic Analysis** - Context detection and arrow selection
- **Code Generation** - Multiple output modes (Simple/Semantic/Bidirectional)
- **Functional Options** - Flexible, composable configuration

All changes maintain CS-compliant architecture with proper separation of concerns.
