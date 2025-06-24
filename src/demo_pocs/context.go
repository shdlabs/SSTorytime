//******************************************************************
//
// Experiment with old CFEngine context gathering approach
// Look at all the ways of grabbing context
//
//******************************************************************

package main

import (
	"fmt"
	"time"
	"os"
        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)


	// Start with the classic time classes

	now := time.Now()
	c,slot := SST.DoNowt(now)
	fmt.Println("TIME_CLASSES",c,"\nSLOT",slot)
	name,_ := os.Hostname()
	fmt.Println("HOST",name)


	// Look at ad hoc text input (small language model)

	input_stream := "/home/mark/Laptop/Work/SST/data_samples/MobyDick.dat"
	input_stream = "../../../org-42/roam/how-i-org.org"

	SST.FractionateTextFile(input_stream)

	fmt.Println("INPUT_STREAM",input_stream)

	for n := SST.N_GRAM_MIN; n < SST.N_GRAM_MAX; n++ {
		for g := range SST.STM_NGRAM_RANK[n] {
			fmt.Println("ng",n,g)
		}
	}

	// When we search for something already in the db, we need to look at 
	// the EntireCone to see what concepts context joins together

	search := "pay"

	ctx_set := SST.GetDBContextsMatchingName(ctx,search)

	fmt.Println("DIRECT CONTEXT SEARCH",ctx_set)


	// Now look for nodes and their orbits

	chap := ""
	nptrs := SST.GetDBNodePtrMatchingName(ctx,search,chap)

	confidence := 2

	for i := range nptrs {
		n := SST.GetDBNodeByNodePtr(ctx,nptrs[i])
		paths,_ := SST.GetEntireConePathsAsLinks(ctx,"any",nptrs[i],confidence)
		for p := range paths {
			for d := range paths[p] {
				if len(paths[p][d].Ctx) > 1 {
					fmt.Println("from",n.S," - ",paths[p][d].Ctx)
				}
			}
		}
	}

	SST.Close(ctx)
}

