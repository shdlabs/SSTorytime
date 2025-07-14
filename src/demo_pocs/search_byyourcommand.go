//******************************************************************
//
// single search string without complex options
//

/*
api for Orbits
get a visualization
notes in chinese
notes about chinese
integrate in math
arrows pe,ep, eh
paths from start to target
a1 to b6
hubs for
stories about (bjorvika)
sequences about 
notes context "not only", "come in"
containing / matching "(),"
*/

//******************************************************************

package main

import (
	"fmt"
	"os"
	"flag"
	"strings"

        SST "SSTorytime"
)

//******************************************************************

type SearchParameters struct {

	Chapter string
	Context []string
	Arrows []string
}

func main() {

	args := GetArgs()

	SST.MemoryInit()

	load_arrows := false
	ctx := SST.Open(load_arrows)

	search := DecodeSearch(args)
		
	Search(ctx,search)


	SST.Close(ctx)
}

//**************************************************************

func Usage() {
	
	fmt.Printf("usage: ByYourCommand <search request>n")
	flag.PrintDefaults()

	os.Exit(2)
}

//**************************************************************

func GetArgs() []string {

	flag.Usage = Usage

	flag.Parse()
	return flag.Args()
}

//******************************************************************

func Search(ctx SST.PoSST, search SST.SearchParameters) {

chaptext string
context []string
searchtext string

	chaptext = strings.TrimSpace(chaptext)
	searchtext = strings.TrimSpace(searchtext)

	fmt.Println("--------------------------------------------------")
	fmt.Println("Looking for relevant nodes by",searchtext)
	fmt.Println("--------------------------------------------------")

	const maxdepth = 3
	
	var start_set []SST.NodePtr
	
	search_items := strings.Split(searchtext," ")
	
	for w := range search_items {
		fmt.Print("Looking for nodes like ",search_items[w],"...")
		start_set = append(start_set,SST.GetDBNodePtrMatchingName(ctx,search_items[w],chaptext)...)
	}

	fmt.Println("   Found possible relevant nodes:",start_set)

	for start := range start_set {

		for sttype := -SST.EXPRESS; sttype <= SST.EXPRESS; sttype++ {

			name :=  SST.GetDBNodeByNodePtr(ctx,start_set[start])

			allnodes := SST.GetFwdConeAsNodes(ctx,start_set[start],sttype,maxdepth)
			
			if len(allnodes) > 1 {
				fmt.Println()
				fmt.Println("    -------------------------------------------")
				fmt.Printf("     Search text MATCH #%d via %s connection\n",start+1,SST.STTypeName(sttype))
				fmt.Printf("     (search %s => hit %s)\n",searchtext,name.S)
				fmt.Println("    -------------------------------------------")

				for l := range allnodes {
					fullnode := SST.GetDBNodeByNodePtr(ctx,allnodes[l])
					fmt.Println("     - SSType",SST.STTypeName(sttype)," cone item: ",fullnode.S,", found in",fullnode.Chap)
				}

				// Conic proper time paths
			
				alt_paths,path_depth := SST.GetFwdPathsAsLinks(ctx,start_set[start],sttype,maxdepth)
				
				if alt_paths != nil {
					
					fmt.Println("\n-- Forward (",SST.STTypeName(sttype),") cone stories ----------------------------------\n")
					
					for p := 0; p < path_depth; p++ {
						SST.PrintLinkPath(ctx,alt_paths,p,"\nStory:","",nil)
					}
				}
				fmt.Printf("     (END %d)\n",start+1)
			}
		}
	}
	
	
	// Now look at the arrow content
	
	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println("checking whether any arrows also match search",searchtext,"(in any context)")
	fmt.Println("--------------------------------------------------")
	
	matching_arrows := SST.GetDBArrowsMatchingArrowName(ctx,searchtext)
	
	relns := SST.GetDBNodeArrowNodeMatchingArrowPtrs(ctx,chaptext,context,matching_arrows)
	
	for r := range relns {
		
		from := SST.GetDBNodeByNodePtr(ctx,relns[r].NFrom)
		to := SST.GetDBNodeByNodePtr(ctx,relns[r].NTo)
		arr := SST.ARROW_DIRECTORY[relns[r].Arr].Long
		wgt := relns[r].Wgt
		actx := relns[r].Ctx
		fmt.Println("   See also: ",from.S,"--(",arr,")->",to.S,"\n       (... wgt",wgt,"in the contexts",actx,")\n")
		
	}
}









