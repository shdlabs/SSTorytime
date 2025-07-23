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
			if strings.Contains(args[a]," ") {
				search_string += fmt.Sprintf("\"%s\"",args[a]) + " "
			} else {
				search_string += args[a] + " "
			}
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
	fmt.Println("searchN4L (1,1) (1,3) (4,4) (3,3) other stuff")
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
	fmt.Println("searchN4L <b5|a2> distance 10")

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

	// Check for Dirac notation

	for arg := range search.Name {

		isdirac,beg,end,cnt := SST.DiracNotation(search.Name[arg])
		
		if isdirac {
			search.Name = nil
			search.From = []string{beg}
			search.To = []string{end}
			search.Context = []string{cnt}
			break
		}
	}

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
	context := search.Context != nil
	chapter := search.Chapter != ""
	pagenr := search.PageNr > 0
	sequence := search.Sequence

	// Now convert strings into NodePointers

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

	// if we have name, (maybe with context, chapter, arrows)

	if name && ! sequence && !pagenr {

		fmt.Println("------------------------------------------------------------------")
		FindOrbits(ctx, nodeptrs, limit)
		return
	}

	if (name && from) || (name && to) {
		fmt.Printf("\nSearch \"%s\" has conflicting parts <to|from> and match strings\n",line)
		os.Exit(-1)
	}

	// Closed path solving, two sets of nodeptrs
	// if we have BOTH from/to (maybe with chapter/context) then we are looking for paths

	if from && to {

		fmt.Println("------------------------------------------------------------------")
		PathSolve(ctx,leftptrs,rightptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
		return
	}

	// Open causal cones, from one of these three

	if name || from || to {

		if sttypes || arrows {
			// from or to or name
			if VERBOSE {
				fmt.Println("CausalCones(ctx,nodeptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)")
			}

			if nodeptrs != nil {
				fmt.Println("------------------------------------------------------------------")
				CausalCones(ctx,nodeptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
				return
			}
			if leftptrs != nil {
				fmt.Println("------------------------------------------------------------------")
				CausalCones(ctx,leftptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
				return
			}
			if rightptrs != nil {
				fmt.Println("------------------------------------------------------------------")
				CausalCones(ctx,rightptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
				return
			}
		}
		

	}

	// if we have page number then we are looking for notes by pagemap

	if (name || chapter) && pagenr {

		var notes []SST.PageMap

		if chapter {
			notes = SST.GetDBPageMap(ctx,search.Chapter,search.Context,search.PageNr)
			ShowNotes(ctx,notes)
			return
		} else {
			for n := range search.Name {
				notes = SST.GetDBPageMap(ctx,search.Name[n],search.Context,search.PageNr)
				ShowNotes(ctx,notes)
			}
			return
		}
	}

	// Look for axial trails following a particular arrow, like _sequence_ 

	if name && sequence || sequence && arrows {
		ShowStories(ctx,search.Arrows,search.Name,search.Chapter,search.Context)
		return
	}

	// Match existing contexts

	if chapter {
		ShowMatchingChapter(ctx,search.Chapter)
		return
	}

	if context {
		ShowMatchingContext(ctx,search.Context)
		return
	}

	// if we have sequence with arrows, then we are looking for sequence context or stories
	// GetNodesStartingStoriesForArrow(ctx PoSST,arrow string) ([]NodePtr,int)

	if arrows || sttypes {
		ShowMatchingArrows(ctx,arrowptrs,sttype)
		return
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
				var a,b int = -1,-1
				fmt.Sscanf(np,"%d,%d",&a,&b)
				if a >= 0 && b >= 0 {
					nptr.Class = a
					nptr.CPtr = SST.ClassedNodePtr(b)
					nodeptrs = append(nodeptrs,nptr)
					current = nil
				} else {
					rest = append(rest,"("+np+")")
					current = nil
				}
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

func PathSolve(ctx SST.PoSST,leftptrs,rightptrs []SST.NodePtr,chapter string,context []string,arrowptrs []SST.ArrowPtr,sttype []int,maxdepth int) {

	var Lnum,Rnum int
	var count int
	var left_paths, right_paths [][]SST.Link

	if leftptrs == nil || rightptrs == nil {
		return
	}

	// Find the path matrix

	var solutions [][]SST.Link
	var ldepth,rdepth int = 1,1

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {

		left_paths,Lnum = SST.GetEntireNCSuperConePathsAsLinks(ctx,"fwd",leftptrs,ldepth,chapter,context)
		right_paths,Rnum = SST.GetEntireNCSuperConePathsAsLinks(ctx,"bwd",rightptrs,rdepth,chapter,context)

		// try the reverse

		if Lnum == 0 || Rnum == 0 {
			left_paths,Lnum = SST.GetEntireNCSuperConePathsAsLinks(ctx,"bwd",leftptrs,ldepth,chapter,context)
			right_paths,Rnum = SST.GetEntireNCSuperConePathsAsLinks(ctx,"fwd",rightptrs,rdepth,chapter,context)
		}

		solutions,_ = SST.WaveFrontsOverlap(ctx,left_paths,right_paths,Lnum,Rnum,ldepth,rdepth)

		if len(solutions) > 0 {

			for s := 0; s < len(solutions); s++ {
				prefix := fmt.Sprintf(" - story path: ")
				PrintConstrainedLinkPath(ctx,solutions,s,prefix,chapter,context,arrowptrs,sttype)
			}
			count++
			break
		}

		if turn % 2 == 0 {
			ldepth++
		} else {
			rdepth++
		}
	}
}

//******************************************************************

func ShowMatchingArrows(ctx SST.PoSST,arrowptrs []SST.ArrowPtr,sttype []int) {

	for a := range arrowptrs {
		adir := SST.GetDBArrowByPtr(ctx,arrowptrs[a])
		fmt.Printf("%3d. (%d) %s -> %s\n",arrowptrs[a],SST.STIndexToSTType(adir.STAindex),adir.Short,adir.Long)
	}

	for st := range sttype {
		adirs := SST.GetDBArrowBySTType(ctx,sttype[st])
		for adir := range adirs {
			fmt.Printf("%3d. (%d) %s -> %s\n",adirs[adir].Ptr,SST.STIndexToSTType(adirs[adir].STAindex),adirs[adir].Short,adirs[adir].Long)
		}
	}
}

//******************************************************************

func ShowMatchingContext(ctx SST.PoSST,s []string) {

	for i := range s {
		res := SST.GetDBContextsMatchingName(ctx,s[i])
		for c := 0; c < len(res); c++ {
			fmt.Printf("%3d. \"%s\"\n",c,res[c])
		}
	}
}

//******************************************************************

func ShowMatchingChapter(ctx SST.PoSST,s string) {

	res := SST.GetDBChaptersMatchingName(ctx,s)
	for c := 0; c < len(res); c++ {
		fmt.Printf("%3d. \"%s\"\n",c,res[c])
	}
}

//******************************************************************

func ShowStories(ctx SST.PoSST,arrows []string,name []string,chapter string,context []string) {

	if arrows == nil {
		arrows = []string{"then"}
	}

	for n := range name {
		for a := range arrows {
			stories := SST.GetSequenceContainers(ctx,arrows[a],name[n],chapter,context)

			for s := range stories {
				// if there is no unique match, the data contain a list of alternatives
				if stories[s].Axis == nil {
					fmt.Printf("%3d. %s\n",s,stories[s].Text)
				} else {
					fmt.Printf("The following story/sequence (%s) \"%s\"\n\n",stories[s].Arrow,stories[s].Text)
					for ev := range stories[s].Axis {
						fmt.Printf("\n%3d. %s\n",ev,stories[s].Axis[ev].Text)
					}
				}
			}
			break
		}
		break
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

// **********************************************************

func ShowNode(ctx SST.PoSST,nptr []SST.NodePtr) string {

	var ret string

	for n := range nptr {
		node := SST.GetDBNodeByNodePtr(ctx,nptr[n])
		ret += fmt.Sprintf("\n    %.30s, ",node.S)
	}

	return ret
}

// **********************************************************

func PrintConstrainedLinkPath(ctx SST.PoSST, cone [][]SST.Link, p int, prefix string,chapter string,context []string,arrows []SST.ArrowPtr,sttype []int) {

	for l := 1; l < len(cone[p]); l++ {
		link := cone[p][l]

		if !ArrowAllowed(ctx,link.Arr,arrows,sttype) {
			return
		}
	}

	SST.PrintLinkPath(ctx,cone,p,prefix,chapter,context)
}

// **********************************************************

func ArrowAllowed(ctx SST.PoSST,arr SST.ArrowPtr, arrlist []SST.ArrowPtr, stlist []int) bool {

	st_ok := false
	arr_ok := false

	staidx := SST.GetDBArrowByPtr(ctx,arr).STAindex
	st := SST.STIndexToSTType(staidx)

	if arrlist != nil {
		for a := range arrlist {
			if arr == arrlist[a] {
				arr_ok = true
				break
			}
		}
	} else {
		arr_ok = true
	}

	if stlist != nil {
		for i := range stlist {
			if stlist[i] == st {
				st_ok = true
				break
			}
		}
	} else {
		st_ok = true
	}

	if st_ok || arr_ok {
		return true
	}

	return false
}

// **********************************************************

func ShowNotes(ctx SST.PoSST,notes []SST.PageMap) {

	var last string
	var lastc string

	for n := 0; n < len(notes); n++ {

		txtctx := SST.ContextString(notes[n].Context)
		
		if last != notes[n].Chapter || lastc != txtctx {

			fmt.Println("\n---------------------------------------------")
			fmt.Println("\nTitle:", notes[n].Chapter)
			fmt.Println("Context:", txtctx)
			fmt.Println("---------------------------------------------\n")

			last = notes[n].Chapter
			lastc = txtctx
		}

		for lnk := 0; lnk < len(notes[n].Path); lnk++ {
			
			text := SST.GetDBNodeByNodePtr(ctx,notes[n].Path[lnk].Dst)
			
			if lnk == 0 {
				fmt.Print("\n",text.S," ")
			} else {
				arr := SST.GetDBArrowByPtr(ctx,notes[n].Path[lnk].Arr)
				fmt.Printf("(%s) %s ",arr.Long,text.S)
			}
		}
	}
}



