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

	var n1,n2,n3,n4,n5,n6 SST.Node
	var lnk12,lnk23,lnk34,lnk25,lnk56 SST.Link
	
	n1.NPtr = SST.NodePtr{ CPtr : 1, Class: SST.LT128}
	n1.S = "Mary had a little lamb"
	n1.Chap = "home and away"

	n2.NPtr = SST.NodePtr{ CPtr : 2, Class: SST.LT128}
	n2.S = "Whose fleece was dull and grey"
	n2.Chap = "home and away"

	n3.NPtr = SST.NodePtr{ CPtr : 3, Class: SST.LT128}
	n3.S = "And everytime she washed it clean"
	n3.Chap = "home and away"

	n4.NPtr = SST.NodePtr{ CPtr : 4, Class: SST.LT128}
	n4.S = "It just went to roll in the hay"
	n4.Chap = "home and away"

	n5.NPtr = SST.NodePtr{ CPtr : 5, Class: SST.LT128}
	n5.S = "And every bar that Mary went"
	n5.Chap = "home and away"

	n6.NPtr = SST.NodePtr{ CPtr : 6, Class: SST.LT128}
	n6.S = "Was hot and loud and gay"
	n6.Chap = "home and away"

	lnk12.Arr = 77
	lnk12.Wgt = 0.34
	lnk12.Ctx = []string{"fairy castles","angel air"}
	lnk12.Dst = n2.NPtr

	lnk23.Arr = 77
	lnk23.Wgt = 0.34
	lnk23.Ctx = []string{"fairy castles","angel air"}
	lnk23.Dst = n3.NPtr

	lnk34.Arr = 77
	lnk34.Wgt = 0.34
	lnk34.Ctx = []string{"fairy castles","angel air"}
	lnk34.Dst = n4.NPtr

	lnk25.Arr = 77
	lnk25.Wgt = 0.34
	lnk25.Ctx = []string{"steamy hot tubs"}
	lnk25.Dst = n5.NPtr

	lnk56.Arr = 77
	lnk56.Wgt = 0.34
	lnk56.Ctx = []string{"steamy hot tubs","lady gaga"}
	lnk56.Dst = n6.NPtr

	sttype := SST.LEADSTO

	n1 = SST.CreateDBNode(ctx, n1)
	n2 = SST.CreateDBNode(ctx, n2)
	SST.AppendDBLinkToNode(ctx,n1.NPtr,lnk12,n2.NPtr,sttype)

	n2 = SST.CreateDBNode(ctx, n2)
	n3 = SST.CreateDBNode(ctx, n3)
	SST.AppendDBLinkToNode(ctx,n2.NPtr,lnk23,n3.NPtr,sttype)

	n3 = SST.CreateDBNode(ctx, n3)
	n4 = SST.CreateDBNode(ctx, n4)
	SST.AppendDBLinkToNode(ctx,n3.NPtr,lnk34,n4.NPtr,sttype)

	n2 = SST.CreateDBNode(ctx, n2)
	n5 = SST.CreateDBNode(ctx, n5)
	SST.AppendDBLinkToNode(ctx,n2.NPtr,lnk25,n5.NPtr,sttype)

	n5 = SST.CreateDBNode(ctx, n5)
	n6 = SST.CreateDBNode(ctx, n6)
	SST.AppendDBLinkToNode(ctx,n5.NPtr,lnk56,n6.NPtr,sttype)

	for depth := 0; depth < 4; depth++ {
		val := SST.GetFwdConeAsNodes(ctx,n1.NPtr,sttype,depth)
		fmt.Println("As NodePtr(s) fwd from",n1,"depth",depth)
		for l := range val {
			fmt.Println("   - Step",val[l])
		}

	}

	for depth := 0; depth < 4; depth++ {
		val := SST.GetFwdConeAsLinks(ctx,n1.NPtr,sttype,depth)
		fmt.Println("As Links fwd from",n1,"depth",depth)
		for l := range val {
			fmt.Println("   - Step",val[l])
		}
	}

	SST.Close(ctx)
}








