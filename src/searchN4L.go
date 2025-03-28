//******************************************************************
//
// SearchN4L - a simple command line search tool
//
// Prepare, e.g.:
// cd examples
// ../src/N4L-db -u chinese*n4l Mary.n4l doors.n4l cluedo.n4l brains.n4l 
//
//******************************************************************

package main

import (
	"fmt"
	"os"
	"strings"
	"flag"

        SST "SSTorytime"
)

//******************************************************************

var (
	CHAPTER string
	SUBJECT string
	CONTEXT []string
	VERBOSE bool
	LIMIT int
)

//******************************************************************

func main() {

	Init()

	load_arrows := true
	ctx := SST.Open(load_arrows)

	Search(ctx,CHAPTER,CONTEXT,SUBJECT,LIMIT)

	SST.Close(ctx)
}

//**************************************************************

func Init() []string {

	flag.Usage = Usage

	verbosePtr := flag.Bool("v", false,"verbose")
	chapterPtr := flag.String("chapter", "any", "a optional string to limit to a chapter/section")
	limitPtr := flag.Int("limit", 20, "an approximate limit on the number of items returned, where applicable")

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		Usage()
		os.Exit(1);
	}

	SUBJECT = args[0]

	for c := 1; c < len(args); c++ {
		CONTEXT = append(CONTEXT,args[c])
	}

	if *verbosePtr {
		VERBOSE = true
	}

	LIMIT = *limitPtr

	if *chapterPtr != "" {
		CHAPTER = *chapterPtr
	}

	SST.MemoryInit()

	return args
}

//******************************************************************

func Search(ctx SST.PoSST, chaptext string,context []string,searchtext string, limit int) {

	SearchByNodes(ctx,chaptext,context,searchtext,limit)
	SearchByArrows(ctx,chaptext,context,searchtext,limit)
	SearchStoryPaths(ctx,chaptext,context,searchtext,limit)
}

//**************************************************************

func SearchByNodes(ctx SST.PoSST, chaptext string,context []string,searchtext string, limit int) {

	chaptext = strings.TrimSpace(chaptext)
	searchtext = strings.TrimSpace(searchtext)

	fmt.Println("--------------------------------------------------")
	fmt.Println("Looking for relevant nodes by",searchtext)
	fmt.Println("--------------------------------------------------")

	var start_set []SST.NodePtr
	
	search_items := strings.Split(searchtext," ")
	
	fmt.Print("Search separately by ")

	for w := range search_items {
		fmt.Print(search_items[w],",..")
		start_set = append(start_set,SST.GetDBNodePtrMatchingName(ctx,chaptext,search_items[w])...)
	}

	fmt.Println("found",len(start_set),"possible relevant nodes:")

	for start := range start_set {

		name :=  SST.GetDBNodeByNodePtr(ctx,start_set[start])

		fmt.Println()
		fmt.Printf("#%d ",start+1)
		fmt.Printf("(search %s => %s)\n",searchtext,name.S)
		fmt.Println("--------------------------------------------")

		SearchPastAndFutureConeBySTType(ctx,start_set[start],SST.EXPRESS,limit) 
		SearchPastAndFutureConeBySTType(ctx,start_set[start],SST.LEADSTO,limit) 
		SearchPastAndFutureConeBySTType(ctx,start_set[start],SST.CONTAINS,limit) 
		SearchPastAndFutureConeBySTType(ctx,start_set[start],SST.NEAR,limit) 
	}
}

//**************************************************************

func SearchByArrows(ctx SST.PoSST, chaptext string,context []string,searchtext string, limit int) {

	// **** Look for meaning in the arrows ***

	var ama map[SST.ArrowPtr][]SST.NodePtr
	var count int

	ama = SST.GetMatroidArrayByArrow(ctx,context,chaptext)

	for arrowptr := range ama {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)

		if SST.MatchesInContext(arr_dir.Long,context) {

			count++

			if count > limit {
				break
			}

			fmt.Println("\nArrow --(",arr_dir.Long,")--> points to a group of nodes with a similar role in the context of",context,"in the chapter",chaptext,"\n")
			
			for n := 0; n < len(ama[arrowptr]); n++ {
				node := SST.GetDBNodeByNodePtr(ctx,ama[arrowptr][n])
				SST.NewLine(n)
				fmt.Print("..  ",node.S,",")
				
			}
			fmt.Println()
			fmt.Println("............................................")
		}
	}

	if count == 0 {
		fmt.Println("\n(No relevant matroid patterns matching by arrow)")
	}
}

//**************************************************************

func SearchStoryPaths(ctx SST.PoSST, chaptext string,context []string,searchtext string, limit int) {

	fmt.Println()
	fmt.Println("Check for story paths of length",limit)
	
	matching_arrows := SST.GetDBArrowsMatchingArrowName(ctx,searchtext)
	
	relns := SST.GetDBNodeArrowNodeMatchingArrowPtrs(ctx,chaptext,context,matching_arrows)
	
	for r := range relns {
		
		from := SST.GetDBNodeByNodePtr(ctx,relns[r].NFrom)
		to := SST.GetDBNodeByNodePtr(ctx,relns[r].NTo)
		arr := SST.ARROW_DIRECTORY[relns[r].Arr].Long
		wgt := relns[r].Wgt
		actx := relns[r].Ctx
		fmt.Println("\n",r,": ",from.S,"--(",arr,")->",to.S,"\n       (... wgt",wgt,"in the contexts",actx,")")
	}

	if len(relns) == 0 {
		fmt.Println("No stories")
	}
}

//**************************************************************

func SearchPastAndFutureConeBySTType(ctx SST.PoSST,node SST.NodePtr,sttype int,limit int) {
		
	// Look for both directions

	allnodes := SST.GetFwdConeAsNodes(ctx,node,sttype,limit)
	
	if len(allnodes) > 1 {
		
		for l := range allnodes {
			fullnode := SST.GetDBNodeByNodePtr(ctx,allnodes[l])
			fmt.Printf(" %s: '%s'\t        (in chapter %s)\n",SST.STTypeName(sttype),fullnode.S,fullnode.Chap)
		}
		SearchStoriesFrom(ctx,node,sttype,limit)
	}

	allnodes = SST.GetFwdConeAsNodes(ctx,node,-sttype,limit)

	if len(allnodes) > 1 {
		for l := range allnodes {
			fullnode := SST.GetDBNodeByNodePtr(ctx,allnodes[l])
			fmt.Printf(" %s: '%s'\t        (in chapter %s)\n",SST.STTypeName(sttype),fullnode.S,fullnode.Chap)
		}
		SearchStoriesFrom(ctx,node,-sttype,limit)
	}

}

//**************************************************************

func SearchStoriesFrom(ctx SST.PoSST,node SST.NodePtr,sttype int,limit int) {

	// Conic proper time paths
	
	alt_paths,num_paths := SST.GetFwdPathsAsLinks(ctx,node,sttype,limit)
	
	if alt_paths != nil {
		for p := 0; p < num_paths; p++ {
			SST.PrintLinkPath(ctx,alt_paths,p,"  Story:",CHAPTER,CONTEXT)
		}
	}
}

//**************************************************************

func Usage() {
	
	fmt.Printf("usage: searchN4L [-v] [-chapter string] subject [context]\n")
	flag.PrintDefaults()
	os.Exit(2)
}








