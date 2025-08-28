package main

import (
        SST "SSTorytime"
	"fmt"
)


func main() {

	ctx := SST.Open(false)

	l := SST.GetTOCStats(ctx)

	for r := range l {
		fmt.Println(l[r])
	}

	var nptr SST.NodePtr
	nptr.Class=2;
	nptr.CPtr=581

	x := SST.GetNPtrStats(ctx,nptr)
	fmt.Println("X",x)

	SST.Close(ctx)
}

