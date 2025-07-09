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
	input := "../../examples/example_data/obama.dat"
	//input := "../../examples/example_data/bede.dat"
	//input := "../../examples/example_data/promisetheory1.dat"
	//input := "../../examples/example_data/Darwin.dat"
	//input := "../../examples/example_data/orgmode.dat"

	SST.MemoryInit()

	_,L := SST.FractionateTextFile(input)  // loads STM_NGRAM*

	ambient,condensed,_ := SST.AssessTextCoherentCoactivation(L,SST.STM_NGRAM_LOCA)

	for n := 1; n < SST.N_GRAM_MAX; n++ {

		var amb []string
		var cond []string

		for ngram := range ambient[n] {
			//if ambient[n][ngram] > 1 {
				amb = append(amb,ngram)
			//}
		}

		for ngram := range condensed[n] {
			//if condensed[n][ngram] > 1 {
				cond = append(cond,ngram)
			//}
		}
		fmt.Println("-------------------------------")

		// Sort by intentionality

		sort.Slice(amb, func(i, j int) bool {
			return SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]]) > SST.StaticIntentionality(L,amb[j],SST.STM_NGRAM_FREQ[n][amb[j]])
		})
		sort.Slice(cond, func(i, j int) bool {
			return SST.StaticIntentionality(L,cond[i],SST.STM_NGRAM_FREQ[n][cond[i]]) > SST.StaticIntentionality(L,cond[j],SST.STM_NGRAM_FREQ[n][cond[j]])
		})

		for i := 0 ; i < 150 && i < len(amb); i++ {
			fmt.Println(n,"ambient",amb[i],"       ",SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]]))
		}

		for i := 0 ; i < 150 && i < len(cond); i++ {
			fmt.Println(n,"condensate",cond[i],"       ",SST.StaticIntentionality(L,cond[i],SST.STM_NGRAM_FREQ[n][cond[i]]))
		}
	}
	
}

