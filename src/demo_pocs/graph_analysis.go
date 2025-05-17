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

	chapter := "notes on chinese"
        context := []string{""}

	sttypes := []int{1,2,3}
	nptrs := SST.GetDBSingletonBySTType(ctx,sttypes,chapter,context)
	PrintNodes(ctx,nptrs)

	adj,nodekey := SST.GetDBAdjacentNodePtrBySTType(ctx,sttypes,chapter,context)

	dim := len(adj)
	sadj := Symmetrize(adj)
	symb := SymbolMatrix(adj)
	ssymb := SymbolMatrix(sadj)

	v := ComputeEVC(sadj)
	PrintVector(v)

	evc,top := GetVecMax(v)

	fmt.Println("Node evc max",top,evc,nodekey[top],SST.GetDBNodeByNodePtr(ctx,nodekey[top]))
	fmt.Println("Got total symbols",dim)

	PrintMatrix(adj,symb,"A")
	PrintMatrix(sadj,ssymb,"sym")

	m2,s2 := SymbolicMultiply(adj,adj,symb,symb)

	PrintMatrix(m2,s2,"A2")
	m3,s3 := SymbolicMultiply(adj,m2,symb,s2)
	PrintMatrix(m3,s3,"A3")
	_,s4 := SymbolicMultiply(adj,m3,symb,s3)
	_,s6 := SymbolicMultiply(m3,m3,s3,s3)

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

func SymbolMatrix(m [][]float32) [][]string {
	
	var symbol [][]string
	dim := len(m)

	for r := 0; r < dim; r++ {

		var srow []string
		
		for c := 0; c < dim; c++ {

			var sym string = ""

			if m[r][c] != 0 {
				sym = fmt.Sprintf("%d*%d",r,c)
			}
			srow = append(srow,sym)
		}
		symbol = append(symbol,srow)
	}
	return symbol
}

//**************************************************************

func SymbolicMultiply(m1,m2 [][]float32,s1,s2 [][]string) ([][]float32,[][]string) {

	var m [][]float32
	var sym [][]string

	dim := len(m1)

	for r := 0; r < dim; r++ {

		var newrow []float32
		var symrow []string

		for c := 0; c < dim; c++ {

			var value float32
			var symbols string

			for j := 0; j < dim; j++ {

				if  m1[r][j] != 0 && m2[j][c] != 0 {
					value += m1[r][j] * m2[j][c]
					symbols += fmt.Sprintf("%s*%s",s1[r][j],s2[j][c])
				}
			}
			newrow = append(newrow,value)
			symrow = append(symrow,symbols)

		}
		m  = append(m,newrow)
		sym  = append(sym,symrow)
	}

	return m,sym
}

//**************************************************************

func GetSparseOccupancy(m [][]float32,dim int) []int {

	var sparse_count = make([]int,dim)

	for r := 0; r < dim; r++ {
		for c := 0; c < dim; c++ {
			sparse_count[r]+= int(m[r][c])
		}
	}

	return sparse_count
}

//**************************************************************

func Symmetrize(m [][]float32) [][]float32 {

	// CAUTION! unless we make a copy, go actually changes the original m!!! :o
	// There is some very weird pathological memory behaviour here .. but this
	// workaround seems to be stable

	var dim int = len(m)
	var symm [][]float32 = make([][]float32,dim)

	for r := 0; r < dim; r++ {
		var row []float32 = make([]float32,dim)
		symm[r] = row
	}
	
	for r := 0; r < dim; r++ {
		for c := r; c < dim; c++ {
			v := m[r][c]+m[c][r]
			symm[r][c] = v
			symm[c][r] = v
		}
	}

	return symm
}

//**************************************************************

func Transpose(matrix [][]float32) [][]float32 {

	var m [][]float32 = matrix

	for r := 0; r < len(m); r++ {
		for c := r; c < len(m); c++ {

			v := m[r][c]
			vt := m[c][r]
			m[r][c] = vt
			m[c][r] = v
		}
	}

	return m
}

//**************************************************************

func MakeInitVector(dim int,init_value float32) []float32 {

	var v = make([]float32,dim)

	for r := 0; r < dim; r++ {
		v[r] = init_value
	}

	return v
}

//**************************************************************

func MatrixOpVector(m [][]float32, v []float32) []float32 {

	var vp = make([]float32,len(m))

	for r := 0; r < len(m); r++ {
		for c := 0; c < len(m); c++ {

			if m[r][c] != 0 {
				vp[r] += m[r][c] * v[c]
			}
		}
	}
	return vp
}

//**************************************************************

func ComputeEVC(adj [][]float32) []float32 {

	v := MakeInitVector(len(adj),1.0)
	vlast := v

	const several = 10

	for i := 0; i < several; i++ {

		v = MatrixOpVector(adj,vlast)

		if CompareVec(v,vlast) < 0.1 {
			break
		}
		vlast = v
	}

	maxval,_ := GetVecMax(v)
	v = NormalizeVec(v,maxval)
	return v
}

//**************************************************************

func GetVecMax(v []float32) (float32,int) {

	var max float32 = -1
	var index int

	for r := range v {
		if v[r] > max {
			max = v[r]
			index = r
		}
	}

	return max,index
}

//**************************************************************

func NormalizeVec(v []float32, div float32) []float32 {

	if div == 0 {
		div = 1
	}

	for r := range v {
		v[r] = v[r] / div
	}

	return v
}

//**************************************************************

func CompareVec(v1,v2 []float32) float32 {

	var max float32 = -1

	for r := range v1 {
		diff := v1[r]-v2[r]

		if diff < 0 {
			diff = -diff
		}

		if diff > max {
			max = diff
		}
	}

	return max
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

		fmt.Printf("( %3.2f )\n",vector[row])
	}
	fmt.Println()
}

//**************************************************************

func PrintMatrix(matrix [][]float32,symbolic [][]string,str string) {

	fmt.Printf("                 DIAG %s \n",str)

	for row := 0; row < len(matrix); row++ {
		for col := 0; col < len(matrix[row]); col++ {
			fmt.Printf("%3.0f ",matrix[row][col])
		}

		fmt.Printf("      %1.1f   ...",matrix[row][row])
		if matrix[row][row] > 0 {
			fmt.Printf("      %s    (loop)\n",symbolic[row][row])
		} else {
			fmt.Println()
		}
	}
	fmt.Println()
}









