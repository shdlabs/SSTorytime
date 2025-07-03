
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

	input := "/home/mark/Laptop/Work/SST/data_samples/MobyDick.dat"
	//input := "/home/mark/Laptop/Work/SST/data_samples/obama.dat"
	//input := "/home/mark/Laptop/Work/SST/data_samples/bede.dat"
	//input := "/home/mark/Laptop/Work/SST/data_samples/pt1.dat"

	SST.MemoryInit()

	_,L := SST.FractionateTextFile(input)

	selections := SST.AssessLongitudinalSignificance(L)

	for n := range selections {

		for ngram := range selections[n] {
			fmt.Println("-",n,selections[n][ngram])
		}
	}

}

