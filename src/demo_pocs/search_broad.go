//******************************************************************
//
// Try out neighbour search for all ST stypes together
//
// Prepare:
// cd examples
// ../src/N4L-db -u chinese.n4l
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
	dbname   = "sstoryline"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	nodeptrs := SST.GetDBNodePtrMatchingName(ctx,"slit","start")

	fmt.Println("Found",nodeptrs)

	for n := range nodeptrs {

		const maxdepth = 7
		context := []string{"physics","slits"}
		chapter := "multi slit"

		alt_paths,path_depth := SST.GetEntireConePathsAsLinks(ctx,"fwd",nodeptrs[n],maxdepth)
		
		if alt_paths != nil {
			
			for p := 0; p < path_depth; p++ {
				SST.PrintLinkPath(ctx,alt_paths,p,"\nStory:",chapter,context)
			}
		}
	}
}







