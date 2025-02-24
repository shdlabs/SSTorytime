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

        ctx := SST.Open()

	fmt.Println("Reset..")

	var n1,n2 SST.Node
	var lnk SST.Link
	
	n1.NPtr = SST.NodePtr{ CPtr : 1, Class: SST.LT128}
	n1.S = "Some crucial and important fact or rumour"
	n1.Chap = "home and away"

	n2.NPtr = SST.NodePtr{ CPtr : 1, Class: SST.N3GRAM}
	n2.S = "Dolly is guity"
	n2.Chap = "home and away"

	lnk.Arr = 77
	lnk.Wgt = 0.34
	lnk.Ctx = []string{"fairy castles","angel air"}
	lnk.Dst = n2.NPtr
	sttype := 2

	SST.CreateDBNode(ctx, n1)
	SST.CreateDBNode(ctx, n2)
	SST.AppendDBLinkToNode(ctx,n1.NPtr,lnk,sttype)

	SST.Close(ctx)
}








