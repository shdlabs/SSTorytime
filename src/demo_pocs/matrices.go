

package main

import (
	"fmt"
	"strings"
	"sort"
)

//**************************************************************

func main () {

	var rows []string

/*        //                  1 2 3 4 5 6 7 8 
	rows = append(rows,"0 0 0 1 0 0 0 0") // 1
	rows = append(rows,"0 0 0 1 0 0 0 0") // 2
	rows = append(rows,"0 0 0 1 0 0 0 0") // 3
	rows = append(rows,"0 0 0 0 1 0 0 0") // 4
	rows = append(rows,"0 0 0 0 0 1 0 0") // 5
	rows = append(rows,"0 0 0 0 0 0 1 1") // 6
	rows = append(rows,"0 0 0 0 0 0 0 0") // 7
	rows = append(rows,"0 0 0 0 0 0 0 0") // 8
*/

/*
        //                  1 2 3 4 5 6 7 8 9 10
	rows = append(rows,"0 1 0 0 0 0 0 0 0 0") // 1
	rows = append(rows,"0 0 1 0 0 0 0 0 0 0") // 2
	rows = append(rows,"0 0 0 1 0 0 0 0 0 0") // 3
	rows = append(rows,"0 0 0 0 1 0 0 0 0 0") // 4
	rows = append(rows,"0 0 0 0 0 1 0 0 0 0") // 5
	rows = append(rows,"0 0 0 0 0 0 1 0 0 0") // 6
	rows = append(rows,"0 0 0 0 0 0 0 1 0 0") // 7
	rows = append(rows,"0 0 0 0 0 0 0 0 1 0") // 8
	rows = append(rows,"0 0 0 0 0 0 0 0 0 1") // 9
 	rows = append(rows,"0 0 0 0 0 0 0 0 0 0") // 10
*/


/*        //                  1 2 3 4 5 6 7 8 9 10
	rows = append(rows,"0 2 0 0 0 0 0 0 0 0") // 1
	rows = append(rows,"0 0 3 0 0 0 0 0 0 0") // 2
	rows = append(rows,"0 0 0 1 0 0 0 0 0 0") // 3
	rows = append(rows,"0 1 0 0 8 0 0 0 0 0") // 4
	rows = append(rows,"0 0 0 0 0 1 0 0 0 0") // 5
	rows = append(rows,"0 0 0 0 0 0 1 0 0 0") // 6
	rows = append(rows,"0 0 0 0 0 0 0 1 0 0") // 7
	rows = append(rows,"0 0 0 0 0 9 0 0 1 1") // 8
	rows = append(rows,"0 0 0 0 0 1 0 0 0 0") // 9
 	rows = append(rows,"0 0 0 0 0 0 0 0 0 0") // 10
*/


        //                  1 2 3 4 5 6 7 8 9 10
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


/*	rows = append(rows,"0 1 0 0 0 0")
	rows = append(rows,"0 0 1 0 0 0")
	rows = append(rows,"0 0 0 1 0 0")
	rows = append(rows,"0 1 0 0 1 0")
	rows = append(rows,"0 0 0 0 0 1")
	rows = append(rows,"0 0 0 0 1 0") */
	
/*	rows = append(rows,"0 1 0 0 0")
	rows = append(rows,"0 0 1 0 0")
	rows = append(rows,"1 0 0 1 0")
	rows = append(rows,"0 0 0 0 1")
	rows = append(rows,"0 0 0 0 0")*/

/*	rows = append(rows,"1 1 0")
	rows = append(rows,"0 0 1")
	rows = append(rows,"0 0 1")*/

/*	rows = append(rows,"0 1 0")
	rows = append(rows,"1 0 1")
	rows = append(rows,"0 1 0")*/


	m,s := Matrix(rows)
	PrintMatrix(m,s,"A")

	//mt := Transpose(m)
	//PrintMatrix(mt,s,"A^T")

	m2,s2 := Mult(m,m,s,s)
	PrintMatrix(m2,s2,"A^2")
	AnalyzePowerMatrix(s2)

	m3,s3 := Mult(m,m2,s,s2)
	PrintMatrix(m3,s3,"A^3")
	AnalyzePowerMatrix(s3)

	m4,s4 := Mult(m,m3,s,s3)
	PrintMatrix(m4,s4,"A^4")
	AnalyzePowerMatrix(s4)

	m6,s6 := Mult(m3,m3,s3,s3)
	PrintMatrix(m6,s6,"A^6")
	AnalyzePowerMatrix(s6)

	fmt.Println(".....................")
	fmt.Println("M^n asymm")

	m12,s12 := Mult(m6,m6,s6,s6)
	PrintMatrix(m12,s12,"M^12")

	v := MakeInitVector(len(m12),1.0)
	v = MatrixOpVector(m12,v)
	PrintVector(v)

	fmt.Println(".....................")

	fmt.Println("EVC plain")
	evc := ComputeEVC("asymm",m)
	PrintVector(evc)

	fmt.Println(".....................")
	sm := Symmetrize(m)

	fmt.Println("EVC symm")
	evc2 := ComputeEVC("symm", sm)
	PrintVector(evc2)

}

//**************************************************************

func Mult(m1,m2 [][]float32,s1,s2 [][]string) ([][]float32,[][]string) {

	var m [][]float32
	var sym [][]string

	if len(m1[0]) != len(m2) {
		fmt.Println("Can't multiply incompatible dimensions")
		return m,sym
	}

	for r := 0; r < len(m1); r++ {

		var newrow []float32
		var symrow []string

		for c := 0; c < len(m2[0]); c++ {

			var value float32
			var symbols string

			for j := 0; j < len(m2); j++ {
				value += m1[r][j] * m2[j][c]
 
				if  m1[r][j] != 0 && m2[j][c] != 0 {
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

func ComputeEVC(s string, adj [][]float32) []float32 {

	v := MakeInitVector(len(adj),1.0)
	vlast := v
	fmt.Println("...........",s,"..........")
	//PrintMatrix(adj,"evc")

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

//**************************************************************

func AnalyzePowerMatrix(symbolic [][]string) {

	var loop = make(map[string]int)

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

		for n := 0; n < len(nodes); n++ {
			members += fmt.Sprintf("(%s)",nodes[n])
		}


		loop[members] = degeneracy
	}

	for m := range loop {
		length := len(strings.Split(m,")("))
		fmt.Println("Loop of length",length,"and degeneracy",loop[m],"with members",m)
	}
}

//**************************************************************

func PrintVector(vector []float32) {

	for row := 0; row < len(vector); row++ {

		fmt.Printf("( %3.2f )\n",vector[row])
	}
	fmt.Println()
}

