
# graph_report tool

There's a few things one would like to know about pretty much every graph.
the `graph_report` tool offers a simple analysis for determining these basic
features, for each link meta-type, chapter by chapter etc.

* *Loops*: graphs that contain loops (cyclic graphs)

* *Sources* and *Sinks*: these are nodes that start and end a path through the graph.
They exchange places if one changes the sign of the link type.

* *Eigenvector centrality*: undirected graphs have a property by virtue of the 
Frobenious-Perron theorem that every undirected graph has a non-negative principal
eigenvector. It ranks the 'connectedness' of nodes, or their importance, by measuring
the amount of 'weight' propagated to each node. If we think of
each node as being a reservoir of 'weight' and each directed arrow as being a gradient, then
all the weight in a directed graph flows to the sinks immediately, leaving all others empty.
In a symmetrized (undirected graph), the flows reach equilibrium and the highest levels settle
in the nodes that are best connected. The capacitance of these nodes tells us something about
their topological connectivity, setting aside the directedness of the arrows.  The graph_report tool
reports the unbiased vector normalization of the principal eigenvector when removing all arrow
directions but preserving their weights.

For example, the report on the "doors.n4l" graph, for leadsto links:

<pre>

$ go run graph_report.go -chapter multi
----------------------------------------------------------------
Analysing chapter "multi slit interference", context [] to path length 6
----------------------------------------------------------------

* TOTAL NODES IN THE SEARCH REGION 13

* TOTAL DIRECTED LINKS = 17 of possible 156 = 0.11 %

* DISTRIBUTION OF NAME TYPE/LENGTHS:
  - single word ngram : 10 / 13
  - two word ngram : 3 / 13


* PROCESS ORIGINS / ROOT DEPENDENCIES / PATH SOURCES for ("+leads to") in multi slit interference

   - NPtr(1,0) -> start



* FINAL END-STATES / PATH SINK NODES for ("+leads to") in multi slit interference

   - NPtr(2,0) -> target 1
   - NPtr(2,1) -> target 2
   - NPtr(2,2) -> target 3

* DIRECTED LOOPS AND CYCLES:

   - Acyclic

* SYMMETRIZED EIGENVECTOR CENTRALITY = FLOW RESERVOIR CAPACITANCE AT EQUILIBRIUM = 

   ( 0.993 ) <- 0 = tram
   ( 0.768 ) <- 1 = target 3
   ( 0.847 ) <- 2 = gate
   ( 0.496 ) <- 3 = bike
   ( 1.000 ) <- 4 = start
   ( 0.787 ) <- 5 = door
   ( 0.940 ) <- 6 = port
   ( 0.678 ) <- 7 = hole
   ( 0.767 ) <- 8 = river
   ( 0.272 ) <- 9 = passage
   ( 0.093 ) <- 10 = target 1
   ( 0.272 ) <- 11 = road
   ( 0.093 ) <- 12 = target 2


* THERE ARE 2 LOCAL MAXIMA IN THE EQUILIBRIUM EVC LANDSCAPE:

  - subregion of maximum 0 consisting of nodes [0 1]
     - where 0 -> tram
     - where 1 -> target 3
  - subregion of maximum 4 consisting of nodes [2 3 4 5 6 7 8 9 10 11 12]
     - where 2 -> gate
     - where 3 -> bike
     - where 4 -> start
     - where 5 -> door
     - where 6 -> port
     - where 7 -> hole
     - where 8 -> river
     - where 9 -> passage
     - where 10 -> target 1
     - where 11 -> road
     - where 12 -> target 2

* HILL-CLIMBING EVC-LAMDSCAPE GRADIENT PATHS:

     - Path node 0 has local maximum at node * 0 *, hop distance 0 along [0]
     - Path node 1 has local maximum at node * 0 *, hop distance 1 along [1 0]
     - Path node 2 has local maximum at node * 4 *, hop distance 1 along [2 4]
     - Path node 3 has local maximum at node * 4 *, hop distance 2 along [3 2 4]
     - Path node 4 has local maximum at node * 4 *, hop distance 0 along [4]
     - Path node 5 has local maximum at node * 4 *, hop distance 1 along [5 4]
     - Path node 6 has local maximum at node * 4 *, hop distance 1 along [6 4]
     - Path node 7 has local maximum at node * 4 *, hop distance 1 along [7 4]
     - Path node 8 has local maximum at node * 4 *, hop distance 2 along [8 6 4]
     - Path node 9 has local maximum at node * 4 *, hop distance 2 along [9 5 4]
     - Path node 10 has local maximum at node * 4 *, hop distance 3 along [10 9 5 4]
     - Path node 11 has local maximum at node * 4 *, hop distance 2 along [11 5 4]
     - Path node 12 has local maximum at node * 4 *, hop distance 3 along [12 11 5 4]

</pre>