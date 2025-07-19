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
	"strconv"
	"strings"

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
	"(1,1), (1,3), (4,4) (3,3) other stuff",
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

	fmt.Println("STARTING EXPRESSION: ",line)
	fmt.Println(" - name:",search.Name)
	fmt.Println(" - from:",search.From)
	fmt.Println(" - to:",search.To)
	fmt.Println(" - chap:",search.Chapter)
	fmt.Println(" - context:",search.Context)
	fmt.Println(" - arrows:",search.Arrows)
	fmt.Println(" - pagenr:",search.PageNr)
	fmt.Println(" - seq:",search.Sequence)

	name := search.Name != nil
	from := search.From != nil
	to := search.To != nil
	chapter := search.Chapter != ""
	context := search.Context != nil
	pagenr := search.PageNr > 0
	sequence := search.Sequence

	arrowptrs,sttype := ArrowPtrFromArrowsNames(ctx,search.Arrows)
	nodeptrs,rest := ParseLiteralNodePtrs(search.Name)

	arrows := arrowptrs != nil
	sttypes := sttype != nil

	fmt.Println("NODEPTRS,REST=",nodeptrs,rest,sttypes)


//	name = rest

	if (name && from) || (name && to) {
		fmt.Println("Search has conflicting parts <to|from> and match strigns")
		os.Exit(-1)
	}

	// if we have name, (maybe with context, chapter, arrows)

	if name && ! sequence {

		// foreach name
		//SST.GetDBNodePtrMatchingNCC(ctx,search.Name,search.Chapter,search.Context,arrowptrs)
	}

	if name && sequence {
	}

	if sttypes {
	}

	// if we only have context then search NodeArrowNode

	if !name && context {

	}

	// if we only have chapter then we are looking for page notes
	// if we have page number then we are looking for notes by pagemap

	if chapter && pagenr && !arrows && !context {
	}

	// if we have from/to (maybe with chapter/context) then we are looking for paths
	// If we have arrows and a name|to|from && maybe chapter|context, then looking for things pointing

	if sttypes {  // from/to
	}

	if from && to {
	}

	// if we have sequence with arrows, then we are looking for sequence context or stories
	// GetNodesStartingStoriesForArrow(ctx PoSST,arrow string) ([]NodePtr,int)

	if sequence && arrows {
	}
}

//******************************************************************

func ParseLiteralNodePtrs(names []string) ([]SST.NodePtr,[]string) {

	var current []rune
	var rest []string
	var nodeptrs []SST.NodePtr

	for n := range names {

		line := []rune(names[n])
		
		for i := 0; i < len(line); i++ {
			
			if line[i] == '(' {
				rs := strings.TrimSpace(string(current))
				if len(rs) > 0 {
					rest = append(rest,string(current))
					current = nil
				}
				continue
			}
			
			if line[i] == ')' {
				np := string(current)
				var nptr SST.NodePtr
				fmt.Sscanf(np,"%d,%d",&nptr.Class,&nptr.CPtr)
				nodeptrs = append(nodeptrs,nptr)
				current = nil
				continue
			}

			current = append(current,line[i])
			
		}
		rs := strings.TrimSpace(string(current))
		if len(rs) > 0 {
			rest = append(rest,rs)
		}
		current = nil
	}

	return nodeptrs,rest
}

//******************************************************************

func ArrowPtrFromArrowsNames(ctx SST.PoSST,arrows []string) ([]SST.ArrowPtr,[]int) {

	var arr []SST.ArrowPtr
	var stt []int

	for a := range arrows {

		// is the entry a number? sttype?

		number, err := strconv.Atoi(arrows[a])
		notnumber := err != nil

		if notnumber {
			arrowptr,_ := SST.GetDBArrowsWithArrowName(ctx,arrows[a])
			if arrowptr != -1 {
				arrdir := SST.GetDBArrowByPtr(ctx,arrowptr)
				fmt.Println("lookup 1arrow",arrdir)
			}
		} else {
			if number < -SST.EXPRESS {
				fmt.Println("Negative arrow value doesn't make sense",number)
			} else if number >= -SST.EXPRESS && number <= SST.EXPRESS {
				stt = append(stt,number)
			} else {
				// whatever remains can only be an arrowpointer
				arrdir := SST.GetDBArrowByPtr(ctx,SST.ArrowPtr(number))
				fmt.Println("lookup 2arrow",arrdir)
			}
		}
	}

	return arr,stt
}

//******************************************************************

func DecodeBoundarySet(s string) []SST.NodePtr {

	var nptrs []SST.NodePtr

	return nptrs

}










