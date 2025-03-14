//******************************************************************
//
// Exploring how to present a search text, with API
//
// Prepare:
// cd examples
// ../src/N4L-db -u chinese.n4l
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

	reader := bufio.NewReader(os.Stdin)

	var context []string

	fmt.Println("\n\nEnter chapter text (e.g. poetry,chinese, etc):")
	chaptext, _ := reader.ReadString('\n')

	for goes := 0; goes < 10; goes ++ {


		fmt.Println("Current context:",context)

		fmt.Println("\n\nEnter newcontext text:")
		cntext, _ := reader.ReadString('\n')

		if cntext != "" {
			w := strings.Split(cntext," ")
			for c := range w {
				context = append(context,strings.TrimSpace(w[c]))
			}
		}
		fmt.Println("\n\nEnter search text:")
		searchtext, _ := reader.ReadString('\n')

		context := []string{"poem"}
		
		Search(ctx,chaptext,context,searchtext)
	}

	SST.Close(ctx)
}

//******************************************************************

func Search(ctx SST.PoSST, chaptext string,context []string,searchtext string) {

	chaptext = strings.TrimSpace(chaptext)
	searchtext = strings.TrimSpace(searchtext)

	// **** Look for meaning in the arrows ***

	var ama map[SST.ArrowPtr][]SST.NodePtr
	var count int

	ama = SST.GetMatroidArrayByArrow(ctx,context,chaptext)

	fmt.Println("--------------------------------------------------")
	fmt.Println("Looking for relevant arrows by",context,chaptext)
	fmt.Println("--------------------------------------------------")
	
	for arrowptr := range ama {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)

		if SST.MatchesInContext(arr_dir.Long,context) {

			count++
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
		fmt.Println("    (No relevant matches)")
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("Looking for relevant nodes by",searchtext)
	fmt.Println("--------------------------------------------------")

	const maxdepth = 5
	
	var start_set []SST.NodePtr
	
	search_items := strings.Split(searchtext," ")
	
	for w := range search_items {
		fmt.Print("Looking for nodes like ",search_items[w],"...")
		start_set = append(start_set,SST.GetDBNodePtrMatchingName(ctx,chaptext,search_items[w])...)
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
			
				alt_paths,path_depth := SST.GetFwdPathsAsLinks(ctx,start_set[start],sttype,maxdepth)
				
				if alt_paths != nil {
					
					fmt.Printf("\n-- Forward",SST.STTypeName(sttype),"cone stories ----------------------------------\n")
					
					for p := 0; p < path_depth; p++ {
						SST.PrintLinkPath(ctx,alt_paths,p,"\nStory:")
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









