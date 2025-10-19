#!/bin/bash

# Scrape multiple Go packages to discover patterns and relationships
# This will help us design the arrow types based on real data

set -e

cd "$(dirname "$0")"

echo "=== Scraping Go Documentation ==="
echo "This will scrape several packages to discover patterns"
echo ""

# Small, well-known packages
packages=(
  "flag:CLI flag parsing"
  "fmt:Formatted I/O"
  "errors:Error handling"
  "context:Context cancellation"
  "io:Basic I/O interfaces"
)

for entry in "${packages[@]}"; do
  IFS=":" read -r pkg desc <<< "$entry"
  echo "📦 Scraping $pkg ($desc)..."
  ./godoc2n4l -v \
    -chapter "$pkg - $desc" \
    -context "golang,stdlib,$pkg" \
    -o "$pkg.n4l" \
    "https://pkg.go.dev/$pkg"
  echo ""
done

echo "=== Scraping complete! ==="
echo ""
echo "Generated files:"
ls -lh *.n4l | awk '{print "  " $9 " (" $5 ")"}'
echo ""
echo "Now examine the files to identify relationships:"
echo "  - Which packages import which?"
echo "  - Which types are used together?"
echo "  - Which functions are examples of concepts?"
echo "  - What prerequisites exist?"
echo ""
echo "Then we'll design arrow types based on what we find!"
