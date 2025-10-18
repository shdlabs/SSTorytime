# N4L Validator Workflow

## Overview

The validator (`n4l-validate`) provides strict syntax checking for N4L files, treating N4L as a compiled language with proper lexical analysis and error reporting.

## Design Decision: Validator on Upload Failure

After testing, we decided on this workflow:

1. **Upload First**: Try uploading the N4L file directly with `N4L -u`
2. **Validate on Failure**: If upload fails, run validator to diagnose the issue
3. **Keep Validator Strict**: The validator correctly identifies syntax errors

### Rationale

- **N4L.go has undefined behavior**: The original uploader sometimes accepts invalid syntax, sometimes rejects it, behavior is inconsistent and undocumented
- **Validator is correct**: Files that pass validation should upload without errors
- **Pragmatic approach**: Don't block legacy files that N4L.go accepts (warnings but successful upload)
- **Future-proof**: New files should pass validation before upload

## Makefile Integration

The `examples/Makefile` implements this workflow:

```makefile
define validate_and_upload
    # Try upload first
    if N4L -u $(flags) $(file); then
        success
    else
        # Upload failed - run validator to diagnose
        validator $(file)
        if validator passed:
            "File passed validation but N4L upload failed"
            "This may be a semantic error or N4L bug"
        else:
            "File has syntax errors - fix these first"
        exit 1
    fi
endef
```

## Known Issues

### 1. String Parsing Inconsistency

**Problem**: N4L.go's string parsing is inconsistent

- Sometimes allows unterminated strings
- Sometimes rejects them
- No clear rules for when it's lenient vs strict

**Example** (SSTorytime.n4l line 131):

```n4l
"    (depends on) create the SST database (e.g.) "$ sudo su -
```

- Validator correctly identifies: `unterminated string`
- N4L.go shows warnings but uploads successfully

**Resolution**: Document this as undefined behavior. Files that pass validation should always upload cleanly.

### 2. Special Characters in Text

**Problem**: N4L parser confuses text content with syntax

- `==` in descriptions triggers arrow parsing
- `=` might be interpreted as annotation
- No escaping mechanism for these characters

**Example** (golang_stdlib_unified.n4l line 82):

```n4l
" (contain) "...err = nil..."
```

- N4L.go reports: `Short word, possible mistake...after annotation =`
- The `= nil` in documentation is confused with syntax

**Workaround**: The scraper replaces `==` with `=` in `escapeN4L()`, but this doesn't solve all cases.

**Long-term fix**: N4L needs proper string escaping or a way to mark text as "literal content".

### 3. Multiline Strings

**Problem**: N4L has no defined multiline string syntax

- Files sometimes contain shell commands spanning multiple lines
- These fail validation correctly
- N4L.go sometimes accepts them (undefined behavior)

**Example**:

```n4l
"$ sudo su -
$ su - postgres
$ psql
```

**Resolution**: Don't use multiline strings. Put each command in separate nodes or use a single-line concatenation.

## Validator Features

### Syntax Checking ✅

- Chapter headers (`-`)
- Context tags (`+`)
- Nodes (quoted and unquoted)
- Arrows with proper formatting
- Ditto support (`"` to repeat previous node)
- Position tracking (line, column, file)

### Error Reporting ✅

- Shows exact location of error
- Displays problematic source line
- Provides hints for fixing
- Proper error messages

### Watch Mode ✅

- `n4l-validate -w file.n4l`
- Monitors file for changes
- Validates on save (300ms debounce)
- Live feedback during editing

## Usage Examples

### Command Line

```bash
# Validate a file
n4l-validate myfile.n4l

# Watch mode (validates on save)
n4l-validate -w myfile.n4l

# Verbose output
n4l-validate -v myfile.n4l
```

### Make targets

```bash
# Validate all files (no upload)
make validate

# Upload single file (validate on failure)
make upload FILE=astronomy.n4l FORCE=1

# Watch a file
make watch FILE=astronomy.n4l

# Full rebuild (validate all, then upload all)
make
```

## Status

**Working**:

- ✅ Lexical analysis (tokenization)
- ✅ Syntax error detection
- ✅ Position tracking
- ✅ Watch mode
- ✅ Makefile integration
- ✅ Error reporting with context

**Pending**:

- ⏳ Parse arrow definitions from `SSTconfig/*.sst`
- ⏳ Semantic validation (arrow usage, node references)
- ⏳ Complete parser implementation
- ⏳ Integration with N4L.go to make it strict

## Recommendations

1. **For new N4L files**: Always validate before uploading

   ```bash
   n4l-validate file.n4l && N4L -u file.n4l
   ```

2. **For existing files**: Upload first, validate only if it fails

   ```bash
   make upload FILE=existing.n4l FORCE=1
   ```

3. **Report N4L bugs**: When validator passes but N4L fails, document the case

4. **Document workarounds**: Add them to this file for future reference

## Future Work

- Make N4L.go use the validator's lexer (unified parsing)
- Add semantic validation phase
- Load arrow definitions dynamically from SSTconfig
- Add LSP server for editor integration
- Generate better error messages for semantic errors
