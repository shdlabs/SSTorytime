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
	"range rover out of its depth",
	"\"range rover\" \"out of its depth\"",
	"from rover range 4",
	"head used as chinese stuff",
	"head context neuro,brain,etc",
	"leg in chapter bodyparts",
	"foot in bodyparts2",
	"visual for prince",	
	"visual of integral",	
	"notes on restaurants in chinese",	
	"notes about brains",
	"notes music writing",
	"page 2 of notes on brains", 
	"notes page 3 brain", 
	"(1,1), (1,3), (4,4) (3,3) other stuff",
	"integrate in math",	
	"arrows pe,ep, eh",
	"arrows 1,-1",
	"forward cone for (bjorvika) range 5",
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
	"paths from start to target limit 5",
	"paths to target3",	
	"a2 to b5 distance 10",
	"to a5",
	"from start",
	"from (1,6)",
	"a1 to b6 arrows then",
	"paths a2 to b5 distance 10",
	"from dog to cat",
        }


//******************************************************************

func main() {

	args := GetArgs()

	fmt.Println(args)

	SST.MemoryInit()

	load_arrows := false
	ctx := SST.Open(load_arrows)

	if len(args) > 0 {
		fmt.Println("args",args,len(args))
		search_string := ""
		for a := 0; a < len(args); a++ {
			search_string += args[a] + " "
		}
		fmt.Println(search_string,"......................")
		search := SST.DecodeSearchField(search_string)
		Search(ctx,search,search_string)
	} else {
		for test := range TESTS {
			fmt.Println("......................")
			search := SST.DecodeSearchField(TESTS[test])
			Search(ctx,search,TESTS[test])
		}
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
	fmt.Println(" - name:",SL(search.Name))
	fmt.Println(" - from:",SL(search.From))
	fmt.Println(" - to:",SL(search.To))
	fmt.Println(" - chap:",search.Chapter)
	fmt.Println(" - context:",SL(search.Context))
	fmt.Println(" - arrows:",SL(search.Arrows))
	fmt.Println(" - pagenr:",search.PageNr)
	fmt.Println(" - range/depth:",search.Range)
	fmt.Println(" - seq:",search.Sequence)

	name := search.Name != nil
	from := search.From != nil
	to := search.To != nil
	chapter := search.Chapter != ""
	context := search.Context != nil
	pagenr := search.PageNr > 0
	sequence := search.Sequence

	arrowptrs,sttype := ArrowPtrFromArrowsNames(ctx,search.Arrows)
	nodeptrs := SolveNodePtrs(ctx,search.Name,search.Chapter,search.Context,arrowptrs)
	leftptrs := SolveNodePtrs(ctx,search.From,search.Chapter,search.Context,arrowptrs)
	rightptrs := SolveNodePtrs(ctx,search.To,search.Chapter,search.Context,arrowptrs)

	arrows := arrowptrs != nil
	sttypes := sttype != nil

	// if we are not looking for paths or cones, depth. range, distance could be actual search terms

	if !from && !to {
	}

	// if we have name, (maybe with context, chapter, arrows)

	if name && ! sequence {

		fmt.Println("AD HOC SELECTION SETS",nodeptrs)
	}

	// RETURN THIS TYPE NOW: []NodePtr for Orbits and Cones, start/end sets
	// or continue to append nodeptrs

	// Next PATHS, which are merged cones
	// Sequences are forward cones

	if (name && from) || (name && to) {
		fmt.Printf("\nSearch \"%s\" has conflicting parts <to|from> and match strings\n",line)
		os.Exit(-1)
	}

	// Path solving, two sets of nodeptrs
	// if we have from/to (maybe with chapter/context) then we are looking for paths
	// If we have arrows and a name|to|from && maybe chapter|context, then looking for things pointing

	if from && to {

		if sttypes {  // from/to
		}

		if arrows {  // from/to
		}

		fmt.Println("PATH BOUNDARY SETS",leftptrs,rightptrs)
	}

	// if we have sequence with arrows, then we are looking for sequence context or stories

	if name && sequence {


	}

	if sttypes {
		//GetEntireCone/Fwd/Bwd
		fmt.Println("FWD CONE")
	}

	// if we only have context then search NodeArrowNode

	if !name && context {
		// GetMatchingContexts(context)
		//notes := SST.GetDBPageMap(CTX,chaptext,context,pagenr)
		fmt.Println("NOTES from context")
	}

	// if we only have chapter then we are looking for page notes
	// if we have page number then we are looking for notes by pagemap

	if chapter && pagenr && !arrows && !context {
	//	GetDBPageMap(ctx PoSST,chap string,cn []string,page int) []PageMap {
		fmt.Println("NOTES BROWSING")
	}


	// if we have sequence with arrows, then we are looking for sequence context or stories
	// GetNodesStartingStoriesForArrow(ctx PoSST,arrow string) ([]NodePtr,int)

	if sequence && arrows {
	}
}

//******************************************************************

func SolveNodePtrs(ctx SST.PoSST,nodenames []string,chap string,cntx []string, arr []SST.ArrowPtr) []SST.NodePtr {

	nodeptrs,rest := ParseLiteralNodePtrs(nodenames)

	for r := range rest {
		nptrs := SST.GetDBNodePtrMatchingNCC(ctx,rest[r],chap,cntx,arr)
		nodeptrs = append(nodeptrs,nptrs...)
	}

	return nodeptrs
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
				arr = append(arr,arrdir.Ptr)
			}
		} else {
			if number < -SST.EXPRESS {
				fmt.Println("Negative arrow value doesn't make sense",number)
			} else if number >= -SST.EXPRESS && number <= SST.EXPRESS {
				stt = append(stt,number)
			} else {
				// whatever remains can only be an arrowpointer
				arrdir := SST.GetDBArrowByPtr(ctx,SST.ArrowPtr(number))
				arr = append(arr,arrdir.Ptr)
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

//******************************************************************

func SL(list []string) string {

	var s string

	s += fmt.Sprint(" [")
	for i := 0; i < len(list); i++ {
		s += fmt.Sprint(list[i],", ")
	}

	s += fmt.Sprint(" ]")

	return s
}








