//
// Scan a document and pick out the sentences that are measured to
// be high in "intentionality" or potential knowledge significance
//

package main

import (
	"fmt"
        SST "SSTorytime"
)

//**************************************************************
// BEGIN
//**************************************************************

func main() {

	const max_class = 100

	//input := "../../examples/example_data/MobyDick.dat"
	//input := "../../examples/example_data/obama.dat"
	//input := "../../examples/example_data/bede.dat"
	//input := "../../examples/example_data/promisetheory1.dat"
	//input := "../../examples/example_data/Darwin.dat"
	input := "../../examples/example_data/orgmode.dat"

	SST.MemoryInit()

	psf,_ := SST.FractionateTextFile(input)
	
	// Rank sentences

	var selections []SST.TextRank

	maxscore := 0.0

	for p := range psf {

		for s := range psf[p] {

			score := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				score += SST.RunningIntent(s,psf[p][s][f])

				text += psf[p][s][f]

				if f < len(psf[p][s])-1 {
					text += ", "
				}
			}

			var this SST.TextRank
			this.Fragment = text
			this.Significance = score
			selections = append(selections,this)
			if score > maxscore {
				maxscore = score
			}
		}
	}

	// Measure relative threshold for percentage of document
	// the lower the threshold, the lower the significance of the document

	const threshold = 0.1
	const parts = 1000
	var cumulative [parts]int
	var total int
	var cutoff float64

	fmt.Println("Summarize approx",threshold*100,"percent\n\n")

	for i := range selections {
		selclass := selections[i].Significance / maxscore * float64(parts-1)
		cumulative[int(selclass)]++
		total++
	}

	// calc the threshold to keep fraction of entries

	cum := 0

	for i := parts-1; i >= 0; i-- {
		cum += cumulative[i]

		if float64(cum)/float64(total) >= threshold {
			cutoff = float64(i)/float64(parts) * maxscore
			break
		}
	}

	// Now print only upper scoring fraction 20%

	printed := 0
	totald := 0

	for i := range selections {
		totald += len(selections[i].Fragment)

		if selections[i].Significance > cutoff {
			printed += len(selections[i].Fragment)
			fmt.Print(i,".")
			SST.ShowText(selections[i].Fragment,100)
			fmt.Println()
		}
	}

	fmt.Println("Fraction of document = ",float64(printed)/float64(totald))

}

