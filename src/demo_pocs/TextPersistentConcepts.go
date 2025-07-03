
//
// Look for longitudinal persistence, which has worked quite well in the past
//

package main

import (
	"fmt"
//	"math"
        SST "SSTorytime"
)

//**************************************************************
// BEGIN
//**************************************************************

func main() {

	const max_class = 100

	//input := "/home/mark/Laptop/Work/SST/data_samples/MobyDick.dat"
	//input := "/home/mark/Laptop/Work/SST/data_samples/obama.dat"
	//input := "/home/mark/Laptop/Work/SST/data_samples/bede.dat"
	//input := "/home/mark/Laptop/Work/SST/data_samples/pt1.dat"
	//input := "/home/mark/Laptop/Work/SST/data_samples/Darwin.dat"
	input := "/home/mark/Laptop/Work/NLnet/SemanticKnowledgeProject/org-42/roam/how-i-org.org"

	SST.MemoryInit()

	_,L := SST.FractionateTextFile(input)

	intentions,context := SST.AssessTextSignificance(L,SST.STM_NGRAM_FREQ,SST.STM_NGRAM_LOCA)

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

