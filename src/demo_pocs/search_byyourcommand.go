//******************************************************************
//
// Replacement for searchN4L
// single search string without complex options
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

var VERBOSE bool = false

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

	SST.MemoryInit()

	load_arrows := false
	ctx := SST.Open(load_arrows)

	if len(args) > 0 {

		search_string := ""
		for a := 0; a < len(args); a++ {
			search_string += args[a] + " "
		}

		search := SST.DecodeSearchField(search_string)

		Search(ctx,search,search_string)
	}

	SST.Close(ctx)
}

//**************************************************************

func Usage() {
	
	fmt.Printf("usage: ByYourCommand <search request>\n\n")
	fmt.Println("searchN4L <mytopic> chapter <mychapter>\n\n")
	fmt.Println("searchN4L range rover out of its depth")
	fmt.Println("searchN4L \"range rover\" \"out of its depth\"")
	fmt.Println("searchN4L from rover range 4")
	fmt.Println("searchN4L head used as \"version control\"")
	fmt.Println("searchN4L head context neuro)brain)etc")
	fmt.Println("searchN4L notes on restaurants in chinese")	
	fmt.Println("searchN4L notes about brains")
	fmt.Println("searchN4L notes music writing")
	fmt.Println("searchN4L page 2 of notes on brains") 
	fmt.Println("searchN4L notes page 3 brain") 
	fmt.Println("searchN4L (1)1)) (1)3)) (4)4) (3)3) other stuff")
	fmt.Println("searchN4L arrows pe)ep) eh")
	fmt.Println("searchN4L arrows 1)-1")
	fmt.Println("searchN4L forward cone for (bjorvika) range 5")
	fmt.Println("searchN4L sequences about fox")	
	fmt.Println("searchN4L context \"not only\"") 
	fmt.Println("searchN4L \"come on down\"")	
	fmt.Println("searchN4L chinese kinds of meat") 
	fmt.Println("searchN4L summary chapter interference")
	fmt.Println("searchN4L paths from arrows pe)ep) eh")
	fmt.Println("searchN4L paths from start to target2 limit 5")
	fmt.Println("searchN4L paths to target3")	
	fmt.Println("searchN4L a2 to b5 distance 10")
	fmt.Println("searchN4L to a5")
	fmt.Println("searchN4L from start")
	fmt.Println("searchN4L from (1)6)")
	fmt.Println("searchN4L a1 to b6 arrows then")
	fmt.Println("searchN4L paths a2 to b5 distance 10")

	flag.PrintDefaults()

	os.Exit(2)
}

//**************************************************************

func GetArgs() []string {

	flag.Usage = Usage
	verbosePtr := flag.Bool("v", false,"verbose")
	flag.Parse()

	if *verbosePtr {
		VERBOSE = true
	}

	return flag.Args()
}

//******************************************************************

func Search(ctx SST.PoSST, search SST.SearchParameters,line string) {

	if VERBOSE {
		fmt.Println("Your starting expression generated this set: ",line,"\n")
		fmt.Println(" - start set:",SL(search.Name))
		fmt.Println(" -      from:",SL(search.From))
		fmt.Println(" -        to:",SL(search.To))
		fmt.Println(" -   chapter:",search.Chapter)
		fmt.Println(" -   context:",SL(search.Context))
		fmt.Println(" -    arrows:",SL(search.Arrows))
		fmt.Println(" -    pagenr:",search.PageNr)
		fmt.Println(" - sequence/story:",search.Sequence)
		fmt.Println(" - limit/range/depth:",search.Range)
		fmt.Println()
	}

	// OPTIONS *********************************************

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
	limit := 0

	if search.Range > 0 {
		limit = search.Range
	} else {
		limit = 5
	}

	// SEARCH SELECTION *********************************************

	fmt.Println("------------------------------------------------------------------")

	// if we have name, (maybe with context, chapter, arrows)

	if name && ! sequence && !pagenr {
		FindOrbits(ctx, nodeptrs, limit)
		return
	}

	if (name && from) || (name && to) {
		fmt.Printf("\nSearch \"%s\" has conflicting parts <to|from> and match strings\n",line)
		os.Exit(-1)
	}

	// Path solving, two sets of nodeptrs
	// if we have from/to (maybe with chapter/context) then we are looking for paths
	// If we have arrows and a name|to|from && maybe chapter|context, then looking for things pointing

	if from && to {

		if sttypes {  // from/to
			fmt.Println("USE GetFwdPathsAsLinks(sttype)")
			fmt.Println("PATH BOUNDARY SETS",leftptrs,rightptrs)
		}

		if arrows {  // from/to
			fmt.Println("SST.GetEntireNCSuperConePathsAsLinks(ctx,FWD,leftptrs,ldepth,chapter,context) AND FILTER")
		}
		fmt.Println("PATH BOUNDARY SETS without arrow constraints",leftptrs,rightptrs)
	}

	// Causal cones, from one of these three

	if name || from || to {

		if VERBOSE {
			fmt.Println("SEARCH: name",nodeptrs,"from",leftptrs,"to",rightptrs,"SSTypes",sttype,"arrow names",arrowptrs)
		}

		if sttypes || arrows {
			// from or to or name
			if VERBOSE {
				fmt.Println("CausalCones(ctx,nodeptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)")
			}

			if nodeptrs != nil {
				CausalCones(ctx,nodeptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
				return
			}
			if leftptrs != nil {
				CausalCones(ctx,leftptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
				return
			}
			if rightptrs != nil {
				CausalCones(ctx,rightptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
				return
			}
		}
		

	}

	// if we have sequence with arrows, then we are looking for sequence context or stories

	if name && pagenr {

	}

	// if we only have context then search NodeArrowNode

	if !name && context {
		// GetMatchingContexts(context)
		fmt.Println("GetDBPageMap(ctx PoSST,chap string,cn []string,page int) []PageMap ")
		fmt.Println("NOTES from context")
	}

	// if we have page number then we are looking for notes by pagemap

	if (name || chapter) && pagenr {
		fmt.Println("GetDBPageMap(ctx PoSST,chap string,cn []string,page int) []PageMap ")
		fmt.Println("NOTES BROWSING")
	}


	// if we have sequence with arrows, then we are looking for sequence context or stories
	// GetNodesStartingStoriesForArrow(ctx PoSST,arrow string) ([]NodePtr,int)

	if arrows {
		fmt.Println("Single links listed by arrow type")
		fmt.Println("RETURN NodeArrowNode for arrows or STType, filter context,name,chapter")
	}

	if name && sequence {
		fmt.Println("STORIES by starting node")
		fmt.Println("Start node pointers and select by arrow/sttype")
	}

	if sequence && arrows {
		fmt.Println("STORIES by arrow type")
		fmt.Println(" GetSequenceContainers(ctx PoSST,arrname string,search,chapter string,context []string) []Story")
	}
}

//******************************************************************

func SolveNodePtrs(ctx SST.PoSST,nodenames []string,chap string,cntx []string, arr []SST.ArrowPtr) []SST.NodePtr {

	nodeptrs,rest := ParseLiteralNodePtrs(nodenames)

	var idempotence = make(map[SST.NodePtr]bool)
	var result []SST.NodePtr

	for n := range nodeptrs {
		idempotence[nodeptrs[n]] = true
	}

	for r := range rest {
		nptrs := SST.GetDBNodePtrMatchingNCC(ctx,rest[r],chap,cntx,arr)
		for n := range nptrs {
			idempotence[nptrs[n]] = true
		}
	}

	for uniqnptr := range idempotence {
		result = append(result,uniqnptr)
	}

	return result
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

//******************************************************************
// SEARCH
//******************************************************************

func FindOrbits(ctx SST.PoSST, nptrs []SST.NodePtr, limit int) {
	
	var count int

	if VERBOSE {
		fmt.Println("First",limit,"orbit result(s):\n")
	}
	for nptr := range nptrs {
		count++
		if count > limit {
			return
		}
		fmt.Print("\n",nptr,": ")
		SST.PrintNodeOrbit(ctx,nptrs[nptr],100)
	}
}

//******************************************************************

func CausalCones(ctx SST.PoSST,nptrs []SST.NodePtr, chap string, context []string,arrows []SST.ArrowPtr, sttype []int,limit int) {
	var total int = 1

	for n := range nptrs {
		for st := range sttype {

			fcone,_ := SST.GetFwdPathsAsLinks(ctx,nptrs[n],sttype[st],limit)

			if fcone != nil {
				fmt.Printf("%d. ",total)
				total += ShowCone(ctx,fcone,sttype[st],chap,context,limit)
			}

			if total > limit {
				return
			}

			bcone,_ := SST.GetFwdPathsAsLinks(ctx,nptrs[n],-sttype[st],limit)

			if bcone != nil {
				fmt.Printf("%d. ",total)
				total += ShowCone(ctx,bcone,sttype[st],chap,context,limit)
			}

			if total > limit {
				return
			}
		}
	}

}

//******************************************************************
// OUTPUT
//******************************************************************

func ShowCone(ctx SST.PoSST,cone [][]SST.Link,sttype int,chap string,context []string,limit int) int {

	if len(cone) < 1 {
		return 0
	}

	if limit <= 0 {
		return 0
	}

	count := 0

	for s := 0; s < len(cone) && s < limit; s++ {
		SST.PrintSomeLinkPath(ctx,cone,s," - ",chap,context,limit)
		count++
	}

	return count
}






