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
	//"strings"
	//"encoding/json"

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := true
	ctx := SST.Open(load_arrows)

	context := []string{""}

	Story(ctx,"notes on chinese",context,1)
	fmt.Println("\n........")
	Story(ctx,"notes on chinese",context,2)
	fmt.Println("\n........")
	Story(ctx,"notes on chinese",context,3)

	SST.Close(ctx)
}

//******************************************************************

func Story(ctx SST.PoSST,chapter string,context []string,page int) {

	var last string

	notes := SST.GetDBPageMap(ctx,chapter,context,page)

	for n := 0; n < len(notes); n++ {

		if last != notes[n].Chapter {
			fmt.Println("\nTitle:", notes[n].Chapter)
			fmt.Println("Context:", notes[n].Context)
			last = notes[n].Chapter
		}

		for lnk := 0; lnk < len(notes[n].Path); lnk++ {
			
			text := SST.GetDBNodeByNodePtr(ctx,notes[n].Path[lnk].Dst)
			
			if lnk == 0 {
				fmt.Print("\n",text.S," ")
			} else {
				arr := SST.GetDBArrowByPtr(ctx,notes[n].Path[lnk].Arr)
				fmt.Printf("(%s) %s ",arr.Long,text.S)
			}
		}
	}
}













