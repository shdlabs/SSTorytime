//
// Scan a document and pick out the sentences that are measured to
// be high in "intentionality" or potential knowledge significance
// using two methods: dynamic running and static posthoc assessment
//

package main

import (
	"os"
	"fmt"
	"sort"
	"flag"
	"strings"
        SST "SSTorytime"
)

var TARGET_PERCENT float64 = 50.0

//**************************************************************
// BEGIN
//**************************************************************

func main() {

	const max_class = 100

	input := GetArgs()

	RipFile2File(input,TARGET_PERCENT)

}

//**************************************************************

func GetArgs() string {

	flag.Usage = Usage

	limitPtr := flag.Float64("%", 50, "approximate percentage of file to skim (overestimates for small values)")

	flag.Parse()
	args := flag.Args()

	TARGET_PERCENT = *limitPtr

	if len(args) != 1 {
		fmt.Println("Missing pure text filename to scan")
		os.Exit(-2)
	} 

	return args[0]
}

//**************************************************************

func Usage() {
	
	fmt.Printf("usage: Text2N4L [-v] [-% percent] filename\n")
	flag.PrintDefaults()

	os.Exit(2)
}

//*******************************************************************

func RipFile2File(filename string,percentage float64){

	SST.MemoryInit()

	psf,L := SST.FractionateTextFile(filename)

	ranking1 := SelectByRunningIntent(psf,L,percentage)
	ranking2 := SelectByStaticIntent(psf,L,percentage)

	selection := MergeText(ranking1,ranking2)

	// save result

	WriteOutput(filename,selection,L,percentage)
}

//*******************************************************************

func WriteOutput(filename string,selection []SST.TextRank,L int, percentage float64) {

	outputfile := filename + "_edit_me.n4l"

	fp, err := os.Create(outputfile)

	if err != nil {
		fmt.Println("Failed to open file for writing: ",outputfile)
		os.Exit(-1)
	}

	defer fp.Close()

	fmt.Fprintf(fp," - Samples from %s\n",filename)

	fmt.Fprintf(fp,"\n# (begin) ************\n")

	fmt.Fprintf(fp,"\n :: _sequence_ , %s::\n", filename)

	var parts = make(map[string]bool)
	
	for i := range selection {
		fmt.Fprintf(fp,"\n@sen%d   %s\n",selection[i].Order,Sanitize(selection[i].Fragment))
		part := PartName(selection[i].Partition,filename)
		fmt.Fprintf(fp,"              \" (is in) %s\n",part)
		parts[part] = true
	}
	
	fmt.Fprintf(fp,"\n# (end) ************\n")
	
	fmt.Fprintf(fp,"\n# Final fraction %.2f of requested %.2f\n",float64(len(selection)*100)/float64(L),percentage)
	
	fmt.Fprintf(fp,"\n# Selected %d samples of %d: ",len(selection),L)
	
	for i := range selection {
		fmt.Fprintf(fp,"%d ",selection[i].Order)
		}
	
	fmt.Fprintf(fp,"\n#\n")



	fmt.Println("Wrote file",outputfile)
}

//*******************************************************************

func PartName(n int,s string) string {

	return fmt.Sprintf("part %d of %s",n,s)
}

//*******************************************************************

func Sanitize(s string) string {

	s = strings.Replace(s,"(","[",-1)
	s = strings.Replace(s,")","]",-1)
	return s
}

//*******************************************************************

func SelectByRunningIntent(psf [][][]string,L int,percentage float64) []SST.TextRank {

	// Rank sentences

	const coherence_length = SST.DUNBAR_30   // approx narrative range or #sentences before new point/topic

	var sentences []SST.TextRank
	var sentence_counter int

	for p := range psf {

		for s := range psf[p] {

			score := 0.0
			text := ""

			for f := 0; f < len(psf[p][s]); f++ {

				score += SST.RunningIntentionality(sentence_counter,psf[p][s][f])

				text += psf[p][s][f]

				if f < len(psf[p][s])-1 {
					text += ", "
				}
			}

			var this SST.TextRank
			this.Fragment = text
			this.Significance = score
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences,this)
			sentence_counter++
		}
	}

	skimmed := OrderAndRank(sentences,percentage)

	return skimmed
}

// ***************************************************

func SelectByStaticIntent(psf [][][]string,L int,percentage float64) []SST.TextRank {

	// Rank sentences

	const coherence_length = SST.DUNBAR_30   // approx narrative range or #sentences before new point/topic

	var sentences []SST.TextRank
	var sentence_counter int

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
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences,this)
			sentence_counter++
		}
	}

	skimmed := OrderAndRank(sentences,percentage)

	return skimmed
}

//*********************************************************************************

func OrderAndRank(sentences []SST.TextRank,percentage float64) []SST.TextRank {

	var selections []SST.TextRank

	// Order by intentionality first to skim cream

	sort.Slice(sentences, func(i, j int) bool {
		return sentences[i].Significance > sentences[j].Significance
	})

	// Measure relative threshold for percentage of document
	// the lower the threshold, the lower the significance of the document

	threshold := percentage / 100.0

	limit := int(threshold * float64(len(sentences)))

	// Skim

	for i := 0; i < limit; i++ {
		selections = append(selections,sentences[i])
	}

	// Order by line number again to restore causal order

	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})

	return selections
}

//*********************************************************************************

func MergeText(one []SST.TextRank,two []SST.TextRank) []SST.TextRank{

	var merge []SST.TextRank
	var already_selected = make(map[int]bool)
	
	for i := range one {
		merge = append(merge,one[i])
		already_selected[one[i].Order] = true
	}

	for i := range two {
		if !already_selected[two[i].Order] {
			merge = append(merge,two[i])
		}
	}

	// Order by line number again to restore causal order

	sort.Slice(merge, func(i, j int) bool {
		return merge[i].Order < merge[j].Order
	})

	return merge
}


