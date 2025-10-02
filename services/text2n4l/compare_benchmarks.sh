#!/bin/bash
# Compare two benchmark results
# Usage: ./compare_benchmarks.sh baseline.txt current.txt

if [ $# -ne 2 ]; then
    echo "Usage: $0 <baseline_file> <current_file>"
    echo "Example: $0 benchmark_baseline.txt benchmark_$(date +%Y%m%d_%H%M%S).txt"
    exit 1
fi

BASELINE="$1"
CURRENT="$2"

if [ ! -f "$BASELINE" ]; then
    echo "Error: Baseline file '$BASELINE' not found"
    exit 1
fi

if [ ! -f "$CURRENT" ]; then
    echo "Error: Current file '$CURRENT' not found"
    exit 1
fi

echo "Benchmark Comparison Report"
echo "==========================="
echo "Baseline: $BASELINE"
echo "Current:  $CURRENT"
echo "Generated: $(date)"
echo ""

# Function to extract metric from benchmark line
extract_metric() {
    local file="$1"
    local benchmark="$2"
    local metric="$3"  # 3=time, 4=memory, 5=allocs
    
    grep "$benchmark" "$file" | head -1 | awk -v col="$metric" '{print $col}' | sed 's/[^0-9.]//g'
}

# Function to calculate percentage change
calc_percentage() {
    local baseline="$1"
    local current="$2"
    
    if [ -z "$baseline" ] || [ -z "$current" ] || [ "$baseline" = "0" ]; then
        echo "N/A"
        return
    fi
    
    echo "scale=2; (($current - $baseline) / $baseline) * 100" | bc -l 2>/dev/null || echo "N/A"
}

# Function to format change with color
format_change() {
    local pct="$1"
    local better_direction="$2"  # "down" means lower is better, "up" means higher is better
    
    if [ "$pct" = "N/A" ]; then
        echo "N/A"
        return
    fi
    
    local sign=""
    if (( $(echo "$pct > 0" | bc -l 2>/dev/null || echo 0) )); then
        sign="+"
    fi
    
    # Determine if this is an improvement
    local is_improvement=false
    if [ "$better_direction" = "down" ] && (( $(echo "$pct < 0" | bc -l 2>/dev/null || echo 0) )); then
        is_improvement=true
    elif [ "$better_direction" = "up" ] && (( $(echo "$pct > 0" | bc -l 2>/dev/null || echo 0) )); then
        is_improvement=true
    fi
    
    if [ "$is_improvement" = true ]; then
        echo "✅ ${sign}${pct}%"
    else
        echo "❌ ${sign}${pct}%"
    fi
}

# Compare key benchmarks
compare_benchmark() {
    local name="$1"
    local benchmark_pattern="$2"
    
    echo "$name:"
    echo "----------------------------------------"
    
    # Extract metrics
    local baseline_time=$(extract_metric "$BASELINE" "$benchmark_pattern" 3)
    local current_time=$(extract_metric "$CURRENT" "$benchmark_pattern" 3)
    local baseline_memory=$(extract_metric "$BASELINE" "$benchmark_pattern" 4)
    local current_memory=$(extract_metric "$CURRENT" "$benchmark_pattern" 4)
    local baseline_allocs=$(extract_metric "$BASELINE" "$benchmark_pattern" 5)
    local current_allocs=$(extract_metric "$CURRENT" "$benchmark_pattern" 5)
    
    # Calculate changes
    local time_change=$(calc_percentage "$baseline_time" "$current_time")
    local memory_change=$(calc_percentage "$baseline_memory" "$current_memory")
    local allocs_change=$(calc_percentage "$baseline_allocs" "$current_allocs")
    
    # Display results
    printf "  %-12s %15s -> %-15s %s\n" "Time:" "${baseline_time:-N/A}" "${current_time:-N/A}" "$(format_change "$time_change" "down")"
    printf "  %-12s %15s -> %-15s %s\n" "Memory:" "${baseline_memory:-N/A}" "${current_memory:-N/A}" "$(format_change "$memory_change" "down")"
    printf "  %-12s %15s -> %-15s %s\n" "Allocations:" "${baseline_allocs:-N/A}" "${current_allocs:-N/A}" "$(format_change "$allocs_change" "down")"
    echo ""
}

# Compare major benchmarks
compare_benchmark "ProcessFile (10% selection)" "BenchmarkProcessFile/PromiseTheory_10pct"
compare_benchmark "RunningIntent (10% selection)" "BenchmarkSelectByRunningIntent/10pct"
compare_benchmark "StaticIntent (10% selection)" "BenchmarkSelectByStaticIntent/10pct"

echo "Legend:"
echo "✅ = Improvement (lower time/memory/allocations)"
echo "❌ = Regression (higher time/memory/allocations)"
echo ""
echo "Note: Lower values are better for time, memory, and allocations"