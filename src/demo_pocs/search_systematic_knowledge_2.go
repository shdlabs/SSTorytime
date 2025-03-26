//******************************************************************
//
// Exploring how to present knowledge systematically, e.g.
// e.g. review/review for an exam!
//  version 2 with a more piecemeal strategy
//
//******************************************************************

package main

import (
	"fmt"
	"strings"

        SST "SSTorytime"
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

	const maxdepth = 8

	fmt.Println("....................................................................................")

	cone,_ := SST.GetFwdPathsAsLinks(ctx,start,1,maxdepth)
	ShowCone(ctx,cone,1,chap,context)

	cone,_ = SST.GetFwdPathsAsLinks(ctx,start,-1,maxdepth)
	ShowCone(ctx,cone,1,chap,context)

	cone,_ = SST.GetFwdPathsAsLinks(ctx,start,2,maxdepth)
	ShowCone(ctx,cone,1,chap,context)

	cone,_ = SST.GetFwdPathsAsLinks(ctx,start,-2,maxdepth)
	ShowCone(ctx,cone,1,chap,context)

	cone,_ = SST.GetFwdPathsAsLinks(ctx,start,3,maxdepth)
	ShowCone(ctx,cone,1,chap,context)

	cone,_ = SST.GetFwdPathsAsLinks(ctx,start,-3,maxdepth)
	ShowCone(ctx,cone,1,chap,context)

	cone,_ = SST.GetFwdPathsAsLinks(ctx,start,0,maxdepth)
	ShowCone(ctx,cone,1,chap,context)
}

//**************************************************************

func ShowCone(ctx SST.PoSST,cone [][]SST.Link,sttype int,chap string,context []string) {

	if len(cone) < 1 {
		return
	}

	for s := 0; s < len(cone); s++ {

		prefix := fmt.Sprintf("\n - %s Word/Phrase ",SST.STTypeName(sttype))
		SST.PrintLinkPath(ctx,cone,s,prefix,chap,context)
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












