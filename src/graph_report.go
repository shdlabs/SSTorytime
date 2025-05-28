//******************************************************************
//
// Study graph properties
// 
//
//******************************************************************

package main

import (
	"fmt"
	"strings"
	"sort"
	"flag"
	"os"
        SST "SSTorytime"
)

var CHAPTER string
var CONTEXT []string
var STTYPES []int
var DEPTH int

//******************************************************************

func main() {

	Init()

	load_arrows := true
	ctx := SST.Open(load_arrows)

	chaps := SST.GetDBChaptersMatchingName(ctx,CHAPTER)

	for chap := range chaps {
		AnalyzeGraph(ctx,chaps[chap],CONTEXT,STTYPES,DEPTH) 
	}

	SST.Close(ctx)
}

//**************************************************************

func Usage() {
	
	fmt.Printf("usage: graph_report [-sttype comma separated L,C,P,N] [-depth integer] [-chapter comma separated string] [context]\n")
	flag.PrintDefaults()

	os.Exit(2)
}

//**************************************************************

func Init() []string {

	flag.Usage = Usage

	chapterPtr := flag.String("chapter", "", "a optional substring to match specific chapters")
	sttypePtr := flag.String("sttype", "+L", "link st-types e.g. L,C,P,N")
	depthPtr := flag.Int("depth", 6, "maximum probe depth for loop detection")

	flag.Parse()
	args := flag.Args()

	CHAPTER = "none"

	if *chapterPtr != "" {
		CHAPTER = *chapterPtr
	}

	if *sttypePtr != "" {

		var sttypes = make(map[int]bool)
		array := strings.Split(*sttypePtr,",")
		for t := range array {
			switch array[t] {
			case "L","+L": 
				sttypes[1] = true
			case "C","+C": 
				sttypes[2] = true
			case "E","+E": 
				sttypes[3] = true
			case "P","+P": 
				sttypes[3] = true
			case "N","+N","-N": 
				sttypes[4] = true
			case "-L": 
				sttypes[-1] = true
			case "-C": 
				sttypes[-2] = true
			case "-E": 
				sttypes[-3] = true
			case "-P": 
				sttypes[-3] = true
			default:
				fmt.Println("Unknown sttype",array[t],"(should be in { L,C,E,N } +/-)")
				os.Exit(-1)
			}
		}

		for t := range sttypes {
			STTYPES = append(STTYPES,t)
		}
	}

	DEPTH = *depthPtr

	SST.MemoryInit()

	return args
}

//******************************************************************

func AnalyzeGraph(ctx SST.PoSST,chapter string,context []string,sttypes []int,depth int) {

	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("Analysing chapter \"%s\", context %v to path length %d\n",chapter,context,depth)
	fmt.Println("----------------------------------------------------------------")

	sources,sinks := SST.GetDBSingletonBySTType(ctx,sttypes,chapter,context)

	fmt.Print("\n\n* PROCESS ORIGINS / ROOT DEPENDENCIES / PATH SOURCES for (")
	for st := range sttypes {
		fmt.Print("\"",SST.STTypeName(sttypes[st]),"\"")
	}
	fmt.Println(") in",chapter)
	fmt.Println("")

	PrintNodes(ctx,sources)

	fmt.Println("")
	fmt.Print("\n\n* FINAL END-STATES / PATH SINK NODES for (")
	for st := range sttypes {
		fmt.Print("\"",SST.STTypeName(sttypes[st]),"\"")
	}
	fmt.Println(") in",chapter)
	fmt.Println("")

	PrintNodes(ctx,sinks)

	adj,nodekey := SST.GetDBAdjacentNodePtrBySTType(ctx,sttypes,chapter,context,false)
	symb := SST.SymbolMatrix(adj)
	sadj := SST.SymmetrizeMatrix(adj)

	fmt.Println("")
	fmt.Println("* DIRECTED LOOP SEARCH:\n")
	fmt.Println("\n")

	// Find power matrices

	an := make([][][]float32,depth+1)
	sn := make([][][]string,depth+1)

	an[1] = adj
	sn[1] = symb

	for power := 2; power <= depth; power++ {

		if power % 2 == 0 {
			an[power],sn[power] = SST.SymbolicMultiply(an[power-1],adj,sn[power-1],symb)
		} else {
			an[power],sn[power] = SST.SymbolicMultiply(an[power-1],adj,sn[power-1],symb)
		}

		loop,memberlist := AnalyzePowerMatrix(ctx,sn[power])

		for m := range loop {
			length := len(strings.Split(m,")("))
			fmt.Println("  Cycle of length",length,"with members",m)
			PrintKeyNodes(ctx,memberlist[m],nodekey)
		}
	}

	fmt.Println("")
	evc := SST.ComputeEVC(sadj)

	fmt.Println("* Symmetrized Eigenvector Centrality = FLOW RESERVOIR CAPACITANCE AT EQUILIBRIUM = \n")

	PrintVector(ctx,evc,nodekey)

	// Now find the undirected graph properties 

	evctop,path := SST.FindGradientFieldTop(sadj,evc)

	fmt.Println("")
	fmt.Println("At directionless equilibrium, there are",len(evctop),"local maxima in the EVC landscape:")

	for index := 0; index < len(evc); index++ {
		fmt.Println("\n  From node",index,"has local maximum at node *",evctop[index],"*, hop distance",len(path[index])-1,"along",path[index])
		PrintKeyNodes(ctx,path[index],nodekey)
	}

}

//**************************************************************

func AnalyzePowerMatrix(ctx SST.PoSST,symbolic [][]string) (map[string]int,map[string][]int) {

	var loop = make(map[string]int)
	var memberlist = make(map[string][]int)

	for r := 0; r < len(symbolic); r++ {

		// check the diagonal

		if len(symbolic[r][r]) == 0 {
			continue
		}

		var distrib = make(map[string]int)
		var nodes []string

		vec := strings.Split(symbolic[r][r],"*")
		
		for i := 0; i < len(vec); i++ {
			distrib[vec[i]]++
		}

		var degeneracy int

		for d := range distrib {
			degeneracy = distrib[d] / 2
			break
		}

		for r := range distrib {
			nodes = append(nodes,r)
		}

		sort.Strings(nodes)
		var members string
		var membindex []int
		var v int

		for n := 0; n < len(nodes); n++ {
			members += fmt.Sprintf("(%s)",nodes[n])
			fmt.Sscanf(nodes[n],"%d",&v)
			membindex = append(membindex,v)
		}

		loop[members] = degeneracy
		memberlist[members] = membindex
	}

	return loop,memberlist
}

//**************************************************************

func PrintNodes(ctx SST.PoSST,nptrs []SST.NodePtr) {

	for n := range nptrs {
		node := SST.GetDBNodeByNodePtr(ctx,nptrs[n])
		fmt.Printf("   - NPtr(%d,%d) -> %s\n",nptrs[n].Class,nptrs[n].CPtr,node.S)
	}
}

//**************************************************************

func PrintKeyNodes(ctx SST.PoSST,m []int,nodekey []SST.NodePtr) {

	for member := range m {
		nptr := nodekey[m[member]]
		node := SST.GetDBNodeByNodePtr(ctx,nptr)
		fmt.Printf("   - where %d -> %s\n",m[member],node.S)
	}
}

//**************************************************************

func PrintVector(ctx SST.PoSST,vector []float32,nodekey []SST.NodePtr) {

	for row := 0; row < len(vector); row++ {
		nptr := nodekey[row]
		node := SST.GetDBNodeByNodePtr(ctx,nptr)
		fmt.Printf("   ( %2.2f ) <- %d = %s\n",vector[row],row,node.S)
	}
	fmt.Println()
}

//**************************************************************

func PrintMatrix(matrix [][]float32,symbolic [][]string,str string) {

	fmt.Printf("                 DIAG %s \n",str)

	for row := 0; row < len(matrix); row++ {
		for col := 0; col < len(matrix[row]); col++ {
			fmt.Printf("%2.0f ",matrix[row][col])
		}

		fmt.Printf(" %1.1f   ...",matrix[row][row])
		if matrix[row][row] > 0 {
			fmt.Printf("      %s    (loop)\n",symbolic[row][row])
		} else {
			fmt.Println()
		}
	}
	fmt.Println()
}

