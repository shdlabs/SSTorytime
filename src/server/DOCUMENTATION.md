# HTTP Server Documentation

## Overview

This HTTP server provides a JSON API for querying and exploring a knowledge graph database (SST - Semantic Spacetime). It supports various types of queries including path finding, contextual search, and narrative sequence exploration.

## Architecture

### Server Initialization

- **Port**: 0.0.0.0:8080 (accessible on all network interfaces)
- **Database**: Single persistent SST connection shared across all requests
- **Static Files**: Embedded at compile time from the `public/` directory
- **Logging**: Structured logging using `slog` at INFO level
- **Shutdown**: Graceful shutdown with 5-second timeout on SIGINT/SIGTERM

### Endpoints

1. **GET/POST /searchN4L** - Main search API
2. **GET /status** - Server health and content discovery
3. **GET /** - Static file serving from embedded public directory

## Key Functions

### Core Handlers

#### `main()`

Entry point that:

- Opens SST database connection
- Sets up embedded filesystem for static files
- Configures HTTP routes with CORS support
- Implements graceful shutdown

#### `EnableCORS(next http.Handler) http.Handler`

Middleware that enables Cross-Origin Resource Sharing:

- Dynamically sets origin from request
- Handles preflight OPTIONS requests
- Supports all common HTTP methods

#### `SearchN4LHandler(w, r)`

Main search request handler supporting:

- **Special Commands**:
  - `\lastnptr` - Update view statistics
  - `\remind` - Show contextual reminders
  - `\help` - Display help documentation
- Direct node access via nclass/ncptr parameters
- Complex search query parsing and delegation

### Search Operations

#### `HandleSearch(search, line, w, r)`

Search orchestrator that routes to specialized handlers based on query type:

1. **Stats Queries** → `ShowStats()`
2. **Chapter/Context TOC** → `ShowChapterContexts()`
3. **Node Neighborhoods** → `HandleOrbit()`
4. **Path Finding** → `HandlePathSolve()`
5. **Causal Cones** → `HandleCausalCones()`
6. **Page-Mapped Notes** → `HandlePageMap()`
7. **Story Sequences** → `HandleStories()`
8. **Arrow Metadata** → `HandleMatchingArrows()`

**Dynamic Limits**:

- Paths/sequences: 30 results (computationally expensive)
- Short terms (<5 chars): 5 results (common words)
- Longer terms: 10 results (more specific)

#### `HandleOrbit(w, r, ctx, search, nptrs, limit)`

Retrieves relationship neighborhoods around nodes:

- Gets all direct connections (orbit) for each node
- Assigns 3D coordinates for visualization
- Separates disconnected nodes spatially in circular pattern
- Returns JSON array of NodeEvents with coordinates

#### `HandleCausalCones(w, r, ctx, nptrs, search, arrows, sttype, limit)`

Explores forward/backward relationship chains:

- **Semantic Types**:
  - 0: Narrative (temporal/sequential)
  - 1: Epistemic (knowledge/belief)
  - 2: Causal (cause/effect)
  - 3: Contains (containment/composition)
- Generates both forward and backward cones
- Stops at result limit to prevent overload

#### `HandlePathSolve(w, r, ctx, leftptrs, rightptrs, search, arrowptrs, sttype, maxdepth)`

Bidirectional path finding algorithm:

- Expands from both start (leftptrs) and end (rightptrs) nodes
- Alternates increasing search depth on each side
- Checks for wavefront overlap (paths found)
- Includes betweenness centrality and super-node detection
- Falls back to direction reversal if initial search fails

#### `HandleStories(w, r, ctx, search, nodeptrs, arrowptrs, sttypes, limit)`

Narrative sequence exploration:

- Follows temporal relationships (defaults to "!then!" arrows)
- Retrieves linear story paths through knowledge graph
- Useful for procedural steps and causal chains

### Context and Statistics

#### `ShowChapterContexts(w, r, ctx, search, limit)`

Table of contents with rich contextual analysis:

- Retrieves all chapters and contexts
- Performs contextual analysis:
  - **IntersectContextParts**: Identifies overlapping keywords
  - **GetContextTokenFrequencies**: Analyzes keyword distribution
  - **ContextIntentAnalysis**: Separates specific vs. ambient context
- Groups contexts into sets, single keywords, and common keywords
- Assigns 3D coordinates for visualization

#### `GetContextSets(dim, clist, adj, xyz)`

Organizes context keywords into related sets:

- Uses adjacency matrix to show keyword co-occurrence
- Returns keywords with relationship indices
- Enables semantic relationship visualization
- **Optimization**: Preallocated slices for better performance

#### `GetContextFragments(clist, ooo)`

Individual context keyword locations:

- Treats keywords independently (no relationships)
- Used for chapter-specific and ambient keywords
- Spaces fragments around origin for visual distribution

#### `UpdateLastSawNPtr(w, r, class, cptr, classifier)`

Tracks node view statistics:

- Updates database with last viewed timestamp
- Updates associated section statistics
- Returns acknowledgment response

#### `ShowStats(w, r, ctx, search, nptrs)`

User engagement analytics:

- Section-level stats (if nptrs is nil)
- Node-level stats (if nptrs provided)
- Includes view counts, timestamps, contexts

### Utility Functions

#### `PackageResponse(ctx, search, kind, jstr)`

Standardized response wrapper:

- Consistent format across all endpoints
- Includes contextual metadata:
  - **Response**: Type identifier
  - **Content**: Actual data (auto-parsed JSON or string)
  - **Time**: Current timestamp
  - **Intent**: User's search intent/context
  - **Ambient**: Environmental context
- Enables temporal and contextual tracking

#### `StatusHandler(w, r)`

Health check and content discovery:

- Returns server operational status
- Lists all available chapters/topics (sorted)
- Provides database connectivity status
- Includes timestamp for monitoring

## Key Changes and Improvements

### Error Handling

- **Before**: Many `w.Write()` calls ignored errors
- **After**: All Write operations check errors with slog logging
- Used `errors.Join()` for cleaner error aggregation

### Logging

- **Before**: Basic fmt.Println debugging
- **After**: Structured logging with slog
- Function entry/exit tracking for debugging
- Contextual error information

### Performance

- **Before**: Repeated slice allocations in loops
- **After**: Preallocated slices in `GetContextSets()` and `GetContextSetsOld()`
- String concatenation uses `strings.Builder`

### Code Quality

- Comprehensive documentation for all functions
- Parameter and return value descriptions
- Use case explanations
- Algorithm documentation

### Response Format

- Intelligent JSON parsing in `PackageResponse()`
- Automatically detects JSON objects/arrays vs strings
- Consistent metadata across all responses

## Semantic Types Reference

The knowledge graph uses four semantic relationship types:

- **0 (Narrative)**: Temporal/sequential relationships
- **1 (Epistemic)**: Knowledge/belief relationships
- **2 (Causal)**: Cause/effect relationships
- **3 (Contains)**: Containment/composition relationships

Each type can be traversed forward or backward (negative values).

## Visualization

The server assigns 3D coordinates (XYZ) for graph visualization:

- **Orbits**: Circular arrangement of disconnected nodes
- **Cones**: Radial expansion from root nodes
- **Context Sets**: Spatial grouping of related keywords
- **Stories**: Linear sequences through narrative space

## Legacy Code

Several functions are kept for backward compatibility but marked for deprecation:

- `SignalHandler()` - Replaced by graceful shutdown in main()
- `GenHeader()` - Replaced by EnableCORS middleware
- `GetContextSetsOld()` - Old version of GetContextSets

Note: `CleanText()` was removed as it's better to use `json.Marshal()` for proper JSON escaping.

## API Usage Examples

### Basic Search

```
POST /searchN4L
name=knowledge graph
```

### Path Finding

```
POST /searchN4L
from=concept A&to=concept B
```

### Chapter Context

```
POST /searchN4L
chapter=Introduction&context=database
```

### Story Sequence

```
POST /searchN4L
name=start node&sequence=true
```

### Server Status

```
GET /status
```
