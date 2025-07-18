//******************************************************************
//
// single search string without complex options
//
// how shall we split a search query into parts to match against?
//
//******************************************************************

package main

import (
	"fmt"
	"os"
	"flag"

        SST "SSTorytime"
)

//******************************************************************

var TESTS = []string{ 
	"head used as chinese stuff",
	"head context neuro",
	"leg in chapter bodyparts",
	"foot in bodyparts2",
	"visual for prince",	
	"visual of integral",	
	"notes on restaurants in chinese",	
	"notes about",	
	"(1,1), (1,3), (4,4) (3,3)",
	"page 2 of notes on brains", 
	"notes page 3 brain", 
	"integrate in math",	
	"arrows pe,ep, eh",
	"arrows 1,-1",
	"forward cone for (bjorvika)",
	"backward sideways cone for (bjorvika)",
	"sequences about fox",	
	"stories about (bjorvika)",	
	"context \"not only\"", 
	"\"come in\"",	
	"containing / matching \"blub blub\"", 
	"chinese kinds of meat", 
	"images prince", 
	"summary chapter interference",
	"showme greetings in norwegian",
	"paths from arrows pe,ep, eh",
	"paths from start to target",	
	"paths to target3",	
	"a2 to b5",
	"to a5",
	"from start",
	"from (1,6)",
	"a1 to b6 arrows then",
        }


//******************************************************************

func main() {

	args := GetArgs()

	fmt.Println(args)

	SST.MemoryInit()

	load_arrows := false
	ctx := SST.Open(load_arrows)

	for test := range TESTS {
		fmt.Println("......................")
		search := SST.DecodeSearchField(TESTS[test])
		Search(ctx,search,TESTS[test])
	}

	SST.Close(ctx)
}

//**************************************************************

func Usage() {
	
	fmt.Printf("usage: ByYourCommand <search request>n")
	flag.PrintDefaults()

	os.Exit(2)
}

//**************************************************************

func GetArgs() []string {

	flag.Usage = Usage
	flag.Parse()
	return flag.Args()
}

//******************************************************************

func Search(ctx SST.PoSST, search SST.SearchParameters,line string) {

	fmt.Println("FROM",line)
	fmt.Println(" - name:",search.Name)
	fmt.Println(" - from:",search.From)
	fmt.Println(" - to:",search.To)
	fmt.Println(" - chap:",search.Chapter)
	fmt.Println(" - context:",search.Context)
	fmt.Println(" - arrows:",search.Arrows)
	fmt.Println(" - pagenr:",search.PageNr)
	fmt.Println(" - seq:",search.Sequence)
}









