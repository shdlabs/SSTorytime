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

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	context := []string{""}
	chapter := "maze"

	arrow := SST.GetDBArrowByName(ctx,"fwd")

	UseGetAppointmentArrayByArrow(ctx,arrow,chapter,context)

	SST.Close(ctx)
}

//******************************************************************

func UseGetAppointmentArrayByArrow(ctx SST.PoSST,arrow SST.ArrowPtr,chapter string,context []string) {

	var ama map[SST.ArrowPtr][]SST.Appointment

	ama = SST.GetAppointedNodesByArrow(ctx,arrow,context,chapter,2)

	fmt.Println(ama)

	fmt.Println("--------------------------------------------------")
	fmt.Println("FEATURE: GetAppointmentArrayByArrow:")
	fmt.Println(" return a map of all the nodes in chap,context that are pointed to by the same type of arrow")
	fmt.Println("--------------------------------------------------")
/*
	for arrowptr := range ama {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)
		fmt.Println("\nArrow --(",arr_dir.Long,")--> points to a group of nodes with a similar role in the context of",context,"in the chapter",chapter,"\n")

		for n := 0; n < len(ama[arrowptr]); n++ {
			node := SST.GetDBNodeByNodePtr(ctx,ama[arrowptr][n])
			NewLine(n)
			fmt.Print("..  ",node.S,",")

		}
		fmt.Println()
		fmt.Println("............................................")
	}*/
}

//******************************************************************
// Tools
//******************************************************************

func NewLine(n int) {

	if n % 8 == 0 {
		fmt.Println()
	}
}

