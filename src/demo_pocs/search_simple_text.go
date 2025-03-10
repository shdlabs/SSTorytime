//******************************************************************
//
// Exploring how to present a search text, with API
//
//******************************************************************

package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"

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

	load_arrows := false
	ctx := SST.Open(load_arrows)

	for goes := 0; goes < 10; goes ++ {

		fmt.Println("\n\nEnter some text:")
		
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		
		Search(ctx,text)
	}

	SST.Close(ctx)
}

//******************************************************************

func Search(ctx SST.PoSST, text string) {

	text = strings.TrimSpace(text)
	fmt.Printf("\n\nSearch text: (%s)\n",text)

	const maxdepth = 8
	sttype := SST.LEADSTO

	// Simplistic text match of a node, and print some stuff

	start_set := SST.GetDBNodePtrMatchingName(ctx,text)

	for start := range start_set {

		allnodes := SST.GetFwdConeAsNodes(ctx,start_set[start],sttype,maxdepth)
		
		for l := range allnodes {
			fullnode := SST.GetDBNodeByNodePtr(ctx,allnodes[l])
			fmt.Println("   - ",fullnode.S,"\t found in",fullnode.Chap)
		}
		
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
							fullnode := SST.GetDBNodeByNodePtr(ctx,paths[p][l].Dst)
							arr := SST.GetDBArrowByPtr(ctx,paths[p][l].Arr)
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








