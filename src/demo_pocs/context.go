//******************************************************************
//
// Experiment with old CFEngine context gathering approach
//
//******************************************************************

package main

import (
	"fmt"
	"time"
        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	now := time.Now()
	c,slot := SST.DoNowt(now)

	fmt.Println("TIME_CLASSES",c,"\nSLOT",slot)


	SST.ContextFromFile("/home/mark/Laptop/Work/SST/data_samples/MobyDick.dat")

	fmt.Println("FILE_SCAN_CLASSES",SST.STM_NGRAM_RANK)

	SST.Close(ctx)
}

