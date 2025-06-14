//******************************************************************
//
// API demo HubJoin function to join a list of nodeptrs to a hub
//
//******************************************************************

package main

import (
	"fmt"
        SST "SSTorytime"
)

//******************************************************************

func main () {


	load_arrows := false
	ctx := SST.Open(load_arrows)

	names := []string{"node1","node2","node3"}
	weights := []float32{0.2, 0.4, 1.0}
	context := []string{"some","context","tags"}

	var nodes []SST.Node
	var nptrs []SST.NodePtr

	for n := range names {

		nodes = append(nodes,SST.Vertex(ctx,names[n],"my chapter"))
		nptrs = append(nptrs,nodes[n].NPtr)
	}


	created := SST.HubJoin(ctx,"","",nptrs,"then",context,weights)

	fmt.Println("Creates hub node",created)

	SST.Close(ctx)
}










