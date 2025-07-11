//
// Scan a document and pick out the sentences that are measured to
// be high in "intentionality" or potential knowledge significance
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

	input := "../../examples/example_data/MobyDick.dat"
	//input := "../../examples/example_data/obama.dat"
	//input := "../../examples/example_data/bede.dat"
	//input := "../../examples/example_data/promisetheory1.dat"
	//input := "../../examples/example_data/Darwin.dat"
	//input := "../../examples/example_data/orgmode.dat"

	set1 := Test1(input)
	set2 := Test2(input)

	var histo = make(map[int]int)

	for i := range set1 {
		histo[set1[i]]++
	}

	for i := range set1 {
		histo[set2[i]]++
	}

	var histo2 = make(map[int]int)

	for a := range histo {
		for i := 0; i < 3; i++ {
			if histo[a] == i {
				histo2[i]++
			}
		}
	}

	fmt.Println("Consistency",100*histo2[2]/(histo2[2]+histo2[1]/2),"%")
}

// ***************************************************

func Test1(input string) []int {

	SST.MemoryInit()

	psf,_ := SST.FractionateTextFile(input)

	// Rank sentences

	var sentences []SST.TextRank
	var selections []SST.TextRank
	var count int

	for p := range psf {

		for s := range psf[p] {

			score := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				score += SST.RunningIntentionality(s,psf[p][s][f])

				text += psf[p][s][f]

				if f < len(psf[p][s])-1 {
					text += ", "
				}
			}

			var this SST.TextRank
			this.Fragment = text
			this.Significance = score
			this.Order = count
			sentences = append(sentences,this)
			count++
		}
	}

	sort.Slice(sentences, func(i, j int) bool {
		return sentences[i].Significance > sentences[j].Significance
	})

	// Measure relative threshold for percentage of document
	// the lower the threshold, the lower the significance of the document

	const threshold = 0.1

	limit := int(threshold * float64(len(sentences)))

	for i := 0; i < limit; i++ {
		selections = append(selections,sentences[i])
	}

	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})

	// Now print only upper scoring fraction 20%

	for i := 0; i < limit; i++ {
		fmt.Print(i,"=",selections[i].Order, ": ")
		SST.ShowText(selections[i].Fragment,100)
		fmt.Println()
	}

	fmt.Println("Fraction of document = ",limit,"->", float64(limit)/float64(len(sentences)))

	var vals []int

	for i := 0; i < limit; i++ {
		vals = append(vals,selections[i].Order)
	}

	return vals
}

// ***************************************************

func Test2(input string) []int {

	SST.MemoryInit()

	psf,L := SST.FractionateTextFile(input)
	
	// Rank sentences

	var sentences []SST.TextRank
	var selections []SST.TextRank
	var count int

	for p := range psf {

		for s := range psf[p] {

			score := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				score += SST.AssessIntent(psf[p][s][f],L,SST.STM_NGRAM_FREQ,1)

				text += psf[p][s][f]

				if f < len(psf[p][s])-1 {
					text += ", "
				}
			}

			var this SST.TextRank
			this.Fragment = text
			this.Significance = score
			this.Order = count
			sentences = append(sentences,this)
			count++
		}
	}

	sort.Slice(sentences, func(i, j int) bool {
		return sentences[i].Significance > sentences[j].Significance
	})

	// Measure relative threshold for percentage of document
	// the lower the threshold, the lower the significance of the document

	const threshold = 0.1

	limit := int(threshold * float64(len(sentences)))

	for i := 0; i < limit; i++ {
		selections = append(selections,sentences[i])
	}

	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})

	// Now print only upper scoring fraction 20%

	for i := 0; i < limit; i++ {
		fmt.Print(i,"=",selections[i].Order, ": ")
		SST.ShowText(selections[i].Fragment,100)
		fmt.Println()
	}

	fmt.Println("Fraction of document = ",limit,"->", float64(limit)/float64(len(sentences)))

	var vals []int

	for i := 0; i < limit; i++ {
		vals = append(vals,selections[i].Order)
	}

	return vals
}
