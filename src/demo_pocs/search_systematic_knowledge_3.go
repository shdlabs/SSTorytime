//******************************************************************
//
// Exploring how to present knowledge systematically, e.g.
// e.g. review/review for an exam!
//  version 3 with axial backbone as a reference to simplify
//
//******************************************************************

package main

import (
	"fmt"
	"strings"

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	context := []string{""}
	arrow := "then"

	Story(ctx,"chinese",context,"fox",arrow)

	SST.Close(ctx)
}

//******************************************************************

func Story(ctx SST.PoSST, chapter string,context []string,searchtext string,arrname string) {

	searchtext = strings.TrimSpace(searchtext)
	stories := SST.GetStoryContainers(ctx,arrname,searchtext,chapter,context)

	for s := range stories {
		fmt.Println(stories[s])
	}
}













