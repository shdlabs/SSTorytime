# Running Local Go Documentation Server

## Quick Start

### Option 1: Using `go doc` (built-in, simple)

For individual packages:

```bash
# View flag package
go doc flag

# View specific function
go doc flag.Parse

# View in browser (if you have godoc installed)
godoc -http=:6060
```

### Option 2: Using `pkgsite` (Official, pkg.go.dev locally)

This runs the **actual pkg.go.dev server** locally!

#### Install pkgsite:

```bash
go install golang.org/x/pkgsite/cmd/pkgsite@latest
```

#### Run the server:

```bash
# Start server on port 8080
pkgsite -http=:8080

# Or specify a different port
pkgsite -http=localhost:6060
```

#### Access it:

```
Open browser: http://localhost:8080
```

This gives you the **exact same interface** as pkg.go.dev but running locally!

### Option 3: Using classic `godoc` (deprecated but still works)

#### Install godoc:

```bash
go install golang.org/x/tools/cmd/godoc@latest
```

#### Run the server:

```bash
# Start on port 6060 (classic default)
godoc -http=:6060

# With index for faster search
godoc -http=:6060 -index
```

#### Access it:

```
Open browser: http://localhost:6060
```

## Recommended: pkgsite

**Use pkgsite** because:

- ✅ Same HTML structure as pkg.go.dev (our scraper already works!)
- ✅ Same selectors, no code changes needed
- ✅ Official tool maintained by Go team
- ✅ Includes all stdlib packages
- ✅ Fast, local, no network needed

## Usage with our scraper

Once running, just change the URL:

```bash
# Instead of:
./godoc2n4l https://pkg.go.dev/flag

# Use:
./godoc2n4l http://localhost:8080/flag

# Or scrape many packages quickly:
for pkg in flag fmt errors context io os net/http; do
  ./godoc2n4l -v -o "$pkg.n4l" "http://localhost:8080/$pkg"
done
```

## Performance Benefits

Local scraping:

- 🚀 **Much faster** - no network latency
- 💾 **Works offline** - no internet needed
- ♻️ **Unlimited requests** - scrape as much as you want
- 😊 **Be nice to Google** - don't hammer their servers

## Verifying it works

```bash
# 1. Start pkgsite
pkgsite -http=:8080 &

# 2. Wait a moment for it to start
sleep 2

# 3. Test with curl
curl -s http://localhost:8080/flag | grep -i "flag package"

# 4. If you see HTML with "Package flag", you're good to go!
```

## Full Example Session

```bash
# Terminal 1: Start the server
pkgsite -http=:8080

# Terminal 2: Use our scraper
cd /home/alex/SSTorytime/cmd/godoc2n4l

# Scrape some packages
./godoc2n4l -v http://localhost:8080/flag
./godoc2n4l -v http://localhost:8080/fmt
./godoc2n4l -v http://localhost:8080/net/http

# Or use the batch script (after updating URLs)
./scrape_samples.sh
```

## Next Steps

1. Install pkgsite: `go install golang.org/x/pkgsite/cmd/pkgsite@latest`
2. Start it: `pkgsite -http=:8080`
3. Update scraper URLs from `pkg.go.dev` to `localhost:8080`
4. Scrape away! 🎉
