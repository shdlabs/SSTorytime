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

	// Contra colliding wavefronts as path integral solver

	const maxdepth = 4
	var ldepth,rdepth int = 1,1
	var left_paths, right_paths [][]SST.Link

	leftptrs := SST.GetDBNodePtrMatchingName(ctx,"slit","start")
	rightptrs := SST.GetDBNodePtrMatchingName(ctx,"slit","target 2")

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {

		for n := range leftptrs {
			left_paths,_ = SST.GetEntireConePathsAsLinks(ctx,"fwd",leftptrs[n],ldepth)
		}
		
		for n := range rightptrs {
			right_paths,_ = SST.GetEntireConePathsAsLinks(ctx,"bwd",rightptrs[n],rdepth)		
		}
		
		j := WaveFrontsOverlap(left_paths,right_paths,maxdepth)
		
		if len(j) > 0 {
			
			fmt.Println("Waves meet....\n")

			for p := 0; p < len(j); p++ {
				SST.PrintLinkPath(ctx,j,p,"  Story:","",nil)
			}
			
		}

		if turn % 2 == 0 {
			ldepth++
		} else {
			rdepth++
		}
	}
}

// **********************************************************

func WaveFrontsOverlap(left_paths,right_paths [][]SST.Link,depth int) [][]SST.Link {
	
	var leftfront,rightfront,joinpts []SST.NodePtr
	var joins [][]SST.Link

	if left_paths != nil && right_paths != nil {

		var lp,rp int

		for lp < len(left_paths) && rp < len(right_paths) {		

			for l := 1; l < len(left_paths[lp]); l++ {
				leftfront = append(leftfront,left_paths[lp][l].Dst)
			}

			for r := 1; r < len(right_paths[rp]); r++ {
				rightfront = append(rightfront,right_paths[rp][r].Dst)
			}
	
			joinpts = NodesOverlap(leftfront,rightfront)
			
			if len(joinpts) > 0 {

				var joinpath []SST.Link

				for l := 0; l < len(left_paths[lp]); l++ {
					joinpath = append(joinpath,left_paths[lp][l])
				}

				for r := len(right_paths[rp])-1; r >= 0; r-- {
					joinpath = append(joinpath,right_paths[rp][r])
				}

				joins = append(joins,joinpath)
				return joins
			}
			lp++
			rp++
		}
	}

	return joins
}

// **********************************************************

func NodesOverlap(left,right []SST.NodePtr) []SST.NodePtr {

	var ret []SST.NodePtr

	for l := range left {
		for r := range right {
			if left[l] == right[r] {
				if !IsIn(ret,left[l]) {
					ret = append(ret,left[l])
				}
			}
		}
	}
	return ret
}

// **********************************************************

func IsIn(set []SST.NodePtr, item SST.NodePtr) bool {

	for i := range set {
		if set[i] == item {
			return true
		}
	}
	return false
}






