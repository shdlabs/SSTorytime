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
	adj,_ := Matrix(rows)
*/

	load_arrows := true
	ctx := SST.Open(load_arrows)

	sttypes := []int{1}

	adj,nodekey := GetDBAdjacentNodePtrBySTType(ctx,sttypes)
	sym := SymbolMatrix(adj)
	dim := len(adj)

	fmt.Println("Got symbols",dim)

	m2,s2 := Mult(adj,adj,sym,sym)
	//PrintMatrix(m2,s2,"A^2")
	AnalyzePowerMatrix(ctx,s2,nodekey)

	m3,s3 := Mult(adj,m2,sym,s2)
	//PrintMatrix(m3,s3,"A^3")
	AnalyzePowerMatrix(ctx,s3,nodekey)

	_,s4 := Mult(adj,m3,sym,s3)
	//PrintMatrix(m4,s4,"A^4")
	AnalyzePowerMatrix(ctx,s4,nodekey)

	_,s6 := Mult(m3,m3,s3,s3)
	//PrintMatrix(m6,s6,"A^6")
	AnalyzePowerMatrix(ctx,s6,nodekey)

	SST.Close(ctx)
}

//**************************************************************

func Matrix(rows []string) ([][]float32,[][]string) {

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

// **************************************************************************

func GetDBAdjacentNodePtrBySTType(ctx SST.PoSST,sttypes []int) ([][]float32,[]SST.NodePtr) {

	// Return a connected adjacency matrix for the subgraph and a lookup table
	// A bit memory intensive, but possibly unavoidable
	
	var qstr,qwhere,qsearch string
	var dim = len(sttypes)

	if dim > 4 {
		fmt.Println("Maximum 4 sttypes in GetDBAdjacentNodePtrBySTType")
		return nil,nil
	}

	for st := 0; st < len(sttypes); st++ {

		stname := SST.STTypeDBChannel(sttypes[st])
		qwhere += fmt.Sprintf("array_length(%s::text[],1) IS NOT NULL",stname)

		if st != dim-1 {
			qwhere += " AND "
		}

		qsearch += "," + stname

	}

	qstr = fmt.Sprintf("SELECT NPtr%s FROM Node WHERE %s",qsearch,qwhere)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY GetDBAdjacentNodePtrBySTType Failed",err)
		return nil,nil
	}

	var linkstr = make ([]string,dim)
	var protoadj = make(map[int][]SST.Link)
	var lookup = make(map[SST.NodePtr]int)
	var rowindex int
	var nodekey []SST.NodePtr
	var counter int

	for row.Next() {		

		var n SST.NodePtr
		var nstr string

		switch dim {

		case 1: err = row.Scan(&nstr,&linkstr[0])
		case 2: err = row.Scan(&nstr,&linkstr[0],&linkstr[1])
		case 3: err = row.Scan(&nstr,&linkstr[0],&linkstr[1],&linkstr[2])
		case 4: err = row.Scan(&nstr,&linkstr[0],&linkstr[1],&linkstr[2],&linkstr[3])

		default:
			fmt.Println("Maximum 4 sttypes in GetDBAdjacentNodePtrBySTType - shouldn't happen")
			row.Close()
			return nil,nil
		}

		if err != nil {
			fmt.Println("Error scanning sql data",err)
			row.Close()
			return nil,nil
		}

		fmt.Sscanf(nstr,"(%d,%d)",&n.Class,&n.CPtr)

		// idempotently gather nptrs into a map, keeping linked nodes close in order

		index,already := lookup[n]

		if already {
			rowindex = index
		} else {
			rowindex = counter
			lookup[n] = counter
			counter++
			nodekey = append(nodekey,n)
		}

		for lnks := range linkstr {

			links := SST.ParseLinkArray(linkstr[lnks])
			
			// we have to go through one by one to avoid duplicates
			// and keep adjacent nodes closer in order
			
			for l := range links {
				
				_,already := lookup[links[l].Dst]
				
				if !already {
					lookup[links[l].Dst] = counter
					counter++
					nodekey = append(nodekey,links[l].Dst)
				}
			}
			protoadj[rowindex] = links // now have a sparse ordered repr		
		}
	}
	
	adj := make([][]float32,counter)

	for r := 0; r < counter; r++ {

		adj[r] = make([]float32,counter)

		row := protoadj[r]

		for l := 0; l < len(row); l++ {
			lnk := row[l]
			c := lookup[lnk.Dst]
			adj[r][c] = lnk.Wgt
		}
	}
	
	row.Close()
	
	return adj,nodekey
	
}

//**************************************************************

func SymbolMatrix(m [][]float32) [][]string {
	
	var symbol [][]string
	
	for r := 0; r < len(m); r++ {
		
		var srow []string
		
		for c := 0; c < len(m); c++ {
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

func Mult(m1,m2 [][]float32,s1,s2 [][]string) ([][]float32,[][]string) {

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

func Symmetrize(matrix [][]float32) [][]float32 {

	var m [][]float32 = matrix

	for r := 0; r < len(m); r++ {
		for c := r; c < len(m); c++ {
			v := m[r][c]+matrix[c][r]
			m[r][c] = v
			m[c][r] = v
		}
	}

	return m
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

	maxval := GetVecMax(v)
	v = NormalizeVec(v,maxval)
	return v
}

//**************************************************************

func GetVecMax(v []float32) float32 {

	var max float32 = -1

	for r := range v {
		if v[r] > max {
			max = v[r]
		}
	}

	return max
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
		PrintNodes(ctx,memberlist[m],nodekey)
	}
}

//**************************************************************

func PrintNodes(ctx SST.PoSST,m []int,nodekey []SST.NodePtr) {

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









