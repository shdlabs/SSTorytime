//
// Scan a document and determine causal separation of ngrams
// 
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
	input := "../../examples/example_data/Darwin.dat"
	//input := "../../examples/example_data/orgmode.dat"

	SST.MemoryInit()

	_,L := SST.FractionateTextFile(input)  // loads STM_NGRAM*

	intent,context,parts := SST.AssessHubFields(L)

	for n := 1; n < SST.N_GRAM_MAX; n++ {

		for ngram := range intent[n] {
			fmt.Println("intent",n,ngram,intent[n][ngram],"of",parts)
		}

		for ngram := range context[n] {
			fmt.Println("context",n,ngram,context[n][ngram],"of",parts)
		}
		fmt.Println("-------------------------------")
	}
}

