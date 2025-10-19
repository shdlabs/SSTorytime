# N4L Validator

A fast, modern N4L validator with watch mode for real-time feedback.

## Quick Start

```bash
# Build
go build -o n4l-validate .

# Install
go install .

# Validate once
./n4l-validate myfile.n4l

# Watch mode - auto-validate on changes
# May become better editor plugin or lsp in the future
./n4l-validate -w myfile.n4l
```

## Features

- ✅ **Fast**: <100µs typical validation time
- 👁️ **Watch Mode**: Auto-validate on file changes
- 📍 **Precise Errors**: Shows exact line, column, and context
- 💡 **Helpful Hints**: Suggests how to fix errors
- 🚫 **No Database**: Pure syntax validation, no upload risk

## Usage

```
Usage: ./n4l-validate [options] <file.n4l>

Validates N4L files without uploading to database

Options:
  -v    Verbose output
  -w    Watch mode - re-validate on file changes
  -watch
        Watch mode - re-validate on file changes
```

## Examples

### Single Validation

```bash
./n4l-validate examples/astronomy.n4l
```

Output:

```
🔍 Validating examples/astronomy.n4l

✓ Validation passed!
```

### Watch Mode

```bash
./n4l-validate -w examples/myfile.n4l
```

Output:

```
👁️  Watch mode enabled - monitoring /path/to/myfile.n4l
Press Ctrl+C to stop

✅ Validation passed! (71.3µs)

👁️  Watching for changes...
```

Every time you save the file, it re-validates automatically!

### Error Example

```
examples/test.n4l:7:25: error: unterminated string
     7 |  Missing closing quote: "this is broken
       |                         ^
  hint: add closing quote

❌ Validation failed (92.4µs)
```

## Watch Mode Benefits

- **No manual re-runs**: Just save your file
- **Instant feedback**: See errors immediately
- **Editor agnostic**: Works with vim, emacs, VS Code, nano, etc.
- **Smart debouncing**: Waits for editor to finish writing (300ms)
- **Performance tracking**: Shows validation time

## Typical Workflow

```bash
# Terminal 1: Start watching
./n4l-validate -w examples/my-notes.n4l

# Terminal 2: Edit the file
vim examples/my-notes.n4l

# Save, see validation results immediately
# Fix errors, save again, repeat
```

## Demo

Run the interactive demo:

```bash
./watch_demo.sh
```

## Documentation

See [Watch Mode Documentation](../../docs/watch-mode.md) for:

- Detailed usage examples
- Editor integration
- Troubleshooting
- Architecture details

## Comparison

### vs Mark's N4L Tool

| Feature         | n4l-validate            | N4L              |
| --------------- | ----------------------- | ---------------- |
| Error location  | Line:col with context   | Just line number |
| Watch mode      | ✅ Yes                  | ❌ No            |
| Database upload | ❌ No (validation only) | ✅ Yes           |
| Speed           | ~50-100µs               | ~10-50ms         |
| Hints           | ✅ Yes                  | ❌ No            |

Use both together:

1. `n4l-validate -w` while editing (fast, safe)
2. `N4L -u` when ready to upload (complete validation)

## Building from Source

```bash
cd /home/alex/SSTorytime/cmd/n4l-validate
go get github.com/fsnotify/fsnotify
go build -o n4l-validate .
```

## Dependencies

- Go 1.24+
- [fsnotify](https://github.com/fsnotify/fsnotify) - File system notifications
- `internal/n4l` - N4L compiler package

## License

Same as SSTorytime project.
