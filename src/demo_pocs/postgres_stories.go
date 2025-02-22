//******************************************************************
//
// Demo of accessing postgres with custom data structures and arrays
// converting to the package library format
//
//******************************************************************

package main

import (
	"fmt"


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

	return

        ctx := SST.Open()

	fmt.Println("Test..")

	var node SST.Node
	var nodeptr SST.NodePtr
	var lnk SST.Link
	
	SST.CreateDBNode(ctx, node)
	SST.AppendDBLinkToNode(ctx,nodeptr,lnk)

	SST.Close(ctx)
}







