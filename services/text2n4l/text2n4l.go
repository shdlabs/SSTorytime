// Package text2n4l provides functionality for extracting high-intentionality
// sentences from text documents and converting them to N4L format.
//
// This package scans documents and identifies sentences that are measured to
// be high in "intentionality" or potential knowledge significance using two
// methods: dynamic running assessment and static post-hoc assessment.

package text2n4l

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	SST "github.com/shdlabs/SSTorytime/services/sstorytime"
)

// sanitizer is a reusable string replacer for cleaning text
var sanitizer = strings.NewReplacer("(", "[", ")", "]")

// ProcessFile analyzes a text file and extracts high-intentionality sentences
// based on the specified percentage threshold.
func ProcessFile(filename string, percentage float64) error {
	SST.MemoryInit()

	fmt.Println("Fractionating file...", filename)
	psf, L := SST.FractionateTextFile(filename)

	fmt.Println("Analyzing longitudinal patterns")
	ranking1 := SelectByRunningIntent(psf, L, percentage)
	fmt.Println("Analyzing statistical patterns")
	ranking2 := SelectByStaticIntent(psf, L, percentage)
	fmt.Println("Merging selections")
	selection := MergeSelections(ranking1, ranking2)

	fmt.Println("Extracting ambient phrases for context")

	// We only want short fragments for context, else we're repeating
	// significant context info from the actual samples

	const minN = 1 // >= N_GRAM_MIN
	const maxN = 3 // <= N_GRAM_MAX

	f, s, ff, ss := SST.ExtractIntentionalTokens(L, selection, minN, maxN)

	return WriteOutput(filename, selection, L, percentage, f, s, ff, ss)
}

// ProcessTextResult holds the result of processing text content
type ProcessTextResult struct {
	N4LContent        string         `json:"n4l_content"`
	TotalSentences    int            `json:"total_sentences"`
	SelectedSentences int            `json:"selected_sentences"`
	FinalFraction     float64        `json:"final_fraction"`
	RequestedFraction float64        `json:"requested_fraction"`
	Selection         []SST.TextRank `json:"selection"`
}

// ProcessTextContent analyzes text content directly and returns the N4L result
// without creating temporary files.
func ProcessTextContent(content string, percentage float64) (*ProcessTextResult, error) {
	SST.MemoryInit()

	// Create a temporary file to use with the existing SST functions
	tempFile, err := ioutil.TempFile("", "text2n4l_*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up temp file
	defer tempFile.Close()

	// Write content to temporary file
	if _, err := tempFile.WriteString(content); err != nil {
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tempFile.Close() // Close before reading

	// Process using existing functions
	psf, L := SST.FractionateTextFile(tempFile.Name())

	ranking1 := SelectByRunningIntent(psf, L, percentage)
	ranking2 := SelectByStaticIntent(psf, L, percentage)
	selection := MergeSelections(ranking1, ranking2)

	const minN = 1 // >= N_GRAM_MIN
	const maxN = 3 // <= N_GRAM_MAX

	f, s, ff, ss := SST.ExtractIntentionalTokens(L, selection, minN, maxN)

	// Generate N4L content as string instead of writing to file
	n4lContent, err := GenerateN4LContent("text_input", selection, L, percentage, f, s, ff, ss)
	if err != nil {
		return nil, fmt.Errorf("failed to generate N4L content: %w", err)
	}

	finalFraction := float64(len(selection)*100) / float64(L)

	return &ProcessTextResult{
		N4LContent:        n4lContent,
		TotalSentences:    L,
		SelectedSentences: len(selection),
		FinalFraction:     finalFraction,
		RequestedFraction: percentage,
		Selection:         selection,
	}, nil
}

// GenerateN4LContent creates N4L formatted content as a string
func GenerateN4LContent(filename string, selection []SST.TextRank, L int, percentage float64, anom_by_part [][]string, ambi_by_part [][]string, all_anom []string, all_ambi []string) (string, error) {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf(" - Samples from %s\n", filename))
	builder.WriteString("\n# (begin) ************\n")

	filealias := strings.Split(filename, ".")[0]
	builder.WriteString(fmt.Sprintf("\n :: _sequence_ , %s::\n", filealias))

	var partcheck = make(map[string]bool)
	var parts []string
	var lastpart string

	for i := range selection {
		context := SpliceSet(ambi_by_part[selection[i].Partition])
		part := PartName(selection[i].Partition, filealias, context)

		// Add context from n = 2,3 fractions
		if part != lastpart {
			if len(context) > 0 {
				builder.WriteString(fmt.Sprintf("\n :: %s ::\n", context))
				lastpart = part
			}
		}

		builder.WriteString(fmt.Sprintf("\n@sen%d   %s\n", selection[i].Order, Sanitize(selection[i].Fragment)))
		builder.WriteString(fmt.Sprintf("              \" (%s) %s\n", SST.INV_CONT_FOUND_IN_L, part))

		// Add intentional context
		for j := range anom_by_part[selection[i].Partition] {
			builder.WriteString(fmt.Sprintf("              \" (%s) %s\n", SST.NEAR_FRAG_L, anom_by_part[selection[i].Partition][j]))
		}

		if !partcheck[part] {
			parts = append(parts, part)
			partcheck[part] = true
		}
	}

	builder.WriteString("\n# (end) ************\n")

	// some stats
	builder.WriteString(fmt.Sprintf("\n# Final fraction %.2f of requested %.2f\n", float64(len(selection)*100)/float64(L), percentage))
	builder.WriteString(fmt.Sprintf("\n# Selected %d samples of %d: ", len(selection), L))

	for i := range selection {
		builder.WriteString(fmt.Sprintf("%d ", selection[i].Order))
	}
	builder.WriteString("\n#\n")

	// document the parts
	builder.WriteString("\n :: themes and topics you might want to annotate/replace ::\n")
	builder.WriteString("\n :: parts, sections ::\n")

	for p := range parts {
		builder.WriteString(fmt.Sprintf("\n %s\n", parts[p]))
		for w := range ambi_by_part[p] {
			builder.WriteString(fmt.Sprintf("  #AMBI %s\n", ambi_by_part[p][w]))
		}

		for w := range anom_by_part[p] {
			builder.WriteString(fmt.Sprintf("   #INTENT %s\n", anom_by_part[p][w]))
		}
	}

	// whole document summary
	for w := range all_ambi {
		builder.WriteString(fmt.Sprintf(" # %s\n", all_ambi[w]))
	}

	for w := range all_anom {
		builder.WriteString(fmt.Sprintf("  # %s\n", all_anom[w]))
	}

	return builder.String(), nil
}

// WriteOutput creates the N4L output file with the selected sentences and context
func WriteOutput(filename string, selection []SST.TextRank, L int, percentage float64, anom_by_part [][]string, ambi_by_part [][]string, all_anom []string, all_ambi []string) error {
	// See AddMandatory() in N4L.go for reserved names (TBD, collect these one day as const)

	outputfile := filename + "_edit_me.n4l"

	fp, err := os.Create(outputfile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputfile, err)
	}
	defer fp.Close()

	// Adaptive buffering based on expected output size
	// Heuristic: Larger input files typically generate larger outputs
	// Small files (< 50KB input): 4KB buffer or no buffer
	// Medium files (50KB-500KB input): 16KB buffer
	// Large files (> 500KB input): 64KB buffer
	var writer interface {
		Write([]byte) (int, error)
	}

	if inputInfo, err := os.Stat(filename); err == nil {
		inputSize := inputInfo.Size()
		if inputSize > 500*1024 { // > 500KB input
			bufWriter := bufio.NewWriterSize(fp, 64*1024) // 64KB buffer for large files
			defer bufWriter.Flush()
			writer = bufWriter
		} else if inputSize > 50*1024 { // 50KB-500KB input
			bufWriter := bufio.NewWriterSize(fp, 16*1024) // 16KB buffer for medium files
			defer bufWriter.Flush()
			writer = bufWriter
		} else { // < 50KB input
			writer = fp // Direct writes for small files
		}
	} else {
		// Fallback: use direct writes if we can't stat the file
		writer = fp
	}

	fmt.Fprintf(writer, " - Samples from %s\n", filename)
	fmt.Fprintf(writer, "\n# (begin) ************\n")

	filealias := strings.Split(filename, ".")[0]
	fmt.Fprintf(writer, "\n :: _sequence_ , %s::\n", filealias)

	var partcheck = make(map[string]bool)
	var parts []string
	var lastpart string

	for i := range selection {
		context := SpliceSet(ambi_by_part[selection[i].Partition])
		part := PartName(selection[i].Partition, filealias, context)

		// Add context from n = 2,3 fractions
		if part != lastpart {
			if len(context) > 0 {
				fmt.Fprintf(writer, "\n :: %s ::\n", context)
				lastpart = part
			}
		}

		fmt.Fprintf(writer, "\n@sen%d   %s\n", selection[i].Order, Sanitize(selection[i].Fragment))
		fmt.Fprintf(writer, "              \" (%s) %s\n", SST.INV_CONT_FOUND_IN_L, part)

		AddIntentionalContext(writer, anom_by_part[selection[i].Partition])

		if !partcheck[part] {
			parts = append(parts, part)
			partcheck[part] = true
		}
	}

	fmt.Fprintf(writer, "\n# (end) ************\n")

	// some stats
	fmt.Fprintf(writer, "\n# Final fraction %.2f of requested %.2f\n", float64(len(selection)*100)/float64(L), percentage)
	fmt.Fprintf(writer, "\n# Selected %d samples of %d: ", len(selection), L)

	for i := range selection {
		fmt.Fprintf(writer, "%d ", selection[i].Order)
	}
	fmt.Fprintf(writer, "\n#\n")

	// document the parts
	fmt.Fprintf(writer, "\n :: themes and topics you might want to annotate/replace ::\n")
	fmt.Fprintf(writer, "\n :: parts, sections ::\n")

	for p := range parts {
		fmt.Fprintf(writer, "\n %s\n", parts[p])
		for w := range ambi_by_part[p] {
			fmt.Fprintf(writer, "  #AMBI %s\n", ambi_by_part[p][w])
		}

		for w := range anom_by_part[p] {
			fmt.Fprintf(writer, "   #INTENT %s\n", anom_by_part[p][w])
		}
	}

	// whole document summary
	for w := range all_ambi {
		fmt.Fprintf(writer, " # %s\n", all_ambi[w])
	}

	for w := range all_anom {
		fmt.Fprintf(writer, "  # %s\n", all_anom[w])
	}

	fmt.Println("Wrote file", outputfile)
	fmt.Printf("Final fraction %.2f of requested %.2f sampled\n", float64(len(selection)*100)/float64(L), percentage)

	return nil
}

// PartName generates a descriptive name for a document partition
func PartName(p int, file string, context string) string {
	// include ambient context in the section name
	return fmt.Sprintf("part %d of %s with %s", p, file, context)
}

// SpliceSet joins a slice of strings with commas
func SpliceSet(ctx []string) string {
	return strings.Join(ctx, ", ")
}

// AddIntentionalContext writes intentional context markers to the output file
func AddIntentionalContext(writer interface{ Write([]byte) (int, error) }, ctx []string) {
	for i := range ctx {
		fmt.Fprintf(writer, "              \" (%s) %s\n", SST.NEAR_FRAG_L, ctx[i])
	}
}

// Sanitize cleans up text by replacing problematic characters
func Sanitize(s string) string {
	return sanitizer.Replace(s)
}

// SelectByRunningIntent analyzes text using dynamic running assessment
func SelectByRunningIntent(psf [][][]string, L int, percentage float64) []SST.TextRank {
	// Rank sentences
	const coherence_length = SST.DUNBAR_30 // approx narrative range or #sentences before new point/topic

	var sentences []SST.TextRank
	var sentence_counter int

	for p := range psf {
		for s := range psf[p] {
			score := 0.0

			// Use strings.Builder for efficient string concatenation
			var builder strings.Builder
			// Pre-allocate capacity (estimate average word length * number of fragments)
			builder.Grow(len(psf[p][s]) * 10)

			for f := 0; f < len(psf[p][s]); f++ {
				score += SST.RunningIntentionality(sentence_counter, psf[p][s][f])

				builder.WriteString(psf[p][s][f])

				if f < len(psf[p][s])-1 {
					builder.WriteString(", ")
				}
			}

			var this SST.TextRank
			this.Fragment = builder.String()
			this.Significance = score
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences, this)
			sentence_counter++
		}
	}

	skimmed := OrderAndRank(sentences, percentage)
	return skimmed
}

// SelectByStaticIntent analyzes text using static post-hoc assessment
func SelectByStaticIntent(psf [][][]string, L int, percentage float64) []SST.TextRank {
	// Rank sentences
	const coherence_length = SST.DUNBAR_30 // approx narrative range or #sentences before new point/topic

	var sentences []SST.TextRank
	var sentence_counter int

	for p := range psf {
		for s := range psf[p] {
			score := 0.0

			// Use strings.Builder for efficient string concatenation
			var builder strings.Builder
			// Pre-allocate capacity (estimate average word length * number of fragments)
			builder.Grow(len(psf[p][s]) * 10)

			for f := 0; f < len(psf[p][s]); f++ {
				score += SST.AssessStaticIntent(psf[p][s][f], L, SST.STM_NGRAM_FREQ, 1)

				builder.WriteString(psf[p][s][f])

				if f < len(psf[p][s])-1 {
					builder.WriteString(", ")
				}
			}

			var this SST.TextRank
			this.Fragment = builder.String()
			this.Significance = score
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences, this)
			sentence_counter++
		}
	}

	skimmed := OrderAndRank(sentences, percentage)
	return skimmed
}

// OrderAndRank sorts sentences by significance and selects the top percentage
func OrderAndRank(sentences []SST.TextRank, percentage float64) []SST.TextRank {
	var selections []SST.TextRank

	// Order by intentionality first to skim cream
	sort.Slice(sentences, func(i, j int) bool {
		return sentences[i].Significance > sentences[j].Significance
	})

	// Measure relative threshold for percentage of document
	// the lower the threshold, the lower the significance of the document
	threshold := percentage / 100.0
	limit := int(threshold * float64(len(sentences)))

	// Skim
	for i := 0; i < limit; i++ {
		selections = append(selections, sentences[i])
	}

	// Order by line number again to restore causal order
	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})

	return selections
}

// MergeSelections combines two sets of text rankings, avoiding duplicates
func MergeSelections(one []SST.TextRank, two []SST.TextRank) []SST.TextRank {
	var merge []SST.TextRank
	var already_selected = make(map[int]bool)

	for i := range one {
		merge = append(merge, one[i])
		already_selected[one[i].Order] = true
	}

	for i := range two {
		if !already_selected[two[i].Order] {
			merge = append(merge, two[i])
		}
	}

	// Order by line number again to restore causal order
	sort.Slice(merge, func(i, j int) bool {
		return merge[i].Order < merge[j].Order
	})

	return merge
}
