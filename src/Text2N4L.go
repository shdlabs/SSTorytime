//
// Scan a document and pick out the sentences that are measured to
// be high in "intentionality" or potential knowledge significance
// using two methods: dynamic running and static posthoc assessment
//

package main

import (
	"fmt"
	"sort"

        SST "SSTorytime"
)

//**************************************************************
// BEGIN
//**************************************************************

func main() {

	const max_class = 100
	const percentage = 32.0

	input := "../examples/example_data/MobyDick.dat"

	SST.MemoryInit()

	psf,L := SST.FractionateTextFile(input)

	ranking1 := SelectByRunningIntent(psf,L,percentage)
	ranking2 := SelectByStaticIntent(psf,L,percentage)

	selection := MergeText(ranking1,ranking2)

	for i := range selection {
		fmt.Println(i,selection[i])
	}

	fmt.Println("FRACTION",float64(len(selection)*100)/float64(L),"of requested",percentage)
}

// ***************************************************

func SelectByRunningIntent(psf [][][]string,L int,percentage float64) []SST.TextRank {

	// Rank sentences

	var sentences []SST.TextRank

	var sentence_counter int

	for p := range psf {

		for s := range psf[p] {

			score := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				score += SST.RunningIntentionality(sentence_counter,psf[p][s][f])

				text += psf[p][s][f]

				if f < len(psf[p][s])-1 {
					text += ", "
				}
			}

			var this SST.TextRank
			this.Fragment = text
			this.Significance = score
			this.Order = sentence_counter
			sentences = append(sentences,this)
			sentence_counter++
		}
	}

	skimmed := OrderAndRank(sentences,percentage)

	return skimmed
}

// ***************************************************

func SelectByStaticIntent(psf [][][]string,L int,percentage float64) []SST.TextRank {

	// Rank sentences

	var sentences []SST.TextRank
	var sentence_counter int

	for p := range psf {

		for s := range psf[p] {

			score := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				score += SST.AssessStaticIntent(psf[p][s][f],L,SST.STM_NGRAM_FREQ,1)

				text += psf[p][s][f]

				if f < len(psf[p][s])-1 {
					text += ", "
				}
			}

			var this SST.TextRank
			this.Fragment = text
			this.Significance = score
			this.Order = sentence_counter
			sentences = append(sentences,this)
			sentence_counter++
		}
	}

	skimmed := OrderAndRank(sentences,percentage)

	return skimmed
}

//*********************************************************************************

func OrderAndRank(sentences []SST.TextRank,percentage float64) []SST.TextRank {

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
		selections = append(selections,sentences[i])
	}

	// Order by line number again to restore causal order

	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})

	return selections
}

//*********************************************************************************

func MergeText(one []SST.TextRank,two []SST.TextRank) []SST.TextRank{

	var merge []SST.TextRank
	var already_selected = make(map[int]bool)
	
	for i := range one {
		merge = append(merge,one[i])
		already_selected[one[i].Order] = true
	}

	for i := range two {
		if !already_selected[two[i].Order] {
			merge = append(merge,two[i])
		}
	}

	// Order by line number again to restore causal order

	sort.Slice(merge, func(i, j int) bool {
		return merge[i].Order < merge[j].Order
	})

	return merge
}


