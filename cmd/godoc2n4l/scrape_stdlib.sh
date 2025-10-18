#!/bin/bash

# Scrape Go standard library packages from LOCAL pkgsite server
# This discovers patterns and relationships in the documentation
# Run this AFTER starting: pkgsite -http=:6060

set -e

cd "$(dirname "$0")"

# Check if pkgsite is running
if ! curl -s http://localhost:6060/ > /dev/null 2>&1; then
  echo "❌ pkgsite is not running on port 6060"
  echo ""
  echo "Start it with:"
  echo "  pkgsite -http=:6060 &"
  echo ""
  exit 1
fi

echo "✓ pkgsite is running on http://localhost:6060"
echo ""
echo "=== Scraping Go Standard Library ==="
echo "This will scrape packages to discover documentation patterns"
echo ""

# Create output directory
mkdir -p scraped_docs

# Packages organized by category to find relationships
declare -A packages

# Basic I/O and types
packages[errors]="Error handling basics"
packages[io]="Basic I/O interfaces"
packages[fmt]="Formatted I/O with printing"
packages[bufio]="Buffered I/O"

# CLI and flags
packages[flag]="Command-line flag parsing"
packages[os]="Operating system interface"

# Context and concurrency
packages[context]="Context for cancellation"
packages[sync]="Synchronization primitives"

# Strings and encoding
packages[strings]="String manipulation"
packages[strconv]="String conversions"
packages[encoding/json]="JSON encoding/decoding"

# Networking (shows dependencies)
packages[net]="Network I/O"
packages[net/http]="HTTP client and server"
packages[net/url]="URL parsing"

count=0
total=${#packages[@]}

for pkg in "${!packages[@]}"; do
  count=$((count + 1))
  desc="${packages[$pkg]}"
  
  echo "[$count/$total] 📦 $pkg - $desc"
  
  # Use pkg name without slashes for filename
  filename=$(echo "$pkg" | tr '/' '_')
  
  ./godoc2n4l \
    -chapter "$pkg - $desc" \
    -context "golang,stdlib,docs-analysis,$(basename $pkg)" \
    -o "scraped_docs/${filename}.n4l" \
    "http://localhost:6060/$pkg" \
    2>&1 | grep -E "(Scraping|functions|types|examples)" || true
    
  echo ""
done

echo "=== Scraping Complete! ==="
echo ""
echo "📊 Generated files:"
ls -lh scraped_docs/*.n4l | awk '{print "   " $9 " (" $5 ")"}'
echo ""
echo "🔍 Now let's analyze the data to find patterns:"
echo ""

# Analysis: Find what data we actually captured
echo "Function counts:"
for file in scraped_docs/*.n4l; do
  name=$(basename "$file" .n4l)
  funcs=$(grep -c "\" (contain) signature" "$file" 2>/dev/null || echo "0")
  types=$(grep -c "\" (contain) kind" "$file" 2>/dev/null || echo "0")
  printf "  %-20s %3d functions, %2d types\n" "$name" "$funcs" "$types"
done

echo ""
echo "📝 Next steps:"
echo "  1. Examine the .n4l files in scraped_docs/"
echo "  2. Look for patterns:"
echo "     - Which packages use which other packages?"
echo "     - Which types appear in multiple packages?"
echo "     - What concepts are shared?"
echo "  3. Design arrow types based on REAL relationships we find"
echo "  4. Create a document describing the patterns"
echo ""
echo "Example analysis commands:"
echo "  # Find all mentions of 'Reader' (shows io.Reader usage)"
echo "  grep -h Reader scraped_docs/*.n4l | sort -u"
echo ""
echo "  # Find packages that mention 'error'"
echo "  grep -l '\" (contain) error' scraped_docs/*.n4l"
echo ""
echo "  # See what net/http imports"
echo "  grep 'imported packages' -A 20 scraped_docs/net_http.n4l"
