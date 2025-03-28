//******************************************************************
//
// SearchN4L - a simple command line search tool
//
// Prepare, e.g.:
// cd examples
// ../src/N4L-db -u chinese*n4l Mary.n4l doors.n4l cluedo.n4l brains.n4l 
//
//******************************************************************

package main

import (
	"fmt"
	"os"
	"strings"
	"flag"

        SST "SSTorytime"
)

//******************************************************************

var (
	ARROWS []string
	CHAPTER string
	SUBJECT string
	CONTEXT []string
	VERBOSE bool
	LIMIT int
)

//******************************************************************

func main() {

	Init()

	load_arrows := true
	ctx := SST.Open(load_arrows)

	Search(ctx,ARROWS,CHAPTER,CONTEXT,SUBJECT,LIMIT)

	SST.Close(ctx)
}


//**************************************************************

func Usage() {
	
	fmt.Printf("usage: searchN4L [-v] [-arrows=] [-chapter string] subject [context]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

//**************************************************************

func Init() []string {

	flag.Usage = Usage

	verbosePtr := flag.Bool("v", false,"verbose")
	chapterPtr := flag.String("chapter", "any", "a optional string to limit to a chapter/section")
	arrowsPtr := flag.String("arrows", "then", "a list of forward/outward arrows to start with")
	limitPtr := flag.Int("limit", 20, "an approximate limit on the number of items returned, where applicable")

	flag.Parse()
	args := flag.Args()

	if *verbosePtr {
		VERBOSE = true
	}

	if *arrowsPtr != "" {
		ARROWS = strings.Split(*arrowsPtr,",")
	}

	LIMIT = *limitPtr

	if *chapterPtr != "" {
		CHAPTER = *chapterPtr
	}

	if len(args) > 0 {
		SUBJECT = args[0]
		
		for c := 1; c < len(args); c++ {
			CONTEXT = append(CONTEXT,args[c])
		}

		if len(ARROWS) == 0 && len(args) < 1 {
			Usage()
			os.Exit(1);
		}
	} 

	if CONTEXT == nil {
		CONTEXT = append(CONTEXT,"")
	}

	if len(ARROWS) == 0 {
		Usage()
		os.Exit(1);
	}
	
	SST.MemoryInit()

	return args
}

//******************************************************************

func Search(ctx SST.PoSST,arrows []string,chapter string,context []string,searchtext string, limit int) {

	fmt.Println()
	fmt.Println("** PROVISIONAL SEARCH TOOL *************************************\n")
	fmt.Println("   Searching in chapter",chapter)
	fmt.Println("   With context",context)
	fmt.Println("   Selected arrows",arrows)
	fmt.Println("   Node filter",searchtext)

	Systematic(ctx,chapter,context,searchtext,arrows)
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

		SST.PrintLinkPath(ctx,cone,s," - ",chap,context)
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








