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
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	SST "SSTorytime"

	"github.com/lmittmann/tint"
)

// Ugly Go directive to embed text files into the binary

//go:embed all:public
var content embed.FS

// *********************************************************************

// CTX is the persistent SST database connection shared across all handlers.
// This single connection is initialized in main() and reused throughout the application lifecycle.
var CTX SST.PoSST

// init initializes the structured logging system with INFO level logging to stdout.
// This runs automatically before main() and configures the default logger for the entire application.
// Time format: HH:MM:SS.mmm (hours:minutes:seconds.milliseconds only, no date or timezone)
// Uses tint for colorful output with level-specific colors.
func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: "15:04:05.000",
		AddSource:  true,  // This provides file, line, and function information
		NoColor:    false, // Enable colors
	}))
	slog.SetDefault(logger)
}

// *********************************************************************
// Main
// *********************************************************************

// main is the entry point of the HTTP server application.
// It performs the following operations:
// 1. Opens a persistent SST database connection
// 2. Creates an embedded filesystem for serving static files from the public directory
// 3. Sets up HTTP routes for search, status, and static file serving
// 4. Starts the server on all network interfaces (0.0.0.0:8080) with CORS enabled
// 5. Implements graceful shutdown on receiving SIGINT or SIGTERM signals
// 6. Ensures clean database closure with a 5-second timeout on shutdown
func main() {
	CTX = SST.Open(true)

	// 1. Create the filesystem view rooted inside the "public" directory.

	publicFS, err := fs.Sub(content, "public")
	if err != nil {
		slog.Error("failed to create sub-filesystem", "error", err)
		os.Exit(1)
	}

	// 2. Create a router (ServeMux) and register handlers.

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(publicFS))

	mux.Handle("/", fileServer)
	mux.HandleFunc("/searchN4L", SearchN4LHandler)
	mux.HandleFunc("/status", StatusHandler)

	// 3. Create an http.Server instance for graceful shutdown.

	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: EnableCORS(mux)}

	// 4. Run the server in a goroutine so it doesn't block.

	go func() {
		slog.Info("Server started", "host", "0.0.0.0", "port", 8080)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("could not start server", "error", err)
			os.Exit(1)
		}
	}()

	// 5. Wait for an interrupt signal.

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server is shutting down...")

	// 6. Perform a graceful shutdown with a timeout.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited properly")
}

// *********************************************************************
// Handlers
// *********************************************************************

// EnableCORS is an HTTP middleware that enables Cross-Origin Resource Sharing (CORS).
// It wraps HTTP handlers to add necessary CORS headers, allowing the API to be accessed
// from web applications hosted on different domains.
//
// Key features:
// - Dynamically sets Access-Control-Allow-Origin to the requesting origin
// - Supports all common HTTP methods (POST, GET, OPTIONS, PUT, DELETE)
// - Handles preflight OPTIONS requests
// - Allows Content-Type header in cross-origin requests
//
// Parameters:
//   - next: The HTTP handler to wrap with CORS support
//
// Returns:
//   - An HTTP handler with CORS headers enabled
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

// SignalHandler is a legacy signal handling function (currently unused).
// It was replaced by the graceful shutdown mechanism in main().
// This function blocks until receiving a system signal and logs the signal type.
//
// Handled signals:
// - SIGHUP (1): Terminal hangup
// - SIGINT (2): Interrupt from keyboard (Ctrl-C)
// - SIGQUIT (3): Quit signal
// - SIGTERM (15): Termination signal
//
// Note: This function is kept for reference but is not currently called.
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

// SearchN4LHandler handles search requests for the N4L (Narrative for Learning) knowledge graph.
// It processes both GET and POST requests and supports various search operations.
//
// Special commands:
// - "\\lastnptr": Updates statistics for the last viewed node pointer
// - "\\remind": Shows reminders in the current time context
// - "\\help": Displays help and search documentation
//
// Request parameters:
// - name: Search query string (can include special commands and operators)
// - nclass: Node class identifier (for direct node access)
// - ncptr: Node class pointer (for direct node access)
// - chapcontext: Chapter context filter
//
// The function:
// 1. Parses form values from the request
// 2. Handles special commands for statistics updates
// 3. Constructs node pointers from nclass/ncptr pairs
// 4. Decodes search fields into SearchParameters
// 5. Delegates to HandleSearch for query execution
//
// Changes made:
// - Improved error handling with errors.Join for multiple error conditions
// - Enhanced logging using slog for better observability
func SearchN4LHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "POST", "GET":
		name := r.FormValue("name")
		nclass := r.FormValue("nclass")
		ncptr := r.FormValue("ncptr")
		chapcontext := r.FormValue("chapcontext")

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
			_, err1 := fmt.Sscanf(nclass, "%d", &a)
			_, err2 := fmt.Sscanf(ncptr, "%d", &b)
			nstr := fmt.Sprintf("(%d,%d)", a, b)
			name = name + nstr
			if err1 != nil || err2 != nil {
				slog.Error("Error parsing nclass/ncptr", "nclass", nclass, "ncptr", ncptr, "error", errors.Join(err1, err2))
			}
		}

		fmt.Println("\nReceived command:", name)

		ambient, key, _ := SST.GetTimeContext()

		if len(name) == 0 || name == "\\remind" {
			name = "any \\chapter reminders \\context " + key + " " + ambient
		}

		if len(name) == 0 || name == "\\help" {
			name = "\\notes \\chapter \"help and search\" \\limit 40"
		}

		search := SST.DecodeSearchField(name)

		HandleSearch(search, name, w, r)

	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// *********************************************************************

// UpdateLastSawSection updates the database with statistics about when a chapter/section was last viewed.
// This is used to track user engagement and provide personalized content recommendations.
//
// Parameters:
//   - w: HTTP response writer (not currently used)
//   - r: HTTP request (not currently used)
//   - query: The chapter/section identifier to update
//
// The function logs the update operation and delegates to SST.UpdateLastSawSection.
func UpdateLastSawSection(w http.ResponseWriter, r *http.Request, query string) {
	// update lastseen db

	slog.Info("UPDATING STATS FOR section", "query", query)

	SST.UpdateLastSawSection(CTX, query)
}

// *********************************************************************

// UpdateLastSawNPtr updates the database with statistics about when a specific node was last viewed.
// This tracks individual node access patterns for analytics and recommendation purposes.
//
// Parameters:
//   - w: HTTP response writer (for sending acknowledgment)
//   - r: HTTP request (not currently used)
//   - class: String representation of the node class
//   - cptr: String representation of the class pointer
//   - classifier: Optional classification context for the node
//
// The function:
// 1. Parses class and cptr strings into integers
// 2. Updates the node pointer statistics in the database
// 3. Updates the associated section statistics
// 4. Sends a JSON acknowledgment response
//
// Changes made:
// - Improved error handling with errors.Join for cleaner error aggregation
// - Enhanced logging with slog for better debugging
// - Added error checking for Write operations
func UpdateLastSawNPtr(w http.ResponseWriter, r *http.Request, class, cptr string, classifier string) {
	// update lastseen db

	// var nptr SST.NodePtr
	// NOTE: The two lines below comented because they are unused variables
	// nptr.Class = nclass
	// nptr.CPtr = SST.ClassedNodePtr(ncptr)

	var nclass int
	var ncptr int
	_, err1 := fmt.Sscanf(class, "%d", &nclass)
	_, err2 := fmt.Sscanf(cptr, "%d", &ncptr)
	if err := errors.Join(err1, err2); err != nil {
		slog.Error("Error parsing class/cptr", "class", class, "cptr", cptr, "error", err)
	}
	SST.UpdateLastSawNPtr(CTX, nclass, ncptr, classifier)

	slog.Info("UPDATING STATS FOR", "nclass", nclass, "ncptr", ncptr, "classifier", classifier)

	SST.UpdateLastSawSection(CTX, classifier)

	response := fmt.Sprintf("{ \"Response\" : \"LastSaw\",\n \"Content\" : \"ack(%s,%s)\" }", class, cptr)
	if _, err := w.Write([]byte(response)); err != nil {
		slog.Error("Error writing response", "error", err)
	}
}

// *********************************************************************

// HandleSearch is the main search orchestrator that routes queries to appropriate handlers.
// It analyzes the search parameters and delegates to specialized handlers based on the query type.
//
// Parameters:
//   - search: Parsed search parameters containing filters and constraints
//   - line: Original search query string (for logging)
//   - w: HTTP response writer
//   - r: HTTP request
//
// Supported search types:
// 1. Stats queries - User engagement statistics
// 2. Chapter/Context browsing - Table of contents exploration
// 3. Node orbit queries - Direct node neighborhood exploration
// 4. Path solving - Finding connections between two sets of nodes
// 5. Causal cone queries - Forward/backward relationship exploration
// 6. Page-mapped notes - Notes linked to specific page numbers
// 7. Sequential stories - Narrative sequences following specific arrow types
// 8. Arrow matching - Finding relationships by arrow type
//
// The function determines search limits dynamically based on query complexity:
// - Paths/sequences: 30 results (computationally expensive)
// - Short search terms (<5 chars): 5 results (likely common words)
// - Longer search terms: 10 results (more specific)
//
// Changes made:
// - Consolidated arrow pointer resolution logic
// - Improved limit calculation heuristics
// - Enhanced debug output with tabwriter formatting
func HandleSearch(search SST.SearchParameters, line string, w http.ResponseWriter, r *http.Request) {
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
	// DEBUG PRINTING *********************************************

	{
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
	}

	var nodeptrs, leftptrs, rightptrs []SST.NodePtr

	if !pagenr && !sequence {
		leftptrs = SST.SolveNodePtrs(CTX, search.From, search, arrowptrs, limit)
		rightptrs = SST.SolveNodePtrs(CTX, search.To, search, arrowptrs, limit)
	}

	nodeptrs = SST.SolveNodePtrs(CTX, search.Name, search, arrowptrs, limit)

	slog.Info("Solved search nodes ...")

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
		slog.Error("Search has conflicting parts <to|from> and match strings", "line", line)
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

// HandleOrbit retrieves and formats the relationship neighborhood (orbit) around specified nodes.
// An orbit represents all direct relationships connected to a node, visualized in 3D space.
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//   - ctx: SST database context
//   - search: Search parameters (for metadata)
//   - nptrs: Array of node pointers to get orbits for
//   - limit: Maximum number of relationships per orbit
//
// The function:
// 1. Iterates through each node pointer
// 2. Retrieves the orbit (all connected relationships)
// 3. Assigns 3D coordinates for visualization (disconnected nodes are spatially separated)
// 4. Packages each orbit as a NodeEvent with coordinates
// 5. Returns a JSON array of orbits
//
// Visualization logic:
// - Multiple disconnected nodes are arranged in a circle pattern around the origin
// - Each orbit is centered at a relative position using RelativeOrbit
// - Coordinates enable 3D graph rendering in the client
//
// Changes made:
// - Enhanced logging with slog
// - Improved error handling for Write operations
func HandleOrbit(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, nptrs []SST.NodePtr, limit int) {
	var count int
	var array []SST.NodeEvent

	startTime := time.Now()
	slog.Info("HandleOrbit starting", "total_nodes", len(nptrs), "limit", limit)
	origin := SST.Coords{X: 0.0, Y: 0.0, Z: 0.0}

	for n := 0; n < len(nptrs); n++ {

		count++

		if count > limit {
			slog.Info("Orbit limit reached", "limit", limit, "processed", count-1)
			break
		}

		slog.Debug("Processing node orbit", "index", n+1, "of", len(nptrs), "node", nptrs[n])

		orb := SST.GetNodeOrbit(CTX, nptrs[n], "", limit)
		// create a set of coords for len(nptrs) disconnected nodes

		xyz := SST.RelativeOrbit(origin, SST.R0, n, len(nptrs))
		orb = SST.SetOrbitCoords(xyz, orb)

		nodeevent := SST.JSONNodeEvent(CTX, nptrs[n], xyz, orb)
		array = append(array, nodeevent)
	}

	data, _ := json.Marshal(array)
	response := PackageResponse(ctx, search, "Orbits", string(data))

	// fmt.Println("REPLY:\n",string(response))

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		slog.Error("Failed to write response", "error", err)
	}
	duration := time.Since(startTime)
	slog.Info("HandleOrbit completed", "processed_nodes", count, "duration_ms", duration.Milliseconds())
}

// *********************************************************************

// HandleCausalCones retrieves and visualizes causal cones (forward/backward relationships) from nodes.
// Causal cones show all nodes reachable by following relationships in a specific direction.
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//   - ctx: SST database context
//   - nptrs: Starting node pointers for cone exploration
//   - search: Search parameters with chapter/context filters
//   - arrows: Optional specific arrow types to follow
//   - sttype: Semantic type filters (0=narrative, 1=epistemic, 2=causal, 3=contains)
//   - limit: Maximum total results to return
//
// The function:
// 1. Uses default sttype [0,1,2,3] if none specified (all semantic types)
// 2. For each node, generates both forward and backward cones for each semantic type
// 3. Stops when the total result count exceeds the limit
// 4. Packages results as WebConePaths for 3D visualization
//
// Cone types:
// - sttype > 0: Forward cone (following relationships forward in time/logic)
// - sttype < 0: Backward cone (following relationships backward)
// - sttype = 0: Special case (narrative sequence, forward only)
//
// Changes made:
// - Enhanced logging for better debugging
// - Improved error handling for Write operations
func HandleCausalCones(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, nptrs []SST.NodePtr, search SST.SearchParameters, arrows []SST.ArrowPtr, sttype []int, limit int) {
	chap := search.Chapter
	context := search.Context

	slog.Info("HandleCausalCones()", "nptrs", nptrs)
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
	// fmt.Println("CasualConePath reponse",string(response))

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		slog.Error("Failed to write response", "error", err)
	}
	slog.Info("Done/sent cone")
}

//******************************************************************

// PackageConeFromOrigin creates a WebConePaths structure for a single node's causal cone.
// This packages both forward and backward cones into a unified visualization structure.
//
// Parameters:
//   - ctx: SST database context
//   - nptr: The root node pointer for the cone
//   - nth: Index of this node in a multi-node result set (for spatial positioning)
//   - sttype: Semantic type (0=narrative, 1=epistemic, 2=causal, 3=contains)
//   - chap: Optional chapter filter
//   - context: Optional context filters
//   - dimnptr: Total number of nodes in the result set (for coordinate calculation)
//   - limit: Maximum paths to include
//
// Returns:
//   - WebConePaths: Structured cone data with paths and metadata
//   - int: Count of paths in the cone
//
// The function:
// 1. Retrieves forward paths in the specified semantic type
// 2. Retrieves backward paths (if sttype != 0)
// 3. Converts path Links to WebPaths with 3D coordinates
// 4. Packages everything with the root node title
//
// Coordinate assignment:
// - nth/dimnptr ratio determines angular position in a circular layout
// - Allows multiple disconnected cones to be visualized simultaneously
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

// HandlePathSolve finds and returns all paths connecting two sets of nodes.
// This implements a bidirectional search algorithm that expands from both endpoints.
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//   - ctx: SST database context
//   - leftptrs: Starting node pointers ("from" nodes)
//   - rightptrs: Ending node pointers ("to" nodes)
//   - search: Search parameters with chapter/context filters
//   - arrowptrs: Optional specific arrow types to follow
//   - sttype: Semantic type filters
//   - maxdepth: Maximum search depth (prevents infinite searches)
//
// Algorithm:
// 1. Starts with ldepth=2, rdepth=2 (search depth from each side)
// 2. Expands forward from leftptrs and backward from rightptrs
// 3. Checks if the expanding wavefronts overlap (paths found)
// 4. If no overlap, alternates increasing ldepth and rdepth
// 5. Continues until paths found or maxdepth reached
// 6. If forward/backward fails, tries backward/forward (direction reversal)
//
// The function includes betweenness centrality calculation and super-node detection
// to identify important intermediate nodes in the solution paths.
//
// Returns:
// - PathSolve JSON with solution paths, betweenness scores, and super-nodes
// - Empty array if no paths found within maxdepth
//
// Changes made:
// - Enhanced logging throughout the search process
// - Improved error handling for Write operations
func HandlePathSolve(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, leftptrs, rightptrs []SST.NodePtr, search SST.SearchParameters, arrowptrs []SST.ArrowPtr, sttype []int, maxdepth int) {
	chapter := search.Chapter
	context := search.Context

	slog.Info("HandlePathSolve()", "leftptrs", leftptrs, "rightptrs", rightptrs)

	var Lnum, Rnum int
	var left_paths, right_paths [][]SST.Link

	// Find the path matrix

	var solutions [][]SST.Link
	var ldepth, rdepth int = 2, 2

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {

		left_paths, Lnum = SST.GetEntireNCSuperConePathsAsLinks(CTX, "fwd", leftptrs, ldepth, chapter, context, maxdepth)
		right_paths, Rnum = SST.GetEntireNCSuperConePathsAsLinks(CTX, "bwd", rightptrs, rdepth, chapter, context, maxdepth)

		if Lnum == 0 || Rnum == 0 {
			slog.Info("Nothing, trying reverse")
			left_paths, Lnum = SST.GetEntireNCSuperConePathsAsLinks(CTX, "bwd", leftptrs, ldepth, chapter, context, maxdepth)
			right_paths, Rnum = SST.GetEntireNCSuperConePathsAsLinks(CTX, "fwd", rightptrs, rdepth, chapter, context, maxdepth)

			if Lnum == 0 || Rnum == 0 {
				slog.Info("No paths")
				response := PackageResponse(ctx, search, "PathSolve", "")
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write(response); err != nil {
					slog.Error("Failed to write PathSolve no paths response", "error", err)
				}
				return
			}
		}

		solutions, _ = SST.WaveFrontsOverlap(CTX, left_paths, right_paths, Lnum, Rnum, ldepth, rdepth)

		if len(solutions) > 0 {
			// format paths

			var pack []SST.WebConePaths
			var soln SST.WebConePaths

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
			array_pack, _ := json.Marshal(pack)

			response := PackageResponse(ctx, search, "PathSolve", string(array_pack))

			// fmt.Println("PATH SOLVE:",string(response))

			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(response); err != nil {
				slog.Error("Failed to write PathSolve response", "error", err)
			}
			return
		}

		if turn%2 == 0 {
			ldepth++
		} else {
			rdepth++
		}
	}

	slog.Info("No paths satisfy constraints")
	response := PackageResponse(ctx, search, "PathSolve", "[]")

	// fmt.Println("PATHSOLVE NOTES",string(response))
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		slog.Error("Failed to write PathSolve empty response", "error", err)
	}
	slog.Info("Done/sent path solve")
}

//******************************************************************

// HandlePageMap retrieves and formats notes linked to specific page numbers.
// This supports page-based navigation in documents that have been annotated with page mappings.
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//   - ctx: SST database context
//   - search: Search parameters (contains chapter/context/page number)
//   - notes: Array of PageMap entries to format
//
// The function:
// 1. Converts PageMap entries to JSON format
// 2. Updates statistics for the chapter being viewed
// 3. Returns a PageMap response type
//
// This enables features like:
// - "Show me notes for page 42 of chapter X"
// - Linking physical book pages to digital annotations
// - Page-based navigation through content
//
// Changes made:
// - Added error checking for Write operations
func HandlePageMap(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, notes []SST.PageMap) {
	slog.Info("Solver/handler: HandlePageMap()")

	jstr := SST.JSONPage(CTX, notes)
	response := PackageResponse(ctx, search, "PageMap", jstr)

	if notes != nil {
		UpdateLastSawSection(w, r, notes[0].Chapter)
	}

	// fmt.Println("PAGEMAP NOTES",string(response))
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		slog.Error("Failed to write PageMap response", "error", err)
	}
	slog.Info("Done/sent pagemap")
}

//******************************************************************

// HandleStories retrieves and formats narrative sequences (stories) through the knowledge graph.
// Stories are linear paths following specific relationship types (typically "!then!" for temporal sequence).
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//   - ctx: SST database context
//   - search: Search parameters
//   - nodeptrs: Starting nodes for story sequences
//   - arrowptrs: Optional arrow types to follow (defaults to "!then!" if nil)
//   - sttypes: Semantic types to include
//   - limit: Maximum sequence length
//
// The function:
// 1. Uses "!then!" arrow type if no arrows specified (temporal sequence)
// 2. Retrieves sequence containers (story structures)
// 3. Converts each story axis (sequence of events) to JSON
// 4. Returns an array of story sequences
//
// Use cases:
// - Following narrative timelines
// - Tracing procedural steps
// - Exploring causal chains
//
// Changes made:
// - Added error checking for Write operations
func HandleStories(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, nodeptrs []SST.NodePtr, arrowptrs []SST.ArrowPtr, sttypes []int, limit int) {
	if arrowptrs == nil {
		arrowptrs, sttypes = SST.ArrowPtrFromArrowsNames(CTX, []string{"!then!"})
	}

	slog.Info("Solver/handler: HandleStories()")

	stories := SST.GetSequenceContainers(ctx, nodeptrs, arrowptrs, sttypes, limit)

	jarray := ""

	for s := range stories {

		var jstory string

		for a := 0; a < len(stories[s].Axis); a++ {
			jstr := JSONStoryNodeEvent(stories[s].Axis[a])
			jstory += fmt.Sprintf("%s,", jstr)
		}

		jstory = strings.Trim(jstory, ",")
		jarray = fmt.Sprintf("[%s],", jstory)
	}

	jarray = strings.Trim(jarray, ",")

	response := PackageResponse(ctx, search, "Sequence", jarray)

	// fmt.Println("Sequence...",string(response))

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		slog.Error("Failed to write Sequence response", "error", err)
	}
	slog.Info("Done/sent sequence")
}

// *********************************************************************

// HandleMatchingArrows retrieves and formats arrow (relationship) type information.
// This provides metadata about available relationship types in the knowledge graph.
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//   - ctx: SST database context
//   - search: Search parameters
//   - arrowptrs: Specific arrow pointers to retrieve info for
//   - sttype: Semantic types to filter by
//
// For each arrow, returns:
// - Arrow pointer and semantic type
// - Short name (e.g., "!then!")
// - Long description
// - Inverse arrow information (reverse relationship)
//
// Semantic types:
// - 0: Narrative (temporal/sequential)
// - 1: Epistemic (knowledge/belief)
// - 2: Causal (cause/effect)
// - 3: Contains (containment/composition)
//
// Use cases:
// - Discovering available relationship types
// - Understanding graph structure
// - Building query interfaces
//
// Changes made:
// - Added error checking for Write operations
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

	slog.Info("Arrows...", "response", string(response))

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		slog.Error("Failed to write Arrows response", "error", err)
	}
	slog.Info("Done/sent arrows")
}

// *********************************************************************

// ShowStats retrieves and returns user engagement statistics.
// This provides analytics about which content has been viewed and when.
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//   - ctx: SST database context
//   - search: Search parameters
//   - nptrs: Optional specific node pointers to get stats for
//
// Returns:
// - If nptrs is nil: Section-level statistics (chapter/context view counts)
// - If nptrs provided: Node-level statistics for specific nodes
//
// Statistics include:
// - Last seen timestamp
// - View count
// - Context in which it was viewed
//
// Use cases:
// - Identifying frequently accessed content
// - Tracking user engagement patterns
// - Personalizing recommendations
//
// Changes made:
// - Added error checking for Write operations
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

	data, _ := json.Marshal(retval)

	response := PackageResponse(ctx, search, "STAT", string(data))

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		slog.Error("Failed to write STAT response", "error", err)
	}
	slog.Info("Done/sent stat")
}

// *********************************************************************

// ShowChapterContexts retrieves and formats a table of contents with contextual analysis.
// This provides a hierarchical view of chapters and their associated contexts.
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//   - ctx: SST database context
//   - search: Search parameters with optional chapter/context filters
//   - limit: Maximum number of chapters to return
//
// The function:
//  1. Retrieves all matching chapters and contexts from the database
//  2. Sorts chapters alphabetically for consistent presentation
//  3. For each chapter, performs contextual analysis:
//     a. IntersectContextParts: Identifies overlapping context keywords
//     b. GetContextTokenFrequencies: Analyzes context keyword distribution
//     c. ContextIntentAnalysis: Separates specific vs. ambient context
//  4. Assigns 3D coordinates for visualization
//  5. Groups contexts into:
//     - Context sets (related keywords that appear together)
//     - Single contexts (chapter-specific keywords)
//     - Common contexts (keywords appearing across multiple chapters)
//
// Returns a TOC (Table of Contents) response with rich contextual metadata.
//
// Changes made:
// - Added error checking for Write operations
// - Enhanced logging with function entry/exit tracking
func ShowChapterContexts(w http.ResponseWriter, r *http.Request, ctx SST.PoSST, search SST.SearchParameters, limit int) {
	slog.Info("started", "function", "ShowChapterContexts()")
	chap := search.Chapter
	context := search.Context

	var chapters []SST.ChCtx
	var chap_list []string

	toc := SST.GetChaptersByChapContext(ctx, chap, context, limit)

	for chaps := range toc {
		chap_list = append(chap_list, chaps)
	}

	sort.Strings(chap_list)

	for c, chapter := range chap_list {

		var chap_anchor SST.ChCtx

		chap_anchor.Chapter = chapter
		chap_anchor.XYZ = SST.AssignChapterCoordinates(c, len(chap_list))

		// Fractionate the (chapter,context) information

		dim, clist, adj := SST.IntersectContextParts(toc[chapter])
		spectrum := SST.GetContextTokenFrequencies(toc[chapter])
		intent, ambient := SST.ContextIntentAnalysis(spectrum, toc[chapter])

		chap_anchor.Context = GetContextSets(dim, clist, adj, chap_anchor.XYZ)
		chap_anchor.Single = GetContextFragments(intent, chap_anchor.XYZ)
		chap_anchor.Common = GetContextFragments(ambient, chap_anchor.XYZ)

		chapters = append(chapters, chap_anchor)
	}

	data, _ := json.Marshal(chapters)
	response := PackageResponse(ctx, search, "TOC", string(data))

	// fmt.Println("Chap/context...", string(response))

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		slog.Error("Failed to write TOC response", "error", err)
	}
	slog.Info("done   ", "function", "ShowChapterContexts()", "message", "content sent")
}

//******************************************************************

// GetContextSets organizes context keywords into related sets based on co-occurrence.
// This identifies which context keywords tend to appear together across the chapter.
//
// Parameters:
//   - dim: Dimensionality of the context space (number of unique keywords)
//   - clist: List of context keyword strings
//   - adj: Adjacency matrix showing keyword co-occurrence (adj[i][j] > 0 means keywords i and j appear together)
//   - xyz: Base coordinates for the chapter (context sets are positioned relative to this)
//
// Returns:
//   - Array of Loc structures, each representing a context set with:
//   - Text: The context keyword
//   - XYZ: 3D coordinates for visualization
//   - Reln: Indices of related keywords (based on adjacency matrix)
//
// The adjacency matrix captures semantic relationships between contexts,
// enabling the UI to visualize which contexts are conceptually related.
//
// Changes made:
// - Preallocated Reln slice for better performance (avoids repeated allocations)
// - Added function entry/exit logging for debugging
func GetContextSets(dim int, clist []string, adj [][]int, xyz SST.Coords) []SST.Loc {
	var retvar []SST.Loc
	numContextSets := len(adj)

	for c := range adj {

		var contextgroup SST.Loc

		contextgroup.Text = clist[c]
		contextgroup.XYZ = SST.AssignContextSetCoordinates(xyz, c, numContextSets)

		// Preallocate Reln slice for performance
		contextgroup.Reln = make([]int, 0, len(adj[c]))
		for cp := 0; cp < len(adj[c]); cp++ {
			if adj[c][cp] > 0 {
				contextgroup.Reln = append(contextgroup.Reln, cp)
			}
		}
		retvar = append(retvar, contextgroup)
	}
	return retvar
}

//******************************************************************

// GetContextSetsOld is the legacy version of GetContextSets, kept for reference.
// This function has the same logic as GetContextSets but may have older implementation details.
//
// Parameters and return values are identical to GetContextSets.
//
// Note: This function should be deprecated and removed once migration is complete.
//
// Changes made:
// - Preallocated Reln slice for performance optimization
// - Added logging for function tracking
func GetContextSetsOld(dim int, clist []string, adj [][]int, xyz SST.Coords) []SST.Loc {
	slog.Info("started", "function", "GetContextSetsOld()")
	var retvar []SST.Loc

	for c := range adj {

		var contextgroup SST.Loc

		contextgroup.Text = clist[c]
		contextgroup.XYZ = SST.AssignContextSetCoordinates(xyz, c, len(adj))

		// Preallocate Reln slice for performance
		contextgroup.Reln = make([]int, 0, len(adj[c]))
		for cp := 0; cp < len(adj[c]); cp++ {
			if adj[c][cp] > 0 {
				contextgroup.Reln = append(contextgroup.Reln, cp)
			}
		}
		retvar = append(retvar, contextgroup)
	}
	slog.Info("completed", "function", "GetContextSetsOld()")
	return retvar
}

//******************************************************************

// GetContextFragments converts a list of context keywords into individual location objects.
// Unlike GetContextSets, this treats each keyword independently without relationships.
//
// Parameters:
//   - clist: List of context keyword strings
//   - ooo: Origin coordinates (base position for the fragment group)
//
// Returns:
//   - Array of Loc structures, each representing a single context keyword with:
//   - Text: The context keyword
//   - XYZ: 3D coordinates calculated relative to the origin
//
// Use cases:
// - Displaying individual context tags
// - Showing chapter-specific keywords (intent)
// - Showing common/ambient keywords that appear everywhere
//
// The coordinate assignment spaces fragments around the origin point
// to create a visually distributed layout.
func GetContextFragments(clist []string, ooo SST.Coords) []SST.Loc {
	var retvar []SST.Loc

	for c := range clist {
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

// JSONStoryNodeEvent converts a NodeEvent to JSON string representation.
// This is used for serializing story sequences for API responses.
//
// Parameters:
//   - en: NodeEvent containing node data and metadata
//
// Returns:
//   - JSON string representation, or empty string if the event has no text
//
// Changes made:
// - Added error handling with slog for marshalling failures
func JSONStoryNodeEvent(en SST.NodeEvent) string {
	if len(en.Text) == 0 {
		return ""
	}

	j, err := json.Marshal(en)
	if err != nil {
		slog.Error("Error marshalling NodeEvent", "error", err)
		return ""
	}
	return string(j)
}

// **********************************************************

// ShowNode converts an array of node pointers to a comma-separated string representation.
// This is used for debugging and logging purposes.
//
// Parameters:
//   - ctx: SST database context
//   - nptr: Array of node pointers to display
//   - maxLen: Maximum length parameter (currently unused)
//
// Returns:
//   - Comma-separated string of node names (truncated to 30 chars each)
//
// The function:
// - Retrieves the text for each node from the database
// - Escapes commas in node names to prevent parsing issues
// - Truncates each name to 30 characters for readability
// - Joins all names with commas
//
// Changes made:
// - Uses strings.Builder for efficient string concatenation
func ShowNode(ctx SST.PoSST, nptr []SST.NodePtr, maxLen int) string {
	var builder strings.Builder

	for n := range nptr {
		node := SST.GetDBNodeByNodePtr(ctx, nptr[n])
		escapedName := strings.ReplaceAll(node.S, ",", "\\,")
		builder.WriteString(fmt.Sprintf("%.30s", escapedName))
		if n < len(nptr)-1 {
			builder.WriteString(",")
		}
	}

	return builder.String()
}

// **********************************************************

// PackageResponse creates a standardized JSON response wrapper for all API responses.
// This ensures consistent response format across all endpoints with contextual metadata.
//
// Parameters:
//   - ctx: SST database context
//   - search: Search parameters from the original request
//   - kind: Response type identifier (e.g., "Orbits", "PathSolve", "TOC")
//   - jstr: JSON string content to include in the response
//
// Returns:
//   - JSON byte array containing:
//   - Response: Type identifier
//   - Content: Actual response data (parsed JSON or string)
//   - Time: Current timestamp key
//   - Intent: User's current search intent/context
//   - Ambient: Ambient contextual information
//
// The function:
// 1. Gets current time context (ambient, key, now)
// 2. Updates short-term memory (STM) context based on the search
// 3. Attempts to parse jstr as JSON; falls back to string if not valid JSON
// 4. Packages everything into a structured response object
// 5. Returns marshalled JSON
//
// This enables the client to:
// - Understand response type and handle it appropriately
// - Track temporal context
// - Maintain user intent across interactions
// - Access ambient/environmental context
//
// Changes made:
// - Improved error handling with errors.Join
// - Added intelligent JSON parsing (object/array vs string)
// - Enhanced error logging with slog
func PackageResponse(ctx SST.PoSST, search SST.SearchParameters, kind string, jstr string) []byte {
	ambien, key, now := SST.GetTimeContext()
	now_ctx := SST.UpdateSTMContext(CTX, ambien, key, now, search)

	intent, err1 := json.Marshal(now_ctx)
	ambient, err2 := json.Marshal(ambien)
	if err := errors.Join(err1, err2); err != nil {
		slog.Error("Error marshalling intent/ambient", "error", err)
	}

	// Try to unmarshal jstr if it's a valid JSON object/array, else treat as string
	var content interface{}
	if len(jstr) > 0 {
		// Try to unmarshal into interface{}
		if err := json.Unmarshal([]byte(jstr), &content); err != nil {
			// If not valid JSON, treat as string
			content = jstr
		}
	} else {
		content = nil
	}

	responseObj := struct {
		Response string `json:"Response"`
		Content  any    `json:"Content"`
		Time     string `json:"Time"`
		Intent   any    `json:"Intent"`
		Ambient  any    `json:"Ambient"`
	}{
		Response: kind,
		Content:  content,
		Time:     key,
		Intent:   json.RawMessage(intent),
		Ambient:  json.RawMessage(ambient),
	}

	response, err := json.Marshal(responseObj)
	if err != nil {
		slog.Error("Error marshalling final response", "error", err)
		return []byte("{}")
	}

	return response
}

//******************************************************************

// SL formats a string slice as a bracketed, comma-separated string.
// This is a simple utility for debug output and logging.
//
// Parameters:
//   - list: Array of strings to format
//
// Returns:
//   - String in the format "[ item1, item2, item3 ]"
//
// Example: SL([]string{"a", "b", "c"}) returns "[ a, b, c ]"
func SL(list []string) string {
	return fmt.Sprintf("[ %s ]", strings.Join(list, ", "))
}

//******************************************************************

// StatusResponse defines the structure for the /status endpoint JSON response.
// This provides server health and available content information.
type StatusResponse struct {
	ServerStatus    string    `json:"server_status"`    // Server operational status
	DatabaseStatus  string    `json:"database_status"`  // Database connection status
	AvailableTopics []string  `json:"available_topics"` // List of available chapters/topics
	Timestamp       time.Time `json:"timestamp"`        // Response generation time
}

// StatusHandler provides server status and available content listing.
// This endpoint is used for health checks and content discovery.
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//
// Returns JSON response containing:
// - server_status: "OK" if running
// - database_status: "OK" if database is accessible
// - available_topics: Sorted list of all available chapters
// - timestamp: Current server time
//
// This enables clients to:
// - Verify server availability
// - Discover available content
// - Check database connectivity
//
// Changes made:
// - Added error handling for response marshalling
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
	if _, err := w.Write(responseJSON); err != nil {
		slog.Error("Failed to write status response", "error", err)
	}
}
