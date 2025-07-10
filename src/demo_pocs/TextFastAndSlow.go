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
	input := "../../examples/example_data/Darwin.dat"
	//input := "../../examples/example_data/orgmode.dat"

	SST.MemoryInit()

	_,L := SST.FractionateTextFile(input)  // loads STM_NGRAM*

	slow,fast,pts := SST.AssessTextFastSlow(L,SST.STM_NGRAM_LOCA)

	var grad_amb [SST.N_GRAM_MAX]map[string]int
	var grad_int [SST.N_GRAM_MAX]map[string]int

	for n := 1; n < SST.N_GRAM_MAX; n++ {
		grad_amb[n] = make(map[string]int)
		grad_int[n] = make(map[string]int)
	}

	for p := 0; p < pts; p++ {

		for n := 1; n < SST.N_GRAM_MAX; n++ {

			var amb []string
			var intent []string

			for ngram := range slow[n][p] {
				amb = append(amb,ngram)
			}
			
			for ngram := range fast[n][p] {
				intent = append(intent,ngram)
			}
			fmt.Println("----- PARTITION ",n,p," --------------------------")
			
			// Sort by intentionality
			
			sort.Slice(amb, func(i, j int) bool {
				return SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]]) > SST.StaticIntentionality(L,amb[j],SST.STM_NGRAM_FREQ[n][amb[j]])
			})
			sort.Slice(intent, func(i, j int) bool {
				return SST.StaticIntentionality(L,intent[i],SST.STM_NGRAM_FREQ[n][intent[i]]) > SST.StaticIntentionality(L,intent[j],SST.STM_NGRAM_FREQ[n][intent[j]])
			})
			
			for i := 0 ; i < 150 && i < len(amb); i++ {
				fmt.Println(n,"slow: ",amb[i],"       ",SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]]))
				grad_amb[n][amb[i]]++
			}
			
			for i := 0 ; i < 150 && i < len(intent); i++ {
				fmt.Println(n,"fast: ",intent[i],"       ",SST.StaticIntentionality(L,intent[i],SST.STM_NGRAM_FREQ[n][intent[i]]))
				grad_int[n][intent[i]]++
			}
		}
	}

	fmt.Println("\n===========================================\n")
	fmt.Println(" SUMMARY")
	fmt.Println("\n===========================================\n")

	for n := 1; n < SST.N_GRAM_MAX; n++ {
		
		for t := range grad_int[n] {
			if grad_int[n][t] > 1 {
				fmt.Println(n, "intentional theme/concept:", t, grad_int[n][t])
			}
		}

		for t := range grad_amb[n] {
			if grad_amb[n][t] > 1 {
				fmt.Println(n, "ambient theme/concept:", t, grad_amb[n][t])
			}
		}
	}

}

