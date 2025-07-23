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

var path [8][]string

//******************************************************************

func main() {

	load_arrows := true
	ctx := SST.Open(load_arrows)

	Solve(ctx)

	SST.Close(ctx)
}

//******************************************************************

func Solve(ctx SST.PoSST) {

	// Contra colliding wavefronts as path integral solver

	const maxdepth = 16
	var ldepth,rdepth int = 1,1
	var Lnum,Rnum int
	var left_paths, right_paths [][]SST.Link

	start_bc := "a7"
	end_bc := "i6"

	leftptrs := SST.GetDBNodePtrMatchingName(ctx,start_bc,"")
	rightptrs := SST.GetDBNodePtrMatchingName(ctx,end_bc,"")

	if leftptrs == nil || rightptrs == nil {
		fmt.Println("No paths available from end points")
		return
	}

	cntx := []string{""}

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {

		left_paths,Lnum = SST.GetEntireNCConePathsAsLinks(ctx,"fwd",leftptrs[0],ldepth,"",cntx)
		xleft_paths,Lnumx := SST.GetEntireConePathsAsLinks(ctx,"fwd",leftptrs[0],ldepth)	
		right_paths,Rnum = SST.GetEntireNCConePathsAsLinks(ctx,"bwd",rightptrs[0],rdepth,"",cntx)	
		xright_paths,Rnumx := SST.GetEntireConePathsAsLinks(ctx,"bwd",rightptrs[0],rdepth)		
		if Lnum != Lnumx {
			fmt.Println("LEFT sizes differ at depth",ldepth,"=",Lnum,Lnumx)
		}
		if Diff(left_paths,xleft_paths) {
			fmt.Println("LEFT SETS differ at depth",ldepth)
		}

		if Rnum != Rnumx {
			fmt.Println("RIGHT sizes differ at depth",rdepth,"=",Rnum,Rnumx)
		}
		if Diff(right_paths,xright_paths) {
			fmt.Println("RIGHT SETS differ at depth",rdepth)
		}

		ldepth++
		rdepth++
	}
}

// **********************************************************

func Diff(left,right [][]SST.Link) bool {

	var LRsplice = make(map[int]int)
	var list string

	// Return coordinate pairs of partial paths to splice

	for path := 0; path < len(left); path++ {
		for l := 0; l < len(left[path]); l++ {
			if left[path][l].Dst != right[path][l].Dst {
				fmt.Println("Mismatch at path",path,l)
				fmt.Println("Mismatch:",left[path][l],right[path][l])
				return true
			}
		}
	}

	if len(list) > 0 {
		fmt.Println("  (i.e. waves impinge",len(LRsplice),"times at: ",list,")\n")
	}

	return false
}







