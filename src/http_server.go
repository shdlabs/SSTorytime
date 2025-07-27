//
// Simple web server lookup
//

package main

import (
	"fmt"
	"net/http"
	"strings"
	"os"
	"encoding/json"

        SST "SSTorytime"
)

var CTX SST.PoSST

// *********************************************************************

func main() {
	
	CTX = SST.Open(true)	
	
	http.HandleFunc("/",PageHandler)
	http.HandleFunc("/searchN4L",SearchN4LHandler)
	fmt.Println("Listening at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// *********************************************************************

func PageHandler(w http.ResponseWriter, r *http.Request) {

	GenHeader(w,r)

	switch r.Method {
	case "GET":

		w.Header().Set("Content-Type", "text/html")
		page,err := os.ReadFile("./page.html")

		localAddr := r.Context().Value(http.LocalAddrContextKey) 
		ipaddr := fmt.Sprintf("%s",localAddr)
		page = []byte(strings.Replace(string(page),"localhost:8080",ipaddr,-1))

		if err != nil {
			fmt.Println("Can't find ./page.html")
			os.Exit(-1)
		}
		w.Write(page)

	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************
// SEARCH
// *********************************************************************

func SearchN4LHandler(w http.ResponseWriter, r *http.Request) {
	
	GenHeader(w,r)
	
	switch r.Method {
	case "POST","GET":
		name := r.FormValue("name")
		nclass := r.FormValue("nclass")
		ncptr := r.FormValue("ncptr")

		if len(nclass) > 0 && len(ncptr) > 0 {
			// direct click on an item
			var a,b int
			fmt.Sscanf(nclass,"%d",&a)
			fmt.Sscanf(ncptr,"%d",&b)
			nstr := fmt.Sprintf("(%d,%d)",a,b)
			name = name + nstr
		}

		if len(name) == 0 {
			name = "semantic spacetime"
		}
		search := SST.DecodeSearchField(name)
		HandleSearch(search,name,w,r)
		
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func HandleSearch(search SST.SearchParameters,line string,w http.ResponseWriter, r *http.Request) {

	// This is analogous to searchN4L

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

	// OPTIONS *********************************************

	name := search.Name != nil
	from := search.From != nil
	to := search.To != nil
	context := search.Context != nil
	chapter := search.Chapter != ""
	pagenr := search.PageNr > 0
	sequence := search.Sequence

	// Now convert strings into NodePointers

	arrowptrs,sttype := SST.ArrowPtrFromArrowsNames(CTX,search.Arrows)
	nodeptrs := SST.SolveNodePtrs(CTX,search.Name,search.Chapter,search.Context,arrowptrs)
	leftptrs := SST.SolveNodePtrs(CTX,search.From,search.Chapter,search.Context,arrowptrs)
	rightptrs := SST.SolveNodePtrs(CTX,search.To,search.Chapter,search.Context,arrowptrs)

	arrows := arrowptrs != nil
	sttypes := sttype != nil
	limit := 0

	if search.Range > 0 {
		limit = search.Range
	} else {
		limit = 10
	}

	// SEARCH SELECTION *********************************************

	if name && ! sequence && !pagenr {
		fmt.Println("HandleOrbits()")
		HandleOrbit(w,r,nodeptrs,limit)
		return
	}

	if (name && from) || (name && to) {
		fmt.Printf("\nSearch \"%s\" has conflicting parts <to|from> and match strings\n",line)
		os.Exit(-1)
	}

	// Closed path solving, two sets of nodeptrs
	// if we have BOTH from/to (maybe with chapter/context) then we are looking for paths

	if from && to {
		HandlePathSolve(w,r,leftptrs,rightptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
		return
	}

	// Open causal cones, from one of these three

	if (name || from || to) && !pagenr && !sequence {

		if nodeptrs != nil {
			HandleCausalCones(w,r,nodeptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
			return
		}
		if leftptrs != nil {
			HandleCausalCones(w,r,leftptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
			return
		}
		if rightptrs != nil {
			HandleCausalCones(w,r,rightptrs,search.Chapter,search.Context,arrowptrs,sttype,limit)
			return
		}
	}

	// if we have page number then we are looking for notes by pagemap

	if (name || chapter) && pagenr {

		var notes []SST.PageMap

		if chapter {
			notes = SST.GetDBPageMap(CTX,search.Chapter,search.Context,search.PageNr)
			HandlePageMap(w,r,notes)
			return
		} else {
			for n := range search.Name {
				notes = SST.GetDBPageMap(CTX,search.Name[n],search.Context,search.PageNr)
				HandlePageMap(w,r,notes)
			}
			return
		}
	}

	// Look for axial trails following a particular arrow, like _sequence_ 

	if name && sequence || sequence && arrows {
		HandleStories(w,r,CTX,search.Arrows,search.Name,search.Chapter,search.Context,limit)
		return
	}

	// Match existing contexts

	if chapter {
		//ShowMatchingChapter(ctx,search.Chapter)
		return
	}

	if context {
		//ShowMatchingContext(ctx,search.Context)
		return
	}

	// if we have sequence with arrows, then we are looking for sequence context or stories

	if arrows || sttypes {
		//ShowMatchingArrows(ctx,arrowptrs,sttype)
		return
	}

	fmt.Println("Didn't find a solver")
}

// *********************************************************************

func HandleOrbit(w http.ResponseWriter, r *http.Request,nptrs []SST.NodePtr,limit int) {

	w.Header().Set("Content-Type", "application/json")

	fmt.Println("HandleOrbit()")
	var count int
	var array string

	for n := 0; n < len(nptrs); n++ {
		count++
		if count > limit {
			return
		}

		orb := SST.GetNodeOrbit(CTX,nptrs[n],"")

		// create a set of coords for len(nptrs) disconnected nodes
		// SetOrbitCoords(ne.XYZ,axis,nth,len(openings)) .. for DISconnected matches

		array += SST.JSONNodeEvent(CTX,nptrs[n],orb)

		if n < len(nptrs)-1 {
			array += ",\n"
		}
	}

	content := fmt.Sprintf("[ %s ]",array)
	response := PackageResponse("Orbits",content)
	
	fmt.Println("REPLY:\n",string(response))

	w.Write(response)
	fmt.Println("Reply Orbit sent")
}

// *********************************************************************

func HandleCausalCones(w http.ResponseWriter, r *http.Request,nptrs []SST.NodePtr, chap string, context []string,arrows []SST.ArrowPtr, sttype []int,limit int) {

	fmt.Println("HandleCausalCones()")
	var total int = 1
	var data string

	if len(sttype) == 0 {
		sttype = []int{0,1,2,3}
	}

	for n := range nptrs {
		for st := range sttype {

			fmt.Println("Cones from",nptrs[n],"sttype",sttype[st])

			jstr,count := PackageConeFromOrigin(nptrs[n],n,sttype[st],chap,context,len(nptrs),limit)

			if count > 0 {
				total += count
				data += jstr
				data += ","
			}

			if total > limit {
				break
			}
		}

		if total > limit {
			break
		}
	}

	data = strings.Trim(data,",")
	array := fmt.Sprintf("[%s]",data)

	response := PackageResponse("ConePaths",array)
	fmt.Println("CasualConePath reponse",string(response))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

//******************************************************************

func PackageConeFromOrigin(nptr SST.NodePtr,nth int,sttype int,chap string,context []string,dimnptr,limit int) (string,int) {
	// Package a JSON object for nptr's causal cone 

	var wpaths [][]SST.WebPath

	fcone,countf := SST.GetFwdPathsAsLinks(CTX,nptr,sttype,limit)
	wpaths = append(wpaths,SST.LinkWebPaths(CTX,fcone,nth,chap,context,dimnptr,limit)...)

	bcone,countb := SST.GetFwdPathsAsLinks(CTX,nptr,-sttype,limit)
	wpaths = append(wpaths,SST.LinkWebPaths(CTX,bcone,nth,chap,context,dimnptr,limit)...)
	
	wstr,err := json.Marshal(wpaths)

	if wpaths == nil {
		return "",0
	}

	if err != nil {
		fmt.Println("Error in PackageConeFromOrigin",err)
		os.Exit(-1)
	}

	jstr := fmt.Sprintf(" { \"NClass\" : %d,\n",nptr.Class)
	jstr += fmt.Sprintf("   \"NCPtr\" : %d,\n",nptr.CPtr)
	jstr += fmt.Sprintf("   \"Title\" : \"%v\",\n",nptr)  // tbd
	jstr += fmt.Sprintf("   \"Paths\" : %s\n}",string(wstr))	

	return jstr,countf + countb
}

//******************************************************************

func HandlePathSolve(w http.ResponseWriter, r *http.Request,leftptrs,rightptrs []SST.NodePtr,chapter string,context []string,arrowptrs []SST.ArrowPtr,sttype []int,maxdepth int) {

	fmt.Println("HandlePathSolve()")

	var Lnum,Rnum int
	var left_paths, right_paths [][]SST.Link

	// Find the path matrix

	var solutions [][]SST.Link
	var ldepth,rdepth int = 1,1

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {

		left_paths,Lnum = SST.GetEntireNCSuperConePathsAsLinks(CTX,"fwd",leftptrs,ldepth,chapter,context)
		right_paths,Rnum = SST.GetEntireNCSuperConePathsAsLinks(CTX,"bwd",rightptrs,rdepth,chapter,context)

		if Lnum == 0 || Rnum == 0 {
			fmt.Println("Nothing, trying reverse")
			left_paths,Lnum = SST.GetEntireNCSuperConePathsAsLinks(CTX,"bwd",leftptrs,ldepth,chapter,context)
			right_paths,Rnum = SST.GetEntireNCSuperConePathsAsLinks(CTX,"fwd",rightptrs,rdepth,chapter,context)
		}

		solutions,_ = SST.WaveFrontsOverlap(CTX,left_paths,right_paths,Lnum,Rnum,ldepth,rdepth)

		if len(solutions) > 0 {

			// format paths
			var jstr string

			jstr += fmt.Sprintf(" { \"NClass\" : %d,\n",solutions[0][0].Dst.Class)
			jstr += fmt.Sprintf("   \"NCPtr\" : %d,\n",solutions[0][0].Dst.CPtr)
			jstr += fmt.Sprintf("   \"Title\" : \"%s\",\n","path solutions")
			jstr += fmt.Sprintf("   \"BTWC\" : [ %s ],\n",SST.BetweenNessCentrality(CTX,solutions))
			jstr += fmt.Sprintf("   \"Supernodes\" : [ %s ],\n",SST.SuperNodes(CTX,solutions,maxdepth))

			var wpaths [][]SST.WebPath
			nth := 1
			dimnptr := 1

			wpaths = append(wpaths,SST.LinkWebPaths(CTX,solutions,nth,chapter,context,dimnptr,maxdepth)...)

			if wpaths == nil {
				break
			}

			wstr,_ := json.Marshal(wpaths)
			jstr += fmt.Sprintf("   \"Paths\" : %s }",string(wstr))

			array_pack := fmt.Sprintf("[%s]",jstr)
			response := PackageResponse("PathSolve",array_pack)
			fmt.Println("PATH SOLVE:",string(response))
			w.Write(response)
			return
		}

		if turn % 2 == 0 {
			ldepth++
		} else {
			rdepth++
		}
	}
	
	fmt.Println("No paths satisfy constraints")
	response := PackageResponse("PathSolve","")
	w.Write(response)
}

//******************************************************************

func HandlePageMap(w http.ResponseWriter, r *http.Request,notes []SST.PageMap) {

	jstr := SST.JSONPage(CTX,notes)
	response := PackageResponse("PageMap",jstr)
	fmt.Println("PAGEMAP NOTES",string(response))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

//******************************************************************

func HandleStories(w http.ResponseWriter, r *http.Request,ctx SST.PoSST,arrows []string,name []string,chapter string,context []string, limit int) {

	if arrows == nil {
		arrows = []string{"then"}
	}

	fmt.Println("Solver/handler: HandleStories()")

	var jarray string

	for n := range name {
		for a := range arrows {
			stories := SST.GetSequenceContainers(ctx,arrows[a],name[n],chapter,context,limit)

			for s := range stories {
				var jstory string

				for a := 0; a < len(stories[s].Axis); a++ {
					jstr := JSONStoryNodeEvent(stories[s].Axis[a])
					jstory += fmt.Sprintf("%s,",jstr)
				}
				jstory = strings.Trim(jstory,",")
				jarray += fmt.Sprintf("[%s],",jstory)
			}
			break
		}
		break
	}

	data := strings.Trim(jarray,",")
	response := PackageResponse("Sequence",data)
	fmt.Println("Sequence...",string(response))
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

// *********************************************************************

func JSONStoryNodeEvent(en SST.NodeEvent) string {

	var jstr string

	if len(en.Text) == 0 {
		return ""
	}

	jstr += fmt.Sprintf("{\"Text\": \"%s\",\n",SST.EscapeString(en.Text))
	jstr += fmt.Sprintf("\"L\": \"%d\",\n",en.L)
	jstr += fmt.Sprintf("\"Chap\": \"%s\",\n",SST.EscapeString(en.Chap))
	jstr += fmt.Sprintf("\"NPtr\": { \"Class\": \"%d\", \"CPtr\" : \"%d\"},\n",en.NPtr.Class,en.NPtr.CPtr)

	var arrays string

	for sti := 0; sti < SST.ST_TOP; sti++ {
		var arr string
		if en.Orbits[sti] != nil {
			js,_ := json.Marshal(en.Orbits[sti])
			arr = fmt.Sprintf("%s,",string(js))
		} else {
			arr = "[],"
		}
		arrays += arr
	}
	arrays = strings.Trim(arrays,",")
	jstr += fmt.Sprintf("\"Orbits\": [%s] }",arrays)
	return jstr
}

// *********************************************************************

func TableOfContents(w http.ResponseWriter, r *http.Request) {

	fmt.Println("TableOfContents handler")

	switch r.Method {
	case "POST","GET":
		chapter := r.FormValue("chapter")
		cntstr := r.FormValue("context")
		context,_ := SST.Str2Array(cntstr)

		toc := SST.JSON_TableOfContents(CTX,chapter,context)
		w.Write([]byte(toc))
		fmt.Println(toc)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func GenHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Add("Vary", "Origin")
}

// *********************************************************************

func CleanText(c string) string {

	c = strings.Replace(c,"{","",-1)
	c = strings.Replace(c,"}","",-1)
	c = strings.Replace(c,","," ",-1)
	c = strings.Replace(c,"\"","\\\"",-1)
	return c
}

// **********************************************************

func ShowNode(ctx SST.PoSST,nptr []SST.NodePtr) string {

	var ret string

	for n := 0; n < len(nptr); n++ {
		node := SST.GetDBNodeByNodePtr(ctx,nptr[n])
		ret += fmt.Sprintf("%.30s",node.S)
		if n < len(nptr)-1 {
			ret += ","
		}
	}

	return ret
}

// **********************************************************

func PackageResponse(kind string, jstr string) []byte {

	response := fmt.Sprintf("{ \"Response\" : \"%s\",\n \"Content\" : %s }",kind,jstr)

	return []byte(response)
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














