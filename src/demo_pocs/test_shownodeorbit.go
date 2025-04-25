//******************************************************************
//
// Exploring how to present node text
//
//******************************************************************

package main

import (
	"fmt"

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	searchtext := "hypo"
	chaptext := "brain"
	context := []string{""}

	fmt.Println("Look for",searchtext,"\n")
	Search(ctx,chaptext,context,searchtext)

	searchtext = "S1"
	chaptext = ""
	context = []string{"physics"}

	fmt.Println("Look for",searchtext,"\n")
	Search(ctx,chaptext,context,searchtext)

	SST.Close(ctx)
}

//******************************************************************

func Search(ctx SST.PoSST, chaptext string,context []string,searchtext string) {
	
	nptrs := SST.GetDBNodePtrMatchingName(ctx,searchtext,chaptext)

	for nptr := range nptrs {
		fmt.Print(nptr,": ")
		SST.PrintNodeOrbit(ctx,nptrs[nptr],100)


	}

}









