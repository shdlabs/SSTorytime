//
// Scan a document and pick out the sentences that are measured to
// be high in "intentionality" or potential knowledge significance
//

package main

import (
	"fmt"
	"sort"
	"math"
        SST "SSTorytime"
)

/*  GNUPLOT

set xlabel "Fraction of whole"
set ylabel "Stability"
set title "Stability of Intentionality"
set term pdf monochrome
set key outside vertical top right
set output "average_intentionalities1.pdf"
plot "averages.in1" with errorbars

set xlabel "Fraction of whole"
set ylabel "Stability"
set title "Convergence to Stability"
set term pdf monochrome
set output "converge_intentionalities1.pdf"
set key outside vertical top right
plot "converge.in1" using 1:2 w line title "Moby Dick", "converge.in1" using 1:3 w line title "Obama", "converge.in1" using 1:4 w line title "Venerable Bede", "converge.in1" using 1:5 w line title "Promises", "converge.in1" using 1:6 w line title "Darwin", "converge.in1" using 1:7 w line title "Notes"   

set xlabel "Fraction of whole"
set ylabel "Stability"
set title "Stability of Intentionality"
set term pdf monochrome
set key outside vertical top right
set output "average_intentionalities2.pdf"
plot "averages.in2" with errorbars

set xlabel "Fraction of whole"
set ylabel "Stability"
set title "Convergence to Stability"
set term pdf monochrome
set output "converge_intentionalities2.pdf"
set key outside vertical top right
plot "converge.in2" using 1:2 w line title "Moby Dick", "converge.in2" using 1:3 w line title "Obama", "converge.in2" using 1:4 w line title "Venerable Bede", "converge.in2" using 1:5 w line title "Promises", "converge.in2" using 1:6 w line title "Darwin", "converge.in2" using 1:7 w line title "Notes"   

*/

//**************************************************************
// BEGIN
//**************************************************************

func main() {

	const max_class = 100

	var input [6]string

	input[0] = "../../examples/example_data/MobyDick.dat"
	input[1] = "../../examples/example_data/obama.dat"
	input[2] = "../../examples/example_data/bede.dat"
	input[3] = "../../examples/example_data/promisetheory1.dat"
	input[4] = "../../examples/example_data/Darwin.dat"
	input[5] = "../../examples/example_data/orgmode.dat"

	var consistency [6][10]float64

	for i := 0; i < 6; i++ {
		for thresh := 1; thresh <= 9; thresh++ {
			threshold := 0.1 * float64(thresh)
			
			set1 := Test1(input[i],threshold)
			set2 := Test2(input[i],threshold)
			
			var histo = make(map[int]int)
			
			for i := range set1 {
				histo[set1[i]]++
			}
			
			for i := range set1 {
				histo[set2[i]]++
			}
			
			var histo2 = make(map[int]int)
			
			for a := range histo {
				for i := 0; i < 3; i++ {
					if histo[a] == i {
						histo2[i]++
					}
				}
			}
			
			consistency[i][thresh] = float64(100*histo2[2])/float64(histo2[2]+histo2[1]/2)
		}
	}

	for thresh := 1; thresh <= 9; thresh++ {

		sum := 0.0
		dev := 0.0

		fmt.Printf("# %.2f ",0.1*float64(thresh))

		for i := 0; i < 6; i++ {

			sum += consistency[i][thresh]
			fmt.Printf("%f ",consistency[i][thresh])
		}

		fmt.Println()

		mean := sum/6.0

		for i := 0; i < 6; i++ {

			dev += (consistency[i][thresh]-mean)*(consistency[i][thresh]-mean)
		}

		stddev := math.Sqrt(dev/6.0)

		fmt.Printf("%.2f %f %f\n",0.1*float64(thresh), mean, stddev)
	}
}

// ***************************************************

func Test1(input string,threshold float64) []int {

	SST.MemoryInit()

	psf,_ := SST.FractionateTextFile(input)

	// Rank sentences

	var sentences []SST.TextRank
	var selections []SST.TextRank
	var count int

	for p := range psf {

		for s := range psf[p] {

			score := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				score += SST.RunningIntentionality(count,psf[p][s][f])

				text += psf[p][s][f]

				if f < len(psf[p][s])-1 {
					text += ", "
				}
			}

			var this SST.TextRank
			this.Fragment = text
			this.Significance = score
			this.Order = count
			sentences = append(sentences,this)
			count++
		}
	}

	sort.Slice(sentences, func(i, j int) bool {
		return sentences[i].Significance > sentences[j].Significance
	})

	// Measure relative threshold for percentage of document
	// the lower the threshold, the lower the significance of the document

	limit := int(threshold * float64(len(sentences)))

	for i := 0; i < limit; i++ {
		selections = append(selections,sentences[i])
	}

	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})

	// Now print only upper scoring fraction 20%

	for i := 0; i < limit; i++ {
		fmt.Print(i,"=",selections[i].Order, ": ")
		SST.ShowText(selections[i].Fragment,100)
		fmt.Println()
	}

	fmt.Println("Fraction of document = ",limit,"->", float64(limit)/float64(len(sentences)))

	var vals []int

	for i := 0; i < limit; i++ {
		vals = append(vals,selections[i].Order)
	}

	return vals
}

// ***************************************************

func Test2(input string,threshold float64) []int {

	SST.MemoryInit()

	psf,L := SST.FractionateTextFile(input)
	
	// Rank sentences

	var sentences []SST.TextRank
	var selections []SST.TextRank
	var count int

	for p := range psf {

		for s := range psf[p] {

			score := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				score += SST.AssessStaticIntent(psf[p][s][f],L,SST.STM_NGRAM_FREQ,1)

				text += psf[p][s][f]

				if f < len(psf[p][s])-1 {
					text += ", "
				}
			}

			var this SST.TextRank
			this.Fragment = text
			this.Significance = score
			this.Order = count
			sentences = append(sentences,this)
			count++
		}
	}

	sort.Slice(sentences, func(i, j int) bool {
		return sentences[i].Significance > sentences[j].Significance
	})

	// Measure relative threshold for percentage of document
	// the lower the threshold, the lower the significance of the document

	limit := int(threshold * float64(len(sentences)))

	for i := 0; i < limit; i++ {
		selections = append(selections,sentences[i])
	}

	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})

	// Now print only upper scoring fraction 20%

	for i := 0; i < limit; i++ {
		fmt.Print(i,"=",selections[i].Order, ": ")
		SST.ShowText(selections[i].Fragment,100)
		fmt.Println()
	}

	fmt.Println("Fraction of document = ",limit,"->", float64(limit)/float64(len(sentences)))

	var vals []int

	for i := 0; i < limit; i++ {
		vals = append(vals,selections[i].Order)
	}

	return vals
}
