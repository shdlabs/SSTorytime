//******************************************************************
//
//  Web server for lookup requests and JSON interface
//
//******************************************************************

package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	SST "github.com/shdlabs/SSTorytime/services/sstorytime"
	"github.com/shdlabs/SSTorytime/services/text2n4l"
)

// Ugly Go directive to embed text files into the binary

//go:embed all:static
var content embed.FS

// *********************************************************************

var CTX SST.PoSST // just one persistent connection

// *********************************************************************
// Main
// *********************************************************************

func main() {

	CTX = SST.Open(true)

	// 1. Create the filesystem view rooted inside the "public" directory.

	publicFS, err := fs.Sub(content, "static")

	if err != nil {
		log.Fatal("failed to create sub-filesystem:", err)
	}

	// 2. Create a router (ServeMux) and register handlers.

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(publicFS))

	mux.Handle("/", fileServer)
	mux.HandleFunc("/searchN4L", SearchN4LHandler)
	mux.HandleFunc("/api/text2n4l/process", Text2N4LHandler)
	mux.HandleFunc("/status", StatusHandler)
	mux.HandleFunc("/debug", DebugHandler) // Debug endpoint to test form parameters

	// 3. Create an http.Server instance for graceful shutdown.

	srv := &http.Server{Addr: "0.0.0.0:8081", Handler: EnableCORS(mux)}

	// 4. Run the server in a goroutine so it doesn't block.

	go func() {
		log.Println("Server starting on http://localhost:8080")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start server: %s\n", err)
		}
	}()

	// 5. Wait for an interrupt signal.

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	// 6. Perform a graceful shutdown with a timeout.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %s\n", err)
	}

	log.Println("Server exited properly")
}

// *********************************************************************
// Error handling and safety functions
// *********************************************************************

// SafeSolveNodePtrs wraps SST.SolveNodePtrs with error handling
func SafeSolveNodePtrs(ctx SST.PoSST, names []string, search SST.SearchParameters, arrowptrs []SST.ArrowPtr, limit int) ([]SST.NodePtr, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in SolveNodePtrs: %v", r)
		}
	}()

	// Handle common problematic cases
	if len(names) > 0 {
		for _, name := range names {
			// Allow "any" when there are context constraints that make the search more specific
			if (name == "any" || name == "%%") && len(search.Context) == 0 {
				return nil, fmt.Errorf("'%s' is too broad a search term and causes database conflicts. Please use more specific terms or try \\help for search guidance", name)
			}
			if len(name) < 2 && name != "%%" {
				return nil, fmt.Errorf("search term '%s' is too short. Please use at least 2 characters", name)
			}
		}
	}

	// Call the actual function with additional error recovery
	var result []SST.NodePtr

	// Wrap the call with recovery for any internal panics
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Internal panic in SST.SolveNodePtrs: %v", r)
				// Don't re-panic, just log the error
			}
		}()
		result = SST.SolveNodePtrs(ctx, names, search, arrowptrs, limit)
	}()

	return result, nil
}

// SendErrorResponse sends a structured error response to the client
func SendErrorResponse(w http.ResponseWriter, errorType, message, query string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	errorResponse := map[string]interface{}{
		"Response":  "ERROR",
		"ErrorType": errorType,
		"Message":   message,
		"Query":     query,
		"Suggestions": []string{
			"Try using more specific search terms",
			"Use \\help for search guidance",
			"Check the documentation at \\notes \\chapter \"help and search\"",
			"Avoid very common words like 'any', 'the', 'a'",
		},
	}

	json.NewEncoder(w).Encode(errorResponse)
}

// *********************************************************************
// Handlers
// *********************************************************************
// *********************************************************************

func EnableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set the Access-Control-Allow-Origin header to the origin of the request.

		origin := r.Header.Get("Origin")

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Browsers send a pre-flight OPTIONS request for CORS. We need to handle it.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// *********************************************************************
// Handlers
// *********************************************************************

func SignalHandler() {

	signal_chan := make(chan os.Signal, 1)

	signal.Notify(signal_chan,
		syscall.SIGHUP,  // 1
		syscall.SIGINT,  // 2 ctrl-c
		syscall.SIGQUIT, // 3
		syscall.SIGTERM) // 15, CTRL-c

	sig := <-signal_chan // block until signal

	switch sig {

	case syscall.SIGHUP:
		fmt.Println("hungup")

	case syscall.SIGINT:
		fmt.Println("Warikomi, cutting in, sandoichi")

	case syscall.SIGTERM:
		fmt.Println("force stop")

	case syscall.SIGQUIT:
		fmt.Println("stop and core dump")

	default:
		fmt.Println("Unknown signal.")
	}
}

// *********************************************************************
// SEARCH
// *********************************************************************

func SearchN4LHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST", "GET":
		name := r.FormValue("name")
		nclass := r.FormValue("nclass")
		ncptr := r.FormValue("ncptr")
		chapcontext := r.FormValue("chapcontext")

		fmt.Printf("DEBUG: Received parameters - name='%s', nclass='%s', ncptr='%s', chapcontext='%s'\n", name, nclass, ncptr, chapcontext)

		if name == "\\lastnptr" {
			if chapcontext != "" && chapcontext != "any" {
				UpdateLastSawSection(w, r, chapcontext)
			}
			UpdateLastSawNPtr(w, r, nclass, ncptr, chapcontext)
			return
		}

		if name == "" && len(nclass) > 0 && len(ncptr) > 0 {
			// direct click on an item
			var a, b int
			fmt.Sscanf(nclass, "%d", &a)
			fmt.Sscanf(ncptr, "%d", &b)
			nstr := fmt.Sprintf("(%d,%d)", a, b)
			name = name + nstr
		}

		fmt.Println("\nReceived command:", name)

		ambient, key, _ := SST.GetTimeContext()

		if name == "\\remind" {
			name = "any \\chapter reminders \\context " + key + " " + ambient
		} else if name == "\\help" {
			name = "\\notes \\chapter \"help and search\" \\limit 40"
		} else if len(name) == 0 {
			// Default query for initial page load - show a welcome/intro
			name = "\\notes \\chapter \"welcome\" \\limit 10"
		}

		search := SST.DecodeSearchField(name)

		HandleSearch(search, name, w, r)

	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

func UpdateLastSawSection(w http.ResponseWriter, r *http.Request, query string) {

	// update lastseen db

	fmt.Println("UPDATING STATS FOR section", query)

	SST.UpdateLastSawSection(CTX, query)
}

// *********************************************************************

func UpdateLastSawNPtr(w http.ResponseWriter, r *http.Request, class, cptr string, classifier string) {

	// update lastseen db

	var nptr SST.NodePtr
	var nclass int
	var ncptr int
	fmt.Sscanf(class, "%d", &nclass)
	fmt.Sscanf(cptr, "%d", &ncptr)
	nptr.Class = nclass
	nptr.CPtr = SST.ClassedNodePtr(ncptr)

	SST.UpdateLastSawNPtr(CTX, nclass, ncptr, classifier)

	fmt.Println("UPDATING STATS FOR", nclass, ncptr, "WITHIN", classifier)

	SST.UpdateLastSawSection(CTX, classifier)

	// Create proper JSON response
	lastSawResponse := map[string]string{
		"Response": "LastSaw",
		"Content":  fmt.Sprintf("ack(%s,%s)", class, cptr),
	}
	responseData, _ := json.Marshal(lastSawResponse)
	w.Write(responseData)

}

// *********************************************************************

func HandleSearch(search SST.SearchParameters, line string, w http.ResponseWriter, r *http.Request) {

	// Add error recovery to prevent server crashes
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in HandleSearch: %v", r)
			errorMsg := fmt.Sprintf("Search error: %v", r)
			SendErrorResponse(w, "SEARCH_ERROR", errorMsg, line)
		}
	}()

	// This is analogous to searchN4L

	// OPTIONS *********************************************

	name := search.Name != nil
	from := search.From != nil
	to := search.To != nil
	context := search.Context != nil
	chapter := search.Chapter != ""
	pagenr := search.PageNr > 0
	sequence := search.Sequence

	// Now convert strings into NodePointers

	arrowptrs, sttype := SST.ArrowPtrFromArrowsNames(CTX, search.Arrows)

	arrows := arrowptrs != nil
	sttypes := sttype != nil
	limit := 0

	if search.Range > 0 {
		limit = search.Range
	} else {
		if from || to || sequence {
			limit = 30 // many paths make hard work
		} else {
			const common_word = 5

			if SST.SearchTermLen(search.Name) < common_word {
				limit = 5
			} else {
				limit = 10
			}
		}
	}

	fmt.Println()
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tabWriter, "start set:\t", SL(search.Name))
	fmt.Fprintln(tabWriter, "from:\t", SL(search.From))
	fmt.Fprintln(tabWriter, "to:\t", SL(search.To))
	fmt.Fprintln(tabWriter, "chapter:\t", search.Chapter)
	fmt.Fprintln(tabWriter, "context:\t", SL(search.Context))
	fmt.Fprintln(tabWriter, "arrows:\t", SL(search.Arrows))
	fmt.Fprintln(tabWriter, "pageNR:\t", search.PageNr)
	fmt.Fprintln(tabWriter, "sequence/story:\t", search.Sequence)
	fmt.Fprintln(tabWriter, "limit/range/depth:\t", limit)
	fmt.Fprintln(tabWriter, "show stats:\t", search.Stats)

	tabWriter.Flush()
	fmt.Println()

	var nodeptrs, leftptrs, rightptrs []SST.NodePtr
	var searchError error

	// Wrap potentially problematic search operations in error handling
	if !pagenr && !sequence {
		leftptrs, searchError = SafeSolveNodePtrs(CTX, search.From, search, arrowptrs, limit)
		if searchError != nil {
			SendErrorResponse(w, "SEARCH_ERROR", fmt.Sprintf("Error searching 'from' terms: %v", searchError), line)
			return
		}
		rightptrs, searchError = SafeSolveNodePtrs(CTX, search.To, search, arrowptrs, limit)
		if searchError != nil {
			SendErrorResponse(w, "SEARCH_ERROR", fmt.Sprintf("Error searching 'to' terms: %v", searchError), line)
			return
		}
	}

	nodeptrs, searchError = SafeSolveNodePtrs(CTX, search.Name, search, arrowptrs, limit)
	if searchError != nil {
		SendErrorResponse(w, "SEARCH_ERROR", fmt.Sprintf("Error searching terms: %v", searchError), line)
		return
	}

	fmt.Println("Solved search nodes ...")

	// SEARCH SELECTION *********************************************

	// Table of contents

	if search.Stats {
		ShowStats(w, r, CTX, search, nodeptrs)
		return
	}

	if (context || chapter) && !name && !sequence && !pagenr && !(from || to) {
		ShowChapterContexts(w, r, CTX, search, limit)
		return
	}

	if name && !sequence && !pagenr {
		HandleOrbit(w, r, CTX, search, nodeptrs, limit)
		return
	}

	if (name && from) || (name && to) {
		fmt.Printf("\nSearch \"%s\" has conflicting parts <to|from> and match strings\n", line)
		os.Exit(-1)
	}

	// Closed path solving, two sets of nodeptrs
	// if we have BOTH from/to (maybe with chapter/context) then we are looking for paths

	if from && to {
		HandlePathSolve(w, r, CTX, leftptrs, rightptrs, search, arrowptrs, sttype, limit)
		return
	}

	// Open causal cones, from one of these three

	if (name || from || to) && !pagenr && !sequence {

		if nodeptrs != nil {
			HandleCausalCones(w, r, CTX, nodeptrs, search, arrowptrs, sttype, limit)
			return
		}
		if leftptrs != nil {
			HandleCausalCones(w, r, CTX, leftptrs, search, arrowptrs, sttype, limit)
			return
		}
		if rightptrs != nil {
			HandleCausalCones(w, r, CTX, rightptrs, search, arrowptrs, sttype, limit)
			return
		}
	}

	// if we have page number then we are looking for notes by pagemap

	if (name || chapter) && pagenr {

		var notes []SST.PageMap

		if chapter {
			notes = SST.GetDBPageMap(CTX, search.Chapter, search.Context, search.PageNr)
			HandlePageMap(w, r, CTX, search, notes)
			return
		} else {
			for n := range search.Name {
				notes = SST.GetDBPageMap(CTX, search.Name[n], search.Context, search.PageNr)
				HandlePageMap(w, r, CTX, search, notes)
			}
			return
		}
	}

	// Look for axial trails following a particular arrow, like _sequence_

	if sequence {
		HandleStories(w, r, CTX, search, nodeptrs, arrowptrs, sttype, limit)
		return
	}

	// if we have sequence with arrows, then we are looking for sequence context or stories

	if arrows || sttypes {
		HandleMatchingArrows(w, r, CTX, search, arrowptrs, sttype)
		return
	}

	fmt.Println("Didn't find a solver")
}

// *********************************************************************

func HandleOrbit(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, nptrs []SST.NodePtr, limit int) {

	// Add panic recovery to prevent server crashes
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in HandleOrbit: %v", r)
			errorMsg := fmt.Sprintf("Error processing node orbit: %v", r)
			SendErrorResponse(w, "ORBIT_ERROR", errorMsg, fmt.Sprintf("%v", search.Name))
		}
	}()

	var count int
	var array []SST.NodeEvent

	origin := SST.Coords{X: 0.0, Y: 0.0, Z: 0.0}

	for n := 0; n < len(nptrs); n++ {

		count++

		if count > limit {
			break
		}

		fmt.Printf("Assembling Node Orbit(%v)\n", nptrs[n])

		// Wrap potentially problematic SST calls with proper error handling
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic processing node %v: %v", nptrs[n], r)
					// Skip this node and continue with the next one
				}
			}()

			orb := SST.GetNodeOrbit(CTX, nptrs[n], "", limit)
			// create a set of coords for len(nptrs) disconnected nodes

			xyz := SST.RelativeOrbit(origin, SST.R0, n, len(nptrs))
			orb = SST.SetOrbitCoords(xyz, orb)

			nodeevent := SST.JSONNodeEvent(CTX, nptrs[n], xyz, orb)
			array = append(array, nodeevent)
		}()
	}

	// Check if we have any valid results
	if len(array) == 0 {
		SendErrorResponse(w, "NO_RESULTS", "No valid results found for the search query. This might be due to database conflicts or overly broad search terms.", fmt.Sprintf("%v", search.Name))
		return
	}

	response := PackageResponse(ctx, search, "Orbits", array)

	//fmt.Println("REPLY:\n",string(response))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println("Reply Orbit sent")
}

// *********************************************************************

func HandleCausalCones(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, nptrs []SST.NodePtr, search SST.SearchParameters, arrows []SST.ArrowPtr, sttype []int, limit int) {

	chap := search.Chapter
	context := search.Context

	fmt.Println("HandleCausalCones()", nptrs)
	var total int = 1

	if len(sttype) == 0 {
		sttype = []int{0, 1, 2, 3}
	}

	var cones []SST.WebConePaths

	for n := range nptrs {
		for st := range sttype {

			subcone, count := PackageConeFromOrigin(ctx, nptrs[n], n, sttype[st], chap, context, len(nptrs), limit)
			cones = append(cones, subcone)

			total += count

			if total > limit {
				break
			}
		}

		if total > limit {
			break
		}
	}

	array, _ := json.Marshal(cones)

	response := PackageResponse(ctx, search, "ConePaths", string(array))
	//fmt.Println("CasualConePath reponse",string(response))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println("Done/sent cone")
}

//******************************************************************

func PackageConeFromOrigin(ctx SST.PoSST, nptr SST.NodePtr, nth int, sttype int, chap string, context []string, dimnptr, limit int) (SST.WebConePaths, int) {

	// Package a JSON object for the nth/dimnptr causal cone , assigning each nth the same width

	var wpaths [][]SST.WebPath

	fcone, count := SST.GetFwdPathsAsLinks(CTX, nptr, sttype, limit, limit)
	wpaths = append(wpaths, SST.LinkWebPaths(CTX, fcone, nth, chap, context, dimnptr, limit)...)

	if sttype != 0 {
		bcone, countb := SST.GetFwdPathsAsLinks(CTX, nptr, -sttype, limit, limit)
		wpaths = append(wpaths, SST.LinkWebPaths(CTX, bcone, nth, chap, context, dimnptr, limit)...)
		count += countb
	}

	var subcone SST.WebConePaths
	subcone.RootNode = nptr
	subcone.Title = SST.GetDBNodeByNodePtr(ctx, nptr).S
	subcone.Paths = wpaths

	return subcone, count
}

//******************************************************************

func HandlePathSolve(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, leftptrs, rightptrs []SST.NodePtr, search SST.SearchParameters, arrowptrs []SST.ArrowPtr, sttype []int, maxdepth int) {

	chapter := search.Chapter
	context := search.Context

	fmt.Println("HandlePathSolve(", leftptrs, ",", rightptrs, ")")

	var Lnum, Rnum int
	var left_paths, right_paths [][]SST.Link

	// Find the path matrix

	var solutions [][]SST.Link
	var ldepth, rdepth int = 2, 2
	var array_pack []SST.WebConePaths

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {

		left_paths, Lnum = SST.GetEntireNCSuperConePathsAsLinks(CTX, "fwd", leftptrs, ldepth, chapter, context, maxdepth)
		right_paths, Rnum = SST.GetEntireNCSuperConePathsAsLinks(CTX, "bwd", rightptrs, rdepth, chapter, context, maxdepth)

		if Lnum == 0 || Rnum == 0 {
			fmt.Println("Nothing, trying reverse")
			left_paths, Lnum = SST.GetEntireNCSuperConePathsAsLinks(CTX, "bwd", leftptrs, ldepth, chapter, context, maxdepth)
			right_paths, Rnum = SST.GetEntireNCSuperConePathsAsLinks(CTX, "fwd", rightptrs, rdepth, chapter, context, maxdepth)

			if Lnum == 0 || Rnum == 0 {
				fmt.Println("No paths")
				response := PackageResponse(ctx, search, "PathSolve", "")
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
				return
			}
		}

		solutions, _ = SST.WaveFrontsOverlap(CTX, left_paths, right_paths, Lnum, Rnum, ldepth, rdepth)

		if len(solutions) > 0 {
			// format paths

			var pack []SST.WebConePaths
			var soln SST.WebConePaths
			var array_pack []SST.WebConePaths
			var the_cone SST.WebConePaths

			soln.RootNode = solutions[0][0].Dst
			soln.Title = "path solutions"
			soln.BTWC = SST.BetweenNessCentrality(CTX, solutions)
			soln.SuperNodes = SST.SuperNodes(CTX, solutions, maxdepth)

			var wpaths [][]SST.WebPath
			nth := 0
			swimlanes := 1

			wpaths = append(wpaths, SST.LinkWebPaths(CTX, solutions, nth, chapter, context, swimlanes, maxdepth)...)

			if wpaths == nil {
				break
			}

			soln.Paths = wpaths
			pack = append(pack, soln)
			the_cone = soln
			array_pack = append(array_pack, the_cone)
		}

		response := PackageResponse(ctx, search, "PathSolve", array_pack)
		//fmt.Println("PATH SOLVE:",string(response))

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return

		if turn%2 == 0 {
			ldepth++
		} else {
			rdepth++
		}
	}

	fmt.Println("No paths satisfy constraints")
	response := PackageResponse(ctx, search, "PathSolve", "[]")

	//fmt.Println("PATHSOLVE NOTES",string(response))
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println("Done/sent path solve")
}

//******************************************************************

func HandlePageMap(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, notes []SST.PageMap) {

	fmt.Println("Solver/handler: HandlePageMap()")

	// Get the JSON string from the existing function
	jstr := SST.JSONPage(CTX, notes)

	// Parse it back into a struct to avoid double JSON encoding
	var pageData map[string]interface{}
	if err := json.Unmarshal([]byte(jstr), &pageData); err != nil {
		log.Printf("Error parsing page data: %v", err)
		pageData = map[string]interface{}{
			"Title":   "Error",
			"Context": "",
			"Notes":   []interface{}{},
		}
	}

	response := PackageResponse(ctx, search, "PageMap", pageData)

	if notes != nil {
		UpdateLastSawSection(w, r, notes[0].Chapter)
	}

	//fmt.Println("PAGEMAP NOTES",string(response))
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println("Done/sent pagemap")
}

//******************************************************************

func HandleStories(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, nodeptrs []SST.NodePtr, arrowptrs []SST.ArrowPtr, sttypes []int, limit int) {

	if arrowptrs == nil {
		arrowptrs, sttypes = SST.ArrowPtrFromArrowsNames(CTX, []string{"!then!"})
	}

	fmt.Println("Solver/handler: HandleStories()")

	stories := SST.GetSequenceContainers(ctx, nodeptrs, arrowptrs, sttypes, limit)

	// Convert stories to proper JSON structure
	var storyArray [][]StoryNodeEvent

	for s := 0; s < len(stories); s++ {
		var story []StoryNodeEvent

		for a := 0; a < len(stories[s].Axis); a++ {
			nodeEventData, err := JSONStoryNodeEvent(stories[s].Axis[a])
			if err != nil || len(nodeEventData) == 0 {
				continue // Skip empty or errored events
			}

			var storyEvent StoryNodeEvent
			if err := json.Unmarshal(nodeEventData, &storyEvent); err == nil {
				story = append(story, storyEvent)
			}
		}

		if len(story) > 0 {
			storyArray = append(storyArray, story)
		}
	}

	response := PackageResponse(ctx, search, "Sequence", storyArray)

	//fmt.Println("Sequence...",string(response))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println("Done/sent sequence")

}

// *********************************************************************

func HandleMatchingArrows(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, arrowptrs []SST.ArrowPtr, sttype []int) {

	fmt.Println("Solver/handler: HandleMatchingArrows()")

	type ArrowList struct {
		ArrPtr  SST.ArrowPtr
		ASTtype int
		Short   string
		Long    string
		InvPtr  SST.ArrowPtr
		ISTtype int
		InvS    string
		InvL    string
	}

	var arrows []ArrowList

	for a := range arrowptrs {
		adir := SST.GetDBArrowByPtr(ctx, arrowptrs[a])
		inv := SST.GetDBArrowByPtr(ctx, SST.INVERSE_ARROWS[arrowptrs[a]])

		var al ArrowList
		al.ArrPtr = arrowptrs[a]
		al.ASTtype = SST.STIndexToSTType(adir.STAindex)
		al.Short = adir.Short
		al.Long = adir.Long
		al.InvPtr = inv.Ptr
		al.ISTtype = SST.STIndexToSTType(inv.STAindex)
		al.InvS = inv.Short
		al.InvL = inv.Long
		arrows = append(arrows, al)
	}

	if arrowptrs == nil {
		for st := range sttype {
			adirs := SST.GetDBArrowBySTType(ctx, sttype[st])
			for adir := range adirs {
				inv := SST.GetDBArrowByPtr(ctx, SST.INVERSE_ARROWS[adirs[adir].Ptr])

				var al ArrowList
				al.ArrPtr = adirs[adir].Ptr
				al.ASTtype = SST.STIndexToSTType(adirs[adir].STAindex)
				al.Short = adirs[adir].Short
				al.Long = adirs[adir].Long
				al.InvPtr = inv.Ptr
				al.ISTtype = SST.STIndexToSTType(inv.STAindex)
				al.InvS = inv.Short
				al.InvL = inv.Long
				arrows = append(arrows, al)
			}
		}
	}

	data, _ := json.Marshal(arrows)
	response := PackageResponse(ctx, search, "Arrows", string(data))

	fmt.Println("Arrows...", string(response))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println("Done/sent arrows")
}

// *********************************************************************

func ShowStats(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, nptrs []SST.NodePtr) {

	var retval []SST.LastSeen

	if nptrs == nil {
		retval = SST.GetLastSawSection(ctx)
	} else {

		for n := range nptrs {
			nptr := SST.GetLastSawNPtr(ctx, nptrs[n])
			retval = append(retval, nptr)
		}
	}

	// Check if any progress tracking data exists
	if len(retval) == 0 {
		// No checkboxes have been checked yet - return friendly guidance
		ambien, key, now := SST.GetTimeContext()
		now_ctx := SST.UpdateSTMContext(CTX, ambien, key, now, search)

		emptyResponse := map[string]interface{}{
			"Response": "GUIDANCE",
			"Content":  "No progress tracking data found. Start exploring content and check the progress boxes on items you've reviewed to build your learning statistics.",
			"Suggestions": []string{
				"Browse chapters using \\notes or \\chapter commands",
				"Check progress boxes on items you've read",
				"Return to \\stats once you've marked some progress",
				"Try \\help for more search guidance",
			},
			"Time":    key,
			"Intent":  now_ctx,
			"Ambient": ambien,
		}

		data, _ := json.Marshal(emptyResponse)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		fmt.Println("Sent guidance for empty stats")
		return
	}

	data, _ := json.Marshal(retval)

	response := PackageResponse(ctx, search, "STAT", string(data))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println("Done/sent stat")

}

// *********************************************************************

func ShowChapterContexts(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, limit int) {

	chap := search.Chapter
	context := search.Context

	fmt.Println("Solver/handler: ShowChapterContexts()")

	var chapters []SST.ChCtx
	var chap_list []string

	toc := SST.GetChaptersByChapContext(ctx, chap, context, limit)

	for chaps := range toc {
		chap_list = append(chap_list, chaps)
	}

	sort.Strings(chap_list)

	for c := 0; c < len(chap_list); c++ {

		var chap_anchor SST.ChCtx

		chap_anchor.Chapter = chap_list[c]
		chap_anchor.XYZ = SST.AssignChapterCoordinates(c, len(chap_list))

		// Fractionate the (chapter,context) information

		dim, clist, adj := SST.IntersectContextParts(toc[chap_list[c]])
		spectrum := SST.GetContextTokenFrequencies(toc[chap_list[c]])
		intent, ambient := SST.ContextIntentAnalysis(spectrum, toc[chap_list[c]])

		chap_anchor.Context = GetContextSets(dim, clist, adj, chap_anchor.XYZ)
		chap_anchor.Single = GetContextFragments(intent, chap_anchor.XYZ)
		chap_anchor.Common = GetContextFragments(ambient, chap_anchor.XYZ)

		chapters = append(chapters, chap_anchor)
	}

	response := PackageResponse(ctx, search, "TOC", chapters)

	fmt.Println("Chap/context...", string(response))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println("Done/sent content")
}

//******************************************************************

func GetContextSets(dim int, clist []string, adj [][]int, xyz SST.Coords) []SST.Loc {

	var retvar []SST.Loc

	for c := 0; c < len(adj); c++ {

		var contextgroup SST.Loc

		contextgroup.Text = clist[c]

		for cp := 0; cp < len(adj[c]); cp++ {
			if adj[c][cp] > 0 {
				contextgroup.Reln = append(contextgroup.Reln, cp)
			}
		}

		contextgroup.XYZ = SST.AssignContextSetCoordinates(xyz, c, len(adj))

		retvar = append(retvar, contextgroup)
	}
	return retvar
}

//******************************************************************

func GetContextFragments(clist []string, ooo SST.Coords) []SST.Loc {

	var retvar []SST.Loc

	for c := 0; c < len(clist); c++ {

		var contextgroup SST.Loc

		contextgroup.Text = clist[c]
		contextgroup.XYZ = SST.AssignFragmentCoordinates(ooo, c, len(clist))

		retvar = append(retvar, contextgroup)
	}
	return retvar
}

// *********************************************************************
// Misc
// *********************************************************************

// StoryNodeEvent represents a JSON-serializable version of SST.NodeEvent
type StoryNodeEvent struct {
	Text    string        `json:"Text"`
	L       int           `json:"L"`
	Chap    string        `json:"Chap"`
	Context string        `json:"Context"`
	NPtr    SST.NodePtr   `json:"NPtr"`
	XYZ     SST.Coords    `json:"XYZ"`
	Orbits  []interface{} `json:"Orbits"`
}

func JSONStoryNodeEvent(en SST.NodeEvent) ([]byte, error) {
	if len(en.Text) == 0 {
		return []byte(""), nil
	}

	// Convert SST.NodeEvent to our JSON-friendly struct
	story := StoryNodeEvent{
		Text:    en.Text,
		L:       en.L,
		Chap:    en.Chap,
		Context: en.Context,
		NPtr:    en.NPtr,
		XYZ:     en.XYZ,
		Orbits:  make([]interface{}, SST.ST_TOP),
	}

	// Handle orbits array
	for sti := 0; sti < SST.ST_TOP; sti++ {
		if en.Orbits[sti] != nil {
			story.Orbits[sti] = en.Orbits[sti]
		} else {
			story.Orbits[sti] = []interface{}{}
		}
	}

	return json.Marshal(story)
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

	c = strings.Replace(c, "{", "", -1)
	c = strings.Replace(c, "}", "", -1)
	c = strings.Replace(c, ",", " ", -1)
	c = strings.Replace(c, "\"", "\\\"", -1)
	return c
}

// **********************************************************

func ShowNode(ctx SST.PoSST, nptr []SST.NodePtr) string {

	var ret string

	for n := range nptr {
		node := SST.GetDBNodeByNodePtr(ctx, nptr[n])
		ret += fmt.Sprintf("%.30s", node.S)
		if n < len(nptr)-1 {
			ret += ","
		}
	}

	return ret
}

// **********************************************************

// APIResponse represents the standard API response format
type APIResponse struct {
	Response string      `json:"Response"`
	Content  interface{} `json:"Content"`
	Time     string      `json:"Time"`
	Intent   interface{} `json:"Intent"`
	Ambient  interface{} `json:"Ambient"`
}

func PackageResponse(ctx SST.PoSST, search SST.SearchParameters, kind string, content interface{}) []byte {
	ambien, key, now := SST.GetTimeContext()
	now_ctx := SST.UpdateSTMContext(CTX, ambien, key, now, search)

	response := APIResponse{
		Response: kind,
		Content:  content,
		Time:     key,
		Intent:   now_ctx,
		Ambient:  ambien,
	}

	data, err := json.Marshal(response)
	if err != nil {
		// Fallback to error response
		errorResp := APIResponse{
			Response: "Error",
			Content:  "Failed to marshal response",
			Time:     key,
			Intent:   nil,
			Ambient:  nil,
		}
		data, _ = json.Marshal(errorResp)
	}

	return data
}

//******************************************************************

func SL(list []string) string {
	if list == nil {
		return " []"
	}
	// Use strings.Join for efficient concatenation and fmt.Sprintf for formatting.
	return fmt.Sprintf(" [%s]", strings.Join(list, ", "))
}

// Text2N4LRequest defines the structure for text processing requests
type Text2N4LRequest struct {
	Text       string  `json:"text"`
	Percentage float64 `json:"percentage"`
}

// Text2N4LResponse defines the structure for text processing responses
type Text2N4LResponse struct {
	N4LContent string `json:"n4l_content"`
	Stats      struct {
		TotalSentences    int     `json:"total_sentences"`
		SelectedSentences int     `json:"selected_sentences"`
		FinalFraction     float64 `json:"final_fraction"`
		RequestedFraction float64 `json:"requested_fraction"`
	} `json:"stats"`
	Error string `json:"error,omitempty"`
}

// Text2N4LHandler processes text input and returns N4L formatted output
func Text2N4LHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var req Text2N4LRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Text2N4LResponse{
			Error: "Invalid JSON request: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate input
	if strings.TrimSpace(req.Text) == "" {
		response := Text2N4LResponse{
			Error: "Text field is required and cannot be empty",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set default percentage if not provided or invalid
	if req.Percentage <= 0 || req.Percentage > 100 {
		req.Percentage = 10.0
	}

	// Process the text
	result, err := text2n4l.ProcessTextContent(req.Text, req.Percentage)
	if err != nil {
		response := Text2N4LResponse{
			Error: "Failed to process text: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Prepare successful response
	response := Text2N4LResponse{
		N4LContent: result.N4LContent,
	}
	response.Stats.TotalSentences = result.TotalSentences
	response.Stats.SelectedSentences = result.SelectedSentences
	response.Stats.FinalFraction = result.FinalFraction
	response.Stats.RequestedFraction = result.RequestedFraction

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// StatusResponse defines the structure for our JSON response.
type StatusResponse struct {
	ServerStatus    string    `json:"server_status"`
	DatabaseStatus  string    `json:"database_status"`
	AvailableTopics []string  `json:"available_topics"`
	Timestamp       time.Time `json:"timestamp"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	toc := SST.GetChaptersByChapContext(CTX, "", nil, 1000) // "" for chapter and nil for context should get all

	var topics []string
	for chapter := range toc {
		topics = append(topics, chapter)
	}
	sort.Strings(topics)

	// Create the response object.
	status := StatusResponse{
		ServerStatus:    "OK",
		DatabaseStatus:  "OK",
		AvailableTopics: topics,
		Timestamp:       time.Now(),
	}

	// Marshal the struct into JSON.
	responseJSON, err := json.Marshal(status)
	if err != nil {
		http.Error(w, "Failed to generate status response", http.StatusInternalServerError)
		return
	}

	// Set the content type and send the response.
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// DebugHandler shows all received parameters for debugging
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
	}

	response := fmt.Sprintf(`{
		"method": "%s",
		"name": "%s",
		"nclass": "%s", 
		"ncptr": "%s",
		"chapcontext": "%s",
		"all_form_values": %q
	}`,
		r.Method,
		r.FormValue("name"),
		r.FormValue("nclass"),
		r.FormValue("ncptr"),
		r.FormValue("chapcontext"),
		r.Form.Encode())

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}
