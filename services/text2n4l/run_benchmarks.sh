#!/bin/bash
# Benchmark runner for text2n4l package
# Usage: ./run_benchmarks.sh [output_file]

set -e

PACKAGE_DIR="/home/alex/SSTorytime/services/text2n4l"
OUTPUT_FILE="${1:-benchmark_$(date +%Y%m%d_%H%M%S).txt}"

echo "Running benchmarks for text2n4l package..."
echo "Output will be saved to: $OUTPUT_FILE"

cd "$PACKAGE_DIR"

# Run comprehensive benchmarks
echo "Running comprehensive benchmarks..."
go test -bench=. -benchmem -count=3 > "$OUTPUT_FILE" 2>&1

# Extract key metrics for quick comparison
echo "" >> "$OUTPUT_FILE"
echo "=== BENCHMARK SUMMARY ===" >> "$OUTPUT_FILE"
echo "Generated: $(date)" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Extract ProcessFile results
echo "ProcessFile Results:" >> "$OUTPUT_FILE"
grep "BenchmarkProcessFile" "$OUTPUT_FILE" | head -3 >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Extract Selection Algorithm results
echo "Selection Algorithm Results:" >> "$OUTPUT_FILE"
grep "BenchmarkSelectBy" "$OUTPUT_FILE" | head -8 >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Extract Memory Usage results
echo "Memory Usage Results:" >> "$OUTPUT_FILE"
grep "BenchmarkMemoryUsage" "$OUTPUT_FILE" | head -3 >> "$OUTPUT_FILE"

echo "Benchmarks completed. Results saved to: $OUTPUT_FILE"

# Show quick summary
echo ""
echo "Quick Summary:"
echo "=============="
echo "ProcessFile (10% selection):"
grep "BenchmarkProcessFile/PromiseTheory_10pct" "$OUTPUT_FILE" | head -1 | awk '{print "  Time: " $3 " | Memory: " $4 " | Allocs: " $5}'

echo "RunningIntent (10% selection):"
grep "BenchmarkSelectByRunningIntent/10pct" "$OUTPUT_FILE" | head -1 | awk '{print "  Time: " $3 " | Memory: " $4 " | Allocs: " $5}'

echo "StaticIntent (10% selection):"
grep "BenchmarkSelectByStaticIntent/10pct" "$OUTPUT_FILE" | head -1 | awk '{print "  Time: " $3 " | Memory: " $4 " | Allocs: " $5}'

echo ""
echo "For detailed results, see: $OUTPUT_FILE"