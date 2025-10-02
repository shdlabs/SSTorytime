package text2n4l

import (
	"fmt"
	"os"
	"testing"

	SST "github.com/shdlabs/SSTorytime/services/sstorytime"
)

// BenchmarkProcessFile benchmarks the full ProcessFile workflow
func BenchmarkProcessFile(b *testing.B) {
	SST.MemoryInit()

	benchmarks := []struct {
		name       string
		inputFile  string
		percentage float64
	}{
		{"PromiseTheory_10pct", "testdata/promisetheory1.dat", 10.0},
		{"PromiseTheory_25pct", "testdata/promisetheory1.dat", 25.0},
		{"PromiseTheory_50pct", "testdata/promisetheory1.dat", 50.0},
		{"MobyDick_10pct", "testdata/MobyDick.dat", 10.0},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := ProcessFile(bm.inputFile, bm.percentage)
				if err != nil {
					b.Fatalf("ProcessFile failed: %v", err)
				}
				// Clean up output file
				os.Remove(bm.inputFile + "_edit_me.n4l")
			}
		})
	}
}

// BenchmarkSelectByRunningIntent benchmarks the dynamic analysis
func BenchmarkSelectByRunningIntent(b *testing.B) {
	SST.MemoryInit()

	// Load test data once
	psf, L := SST.FractionateTextFile("testdata/promisetheory1.dat")

	benchmarks := []struct {
		name       string
		percentage float64
	}{
		{"10pct", 10.0},
		{"25pct", 25.0},
		{"50pct", 50.0},
		{"90pct", 90.0},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = SelectByRunningIntent(psf, L, bm.percentage)
			}
		})
	}
}

// BenchmarkSelectByStaticIntent benchmarks the static analysis
func BenchmarkSelectByStaticIntent(b *testing.B) {
	SST.MemoryInit()

	// Load test data once
	psf, L := SST.FractionateTextFile("testdata/promisetheory1.dat")

	benchmarks := []struct {
		name       string
		percentage float64
	}{
		{"10pct", 10.0},
		{"25pct", 25.0},
		{"50pct", 50.0},
		{"90pct", 90.0},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = SelectByStaticIntent(psf, L, bm.percentage)
			}
		})
	}
}

// BenchmarkMergeSelections benchmarks the merge operation with varying sizes
func BenchmarkMergeSelections(b *testing.B) {
	// Create test data of varying sizes
	sizes := []int{10, 50, 100, 500}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test selections
			selection1 := make([]SST.TextRank, size/2)
			selection2 := make([]SST.TextRank, size/2)

			for i := 0; i < size/2; i++ {
				selection1[i] = SST.TextRank{
					Fragment:     fmt.Sprintf("Fragment %d", i*2),
					Order:        i * 2,
					Significance: float64(i) * 0.1,
				}
				selection2[i] = SST.TextRank{
					Fragment:     fmt.Sprintf("Fragment %d", i*2+1),
					Order:        i*2 + 1,
					Significance: float64(i) * 0.1,
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = MergeSelections(selection1, selection2)
			}
		})
	}
}

// BenchmarkOrderAndRank benchmarks the sorting and selection operation
func BenchmarkOrderAndRank(b *testing.B) {
	sizes := []int{100, 500, 1000, 5000}
	percentages := []float64{10.0, 50.0, 90.0}

	for _, size := range sizes {
		for _, pct := range percentages {
			name := fmt.Sprintf("size_%d_pct_%.0f", size, pct)
			b.Run(name, func(b *testing.B) {
				// Generate test sentences
				sentences := make([]SST.TextRank, size)
				for i := 0; i < size; i++ {
					sentences[i] = SST.TextRank{
						Fragment:     fmt.Sprintf("Fragment %d", i),
						Order:        i,
						Significance: float64(i%100) * 0.01, // Vary significance
					}
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = OrderAndRank(sentences, pct)
				}
			})
		}
	}
}

// BenchmarkSanitize benchmarks text sanitization
func BenchmarkSanitize(b *testing.B) {
	testStrings := []struct {
		name string
		text string
	}{
		{"Short", "This is a test (with parentheses) for sanitization"},
		{"Medium", "This is a longer test string (with multiple) sets of (parentheses) that need to be (sanitized) properly"},
		{"Long", generateLongText(1000)},
		{"VeryLong", generateLongText(10000)},
	}

	for _, ts := range testStrings {
		b.Run(ts.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Sanitize(ts.text)
			}
		})
	}
}

// BenchmarkSpliceSet benchmarks string joining
func BenchmarkSpliceSet(b *testing.B) {
	sizes := []int{5, 20, 100, 500}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test slice
			testSlice := make([]string, size)
			for i := 0; i < size; i++ {
				testSlice[i] = fmt.Sprintf("item%d", i)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = SpliceSet(testSlice)
			}
		})
	}
}

// BenchmarkWriteOutput benchmarks the file writing operation
func BenchmarkWriteOutput(b *testing.B) {
	SST.MemoryInit()

	// Prepare test data
	selection := []SST.TextRank{
		{Fragment: "Test sentence 1", Order: 1, Significance: 0.8, Partition: 0},
		{Fragment: "Test sentence 2", Order: 2, Significance: 0.7, Partition: 0},
		{Fragment: "Test sentence 3", Order: 3, Significance: 0.6, Partition: 1},
	}

	anom_by_part := [][]string{
		{"keyword1", "keyword2"},
		{"keyword3", "keyword4"},
	}

	ambi_by_part := [][]string{
		{"context1", "context2"},
		{"context3", "context4"},
	}

	all_anom := []string{"global_keyword1", "global_keyword2"}
	all_ambi := []string{"global_context1", "global_context2"}

	tempFile := "testdata/benchmark_temp"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := WriteOutput(tempFile, selection, 100, 10.0, anom_by_part, ambi_by_part, all_anom, all_ambi)
		if err != nil {
			b.Fatalf("WriteOutput failed: %v", err)
		}
		// Clean up
		os.Remove(tempFile + "_edit_me.n4l")
	}
}

// BenchmarkMemoryUsage benchmarks memory allocation patterns
func BenchmarkMemoryUsage(b *testing.B) {
	SST.MemoryInit()

	b.Run("ProcessFile_Memory", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := ProcessFile("testdata/promisetheory1.dat", 10.0)
			if err != nil {
				b.Fatalf("ProcessFile failed: %v", err)
			}
			os.Remove("testdata/promisetheory1.dat_edit_me.n4l")
		}
	})
}

// Helper function to generate long text for benchmarking
func generateLongText(length int) string {
	text := "This is a test string (with parentheses) that gets repeated. "
	result := ""
	for len(result) < length {
		result += text
	}
	return result[:length]
}
