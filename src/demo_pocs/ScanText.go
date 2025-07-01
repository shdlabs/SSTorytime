
//
// transform random text or book, suggesting arrow hints for N4LParser
//   Conclusion : this approach is misguided. Need something more authoritative 
//

package main

import (
	"fmt"
	"math"
        SST "SSTorytime"
)

//**************************************************************
// BEGIN
//**************************************************************

func main() {

	const max_class = 100
	var freq_dist [10][max_class]int

	input := "/home/mark/Laptop/Work/SST/data_samples/MobyDick.dat"

	SST.MemoryInit()

	psf,L := SST.FractionateTextFile(input)
	
	for n := range SST.STM_NGRAM_FREQ {
		
		maxf := 0.0
		maxI := 0.0
		
		for ngram := range SST.STM_NGRAM_FREQ[n] {

			freq := SST.STM_NGRAM_FREQ[n][ngram]
			valueI := SST.Intentionality(n,L,ngram,freq)
			
			if freq > maxf {
				maxf = freq
			}

			if valueI > maxI {
				maxI = valueI
			}

			class := int(valueI) / 50

			if class < max_class {
				freq_dist[n][class]++
			}
		}
		fmt.Println("N",n,"f=",maxf,"I=",maxI,"of",L)
	}

// plot

	for f := 1; f < 50; f++ {

		fmt.Printf("%f ",math.Log(float64(f)))
		for n := 1; n < 5; n++ {
			fmt.Printf("%f ",math.Log(float64(1+freq_dist[n][f])))
		}
		fmt.Println()
	}

	freq := SST.STM_NGRAM_FREQ[1]["you"]
	valueI := SST.Intentionality(1,L,"you",freq)
	fmt.Println("you = ",freq,valueI)

	freq = SST.STM_NGRAM_FREQ[1]["Ahab"]
	valueI = SST.Intentionality(1,L,"Ahab",freq)
	fmt.Println("Ahab = ",freq,valueI)

	freq = SST.STM_NGRAM_FREQ[1]["Ahab"]
	valueI = SST.Intentionality(1,L,"Ahab",freq)
	fmt.Println("Ahab = ",freq,valueI)

	freq = SST.STM_NGRAM_FREQ[3]["desires to paint"]
	valueI = SST.Intentionality(3,L,"desires to paint",freq)
	fmt.Println("He desires to paint = ",freq,valueI)

	freq = SST.STM_NGRAM_FREQ[1]["whale-boat"]
	valueI = SST.Intentionality(1,L,"whale-boat",freq)
	fmt.Println("whaling ship = ",freq,valueI)

	// Rank sentences

	sentence := 0

	for p := range psf {

		for s := range psf[p] {

			score := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				score += SST.AssessIntent(psf[p][s][f],L,SST.STM_NGRAM_FREQ,1)

				text += psf[p][s][f]

				if f < len(psf[p][s])-1 {
					text += ", "
				} else {
					text += ". "
				}
			}

			sentence++
			fmt.Println(sentence,score,text,"\n")
		}
	}

	// Now print upper fraction 20% say

}

