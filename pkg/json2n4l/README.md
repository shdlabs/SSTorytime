# json2n4l - JSON to N4L Converter

A tool that converts JSON files to N4L (Narrative for Learning) format, using semantic arrows from SSTconfig to represent JSON structure relationships.

## Overview

The `json2n4l` package maps JSON structures to knowledge graph representations:

- **JSON Objects** → Nodes with `!has-key!` arrows to their properties
- **JSON Arrays** → Nodes with `!has-element!` arrows to indexed items
- **Key-Value Pairs** → `!has-value!` relationships
- **Nested Structures** → `!contain!` and `!belong!` bidirectional relationships
- **Array Items** → `!in-array!` back-references with index annotations

## Installation

```bash
cd /home/alex/SSTorytime/cmd/json2n4l
go build -o json2n4l
```

## Usage

### Basic Conversion

```bash
json2n4l input.json
```

This creates `input.n4l` with default settings.

### Specify Output File

```bash
json2n4l input.json output.n4l
```

### With Options

```bash
json2n4l data.json -pretty -types -chapter "API Response" -context "api,json,rest"
```

## Options

| Option      | Type   | Description                                                |
| ----------- | ------ | ---------------------------------------------------------- |
| `-chapter`  | string | Chapter name for the N4L document (default: "JSON Import") |
| `-context`  | string | Comma-separated context tags                               |
| `-types`    | bool   | Include type annotations (object, array, string, etc.)     |
| `-pretty`   | bool   | Add comments for readability                               |
| `-maxdepth` | int    | Maximum nesting depth (0 = unlimited)                      |
| `-root`     | string | Root node name (default: filename)                         |

## Example

Given this JSON:

```json
{
  "user": {
    "name": "Alice",
    "age": 30,
    "roles": ["admin", "user"]
  }
}
```

Running:

```bash
json2n4l example.json -pretty -types
```

Produces this N4L:

```
\chapter{JSON Import}

# Converted from JSON to N4L format
# Structure encoded using semantic arrows

# JSON Object: example
example
  \anno{type: object}
  !has-key! user

  # JSON Object: user
  user
    \anno{type: object}
    !has-key! name

    name
      !has-value! "Alice"
      \anno{type: string}
      !belong! user

    !has-key! age

    age
      !has-value! 30
      \anno{type: number}
      !belong! user

    !has-key! roles

    # JSON Array: roles [2 items]
    roles
      \anno{type: array, length: 2}
      !has-element! "roles[0]"
        \anno{index: 0}

      roles[0]
        !has-value! "admin"
        \anno{type: string}
        !in-array! roles

      !has-element! "roles[1]"
        \anno{index: 1}

      roles[1]
        !has-value! "user"
        \anno{type: string}
        !in-array! roles

      !belong! user

    !belong! example
```

## Semantic Arrows Used

The converter uses arrows defined in SSTconfig:

- `!has-key!` - Object contains key (from arrows-CN-2.sst: contains)
- `!has-element!` - Array contains element (containment)
- `!has-value!` - Node has literal value
- `!belong!` - Value belongs to parent object (inverse of contains)
- `!in-array!` - Element is part of array (inverse of has-element)

## Use Cases

1. **API Documentation**: Convert API responses to navigable knowledge graphs
2. **Configuration Analysis**: Understand complex config file structures
3. **Data Exploration**: Browse JSON data with semantic relationships
4. **Schema Visualization**: See how JSON schemas relate concepts
5. **Integration**: Import JSON data into SST knowledge base

## Programmatic Usage

```go
import "github.com/markburgess/SSTorytime/pkg/json2n4l"

// Simple conversion
config := json2n4l.Config{
    InputFile:  "data.json",
    OutputFile: "output.n4l",
}
converter := json2n4l.NewConverter(config)
err := converter.Convert()

// With options
err = json2n4l.ConvertFile(
    "input.json",
    "output.n4l",
    json2n4l.WithChapter("My Data"),
    json2n4l.WithContext("api", "production"),
    json2n4l.WithTypes(),
    json2n4l.WithPrettyPrint(),
    json2n4l.WithMaxDepth(5),
)

// Convert string
converter := json2n4l.NewConverter(json2n4l.Config{
    RootName: "mydata",
})
n4lString, err := converter.ConvertString(jsonString)
```

## Features

✅ Handles all JSON types (object, array, string, number, boolean, null)  
✅ Preserves nested structure with semantic relationships  
✅ Bidirectional arrows (parent↔child relationships)  
✅ Array indexing with annotations  
✅ Type information (optional)  
✅ Pretty printing with comments  
✅ Depth limiting for large structures  
✅ Proper N4L escaping

## Future Enhancements

- [ ] Custom arrow mapping configuration
- [ ] JSON Schema integration
- [ ] Streaming for large files
- [ ] Reverse conversion (N4L → JSON)
- [ ] Template-based output formatting
- [ ] Compression for repeated structures

## Testing

```bash
cd /home/alex/SSTorytime/pkg/json2n4l
go test -v
```

## Contributing

The json2n4l converter is part of the SSTorytime project. See the main README for contribution guidelines.
