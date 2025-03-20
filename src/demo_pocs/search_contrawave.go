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

	const maxdepth = 10
	var ldepth,rdepth int = 1,1
	var Lnum,Rnum int
	var count int
	var left_paths, right_paths [][]SST.Link

	start_bc := "A1"
	end_bc := "B6"

	leftptrs := SST.GetDBNodePtrMatchingName(ctx,"",start_bc)
	rightptrs := SST.GetDBNodePtrMatchingName(ctx,"",end_bc)

	if leftptrs == nil || rightptrs == nil {
		fmt.Println("No paths available from end points")
		return
	}

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {

		left_paths,Lnum = SST.GetEntireConePathsAsLinks(ctx,"fwd",leftptrs[0],ldepth)		
		right_paths,Rnum = SST.GetEntireConePathsAsLinks(ctx,"bwd",rightptrs[0],rdepth)		
		
		solutions := WaveFrontsOverlap(ctx,left_paths,right_paths,Lnum,Rnum,ldepth,rdepth)

		if len(solutions) > 0 {
			fmt.Println("Path solution",count,"from",start_bc,"to",end_bc,"with lengths",ldepth,-rdepth)

			for s := 0; s < len(solutions); s++ {
				SST.PrintLinkPath(ctx,solutions,s,"  Story:","",nil)
			}
			count++
		}

		if turn % 2 == 0 {
			ldepth++
		} else {
			rdepth++
		}
	}
}

// **********************************************************

func WaveFrontsOverlap(ctx SST.PoSST,left_paths,right_paths [][]SST.Link,Lnum,Rnum,ldepth,rdepth int) [][]SST.Link {

	// The wave front consists of Lnum and Rnum points left_paths[len()-1].
	// Any of the

	var solutions [][]SST.Link

	// Start expanding the waves from left and right, one step at a time, alternately

	leftfront := WaveFront(left_paths,Lnum)
	rightfront := WaveFront(right_paths,Rnum)
	
	join_match := NodesOverlap(ctx,leftfront,rightfront)
	
	for join := range join_match {

		lp := join
		rp := join_match[join]

		var LRsplice []SST.Link		

		LRsplice = LeftJoin(LRsplice,left_paths[lp])
		LRsplice = RightComplementJoin(LRsplice,right_paths[rp])
		solutions = append(solutions,LRsplice)
		return solutions
	}


	return nil
}

// **********************************************************

func WaveFront(path [][]SST.Link,num int) []SST.NodePtr {

	// assemble the cross cutting nodeptrs of the wavefronts

	var front []SST.NodePtr

	for l := 0; l < num; l++ {
		front = append(front,path[l][len(path[l])-1].Dst)
	}

	return front
}

// **********************************************************

func NodesOverlap(ctx SST.PoSST,left,right []SST.NodePtr) map[int]int {

	var LRsplice = make(map[int]int)

	for l := range left {
		for r := range right {
			if left[l] == right[r] {
				LRsplice[l] = r
			}
		}
	}
	return LRsplice
}

// **********************************************************

func LeftJoin(LRsplice,seq []SST.Link) []SST.Link {

	for i := 0; i < len(seq); i++ {

		LRsplice = append(LRsplice,seq[i])
	}

	return LRsplice
}

// **********************************************************

func RightComplementJoin(LRsplice,seq []SST.Link) []SST.Link {

	// len(seq)-1 matches the last node of right join

	for j := len(seq)-2; j >= 0; j-- {

		link := seq[j]
		link.Arr = SST.INVERSE_ARROWS[link.Arr]

		LRsplice = append(LRsplice,link)
	}

	return LRsplice
}

// **********************************************************

func ShowArray(ctx SST.PoSST,s string, a []SST.NodePtr) {

	fmt.Print(s,": ")

	for m := range a {
		n := SST.GetDBNodeByNodePtr(ctx,a[m])
		fmt.Print(n.S,",")
	}
	fmt.Println()
}




