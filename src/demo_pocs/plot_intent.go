//
// Test intent assessment
//

package main

import (
	"fmt"
	"math"
)

//**************************************************************
// BEGIN
//**************************************************************

func main() {

	L := 100 // sentences

	for f := 0; f < 100; f++ {

		I := Intentionality(1,L,f)
		fmt.Printf("%d %f\n",f,I)
	}

}

//**************************************************************

func Intentionality(n,L int, freq int) float64 {

	// Compute the effective intent of a string s at a position count
	// within a document of many sentences. The weighting due to
	// inband learning uses an exponential deprecation based on
	// SST scales (see "leg" meaning).

	work := 10.0

	// measure occurrences relative to total length L in sentences

	phi := float64(freq)
	phi_0 := float64(L/20)

	// How often is too often for a concept? density/efficiency

	const rho = 1.0/20

	crit := phi/phi_0 - rho

	meaning := phi * work / (1.0 + math.Exp(crit))

	return meaning
}

