//
// Scan a document and see if we can pick out the
// sentences that are worth remembering .. combining the longitidunal and statistical
//

package main

import (
	"fmt"
	"strings"
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

	SST.MemoryInit()

	psf,L := SST.FractionateTextFile(input)

	//intentions,context
	intentions,_ := SST.AssessTextSignificance(L,SST.STM_NGRAM_FREQ,SST.STM_NGRAM_LOCA)

/*

	for n := range intentions {
		for ngram := range intentions[n] {
			fmt.Println("-Intended: ",n,intentions[n][ngram].Fragment,intentions[n][ngram].Significance)
		}
	}
	for n := range context {
		for ngram := range context[n] {
			fmt.Println("-Context: ",n,context[n][ngram].Fragment,context[n][ngram].Significance)
		}
	}
*/

	var selections []SST.TextRank


	for p := range psf {

		for s := range psf[p] {

			keep := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				for n := 1; n < SST.N_GRAM_MAX; n++ {
					for ngram := range intentions[n] {
						if strings.Contains(psf[p][s][f],intentions[n][ngram].Fragment) {
							keep = 1.0
						}
					}
				}

				text += SST.CleanNgram(psf[p][s][f])

				if f < len(psf[p][s])-1 {
					text += ", "
				}
			}
			
			var this SST.TextRank
			this.Fragment = text
			this.Significance = keep
			selections = append(selections,this)
		}
	}

	// Now print only upper scoring fraction 20%

	printed := 0
	totald := 0

	for i := range selections {
		totald += len(selections[i].Fragment)

		if selections[i].Significance > 0 {
			printed += len(selections[i].Fragment)
			fmt.Print(i,".")
			SST.ShowText(selections[i].Fragment,100)
			fmt.Println()
		}
	}

	fmt.Println("Fraction of document = ",float64(printed)/float64(totald))
}

