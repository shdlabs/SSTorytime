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

	fmt.Println("Node causal forward cone doors.N4l")

	const maxdepth = 8
	sttype := SST.LEADSTO

	levels := make([][]SST.NodePtr,maxdepth)

	// Get the start node

	start_set := SST.GetNodePtrMatchingName(ctx,"start")

	for start := range start_set {

		fmt.Println(" ---------------------------------")
		fmt.Println(" - Total forward cone from: ",start_set[start])
		fmt.Println(" ---------------------------------")

		val := SST.GetFwdConeAsNodes(ctx,start_set[start],sttype,maxdepth)
		
		for l := range val {
			fmt.Println("   - node",val[l])
		}
		
		for depth := 0; depth < maxdepth; depth++ {
			
			fmt.Println(" ---------------------------------")
			fmt.Println(" - Cone layers ",depth," from: ",start_set[start])
			fmt.Println(" ---------------------------------")
			
			levels[depth] = make([]SST.NodePtr,0)
			
			val := SST.GetFwdConeAsNodes(ctx,start_set[start],sttype,depth)
			
			for l := range val {
				if IsNew(val[l],levels) {
					levels[depth] = append(levels[depth],val[l])
				}
			}			
			fmt.Println("level",depth,levels[depth])
		}
		
		fmt.Println("Link proper time normal paths:")
			
		for start := range start_set {
			
			for depth := 0; depth < maxdepth; depth++ {
				
				fmt.Println("Searching paths of length",depth,"/",maxdepth,"from",start_set[start])
				
				paths := SST.GetFwdPathsAsLinks(ctx,start_set[start],sttype,depth)
				
				for p := range paths {
					
					if len(paths[p]) > 1 {
						
						fmt.Println("    Path",p," len",len(paths[p]))
						
						for l := 0; l < len(paths[p]); l++ {
							fmt.Print(" --> ",paths[p][l].Dst," weight",paths[p][l].Wgt)
						}
						
						fmt.Println()
					}
				}
			}
		}
	}		
	SST.Close(ctx)
}

//******************************************************************

func IsNew(nptr SST.NodePtr,levels [][]SST.NodePtr) bool {

	for l := range levels {
		for e := range levels[l] {
			if levels[l][e] == nptr {
				return false
			}
		}
	}
	return true
}








