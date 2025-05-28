
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
their topological connectivity, setting aside the directedness of the arrows. 

For example, the report on the "LoopyLoo.n4l" graph, for leadsto links:

<pre>

$ go run graph_report.go -chapter "loop" -sttype=L 
----------------------------------------------------------------
Analysing chapter "loop test", context [] to path length 6
----------------------------------------------------------------

* PROCESS ORIGINS / ROOT DEPENDENCIES / PATH SOURCES for ("+leads to") in loop test

   - NPtr(1,3161) -> L0

* FINAL END-STATES / PATH SINK NODES for ("+leads to") in loop test

   - NPtr(1,3170) -> L9

* DIRECTED LOOP SEARCH:

  Cycle of length 3 with members (1)(4)(5)
   - where 1 -> L1
   - where 4 -> L2
   - where 5 -> L3
  Cycle of length 3 with members (2)(3)(7)
   - where 2 -> L6
   - where 3 -> L7
   - where 7 -> L5
  Cycle of length 4 with members (2)(3)(7)(8)
   - where 2 -> L6
   - where 3 -> L7
   - where 7 -> L5
   - where 8 -> L8
  Cycle of length 3 with members (2)(3)(7)
   - where 2 -> L6
   - where 3 -> L7
   - where 7 -> L5
  Cycle of length 3 with members (1)(4)(5)
   - where 1 -> L1
   - where 4 -> L2
   - where 5 -> L3

* Symmetrized Eigenvector Centrality = FLOW RESERVOIR CAPACITANCE AT EQUILIBRIUM = 

   ( 0.09 ) <- 0 = L0
   ( 0.25 ) <- 1 = L1
   ( 0.70 ) <- 2 = L6
   ( 0.96 ) <- 3 = L7
   ( 0.23 ) <- 4 = L2
   ( 0.35 ) <- 5 = L3
   ( 0.49 ) <- 6 = L4
   ( 1.00 ) <- 7 = L5
   ( 0.70 ) <- 8 = L8
   ( 0.34 ) <- 9 = L9

At directionless equilibrium, there are 10 local maxima in the EVC landscape:

  From node 0 has local maximum at node * 7 *, hop distance 4 along [0 1 5 6 7]
   - where 0 -> L0
   - where 1 -> L1
   - where 5 -> L3
   - where 6 -> L4
   - where 7 -> L5

  From node 1 has local maximum at node * 7 *, hop distance 3 along [1 5 6 7]
   - where 1 -> L1
   - where 5 -> L3
   - where 6 -> L4
   - where 7 -> L5

  From node 2 has local maximum at node * 7 *, hop distance 1 along [2 7]
   - where 2 -> L6
   - where 7 -> L5

  From node 3 has local maximum at node * 7 *, hop distance 1 along [3 7]
   - where 3 -> L7
   - where 7 -> L5

  From node 4 has local maximum at node * 7 *, hop distance 3 along [4 5 6 7]
   - where 4 -> L2
   - where 5 -> L3
   - where 6 -> L4
   - where 7 -> L5

  From node 5 has local maximum at node * 7 *, hop distance 2 along [5 6 7]
   - where 5 -> L3
   - where 6 -> L4
   - where 7 -> L5

  From node 6 has local maximum at node * 7 *, hop distance 1 along [6 7]
   - where 6 -> L4
   - where 7 -> L5

  From node 7 has local maximum at node * 7 *, hop distance 0 along [7]
   - where 7 -> L5

  From node 8 has local maximum at node * 7 *, hop distance 1 along [8 7]
   - where 8 -> L8
   - where 7 -> L5

  From node 9 has local maximum at node * 7 *, hop distance 2 along [9 3 7]
   - where 9 -> L9
   - where 3 -> L7
   - where 7 -> L5

</pre>