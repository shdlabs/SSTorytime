package main

import (
        SST "SSTorytime"
	"fmt"
)


func main() {

	ctx := SST.Open(false)

	l := SST.GetTOCStats(ctx)

	for r := range l {
		fmt.Println("Sec",l[r].Section,"LAST",l[r].Last,"pdel",l[r].Pdelta,"ndel",l[r].Ndelta)
	}

	SST.Close(ctx)
}

