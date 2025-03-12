//******************************************************************
//
// Exploring how to present a search text, with API
//
// Prepare:
// cd examples
// ../src/N4L-db -u Mary.n4l, e.g. try type Mary example, type 1
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
	dbname   = "newdb"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	x := SST.GetMatroidArrayByArrow(ctx)
	fmt.Println("GetMatroidArrayByArrow",x)

	y := SST.GetMatroidArrayBySSType(ctx)
	fmt.Println("GetMatroidArrayBySTType",y)

	z := SST.GetMatroidHistogramByArrow(ctx)
	fmt.Println("GetMatroidHistogramByArrow",z)

	a := SST.GetMatroidHistogramBySSType(ctx)
	fmt.Println("GetMatroidHistogramBySTType",a)

	b := SST.GetMatroidNodesByArrow(ctx)
	fmt.Println("GetMatroidNodesByArrow",b)

	c := SST.GetMatroidNodesBySTType(ctx)
	fmt.Println("GetMatroidNodesBySTType",c)

	SST.Close(ctx)
}








