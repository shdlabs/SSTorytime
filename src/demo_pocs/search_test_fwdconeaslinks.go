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
//	"strings"

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	searchtext := "(limian)"
	chaptext := "chinese"
	context := []string{"city","exercise"}

	Search(ctx,chaptext,context,searchtext)

	SST.Close(ctx)
}

//******************************************************************

func Search(ctx SST.PoSST, chaptext string,context []string,searchtext string) {

	start_set := SST.GetDBNodePtrMatchingName(ctx,chaptext,searchtext)
	
	for sttype := -SST.EXPRESS; sttype <= SST.EXPRESS; sttype++ {
		
		name :=  SST.GetDBNodeByNodePtr(ctx,start_set[0])
		
		alt_paths,path_depth := SST.GetFwdPathsAsLinks(ctx,start_set[0],sttype,2)
		
		if alt_paths != nil {
			
			fmt.Println("\n-------\n",SST.STTypeName(sttype),"\n NPTR=",start_set[0]," with NAME",name,"\n-----")
			
			for p := 0; p < path_depth; p++ {

				SST.PrintLinkPath(ctx,alt_paths,p,"\nStory:","",nil)
			}
		}
	}
}






