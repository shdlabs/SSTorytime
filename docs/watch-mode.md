# N4L Validator Watch Mode

Real-time validation for N4L files while editing.

## Installation

The validator is built with watch mode support:

```bash
cd /home/alex/SSTorytime/cmd/n4l-validate
go build -o n4l-validate .
```

## Usage

### Basic Validation (Single Run)

```bash
./n4l-validate myfile.n4l
```

Output:

```
🔍 Validating myfile.n4l

✓ Validation passed!
```

### Watch Mode

```bash
./n4l-validate -w myfile.n4l
```

or

```bash
./n4l-validate --watch myfile.n4l
```

Output:

```
👁️  Watch mode enabled - monitoring /path/to/myfile.n4l
Press Ctrl+C to stop

✅ Validation passed! (71.317µs)

👁️  Watching for changes...
```

### Verbose Watch Mode

```bash
./n4l-validate -w -v myfile.n4l
```

Shows detailed information on each validation:

```
✅ Validation passed! (85.2µs)
  Chapter: My Chapter
  Contexts: [context1 context2]
  Nodes: 15
  Core arrows: 6

👁️  Watching for changes...
```

## How It Works

1. **Initial Validation**: Validates file immediately when started
2. **File Monitoring**: Uses `fsnotify` to watch for file system events
3. **Debouncing**: Waits 300ms after last change before re-validating
   - Prevents multiple validations while editor is saving
   - Ensures file is fully written before validation
4. **Instant Feedback**: Shows validation results immediately after changes
5. **Continuous**: Keeps watching until you press Ctrl+C

## Use Cases

### 1. Live Editing Feedback

Open your N4L file in any text editor, start watch mode in a terminal:

```bash
# Terminal 1: Start validator
cd /home/alex/SSTorytime/examples
/home/alex/SSTorytime/cmd/n4l-validate/n4l-validate -w myfile.n4l

# Terminal 2: Edit the file
vim myfile.n4l  # or nano, emacs, vscode, etc.
```

Every time you save, the validator runs automatically!

### 2. Development Workflow

```bash
# Start watching your N4L file
./n4l-validate -w examples/my-notes.n4l &

# Edit, save, see immediate results
# No need to manually run validator each time
```

### 3. Debugging Syntax Errors

Watch mode shows exactly where errors are:

```
📝 File changed - re-validating...

examples/test.n4l:7:25: error: unterminated string
     7 |  Missing closing quote: "this is broken
       |                         ^
  hint: add closing quote

❌ Validation failed (92.4µs)

👁️  Watching for changes...
```

Fix the error, save again, and see success immediately:

```
📝 File changed - re-validating...

✅ Validation passed! (78.1µs)

👁️  Watching for changes...
```

## Features

### ✅ What Watch Mode Provides

- **Real-time Feedback**: See errors as soon as you save
- **No Manual Re-runs**: Automatic validation on every change
- **Performance Tracking**: Shows validation time (usually <100µs)
- **Editor Agnostic**: Works with any text editor
- **Smart Debouncing**: Waits for editor to finish writing
- **Clear Status**: Visual indicators (✅ pass, ❌ fail)

### ⏱️ Performance

Validation is extremely fast:

- Typical validation: 50-100 microseconds
- Large files (1000+ lines): Still under 1ms
- No database access: Purely syntax/lexical checks

### 🔔 What Triggers Re-validation

- File write operations (save in editor)
- Manual updates with command-line tools
- Any modification detected by file system

### 🚫 What Doesn't Trigger Re-validation

- Just opening the file
- Moving/renaming (validation stops if file is moved)
- Changes to other files

## Comparison with Traditional Workflow

### Without Watch Mode ❌

```bash
vim myfile.n4l
# ... edit ...
# :wq
./n4l-validate myfile.n4l
# See error, fix it
vim myfile.n4l
# ... edit ...
# :wq
./n4l-validate myfile.n4l
# Repeat...
```

### With Watch Mode ✅

```bash
# Terminal 1
./n4l-validate -w myfile.n4l

# Terminal 2
vim myfile.n4l
# ... edit, save (automatic validation!)
# ... edit, save (automatic validation!)
# ... edit, save (automatic validation!)
```

## Tips & Tricks

### 1. Split Screen Setup

```bash
# Use tmux or split terminal
# Left: Editor
# Right: Watch mode validator
```

### 2. Watch Multiple Files

```bash
# Terminal 1
./n4l-validate -w file1.n4l

# Terminal 2
./n4l-validate -w file2.n4l
```

### 3. Background Validation

```bash
# Run in background (careful with output!)
./n4l-validate -w myfile.n4l > validation.log 2>&1 &

# Check logs
tail -f validation.log
```

### 4. Validation Before Upload

```bash
# Keep validator running while working
./n4l-validate -w examples/myfile.n4l

# In another terminal, when validation passes:
cd examples
../src/N4L -u myfile.n4l
```

## Demo Script

Run the interactive demo:

```bash
cd /home/alex/SSTorytime/cmd/n4l-validate
./watch_demo.sh
```

This will:

1. Create a test file
2. Start watch mode
3. Show you what to edit
4. Display validation results in real-time

## Architecture

### Debouncing Logic

```
File Change Event → Wait 300ms → No more changes? → Validate
                        ↑                               ↓
                        └───────── New change? ─────────┘
                                  (Reset timer)
```

This prevents:

- Multiple validations during one save operation
- Validation of partially-written files
- Excessive CPU usage on rapid changes

### Error Handling

- **File deleted**: Watcher stops, error message shown
- **File moved**: Watcher stops (can't follow moves)
- **Permission denied**: Error shown, watcher continues
- **Syntax errors**: Displayed with context, watcher continues

## Integration with Editors

### VS Code

```json
// In your workspace settings
{
  "terminal.integrated.shellArgs.linux": [
    "-c",
    "/home/alex/SSTorytime/cmd/n4l-validate/n4l-validate -w ${file}"
  ]
}
```

### Vim/Neovim

```vim
" Add to .vimrc
" Open validator in split on save
autocmd BufWritePost *.n4l !n4l-validate %
```

### Emacs

```elisp
;; Add to .emacs
(add-hook 'after-save-hook
  (lambda ()
    (when (string-match "\\.n4l$" (buffer-file-name))
      (compile (concat "n4l-validate " (buffer-file-name))))))
```

## Future Enhancements

Planned features:

- [ ] Watch multiple files simultaneously
- [ ] Integration with VS Code extension
- [ ] LSP (Language Server Protocol) support
- [ ] Auto-fix suggestions
- [ ] Watch entire directory
- [ ] Git-aware watching (only tracked files)
- [ ] Notification system (desktop alerts)
- [ ] Web UI for validation results

## Troubleshooting

### "No such file or directory"

- Check file path is correct
- Use absolute paths or ensure you're in right directory

### "Permission denied"

- Ensure file is readable: `chmod +r myfile.n4l`
- Check directory permissions

### Changes not detected

- Some editors use atomic saves (write temp file, rename)
- Try enabling "safe write" in your editor
- Check if `inotify` is available: `cat /proc/sys/fs/inotify/max_user_watches`

### Too many validations

- Debounce delay is 300ms by default
- If needed, can be increased in code
- Some editors auto-save frequently - consider disabling

## Exit Codes

- `0`: Validation passed (single run mode)
- `1`: Validation failed or file error
- `130`: Interrupted by user (Ctrl+C in watch mode)

## See Also

- [Validator Comparison Report](../docs/validator-comparison.md)
- [N4L Language Specification](../docs/N4L.md)
- [Tutorial](../docs/Tutorial.md)
