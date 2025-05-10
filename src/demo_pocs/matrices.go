

package main

import (
	"fmt"
	"strings"
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


        //                  1 2 3 4 5 6 7 8 9 10
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


/*
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
*/

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


	m := Matrix(rows)
	PrintMatrix(m)

	mt := Transpose(m)
	PrintMatrix(mt)

	m2 := Mult(mt,mt)
	//PrintMatrix(m2)

	m3 := Mult(mt,m2)
	//PrintMatrix(m3)

	//m4 := Mult(m,m3)
	//PrintMatrix(m4)

	m6 := Mult(m3,m3)
	//PrintMatrix(m6)

	fmt.Println(".....................")
	fmt.Println("M^n asymm")

	m12 := Mult(m6,m6)
	PrintMatrix(m12)
	v := MakeInitVector(len(m12),1.0)
	v = MatrixOpVector(m12,v)
	PrintVector(v)

	fmt.Println(".....................")

	fmt.Println("EVC plain")
	evc := ComputeEVC("asymm",mt)
	PrintVector(evc)

	fmt.Println(".....................")
	s := Symmetrize(m)
	fmt.Println("symm")
	PrintMatrix(s)

	fmt.Println("EVC symm")
	evc2 := ComputeEVC("symm", s)
	PrintVector(evc2)

}

//**************************************************************

func Mult(m1,m2 [][]float64) [][]float64 {

	var m [][]float64

	if len(m1[0]) != len(m2) {
		fmt.Println("Can't multiply incompatible dimensions")
		return m
	}

	for r := 0; r < len(m1); r++ {

		var newrow []float64

		for c := 0; c < len(m2[0]); c++ {
			var value float64
			for j := 0; j < len(m2); j++ {
				value += m1[r][j] * m2[j][c]
			}
			newrow = append(newrow,value)
		}
		m  = append(m,newrow)
	}

	return m
}

//**************************************************************

func Matrix(rows []string) [][]float64 {

	var matrix [][]float64

	for r := 0; r < len(rows); r++ {

		var row []float64

		cols := strings.Split(rows[r]," ")

		for c := 0; c < len(cols); c++ {
			var value float64
			fmt.Sscanf(cols[c],"%f",&value)
			row = append(row,value)
		}
		matrix = append(matrix,row)
	}

	return matrix
}

//**************************************************************

func Symmetrize(matrix [][]float64) [][]float64 {

	var m [][]float64 = matrix

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

func Transpose(matrix [][]float64) [][]float64 {

	var m [][]float64 = matrix

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

func MakeInitVector(dim int,init_value float64) []float64 {

	var v = make([]float64,dim)

	for r := 0; r < dim; r++ {
		v[r] = init_value
	}

	return v
}

//**************************************************************

func MatrixOpVector(m [][]float64, v []float64) []float64 {

	var vp = make([]float64,len(m))

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

func ComputeEVC(s string, adj [][]float64) []float64 {

	v := MakeInitVector(len(adj),1.0)
	vlast := v
	fmt.Println("...........",s,"..........")
	PrintMatrix(adj)

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

func GetVecMax(v []float64) float64 {

	var max float64 = -1

	for r := range v {
		if v[r] > max {
			max = v[r]
		}
	}

	return max
}

//**************************************************************

func NormalizeVec(v []float64, div float64) []float64 {

	if div == 0 {
		div = 1
	}

	for r := range v {
		v[r] = v[r] / div
	}

	return v
}

//**************************************************************

func CompareVec(v1,v2 []float64) float64 {

	var max float64 = -1

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

func PrintMatrix(matrix [][]float64) {

	fmt.Printf("                 DIAG  \n")

	for row := 0; row < len(matrix); row++ {
		for col := 0; col < len(matrix[row]); col++ {
			fmt.Printf("%3.0f ",matrix[row][col])
		}

		fmt.Printf("      %1.1f    ...\n",matrix[row][row])
	}
	fmt.Println()
}
//**************************************************************

func PrintVector(vector []float64) {

	for row := 0; row < len(vector); row++ {

		fmt.Printf("( %3.2f )\n",vector[row])
	}
	fmt.Println()
}

