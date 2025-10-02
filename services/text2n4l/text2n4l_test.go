package text2n4l

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	SST "github.com/shdlabs/SSTorytime/services/sstorytime"
)

func TestSanitize(t *testing.T) {
	input := "This is a test (with parentheses) for sanitization"
	expected := "This is a test [with parentheses] for sanitization"
	result := Sanitize(input)

	if result != expected {
		t.Errorf("Sanitize() = %q, want %q", result, expected)
	}
}

func TestSpliceSet(t *testing.T) {
	input := []string{"apple", "banana", "cherry"}
	expected := "apple, banana, cherry"
	result := SpliceSet(input)

	if result != expected {
		t.Errorf("SpliceSet() = %q, want %q", result, expected)
	}
}

func TestPartName(t *testing.T) {
	result := PartName(1, "testfile", "context info")
	expected := "part 1 of testfile with context info"

	if result != expected {
		t.Errorf("PartName() = %q, want %q", result, expected)
	}
}

// TestProcessFileGolden is our golden test that ensures refactoring doesn't break functionality
func TestProcessFileGolden(t *testing.T) {
	// Initialize memory once for all tests
	SST.MemoryInit()

	testCases := []struct {
		name                    string
		inputFile               string
		percentage              float64
		expectedSentenceCount   int
		expectedSentenceNumbers []int
	}{
		{
			name:                    "promisetheory1_10percent",
			inputFile:               "testdata/promisetheory1.dat",
			percentage:              10.0,
			expectedSentenceCount:   58,                                        // Based on our observed output
			expectedSentenceNumbers: []int{2, 4, 5, 8, 14, 16, 18, 20, 22, 25}, // First 10 core sentence numbers that should be selected
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Process the file
			err := ProcessFile(tc.inputFile, tc.percentage)
			if err != nil {
				t.Fatalf("ProcessFile failed: %v", err)
			}

			// The ProcessFile function creates a file with "_edit_me.n4l" suffix
			actualOutput := tc.inputFile + "_edit_me.n4l"
			defer os.Remove(actualOutput) // Clean up after test

			// Read the actual output
			actualBytes, err := os.ReadFile(actualOutput)
			if err != nil {
				t.Fatalf("Failed to read actual output file: %v", err)
			}

			// Check that we got the expected number of selected sentences
			// Count @sen occurrences
			sentenceCount := bytes.Count(actualBytes, []byte("@sen"))
			if sentenceCount != tc.expectedSentenceCount {
				t.Errorf("Expected %d sentences, got %d", tc.expectedSentenceCount, sentenceCount)
			}

			// Check that specific important sentences are included
			for _, senNum := range tc.expectedSentenceNumbers {
				expectedSentenceMarker := fmt.Sprintf("@sen%d", senNum)
				if !bytes.Contains(actualBytes, []byte(expectedSentenceMarker)) {
					t.Errorf("Expected sentence %s not found in output", expectedSentenceMarker)
				}
			}

			// Check basic structure elements
			structureChecks := []string{
				"# (begin) ************",
				"# (end) ************",
				"_sequence_",
				"extract/quote from",
			}

			for _, check := range structureChecks {
				if !bytes.Contains(actualBytes, []byte(check)) {
					t.Errorf("Expected structure element %q not found in output", check)
				}
			}

			// Check that the final fraction is reported and reasonable
			if !bytes.Contains(actualBytes, []byte("Final fraction")) {
				t.Error("Expected 'Final fraction' summary not found")
			}

			// Verify file header contains our input filename
			if !bytes.Contains(actualBytes, []byte(tc.inputFile)) {
				t.Errorf("Expected input filename %s not found in output header", tc.inputFile)
			}

			t.Logf("Golden test passed: %d sentences selected, structure intact", sentenceCount)
		})
	}
}

// TestIndividualFunctions tests the individual functions for more granular testing
func TestSelectByRunningIntent(t *testing.T) {
	SST.MemoryInit()

	// Create test data
	testText := [][][]string{
		{
			{"This", "is", "a", "test", "sentence."},
			{"Another", "test", "sentence", "follows."},
		},
		{
			{"This", "is", "part", "two", "content."},
		},
	}

	result := SelectByRunningIntent(testText, 3, 50.0)

	// Basic sanity checks
	if len(result) == 0 {
		t.Error("SelectByRunningIntent returned no results")
	}

	// Verify that results are ordered by line number
	for i := 1; i < len(result); i++ {
		if result[i-1].Order > result[i].Order {
			t.Error("Results are not ordered by line number")
		}
	}
}

func TestSelectByStaticIntent(t *testing.T) {
	SST.MemoryInit()

	// Create test data
	testText := [][][]string{
		{
			{"This", "is", "a", "test", "sentence."},
			{"Another", "test", "sentence", "follows."},
		},
		{
			{"This", "is", "part", "two", "content."},
		},
	}

	result := SelectByStaticIntent(testText, 3, 50.0)

	// Basic sanity checks
	if len(result) == 0 {
		t.Error("SelectByStaticIntent returned no results")
	}

	// Verify that results are ordered by line number
	for i := 1; i < len(result); i++ {
		if result[i-1].Order > result[i].Order {
			t.Error("Results are not ordered by line number")
		}
	}
}

func TestMergeSelections(t *testing.T) {
	// Create test data
	selection1 := []SST.TextRank{
		{Fragment: "First", Order: 1, Significance: 0.8},
		{Fragment: "Third", Order: 3, Significance: 0.6},
	}

	selection2 := []SST.TextRank{
		{Fragment: "Second", Order: 2, Significance: 0.7},
		{Fragment: "Third", Order: 3, Significance: 0.6}, // Duplicate
		{Fragment: "Fourth", Order: 4, Significance: 0.5},
	}

	result := MergeSelections(selection1, selection2)

	// Should have 4 unique entries (duplicate "Third" should be removed)
	if len(result) != 4 {
		t.Errorf("Expected 4 unique entries, got %d", len(result))
	}

	// Should be ordered by Order field
	expectedOrder := []int{1, 2, 3, 4}
	for i, item := range result {
		if item.Order != expectedOrder[i] {
			t.Errorf("Expected order %d at position %d, got %d", expectedOrder[i], i, item.Order)
		}
	}
}
