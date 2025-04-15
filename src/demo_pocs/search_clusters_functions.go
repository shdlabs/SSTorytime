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
	dbname   = "sstoryline"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	UseGetAppointmentArrayByArrow(ctx)
	UseGetAppointmentArrayBySSType(ctx)
	UseGetAppointmentHistogramByArrow(ctx)
	UseGetAppointmentHistogramBySSType(ctx)
	UseGetAppointmentNodesByArrow(ctx)
	UseGetAppointmentNodesBySTType(ctx)

	SST.Close(ctx)
}

//******************************************************************

func UseGetAppointmentArrayByArrow(ctx SST.PoSST) {

	context := []string{"any"}
	chapter := "any"

	var ama map[SST.ArrowPtr][]SST.NodePtr

	ama = SST.GetAppointmentArrayByArrow(ctx,context,chapter)

	fmt.Println("--------------------------------------------------")
	fmt.Println("FEATURE: GetAppointmentArrayByArrow:")
	fmt.Println("--------------------------------------------------")

	for arrowptr := range ama {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)
		fmt.Println("\nArrow --(",arr_dir.Long,")--> points to a group of nodes with a similar role in the context of",context,"in the chapter",chapter,"\n")

		for n := 0; n < len(ama[arrowptr]); n++ {
			node := SST.GetDBNodeByNodePtr(ctx,ama[arrowptr][n])
			NewLine(n)
			fmt.Print("..  ",node.S,",")

		}
		fmt.Println()
		fmt.Println("............................................")
	}
}

//******************************************************************

func UseGetAppointmentArrayBySSType(ctx SST.PoSST) {

	var ams map[int][]SST.NodePtr

	ams = SST.GetAppointmentArrayBySSType(ctx)

	fmt.Println("--------------------------------------------------")
	fmt.Println("FEATURE: GetAppointmentArrayBySTType:")
	fmt.Println("--------------------------------------------------")

	for sttype := range ams {

		fmt.Println("\nArrow class --(",SST.STTypeName(sttype),")--> acts as a type/interpretation correlator of the following group by pointing/pointed to:\n")

		for n := 0; n < len(ams[sttype]); n++ {
			node := SST.GetDBNodeByNodePtr(ctx,ams[sttype][n])
			NewLine(n)
			fmt.Print("..  ",node.S,",")
		}
		fmt.Println()
		fmt.Println("............................................")
	}
}

//******************************************************************

func UseGetAppointmentHistogramByArrow(ctx SST.PoSST) {

	var ha map[SST.ArrowPtr]int
	ha = SST.GetAppointmentHistogramByArrow(ctx)
	fmt.Println("*****************************************************")
	fmt.Println("FEATURE: GetAppointmentHistogramByArrow")
	fmt.Println("*****************************************************")

	fmt.Println("The relative prevalence of things pointed out by the graph relations:\n")

	for arrowptr := range ha {
		arr_dir := SST.GetDBArrowByPtr(ctx,arrowptr)
		fmt.Println("    - Arrow -(",arr_dir.Long,")-> selects nodes with frequency",ha[arrowptr])
	}

}

//******************************************************************

func UseGetAppointmentHistogramBySSType(ctx SST.PoSST) {

	var hs map[int]int
	hs = SST.GetAppointmentHistogramBySSType(ctx)
	fmt.Println("*****************************************************")
	fmt.Println("FEATURE: GetAppointmentHistogramBySTType")
	fmt.Println("*****************************************************")

	fmt.Println("The relative prevalence of things pointed out by spacetime process:\n")

	for sttype := range hs {
		fmt.Println("   - Arrow class -(",SST.STTypeName(sttype),")-> selects nodes with frequency",hs[sttype])
		fmt.Println("............................................")
	}
}

//******************************************************************

func UseGetAppointmentNodesByArrow(ctx SST.PoSST) {

	var ma []SST.ArrowAppointment
	ma = SST.GetAppointmentNodesByArrow(ctx)

	fmt.Println("*****************************************************")
	fmt.Println("FEATURE: GetAppointmentNodesByArrow")
	fmt.Println("*****************************************************")

	for a := range ma {
		from := SST.GetDBNodeByNodePtr(ctx,ma[a].NFrom)
		arrow := SST.GetDBArrowByPtr(ctx,ma[a].Arr)
		fmt.Println("  - Node",from.S,"acts as a matroidal hub connecting pointing/pointed to all these:")
		for n := range ma[a].NTo {
			to := SST.GetDBNodeByNodePtr(ctx,ma[a].NTo[n])
			NewLine(n)
			fmt.Print("      ...",to.S,",")
		}
		fmt.Println("    with arrow type:  --(",arrow.Long,")-->\n")
		fmt.Println("............................................")
	}
}

//******************************************************************

func UseGetAppointmentNodesBySTType(ctx SST.PoSST) {

	var ms []SST.STTypeAppointment
	ms = SST.GetAppointmentNodesBySTType(ctx)
	fmt.Println("*****************************************************")
	fmt.Println("FEATURE: GetAppointmentNodesBySTType")
	fmt.Println("*****************************************************")

	for s := range ms {
		from := SST.GetDBNodeByNodePtr(ctx,ms[s].NFrom)
		sttype := SST.STTypeName(ms[s].STType)
		fmt.Println("  -Node",from.S,"acts as a matroidal hub pointing/pointed to all these:")
		for n := range ms[s].NTo {
			to := SST.GetDBNodeByNodePtr(ctx,ms[s].NTo[n])
			NewLine(n)
			fmt.Print("     ...",to.S,",")
		}
		fmt.Println("   with arrow class:  --(",sttype,")-->\n")
		fmt.Println("............................................")
	}
}

//******************************************************************
// Tools
//******************************************************************

func NewLine(n int) {

	if n % 8 == 0 {
		fmt.Println()
	}
}







