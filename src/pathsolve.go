//******************************************************************
//
// Find <end|start> transition matrix and calculate symmetries
//
//
//******************************************************************

package main

import (
	"fmt"
	"flag"
	"os"
	"sort"
	"strings"

        SST "SSTorytime"
)

//******************************************************************

var (
	BEGIN   string
	END     string
	CHAPTER string
	CONTEXT string
	VERBOSE bool
	FWD     string
	BWD     string
)

//******************************************************************

func main() {

	Init()

	load_arrows := true
	ctx := SST.Open(load_arrows)

	PathSolve(ctx,CHAPTER,CONTEXT,BEGIN,END)

}

//**************************************************************

func Usage() {
	
	fmt.Printf("usage: PathSolve [-v] -begin <string> -end <string> [-chapter string] subject [context]\n")
	flag.PrintDefaults()

	os.Exit(2)
}

//**************************************************************

func Init() []string {

	flag.Usage = Usage

	verbosePtr := flag.Bool("v", false,"verbose")
	chapterPtr := flag.String("chapter", "", "a optional string to limit to a chapter/section")
	beginPtr := flag.String("begin", "", "a string match start/begin set")
	endPtr := flag.String("end", "", "a string to match final end set")
	dirPtr := flag.Bool("bwd", false, "reverse search direction")

	flag.Parse()
	args := flag.Args()

	if *verbosePtr {
		VERBOSE = true
	}

	CHAPTER = ""

	if *dirPtr {
		FWD = "bwd"
		BWD = "fwd"
	} else {
		BWD = "bwd"
		FWD = "fwd"
	}

	if *beginPtr != "" {
		BEGIN = *beginPtr
	} 

	if *endPtr != "" {
		END = *endPtr
	}

	if *dirPtr {
		FWD = "bwd"
		BWD = "fwd"
	} else {
		BWD = "bwd"
		FWD = "fwd"
	}

	if *chapterPtr != "" {
		CHAPTER = *chapterPtr
	}

	if len(args) > 0 {
		isdirac,beg,end,cnt := SST.DiracNotation(args[0])

		if isdirac {
			BEGIN = beg
			END = end
			CONTEXT = cnt
		}
	} 

	SST.MemoryInit()

	return args
}

//******************************************************************

func PathSolve(ctx SST.PoSST, chapter,cntext,begin, end string) {

	const maxdepth = 15
	var Lnum,Rnum int
	var count int
	var left_paths, right_paths [][]SST.Link

	start_bc := []string{begin}
	end_bc := []string{end}
	context := strings.Split(cntext,",")

	var leftptrs,rightptrs []SST.NodePtr

	for n := range start_bc {
		leftptrs = append(leftptrs,SST.GetDBNodePtrMatchingName(ctx,start_bc[n],chapter)...)
	}

	for n := range end_bc {
		rightptrs = append(rightptrs,SST.GetDBNodePtrMatchingName(ctx,end_bc[n],chapter)...)
	}

	if leftptrs == nil || rightptrs == nil {
		fmt.Println("No paths available from end points",begin,"TO",end,"in chapter",chapter)
		return
	}

	fmt.Printf("\n\n Paths < end_set= {%s} | {%s} = start set>\n\n",ShowNode(ctx,rightptrs),ShowNode(ctx,leftptrs))

	// Find the path matrix

	var solutions [][]SST.Link
	var ldepth,rdepth int = 1,1
	var betweenness = make(map[string]int)

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {

		left_paths,Lnum = SST.GetEntireNCSuperConePathsAsLinks(ctx,FWD,leftptrs,ldepth,chapter,context)
		right_paths,Rnum = SST.GetEntireNCSuperConePathsAsLinks(ctx,BWD,rightptrs,rdepth,chapter,context)
		solutions,_ = SST.WaveFrontsOverlap(ctx,left_paths,right_paths,Lnum,Rnum,ldepth,rdepth)

		if len(solutions) > 0 {

			for s := 0; s < len(solutions); s++ {
				prefix := fmt.Sprintf(" - story path: ")
				SST.PrintLinkPath(ctx,solutions,s,prefix,"",nil)
				betweenness = TallyPath(ctx,solutions[s],betweenness)
			}
			count++
			break
		}

		if turn % 2 == 0 {
			ldepth++
		} else {
			rdepth++
		}
	}

	if len(solutions) == 0 {
		fmt.Println("No paths satisfy constraints",context," between end points",begin,"TO",end,"in chapter",chapter)
		os.Exit(-1)
	}

	// Calculate the node layer sets S[path][depth]

	fmt.Println(" *\n *\n * PATH ANALYSIS: into node flow equivalence groups\n *\n *\n")

	var supernodes [][]SST.NodePtr

	for depth := 0; depth < maxdepth*2; depth++ {

		for p_i := 0; p_i < len(solutions); p_i++ {

			if depth == len(solutions[p_i])-1 {
				supernodes = SST.Together(supernodes,solutions[p_i][depth].Dst,solutions[p_i][depth].Dst)
			}

			if depth > len(solutions[p_i])-1 {
				continue
			}

			supernodes = SST.Together(supernodes,solutions[p_i][depth].Dst,solutions[p_i][depth].Dst)

			for p_j := p_i+1; p_j < len(solutions); p_j++ {

				if depth < 1 || depth > len(solutions[p_j])-2 {
					break
				}

				if solutions[p_i][depth-1].Dst == solutions[p_j][depth-1].Dst && 
				   solutions[p_i][depth+1].Dst == solutions[p_j][depth+1].Dst {
					   supernodes = SST.Together(supernodes,solutions[p_i][depth].Dst,solutions[p_j][depth].Dst)
				}
			}
		}		
	}

	// *** Summarize paths

	for g := range supernodes {
		fmt.Print("\n    - Super node ",g," = {")
		for n := range supernodes[g] {
			node :=SST.GetDBNodeByNodePtr(ctx,supernodes[g][n])
			fmt.Print(node.S,",")
		}
		fmt.Println("}")
	}

	fmt.Println(" *\n *\n * FLOW IMPORTANCE:\n *\n *\n")

	var inv = make(map[int][]string)
	var order []int

	for key := range betweenness {
		inv[betweenness[key]] = append(inv[betweenness[key]],key)
	}

	for key := range inv {
		order = append(order,key)
	}

	sort.Ints(order)

	for key := len(order)-1; key >= 0; key-- {
		fmt.Printf("\n    -Rank (betweenness centrality): %.2f - ",float64(order[key])/float64(len(solutions)))
		for el := range inv[order[key]] {
			fmt.Print(inv[order[key]][el],",")
		}
		fmt.Println()
	}
	
}

// **********************************************************

func TallyPath(ctx SST.PoSST,path []SST.Link,between map[string]int) map[string]int {

	// count how often each node appears in the different path solutions

	for leg := range path {
		n := SST.GetDBNodeByNodePtr(ctx,path[leg].Dst)
		between[n.S]++
	}

	return between
}

// **********************************************************

func ShowNode(ctx SST.PoSST,nptr []SST.NodePtr) string {

	var ret string

	for n := range nptr {
		node := SST.GetDBNodeByNodePtr(ctx,nptr[n])
		ret += fmt.Sprintf("%.30s, ",node.S)
	}

	return ret
}







