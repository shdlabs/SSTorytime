//******************************************************************
//
// Find all nodes that start stories with a named arrow type
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

	chapter := ""
        context := []string{"WaVEs"}
	arrow := "freq"

        matches1 := SST.GetNCCNodesStartingStoriesForArrow(ctx,arrow,chapter,context)

        for p := range matches1 {

                n := SST.GetDBNodeByNodePtr(ctx,matches1[p])

                fmt.Println("Context story start with",n.S,"in",n.Chap)
        }

        matches2 := SST.GetNodesStartingStoriesForArrow(ctx,arrow)

        for p := range matches2 {

                n := SST.GetDBNodeByNodePtr(ctx,matches2[p])

                fmt.Println("Plain story start with",n.S,"in",n.Chap)
        }

	SST.Close(ctx)
}

