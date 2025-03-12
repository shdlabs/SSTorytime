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

	Analyze(ctx)

	for goes := 0; goes < 10; goes ++ {

		fmt.Println("\n\nEnter some text:")
		
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		
		Search(ctx,text)
	}

	SST.Close(ctx)
}

//******************************************************************

func Analyze(ctx SST.PoSST) {

	fmt.Println("------------------------------------------------")
	fmt.Println("Arrow that form correlators between node groups")
	fmt.Println("------------------------------------------------\n")

	var ama map[SST.ArrowPtr][]SST.NodePtr

	ama = SST.GetMatroidArrayByArrow(ctx)

	fmt.Println("FEATURE: GetMatroidArrayByArrow: raw",ama)

	for arrowptr := range ama {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)
		fmt.Println("\nArrow --(",arr_dir.Long,")--> acts as a type/meaning correlator to classify the following group:\n")
		for n := 0; n < len(ama[arrowptr]); n++ {
			node := SST.GetDBNodeByNodePtr(ctx,ama[arrowptr][n])
			fmt.Println("  Node",ama[arrowptr][n]," = ",node.S)
		}
	}

	fmt.Println("\n--------------------------------------------------------------\n")

	var ams map[int][]SST.NodePtr

	ams = SST.GetMatroidArrayBySSType(ctx)

	fmt.Println("FEATURE: GetMatroidArrayBySTType: raw",ams)

	for sttype := range ams {

		fmt.Println("\nArrow class --(",SST.STTypeName(sttype),")--> acts as a type/interpretation correlator of the following group:\n")

		for n := 0; n < len(ams[sttype]); n++ {
			node := SST.GetDBNodeByNodePtr(ctx,ams[sttype][n])
			fmt.Println("  Node",ams[sttype][n]," = ",node.S)
		}
	}

	fmt.Println("\n--------------------------------------------------------------\n")

	var ha map[SST.ArrowPtr]int
	ha = SST.GetMatroidHistogramByArrow(ctx)
	fmt.Println("FEATURE: GetMatroidHistogramByArrow",ha,"\n")

	for arrowptr := range ha {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)
		fmt.Println("Arrow -(",arr_dir.Long,")-> occurs with frequency",ha[arrowptr])
	}

	fmt.Println("\n--------------------------------------------------------------\n")

	var hs map[int]int
	hs = SST.GetMatroidHistogramBySSType(ctx)
	fmt.Println("FEATURE: GetMatroidHistogramBySTType",hs,"\n")

	for sttype := range hs {
		fmt.Println("Arrow class -(",SST.STTypeName(sttype),")-> occurs with frequency",hs[sttype])
	}


	fmt.Println("\n--------------------------------------------------------------\n")

	var ma []SST.ArrowMatroid
	ma = SST.GetMatroidNodesByArrow(ctx)

	fmt.Println("FEATURE: GetMatroidNodesByArrow",ma,"\n")

	for a := range ma {
		from := SST.GetDBNodeByNodePtr(ctx,ma[a].NFrom)
		arrow := SST.GetDBArrowByPtr(ctx,ma[a].Arr)
		fmt.Println("  - Node",from.S,"acts as a matroidal hub connecting:")
		for n := range ma[a].NTo {
			to := SST.GetDBNodeByNodePtr(ctx,ma[a].NTo[n])
			fmt.Println("   ... Node",to.S)
		}
		fmt.Println("    with arrow:  --(",arrow,")-->\n")
	}

	fmt.Println("\n--------------------------------------------------------------\n")

	var ms []SST.STTypeMatroid
	ms = SST.GetMatroidNodesBySTType(ctx)
	fmt.Println("FEATURE: GetMatroidNodesBySTType",ms,"\n")

	for s := range ms {
		from := SST.GetDBNodeByNodePtr(ctx,ms[s].NFrom)
		sttype := SST.STTypeName(ms[s].STType)
		fmt.Println("  -Node",from.S,"acts as a matroidal hub connecting:")
		for n := range ms[s].NTo {
			to := SST.GetDBNodeByNodePtr(ctx,ms[s].NTo[n])
			fmt.Println("   ... Node",to.S)
		}
		fmt.Println("   with arrow class:  --(",sttype,")-->\n")
	}

}

//******************************************************************

func Search(ctx SST.PoSST, text string) {

	text = strings.TrimSpace(text)

	const maxdepth = 5

	sttype := SST.LEADSTO
	fmt.Print("Choose a search type: ")

	for t := SST.NEAR; t <= SST.EXPRESS; t++ {
		fmt.Print(t,"=",SST.STTypeName(t),", ")
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
			fmt.Println("   - Fwd ",SST.STTypeName(sttype)," cone item: ",fullnode.S,", found in",fullnode.Chap)
		}

		alt_paths,path_depth := SST.GetFwdPathsAsLinks(ctx,start_set[start],sttype,maxdepth)
			
		if alt_paths != nil {
			
			fmt.Printf("\n-- Forward",SST.STTypeName(sttype),"cone stories ----------------------------------\n")
			
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








