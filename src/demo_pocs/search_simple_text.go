//******************************************************************
//
// Exploring how to present a search text, with API
//
// Prepare:
// cd examples
// ../src/N4L-db -u Mary.n4l, e.g. try type Mary example, type 1
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
	fmt.Print("Choose a search type: ")

	for t := SST.NEAR; t <= SST.EXPRESS; t++ {
		fmt.Print(t,"=",SST.SST_NAMES[t],", ")
	}

	fmt.Scanf("%d",&sttype)

	var start_set []SST.NodePtr

	search_items := strings.Split(text," ")

	for w := range search_items {
		start_set = append(start_set,SST.GetDBNodePtrMatchingName(ctx,search_items[w])...)
	}

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

	fmt.Println("\nLooking at relations...\n")

	matching_arrows := SST.GetDBArrowsMatchingArrowName(ctx,text)

	relns := SST.GetDBNodeArrowNodeMatchingArrowPtrs(ctx,matching_arrows)

	for r := range relns {

		from := SST.GetDBNodeByNodePtr(ctx,relns[r].NFrom)
		to := SST.GetDBNodeByNodePtr(ctx,relns[r].NFrom)
		//st := relns[r].STType
		arr := SST.ARROW_DIRECTORY[relns[r].Arr].Long
		wgt := relns[r].Wgt
		actx := relns[r].Ctx
		fmt.Println("See also: ",from.S,"--(",arr,")->",to.S,"\n       (... wgt",wgt,"in the contexts",actx,")\n")

	}
}










