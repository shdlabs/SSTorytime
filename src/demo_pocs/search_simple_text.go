//******************************************************************
//
// Exploring how to present a search text, with API
//
//******************************************************************

package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"

        SST "SSTorytime"
)

//******************************************************************

const (
	host     = "localhost"
	port     = 5432
	user     = "sstoryline"
	password = "sst_1234"
	dbname   = "newdb"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	for goes := 0; goes < 10; goes ++ {

		fmt.Println("\n\nEnter some text:")
		
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		
		Search(ctx,text)
	}

	SST.Close(ctx)
}

//******************************************************************

func Search(ctx SST.PoSST, text string) {

	text = strings.TrimSpace(text)

	const maxdepth = 5
	sttype := SST.LEADSTO

	// Simplistic text match of a node, and print some stuff

	start_set := SST.GetDBNodePtrMatchingName(ctx,text)

	for start := range start_set {

		name :=  SST.GetDBNodeByNodePtr(ctx,start_set[start])

		fmt.Println()
		fmt.Println("-------------------------------------------")
		fmt.Printf(" SEARCH MATCH %d: (%s -> %s)\n",start,text,name.S)
		fmt.Println("-------------------------------------------")

		allnodes := SST.GetFwdConeAsNodes(ctx,start_set[start],sttype,maxdepth)
		
		for l := range allnodes {
			fullnode := SST.GetDBNodeByNodePtr(ctx,allnodes[l])
			fmt.Println("   - Fwd ",SST.SST_NAMES[sttype]," cone item: ",fullnode.S,", found in",fullnode.Chap)
		}

		alt_paths,path_depth := SST.GetFwdPathsAsLinks(ctx,start_set[start],sttype,maxdepth)
			
		if alt_paths != nil {
			
			fmt.Printf("\n-- Forward",SST.SST_NAMES[sttype],"cone stories ----------------------------------\n")
			
			for p := 0; p < path_depth; p++ {
				SST.PrintLinkPath(ctx,alt_paths,p,"\nStory:")
			}
		}
		fmt.Printf("     (END %d)\n",start)
	}

	// Now look at the arrow content

	matching_arrows := SST.GetDBNodeArrowNodeMatchingArrowName(ctx,text)

	fmt.Println(matching_arrows)

}

//******************************************************************

func IsNew(nptr SST.NodePtr,levels [][]SST.NodePtr) bool {

	for l := range levels {
		for e := range levels[l] {
			if levels[l][e] == nptr {
				return false
			}
		}
	}
	return true
}








