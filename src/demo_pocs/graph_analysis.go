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
        SST "SSTorytime"
)

//******************************************************************

func main() {

/*
	var rows []string
	rows = append(rows,"0 1 0 0 0 0 0 0 0 0") // 1
	rows = append(rows,"0 0 1 0 0 0 0 0 0 0") // 2
	rows = append(rows,"0 0 0 1 0 0 0 0 0 0") // 3
	rows = append(rows,"0 1 0 0 1 0 0 0 0 0") // 4
	rows = append(rows,"0 0 0 0 0 1 0 0 0 0") // 5
	rows = append(rows,"0 0 0 0 0 0 1 0 0 0") // 6
	rows = append(rows,"0 0 0 0 0 0 0 1 0 0") // 7
	rows = append(rows,"0 0 0 0 0 1 0 0 1 1") // 8
	rows = append(rows,"0 0 0 0 0 1 0 0 0 0") // 9
 	rows = append(rows,"0 0 0 0 0 0 0 0 0 0") // 10
	adj,_ := Rows2Matrix(rows)
	var nodekey []SST.NodePtr
*/



	load_arrows := true
	ctx := SST.Open(load_arrows)

	chapter := "SSTorytime in N4L"

        context := []string{""}

	sttypes := []int{1,2,3}
	sources,sinks := SST.GetDBSingletonBySTType(ctx,sttypes,chapter,context)

	fmt.Println("---------------------------------")
	fmt.Println("\n\nSOURCES types",sttypes)
	fmt.Println("---------------------------------")

	PrintNodes(ctx,sources)

	fmt.Println("---------------------------------")
	fmt.Println("\n\nSINKS types",sttypes)
	fmt.Println("---------------------------------")

	PrintNodes(ctx,sinks)

	adj,nodekey := SST.GetDBAdjacentNodePtrBySTType(ctx,sttypes,chapter,context)

	at := SST.TransposeMatrix(adj)
	asymb := SST.SymbolMatrix(at)

	PrintMatrix(at,asymb,"A^T")

	sadj := SST.SymmetrizeMatrix(adj)
	symb := SST.SymbolMatrix(adj)
	ssymb := SST.SymbolMatrix(sadj)

	fmt.Println("---------------------------------")
	evc := SST.ComputeEVC(sadj)
	fmt.Print("EVC = capacitance at equilibrium = \n")

	PrintVector(evc)

	fmt.Println("---------------------------------")
	fmt.Println("Graph equilibrium EVC landscape:")
	fmt.Println("---------------------------------")

	evctop,path := SST.FindGradientFieldTop(sadj,evc)

	fmt.Println("\nThere are",len(evctop),"local maxima in the EVC landscape:\n")
	fmt.Println("\n  - Gradient paths:\n")

	for index := 0; index < len(evc); index++ {
		fmt.Println("--\n   From node",index,"has local max",evctop[index],"hop distance",len(path[index])-1,"along",path[index],"\n")
		for p := 0; p < len(path[index]); p++ {
			fmt.Printf("   %d = %.40s\n",path[index][p],SST.GetDBNodeByNodePtr(ctx,nodekey[path[index][p]]).S)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Println("Loop search")
	fmt.Println("---------------------------------\n")

	PrintMatrix(adj,symb,"A")
	PrintMatrix(sadj,ssymb,"sym")

	m2,s2 := SST.SymbolicMultiply(adj,adj,symb,symb)

	//SST.PrintMatrix(m2,s2,"A2")
	m3,s3 := SST.SymbolicMultiply(adj,m2,symb,s2)
	//SST.PrintMatrix(m3,s3,"A3")
	_,s4 := SST.SymbolicMultiply(adj,m3,symb,s3)
	_,s6 := SST.SymbolicMultiply(m3,m3,s3,s3)

	AnalyzePowerMatrix(ctx,s2,nodekey)
	AnalyzePowerMatrix(ctx,s3,nodekey)
	AnalyzePowerMatrix(ctx,s4,nodekey)
	AnalyzePowerMatrix(ctx,s6,nodekey)


	SST.Close(ctx)
}

//**************************************************************

func Rows2Matrix(rows []string) ([][]float32,[][]string) {

	var matrix [][]float32
	var symbol [][]string

	for r := 0; r < len(rows); r++ {

		var row []float32
		var srow []string

		cols := strings.Split(rows[r]," ")

		for c := 0; c < len(cols); c++ {

			var value float32
			var sym string = ""

			fmt.Sscanf(cols[c],"%f",&value)
			row = append(row,value)
			if value != 0 {
				sym = fmt.Sprintf("%d*%d",r,c)
			}
			srow = append(srow,sym)
		}
		matrix = append(matrix,row)
		symbol = append(symbol,srow)
	}

	return matrix,symbol
}

//**************************************************************

func AnalyzePowerMatrix(ctx SST.PoSST,symbolic [][]string,nodekey []SST.NodePtr) {

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

	for m := range loop {
		length := len(strings.Split(m,")("))
		fmt.Println("Loop of length",length,"and degeneracy",loop[m],"with members",m)
		PrintKeyNodes(ctx,memberlist[m],nodekey)
	}
}

//**************************************************************

func PrintNodes(ctx SST.PoSST,nptrs []SST.NodePtr) {

	for n := range nptrs {
		node := SST.GetDBNodeByNodePtr(ctx,nptrs[n])
		fmt.Printf("(%d,%d) = %s\n",nptrs[n].Class,nptrs[n].CPtr,node.S)
	}
}

//**************************************************************

func PrintKeyNodes(ctx SST.PoSST,m []int,nodekey []SST.NodePtr) {

	for member := range m {
		nptr := nodekey[m[member]]
		node := SST.GetDBNodeByNodePtr(ctx,nptr)
		fmt.Printf("   %d = %s\n",m[member],node.S)
	}
}

//**************************************************************

func PrintVector(vector []float32) {

	for row := 0; row < len(vector); row++ {

		fmt.Printf("   ( %2.2f )\n",vector[row])
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









