//
// Scan a document and pick out the n-grams that are persistent
// but relatively anomalous, so that they stand out intentionally
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
	input := "../../examples/example_data/obama.dat"
	//input := "../../examples/example_data/bede.dat"
	//input := "../../examples/example_data/promisetheory1.dat"
	//input := "../../examples/example_data/Darwin.dat"
	//input := "../../examples/example_data/orgmode.dat"

	SST.MemoryInit()

	_,L := SST.FractionateTextFile(input)

	intentions,context := SST.AssessTextAnomalies(L,SST.STM_NGRAM_FREQ,SST.STM_NGRAM_LOCA)

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

}

