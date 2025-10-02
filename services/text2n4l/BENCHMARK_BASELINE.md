# Text2N4L Performance Baseline Report
Generated: October 2, 2025
Go Version: $(go version)
Platform: Linux
Test Document: testdata/promisetheory1.dat

## Performance Summary

### Full ProcessFile Workflow
| Test Case              | Time/Op (ms) | Memory/Op (MB) | Allocs/Op |
|------------------------|--------------|-----------------|-----------|
| PromiseTheory_10pct    | 1,025.6      | 551.6          | 7,803,447 |
| PromiseTheory_25pct    | 1,009.1      | 551.6          | 7,805,147 |
| PromiseTheory_50pct    | 1,061.7      | 551.2          | 7,807,285 |

**Analysis**: Full processing takes ~1 second, with consistent ~550MB memory usage and ~7.8M allocations.

### Selection Algorithms (Core Performance)
| Algorithm               | Time/Op (ms) | Memory/Op (MB) | Allocs/Op |
|------------------------|--------------|-----------------|-----------|
| RunningIntent_10pct    | 457.2        | 183.2          | 2,595,355 |
| RunningIntent_25pct    | 376.8        | 182.8          | 2,595,265 |
| RunningIntent_50pct    | 335.7        | 183.2          | 2,595,319 |
| RunningIntent_90pct    | 338.9        | 182.9          | 2,595,276 |
| StaticIntent_10pct     | 331.2        | 183.2          | 2,595,295 |
| StaticIntent_25pct     | 348.8        | 182.7          | 2,595,241 |
| StaticIntent_50pct     | 330.2        | 183.1          | 2,595,281 |
| StaticIntent_90pct     | 332.2        | 183.0          | 2,595,268 |

**Analysis**: 
- Running intent analysis: 335-457ms (varies by percentage)
- Static intent analysis: 330-349ms (more consistent)
- Both use ~183MB and ~2.6M allocations
- Selection percentage affects running intent more than static intent

## Performance Hotspots Identified

1. **Memory Allocation**: 7.8M allocations per full run suggests significant GC pressure
2. **Algorithm Variance**: Running intent shows more performance variation (22% spread vs 6%)
3. **Memory Usage**: 551MB for full processing, 183MB for each selection algorithm
4. **Consistent Allocation Pattern**: ~2.6M allocations per algorithm regardless of percentage

## Refactoring Targets

### High Priority
1. **Reduce allocations**: 7.8M allocations indicate excessive object creation
2. **Memory pooling**: Consistent allocation patterns suggest reusable objects
3. **String handling**: Text processing likely generates many temporary strings

### Medium Priority
1. **Algorithm optimization**: Running intent variance suggests optimization opportunity
2. **Early termination**: Selection percentage doesn't scale linearly with time
3. **IO optimization**: File operations may benefit from buffering

### Low Priority
1. **Merge operation**: Relatively small compared to selection algorithms
2. **Utility functions**: Sanitize, SpliceSet show minimal overhead

## Success Criteria for Refactoring

### Performance Goals
- **Reduce total runtime**: Target 20% improvement (800ms from 1000ms)
- **Reduce memory allocations**: Target 50% reduction (4M from 8M allocs)
- **Improve memory efficiency**: Target 20% reduction (440MB from 550MB)
- **Reduce variance**: Make running intent as consistent as static intent

### Quality Goals
- Maintain functional correctness (golden test must pass)
- Preserve algorithm accuracy (same sentences selected)
- Keep code readability and maintainability

## Benchmark Infrastructure

All benchmarks use:
- Consistent test data (promisetheory1.dat)
- Memory allocation tracking (-benchmem)
- Multiple iterations for statistical significance
- Separate tests for each algorithm component

### Files Created
- `text2n4l_bench_test.go`: Comprehensive benchmark suite
- `benchmark_baseline.txt`: Full baseline results
- `benchmark_results_clean.txt`: ProcessFile focused results
- `benchmark_selection.txt`: Selection algorithm results

## Next Steps

1. **Run refactoring iterations**: Each change should be benchmarked
2. **Compare results**: Use `benchcmp` or similar tools for analysis
3. **Monitor regression**: Golden test ensures functional correctness
4. **Document improvements**: Track performance gains and code quality

This baseline establishes our starting point for performance optimization. Any refactoring should improve these metrics while maintaining correctness.