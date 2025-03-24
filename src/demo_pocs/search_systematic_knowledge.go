//******************************************************************
//
// Exploring how to present knowledge systematically, e.g.
// e.g. review/review for an exam!
//
//******************************************************************

package main

import (
	"fmt"
	"strings"

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

	context := []string{""}
	arrows := []string{"pe","ph"} // Always start with pinyin

	Systematic(ctx,"chinese",context,"",arrows)

	SST.Close(ctx)
}

//******************************************************************

func Systematic(ctx SST.PoSST, chaptext string,context []string,searchtext string,arrnames []string) {

	chaptext = strings.TrimSpace(chaptext)
	searchtext = strings.TrimSpace(searchtext)

	var arrows []SST.ArrowPtr

	for a := range arrnames {
		arr := SST.GetDBArrowByName(ctx,arrnames[a])
		arrows = append(arrows,arr)
	}

	nodes := SST.GetDBNodeContextsMatchingArrow(ctx,chaptext,context,searchtext,arrows)

	var prev string
	var header []string

	for cntxt := range nodes {
				
		for n := 0; n < len(nodes[cntxt]); n++ {

			result := SST.GetDBNodeByNodePtr(ctx,nodes[cntxt][n])

			if cntxt != prev {
				prev = cntxt
				header = SST.ParseSQLArrayString(cntxt)
				Header(header,result.Chap)
			}

			SearchStoryPaths(ctx,result.S,result.NPtr,arrows,result.Chap,context)
		}
	}
}

//**************************************************************

func SearchStoryPaths(ctx SST.PoSST,name string,start SST.NodePtr, arrows []SST.ArrowPtr,chap string,context []string) {

	cone,_ := SST.GetEntireConePathsAsLinks(ctx,"any",start,8)

	var done = make(map[SST.NodePtr]bool)

	if len(cone) < 1 {
		return
	}

	fmt.Println("....................................................................................")

	for s := 0; s < len(cone); s++ {

		if done[cone[s][0].Dst] {
			continue
		} else {
			done[cone[s][0].Dst] = true
		}

		prefix := fmt.Sprintf("\n - Word/Phrase ")
		SST.PrintLinkPath(ctx,cone,s,prefix,chap,context)
	}

}

//**************************************************************

func SearchStoryMatroids(ctx SST.PoSST,name string,start SST.NodePtr, arrows []SST.ArrowPtr) {


	var ams map[int][]SST.NodePtr

	ams = SST.GetMatroidArrayBySSType(ctx)

	for sttype := range ams {

		fmt.Println("\nArrow class --(",SST.STTypeName(sttype),")--> acts as a type/interpretation correlator of the following group by pointing/pointed to:\n")

		for n := 0; n < len(ams[sttype]); n++ {
			node := SST.GetDBNodeByNodePtr(ctx,ams[sttype][n])
			//NewLine(n)
			fmt.Print("..  ",node.S,",")
		}
		fmt.Println()
		fmt.Println("............................................")
	}

}

//**************************************************************

func Header(h []string,chap string) {

	if len(h) == 0 {
		return
	}

	fmt.Println("\n\n============================================================")
	fmt.Println("   In chapter: \"",chap,"\"\n")

	for s := range h {
		fmt.Println("   ::",h[s],"::")
	}

	fmt.Println("\n============================================================")
}

//**************************************************************

func Box(a ...interface{}) {

	fmt.Println("\n------------------------------------")
	fmt.Println(a...)
	fmt.Println("------------------------------------\n")
}










