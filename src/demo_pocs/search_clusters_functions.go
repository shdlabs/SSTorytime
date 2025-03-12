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

	UseGetMatroidArrayByArrow(ctx)
	UseGetMatroidArrayBySSType(ctx)
	UseGetMatroidHistogramByArrow(ctx)
	UseGetMatroidHistogramBySSType(ctx)
	UseGetMatroidNodesByArrow(ctx)
	UseGetMatroidNodesBySTType(ctx)

	SST.Close(ctx)
}

//******************************************************************

func UseGetMatroidArrayByArrow(ctx SST.PoSST) {

	var ama map[SST.ArrowPtr][]SST.NodePtr

	ama = SST.GetMatroidArrayByArrow(ctx)

	fmt.Println("FEATURE: GetMatroidArrayByArrow: raw",ama)

	for arrowptr := range ama {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)
		fmt.Println("\nArrow --(",arr_dir.Long,")--> acts as a type/meaning correlator to classify the following group:\n")
		for n := 0; n < len(ama[arrowptr]); n++ {
			node := SST.GetDBNodeByNodePtr(ctx,ama[arrowptr][n])
			fmt.Println("  Node",ama[arrowptr][n]," = ",node.S)
		}
	}
}

//******************************************************************

func UseGetMatroidArrayBySSType(ctx SST.PoSST) {

	var ams map[int][]SST.NodePtr

	ams = SST.GetMatroidArrayBySSType(ctx)

	fmt.Println("FEATURE: GetMatroidArrayBySTType: raw",ams)

	for sttype := range ams {

		fmt.Println("\nArrow class --(",SST.STTypeName(sttype),")--> acts as a type/interpretation correlator of the following group:\n")

		for n := 0; n < len(ams[sttype]); n++ {
			node := SST.GetDBNodeByNodePtr(ctx,ams[sttype][n])
			fmt.Println("  Node",ams[sttype][n]," = ",node.S)
		}
	}
}

//******************************************************************

func UseGetMatroidHistogramByArrow(ctx SST.PoSST) {

	var ha map[SST.ArrowPtr]int
	ha = SST.GetMatroidHistogramByArrow(ctx)
	fmt.Println("FEATURE: GetMatroidHistogramByArrow",ha,"\n")

	for arrowptr := range ha {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)
		fmt.Println("Arrow -(",arr_dir.Long,")-> occurs with frequency",ha[arrowptr])
	}

}

//******************************************************************

func UseGetMatroidHistogramBySSType(ctx SST.PoSST) {

	var hs map[int]int
	hs = SST.GetMatroidHistogramBySSType(ctx)
	fmt.Println("FEATURE: GetMatroidHistogramBySTType",hs,"\n")

	for sttype := range hs {
		fmt.Println("Arrow class -(",SST.STTypeName(sttype),")-> occurs with frequency",hs[sttype])
	}
}

//******************************************************************

func UseGetMatroidNodesByArrow(ctx SST.PoSST) {

	var ma []SST.ArrowMatroid
	ma = SST.GetMatroidNodesByArrow(ctx)

	fmt.Println("FEATURE: GetMatroidNodesByArrow",ma,"\n")

	for a := range ma {
		from := SST.GetDBNodeByNodePtr(ctx,ma[a].NFrom)
		arrow := SST.GetDBArrowByPtr(ctx,ma[a].Arr)
		fmt.Println("  - Node",from.S,"acts as a matroidal hub connecting:")
		for n := range ma[a].NTo {
			to := SST.GetDBNodeByNodePtr(ctx,ma[a].NTo[n])
			fmt.Println("   ... Node",to.S)
		}
		fmt.Println("    with arrow:  --(",arrow,")-->\n")
	}
}

//******************************************************************

func UseGetMatroidNodesBySTType(ctx SST.PoSST) {

	var ms []SST.STTypeMatroid
	ms = SST.GetMatroidNodesBySTType(ctx)
	fmt.Println("FEATURE: GetMatroidNodesBySTType",ms,"\n")

	for s := range ms {
		from := SST.GetDBNodeByNodePtr(ctx,ms[s].NFrom)
		sttype := SST.STTypeName(ms[s].STType)
		fmt.Println("  -Node",from.S,"acts as a matroidal hub connecting:")
		for n := range ms[s].NTo {
			to := SST.GetDBNodeByNodePtr(ctx,ms[s].NTo[n])
			fmt.Println("   ... Node",to.S)
		}
		fmt.Println("   with arrow class:  --(",sttype,")-->\n")
	}
}







