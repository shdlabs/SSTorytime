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

	f,s,ff,ss := ExtractIntentionalTokens(L)

	fmt.Println("intentional fast by partition",f)
	fmt.Println("ambient slow by partition",s)
	fmt.Println("intentional fast summary",ff)
	fmt.Println("ambient slow summary",ss)
}

//**************************************************************

func ExtractIntentionalTokens(L int) ([][]string,[][]string,[]string,[]string) {

	slow,fast,doc_parts := SST.AssessTextFastSlow(L,SST.STM_NGRAM_LOCA)

	var grad_amb [SST.N_GRAM_MAX]map[string]float64
	var grad_oth [SST.N_GRAM_MAX]map[string]float64

	// returns

	var fastparts = make([][]string,doc_parts)
	var slowparts = make([][]string,doc_parts)
	var fastwhole []string
	var slowwhole []string

	for n := 1; n < SST.N_GRAM_MAX; n++ {
		grad_amb[n] = make(map[string]float64)
		grad_oth[n] = make(map[string]float64)
	}

	for p := 0; p < doc_parts; p++ {

		for n := 1; n < SST.N_GRAM_MAX; n++ {

			var amb []string
			var other []string

			for ngram := range fast[n][p] {
				other = append(other,ngram)
			}

			for ngram := range slow[n][p] {
				amb = append(amb,ngram)
			}
			
			// Sort by intentionality

			sort.Slice(amb, func(i, j int) bool {
				ambi :=	SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]])
				ambj := SST.StaticIntentionality(L,amb[j],SST.STM_NGRAM_FREQ[n][amb[j]])
				return ambi > ambj
			})

			sort.Slice(other, func(i, j int) bool {
				inti := SST.StaticIntentionality(L,other[i],SST.STM_NGRAM_FREQ[n][other[i]])
				intj := SST.StaticIntentionality(L,other[j],SST.STM_NGRAM_FREQ[n][other[j]])
				return inti > intj
			})
			
			for i := 0 ; i < 150 && i < len(amb); i++ {
				v := SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]])
				slowparts[p] = append(slowparts[p],amb[i])
				grad_amb[n][amb[i]] += v
			}
			
			for i := 0 ; i < 150 && i < len(other); i++ {
				v := SST.StaticIntentionality(L,other[i],SST.STM_NGRAM_FREQ[n][other[i]])
				fastparts[p] = append(fastparts[p],other[i])
				grad_oth[n][other[i]] += v
			}
		}
	}
	
	// Summary ranking of whole doc
	
	for n := 1; n < SST.N_GRAM_MAX; n++ {
		
		var amb []string
		var other []string
				
		// there is possible overlap

		for ngram := range grad_oth[n] {
			_,dup := grad_amb[n][ngram]
			if dup {
				continue
			}
			other = append(other,ngram)
		}

		for ngram := range grad_amb[n] {
			amb = append(amb,ngram)
		}

		// Sort by intentionality
		
		sort.Slice(amb, func(i, j int) bool {
			ambi := SST.StaticIntentionality(L,amb[i],SST.STM_NGRAM_FREQ[n][amb[i]])
			ambj := SST.StaticIntentionality(L,amb[j],SST.STM_NGRAM_FREQ[n][amb[j]])
			return ambi > ambj
		})
		sort.Slice(other, func(i, j int) bool {
			inti := SST.StaticIntentionality(L,other[i],SST.STM_NGRAM_FREQ[n][other[i]])
			intj := SST.StaticIntentionality(L,other[j],SST.STM_NGRAM_FREQ[n][other[j]])
			return inti > intj
		})
		
		for i := 0 ; i < 150 && i < len(amb); i++ {
			slowwhole = append(slowwhole,amb[i])
		}

		for i := 0 ; i < 150 && i < len(other); i++ {
			fastwhole = append(fastwhole,other[i])
		}
		fmt.Println()
	}	

	return fastparts,slowparts,fastwhole,slowwhole
}

	