//
// Scan a document and determine causal separation of ngrams
// 
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

	//input := "../../examples/example_data/MobyDick.dat"
	//input := "../../examples/example_data/obama.dat"
	//input := "../../examples/example_data/bede.dat"
	//input := "../../examples/example_data/promisetheory1.dat"
	//input := "../../examples/example_data/Darwin.dat"
	input := "../../examples/example_data/orgmode.dat"

	SST.MemoryInit()

	_,L := SST.FractionateTextFile(input)  // loads STM_NGRAM*

	common,localized,_ := SST.AssessHubFields(L,SST.STM_NGRAM_LOCA)

	for n := 1; n < SST.N_GRAM_MAX; n++ {

		var com []string
		var loc []string

		for ngram := range common[n] {
			com = append(com,ngram)
		}

		for ngram := range localized[n] {
			loc = append(loc,ngram)
		}
		fmt.Println("-------------------------------")

		// Sort by intentionality

		sort.Slice(com, func(i, j int) bool {
			return SST.Intentionality(L,com[i],SST.STM_NGRAM_FREQ[n][com[i]]) > SST.Intentionality(L,com[j],SST.STM_NGRAM_FREQ[n][com[j]])
		})
		sort.Slice(loc, func(i, j int) bool {
			return SST.Intentionality(L,loc[i],SST.STM_NGRAM_FREQ[n][loc[i]]) > SST.Intentionality(L,loc[j],SST.STM_NGRAM_FREQ[n][loc[j]])
		})

		for i := range com {
			fmt.Println(n,"common",com[i],SST.Intentionality(L,com[i],SST.STM_NGRAM_FREQ[n][com[i]]))
		}

		for i := range loc {
			fmt.Println(n,"local",loc[i])
		}
	}
	
}

