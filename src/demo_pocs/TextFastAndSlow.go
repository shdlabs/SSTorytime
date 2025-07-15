//
// Scan a document and determine rate separation of ngrams
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

	slow,fast,pts := SST.AssessTextFastSlow(L,SST.STM_NGRAM_LOCA)

	var grad_amb [SST.N_GRAM_MAX]map[string]float64
	var grad_int [SST.N_GRAM_MAX]map[string]float64

	for n := 1; n < SST.N_GRAM_MAX; n++ {
		grad_amb[n] = make(map[string]float64)
		grad_int[n] = make(map[string]float64)
	}

	for p := 0; p < pts; p++ {

		for n := 1; n < SST.N_GRAM_MAX; n++ {

			var amb []string
			var intent []string

			for ngram := range fast[n][p] {
				intent = append(intent,ngram)
			}

			for ngram := range slow[n][p] {
				amb = append(amb,ngram)
			}
			
			fmt.Println("----- PARTITION ",p," --------------------------")
			
			// Sort by intentionality
			
			sort.Slice(amb, func(i, j int) bool {
				return SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]]) > SST.StaticIntentionality(L,amb[j],SST.STM_NGRAM_FREQ[n][amb[j]])
			})
			sort.Slice(intent, func(i, j int) bool {
				return SST.StaticIntentionality(L,intent[i],SST.STM_NGRAM_FREQ[n][intent[i]]) > SST.StaticIntentionality(L,intent[j],SST.STM_NGRAM_FREQ[n][intent[j]])
			})
			
			for i := 0 ; i < 150 && i < len(amb); i++ {
				v := SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]])
				fmt.Println(n,"slow: ",amb[i],"       ",v)
				grad_amb[n][amb[i]] += v
			}
			
			for i := 0 ; i < 150 && i < len(intent); i++ {
				v := SST.StaticIntentionality(L,intent[i],SST.STM_NGRAM_FREQ[n][intent[i]])
				fmt.Println(n,"fast: ",intent[i],"       ",v)
				grad_int[n][intent[i]] += v
			}
		}
	}
	
	fmt.Println("\n===========================================\n")
	fmt.Println(" SUMMARY")
	fmt.Println("\n===========================================\n")
	
	for n := 1; n < SST.N_GRAM_MAX; n++ {
		
		var amb []string
		var intent []string
				
		// there is possible overlap

		for ngram := range grad_int[n] {
			_,dup := grad_amb[n][ngram]
			if dup {
				continue
			}
			intent = append(intent,ngram)
		}

		for ngram := range grad_amb[n] {
			amb = append(amb,ngram)
		}

		// Sort by intentionality
		
		sort.Slice(amb, func(i, j int) bool {
			return SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]]) > SST.StaticIntentionality(L,amb[j],SST.STM_NGRAM_FREQ[n][amb[j]])
		})
		sort.Slice(intent, func(i, j int) bool {
			return SST.StaticIntentionality(L,intent[i],SST.STM_NGRAM_FREQ[n][intent[i]]) > SST.StaticIntentionality(L,intent[j],SST.STM_NGRAM_FREQ[n][intent[j]])
		})
		
		for i := 0 ; i < 150 && i < len(amb); i++ {
			v := SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]]) / float64(n)
			fmt.Println(n,"slow v context: ",amb[i],"       ",v)
		}
		fmt.Println()
		for i := 0 ; i < 150 && i < len(intent); i++ {
			v := SST.StaticIntentionality(L,intent[i],SST.STM_NGRAM_FREQ[n][intent[i]]) / float64(n)
			fmt.Println(n,"fast v intentional: ",intent[i],"       ",v)
		}
		fmt.Println()
	}	
}

	