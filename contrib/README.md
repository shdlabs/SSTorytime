# SSTorytime Contrib Tools

This directory contains contributed tools and utilities for working with SSTorytime and N4L (Narrative for Language) files.

## Tools

### n4l-validate

Validator and lexer for N4L files. Validates syntax, structure, and semantic correctness of N4L documents.

**Location:** `tools/n4l-validate/`  
**Documentation:** `docs/n4l-validate.md`  
**Build:** `cd tools/n4l-validate && go build`

### godoc2n4l

Extracts Go documentation and converts it to N4L format. Can scrape entire Go stdlib or individual packages.

**Location:** `tools/godoc2n4l/`  
**Documentation:** `docs/godoc2n4l.md`  
**Build:** `cd tools/godoc2n4l && go build`  
**Usage:** See `tools/godoc2n4l/scrape_unified.sh` for batch processing

### json2n4l

Converts JSON data to N4L format, enabling integration with JSON-based systems.

**Location:** `tools/json2n4l/`  
**Build:** `cd tools/json2n4l && go build`

## Documentation

All tool documentation is centralized in the `docs/` directory:

- `docs/n4l-validate.md` - N4L validator documentation
- `docs/godoc2n4l.md` - Go documentation scraper guide

## Building All Tools

```bash
# From the contrib directory
for tool in tools/*/; do
    if [ -f "$tool/go.mod" ] || [ -f "$tool/main.go" ]; then
        echo "Building $(basename $tool)..."
        (cd "$tool" && go build)
    fi
done
```

## Internal Libraries

Tools may include internal libraries in their respective directories:

- `tools/n4l-validate/internal/` - Lexer, parser, and validator components

## Contributing

When adding new tools:

1. Place the tool in `tools/<toolname>/`
2. Add documentation to `docs/<toolname>.md`
3. Update this README with tool description
4. Follow Go project conventions (go.mod, main.go, etc.)
