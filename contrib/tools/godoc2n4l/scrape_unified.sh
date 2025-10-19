#!/bin/bash

# Scrape Go standard library into ONE unified N4L document
# This allows us to find cross-package relationships

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
echo "=== Scraping Go Standard Library into ONE unified document ==="
echo ""

OUTPUT_FILE="golang_stdlib_unified.n4l"

# Start with chapter header and context
cat > "$OUTPUT_FILE" << 'EOF'
- Go Standard Library Documentation

+ golang
+ stdlib
+ documentation
+ reference

EOF

# Get ALL standard library packages
echo "📚 Discovering all standard library packages..."
readarray -t packages < <(go list std | grep -v "^cmd/" | grep -v "^internal/" | grep -v "^vendor/")

total=${#packages[@]}
count=0

echo "✓ Found $total stdlib packages to scrape"
echo ""
echo "Scraping into unified document..."
echo ""

for pkg in "${packages[@]}"; do
  count=$((count + 1))
  echo "[$count/$total] 📦 Scraping $pkg..."
  
  # Scrape to temp file
  filename=$(echo "$pkg" | tr '/' '_')
  temp_file="/tmp/godoc_${filename}.n4l"
  
  ./godoc2n4l \
    -chapter "TEMP" \
    -context "TEMP" \
    -o "$temp_file" \
    "http://localhost:6060/$pkg" 2>&1 | grep -E "(functions|types|examples)" || true
  
  # Extract package content (skip chapter and context lines)
  # The individual scraper now creates the package node with proper formatting
  pkg_name=$(basename "$pkg")
  echo "" >> "$OUTPUT_FILE"
  
  # Extract everything after the context tags
  # Skip first 3 lines (chapter, empty, context tags)
  sed -n '/^$/,$p' "$temp_file" | \
    grep -v '^-' | \
    grep -v '^+' | \
    sed '/^$/d' \
    >> "$OUTPUT_FILE"
  
  rm -f "$temp_file"
done

echo ""
echo "=== Complete! ==="
echo ""
echo "📄 Generated: $OUTPUT_FILE"
ls -lh "$OUTPUT_FILE"
echo ""

# Count what we captured
funcs=$(grep -c "\" (contain) signature" "$OUTPUT_FILE" 2>/dev/null || echo "0")
types=$(grep -c "\" (contain) kind" "$OUTPUT_FILE" 2>/dev/null || echo "0")
pkgs=$(grep -c "\" (contain) \"package type\"" "$OUTPUT_FILE" 2>/dev/null || echo "0")

echo "📊 Statistics:"
echo "  Packages: $pkgs"
echo "  Functions: $funcs"
echo "  Types: $types"
echo ""

# Show sample of structure
echo "📋 Document structure preview:"
head -50 "$OUTPUT_FILE"
echo ""
echo "..."
echo ""

echo "🔍 Now you can:"
echo "  1. View the file: cat $OUTPUT_FILE | less"
echo "  2. Upload to N4L: ../../src/N4L -u -force $OUTPUT_FILE"
echo "  3. Search for relationships: grep 'io.Reader' $OUTPUT_FILE"
echo "  4. Find cross-package references"
echo ""
echo "💡 With everything in ONE graph, we can discover:"
echo "  - Which packages use which (import relationships)"
echo "  - Which types appear in multiple packages (io.Reader, error, context.Context)"
echo "  - Which functions are related across packages"
echo "  - Learning paths (basic → intermediate → advanced)"
