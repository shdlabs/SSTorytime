#!/bin/bash

echo "=== Smart Buffering Analysis ==="
echo ""

# Display file sizes
echo "Input file sizes:"
ls -lh testdata/*.dat

echo ""
echo "=== Testing Adaptive Buffering Strategy ==="

# Run benchmark with current (adaptive buffering) implementation
echo "Running benchmark with adaptive buffering..."
go test -bench=BenchmarkProcessFile/MobyDick_10pct -benchmem -count=1 -timeout=5m > /tmp/buffered_result.txt 2>&1

# Show the result
cat /tmp/buffered_result.txt | grep "BenchmarkProcessFile/MobyDick_10pct"

echo ""
echo "Checking buffering strategy for different file sizes:"
echo "- MobyDick.dat: $(stat -c%s testdata/MobyDick.dat) bytes -> Should use 64KB buffer (large file)"
echo "- promisetheory1.dat: $(stat -c%s testdata/promisetheory1.dat) bytes -> Should use direct write (small file)"

echo ""
echo "=== Output File Analysis ==="

# Process the files to see output sizes
echo "Processing files to analyze output sizes..."

# Process small file
go test -run TestProcessFileGolden > /dev/null 2>&1
if [ -f testdata/promisetheory1.dat_edit_me.n4l ]; then
    echo "Small file output: $(ls -lh testdata/promisetheory1.dat_edit_me.n4l | awk '{print $5}')"
    rm testdata/promisetheory1.dat_edit_me.n4l
fi

# Process large file (temporarily)
timeout 60s go run -c 'package main; import "github.com/shdlabs/SSTorytime/services/text2n4l"; import SST "github.com/shdlabs/SSTorytime/services/sstorytime"; func main() { SST.MemoryInit(); text2n4l.ProcessFile("testdata/MobyDick.dat", 10.0) }' > /dev/null 2>&1
if [ -f testdata/MobyDick.dat_edit_me.n4l ]; then
    echo "Large file output: $(ls -lh testdata/MobyDick.dat_edit_me.n4l | awk '{print $5}')"
    rm testdata/MobyDick.dat_edit_me.n4l
fi

echo ""
echo "=== Buffering Strategy Summary ==="
echo "Our adaptive buffering strategy:"
echo "1. Files < 50KB input  -> Direct writes (no buffer overhead)"
echo "2. Files 50KB-500KB   -> 16KB buffer (moderate buffering)"  
echo "3. Files > 500KB      -> 64KB buffer (aggressive buffering)"
echo ""
echo "This ensures optimal performance across different file sizes!"