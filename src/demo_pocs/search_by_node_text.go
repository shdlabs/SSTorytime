//******************************************************************
//
// Demo of accessing SST in postgres
// Test of future (causal) cone and independent path retrieval
// (Add text to NPtr and Link reference)
//
// Prepare:
// cd examples
// ../src/N4L-db -u doors.n4l
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

		allnodes := SST.GetFwdConeAsNodes(ctx,start_set[start],sttype,maxdepth)
		
		for l := range allnodes {
			fullnode := SST.GetNodeByNodePtr(ctx,allnodes[l])
			fmt.Println("   - ",fullnode.S,"\tin",fullnode.Chap)
		}
		
		for depth := 0; depth < maxdepth; depth++ {
			
			fmt.Println(" ---------------------------------")
			fmt.Println(" - Cone layers ",depth," from: ",start_set[start])
			fmt.Println(" ---------------------------------")
			
			levels[depth] = make([]SST.NodePtr,0)
			
			allnodes := SST.GetFwdConeAsNodes(ctx,start_set[start],sttype,depth)
			
			for l := range allnodes {
				if IsNew(allnodes[l],levels) {
					levels[depth] = append(levels[depth],allnodes[l])
					fullnode := SST.GetNodeByNodePtr(ctx,allnodes[l])
					fmt.Println("   - Level",depth,fullnode.S,"\tin",fullnode.Chap)
				}
			}
		}
		
		fmt.Println("Link proper time normal paths:")
			
		for start := range start_set {
			
			for depth := 0; depth < maxdepth; depth++ {
				
				paths := SST.GetFwdPathsAsLinks(ctx,start_set[start],sttype,depth)

				if paths == nil {
					continue
				}

				fmt.Println("Searching paths of length",depth,"/",maxdepth,"from origin",start_set[start])
								
				for p := range paths {
					
					if len(paths[p]) > 1 {
						
						fmt.Print(" Path",p,"/",len(paths[p]),": ")
						
						for l := 0; l < len(paths[p]); l++ {
							fullnode := SST.GetNodeByNodePtr(ctx,paths[p][l].Dst)
							arr := SST.GetArrowByPtr(ctx,paths[p][l].Arr)
							fmt.Print(fullnode.S)
							if l < len(paths[p])-1 {
								fmt.Print("  -(",arr.Long,";",paths[p][l].Wgt,")->  ")
							}
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








