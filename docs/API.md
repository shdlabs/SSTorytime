

# An API for interacting with the SST graph

Once data have been entered into a SSToryline database, we want to be able to extract it again.
It's possible to create tools for this, but ultimately any set of tools will tend to limit the user.
A user's imagination should be the only limit. 

Many specialized graph databases offer graph languages, but they
expose an important problem with Domain Specific Languages, which is
that by trying to make simple things easy, they make less-simple
things hard. The most well known standard for data (Structured Query
Language, or SQL) is itself a Domain Specific Language with exactly
these problems. However, in Open Source Postgres there are plenty of
extensions that make it possible to overcome the limitations of SQL.

*This project uses Postgres because of that compromise between a well known
standard, and a battle-tested and extensible data platform.*

You will find examples of using Go(lang) code to write custom scripts
that interact with the database through the Go API [here](https://github.com/markburgess/SSTorytime/tree/main/src/demo_pocs).

## Basic queries from SQL

Using perfectly standard SQL, you can interrogate the database established by N4L or the low level API
functions.

### Tables

* To show the different tables:
<pre>
$ psql newdb

newdb=# \dt
              List of relations
 Schema |      Name      | Type  |   Owner    
--------+----------------+-------+------------
 public | arrowdirectory | table | sstoryline
 public | arrowinverses  | table | sstoryline
 public | node           | table | sstoryline
 public | nodearrownode  | table | sstoryline
(4 rows)

</pre>
* To query these, we look at the members:
<pre>
newdb=# \d node
                Table "public.node"
 Column |  Type   | Collation | Nullable | Default 
--------+---------+-----------+----------+---------
 nptr   | nodeptr |           |          | 
 l      | integer |           |          | 
 s      | text    |           |          | 
 chap   | text    |           |          | 
 im3    | link[]  |           |          | 
 im2    | link[]  |           |          | 
 im1    | link[]  |           |          | 
 in0    | link[]  |           |          | 
 il1    | link[]  |           |          | 
 ic2    | link[]  |           |          | 
 ie3    | link[]  |           |          | 
Indexes:
    "node_chan_l_s_idx" btree (((nptr).chan), l, s)

</pre>

### Nodes

Now try:
<pre>
newdb=# select S,chap from Node limit 10;
     s      |       chap       
------------+------------------
 please     | notes on chinese
 yes        | notes on chinese
 请          | notes on chinese
 qǐng       | notes on chinese
 thankyou   | notes on chinese
 Méiyǒu     | notes on chinese
 谢谢        | notes on chinese
 xièxiè     | notes on chinese
 是的        | notes on chinese
 请在这里等    | notes on chinese
(10 rows)

</pre>

* An alternative view of relations is provided by NodeArrowNode:
<pre>
newdb=# select *  from NodeArrowNode LIMIT 10;
 nfrom | sttype | arr | wgt |              ctx              |   nto   
-------+--------+-----+-----+-------------------------------+---------
 (1,0) |     -1 |  69 |   1 | {please,"thank you",thankyou} | (1,1)
 (1,1) |     -1 |  67 |   1 | {thankyou,please,"thank you"} | (1,2)
 (1,1) |      1 |  68 |   1 | {thankyou,please,"thank you"} | (1,0)
 (1,1) |      1 |  68 |   1 | {news,online}                 | (2,291)
 (1,2) |      1 |  66 |   1 | {thankyou,please,"thank you"} | (1,1)
 (1,3) |     -1 |  69 |   1 | {please,"thank you",thankyou} | (1,4)
 (1,4) |     -1 |  67 |   1 | {please,"thank you",thankyou} | (1,5)
 (1,4) |      1 |  68 |   1 | {please,"thank you",thankyou} | (1,3)
 (1,5) |      1 |  66 |   1 | {please,"thank you",thankyou} | (1,4)
 (1,6) |     -1 |  67 |   1 | {please,"thank you",thankyou} | (4,0)
(10 rows)

</pre>

Notice how nodes (`nfrom`,`nto`,`nptr? ) and arrows (`arr`) are represented by pointer references
that are integers. When working with the graph, we often don't need to know the names
of things, we can get away with deferring the lookup of the actual data until we find what we're
looking for. That information can be cached so as to minimize the data transferred over the wire.

<pre>
newdb=# select S from Node where NPtr=(1,5);
   s    
--------
 xièxiè
(1 row)

</pre>

### Links and Arrows

A link is a composite relation that involves an arrow (pointer), a context,
and a destination node. Links are anchored to their origin nodes in the `Node` table
in the six columns `im3`, `im2`, `im1`, `in0`, `il1`, `ic2`, `ie3`.  
To find the links of type `Leads to':
<pre>
newdb=# select Il1 from Node where NPtr=(1,5);
                                       il1                                        
----------------------------------------------------------------------------------
 {"(66,1,\"{ \"\"please\"\", \"\"thank you\"\", \"\"thankyou\"\" }\",\"(1,4)\")"}
(1 row)

</pre>

Arrows are defined for each arrow pointer in the arrow directory:

<pre>
newdb=# select * from arrowdirectory limit 10;
 staindex |         long         | short | arrptr 
----------+----------------------+-------+--------
        4 | leads to             | lt    |      0
        2 | arriving from        | af    |      1
        4 | forwards             | fwd   |      2
        2 | backwards            | bwd   |      3
        4 | affects              | aff   |      4
        2 | affected by          | baff  |      5
        4 | causes               | cf    |      6
        2 | is caused by         | cb    |      7
        4 | used for             | for   |      8
        2 | is a possible use of | use   |      9
(10 rows)

</pre>

## The Go(lang) interfaces

The SSToryline package tries to make querying the data structures easy, by providing
generic scriptable functions that can be used easily in Go.

The open a database connection, to make any kind of query, with the help of the SSToryline package:
<pre>

package main

import (
	"fmt"
        SST "SSTorytime"
)

//******************************************************************

const (
	host     = "localhost"
	port     = 5432
	user     = "sstoryline"
	password = "sst_1234"
	dbname   = "newdb"
)

//******************************************************************

func main() {

	load_arrows := false

	ctx := SST.Open(load_arrows)

	row,err := ctx.DB.Query("SELECT NFrom,Arr,NTo FROM NodeArrowNode LIMIT 10")

	var a,c string	
	var b int

	for row.Next() {		
		err = row.Scan(&a,&b,&c)
		fmt.Println(a,b,c)
	}
	
	row.Close()

	SST.Close(ctx)
}

</pre>

## Low level wrapper functions

In general, you will want to use the special functions written for
querying the data.  These return data into Go structures directly,
performing all the marshalling and de-marshalling. The following are
basic workhorses. You will not normally use these.
For example, [see demo](https://github.com/markburgess/SSTorytime/blob/main/src/demo_pocs/postgres_stories.go).

* `CreateDBNode(ctx PoSST, n Node) Node` - idempotently create a node and return its database pointer structure.
* `AppendDBLinkToNode(ctx PoSST, n1ptr NodePtr, lnk Link, sttype int) bool` - idempotently attach an outgoing link to a Node.
* `UploadArrowToDB(ctx PoSST,arrow ArrowPtr)` - define an arrow in the arrow directory.
* `CreateDBNodeArrowNode(ctx PoSST, org NodePtr, dst Link, sttype int) bool` - Create a NodeArrowNode reference.

## Basic retrieval functions

* `GetDBNodePtrMatchingName(ctx PoSST,s string) []NodePtr` - returns a list of node pointers to names matching the input string as a substring

* `GetDBNodePtrMatchingNCC(ctx PoSST,chap,nm string ,cn []string) []NodePtr` - returns a list of node pointers to nodes matching the combined input string, filtered by chapter and context strategy matches. *This is currently a simple substring match, but is intended to offer more sophisticated fuzzy matching in future.* [See demo](https://github.com/markburgess/SSTorytime/blob/main/src/demo_pocs/search_by_name.go)

* `GetDBNodeByNodePtr(ctx PoSST,db_nptr NodePtr) Node` - Return the full node details from its pointer.

* `GetDBArrowsMatchingArrowName(ctx PoSST,s string) []ArrowPtr` - Find a list of arrows matching the given name as a substring.

* `GetDBNodeArrowNodeMatchingArrowPtrs(ctx PoSST,arrows []ArrowPtr) []NodeArrowNode` - Get a list of NodeArrowNode relations that involve the given arrow pointer type.

* `GetDBArrowByName(ctx PoSST,name string) ArrowPtr` - Return the arrow pointer for the given exact name.

* `GetDBArrowByPtr(ctx PoSST,arrowptr ArrowPtr) ArrowDirectory` - Return the arrow definition for a given arrow pointer.

## Future cone functions

The future cone of graph, from a starting node, is the set of all connected nodes following a particular
meta-type class (one of the four semantic spacetime meta-types). Think of this as the unfolding cone of influence, expanding from a root-cause, to integer link depth. 

This function will most likely be used when making inferences, or trying to connect the dots between the steps in a storyline.

For examples, [see demo](https://github.com/markburgess/SSTorytime/blob/main/src/demo_pocs/postgres_stories.go).

* `GetFwdConeAsNodes(ctx PoSST, start NodePtr, sttype,depth int) []NodePtr` - return a set of nodes in the forward cone, for links of given STType to fixed depth. By increasing the depth, one can slice the spatial hypersurfaces of the graph, circumferentially in expanding rings around the starting node.

* `GetFwdConeAsLinks(ctx PoSST, start NodePtr, sttype,depth int) []Link` - as above, but returning data as Link structures that contain both arrow and node information.

* `GetFwdPathsAsLinks(ctx PoSST, start NodePtr, sttype,depth int) ([][]Link,int)` - return the set of `proper time' paths expanding perpendicularly to the cirumferential/spatial layers. This is a form of path integral representation.

* `PrintLinkPath(ctx PoSST, alt_paths [][]Link, p int, prefix string)` - display a path structure returned above.


## Matroid Analysis Functions

For examples, [see demo](https://github.com/markburgess/SSTorytime/blob/main/src/demo_pocs/search_clusters_functions.go) and [example](https://github.com/markburgess/SSTorytime/blob/main/src/demo_pocs/search_clusters.go).

These functions will most likely be used during browsing of data, when getting a feel for the size and shape of the data.

* `GetMatroidArrayByArrow(ctx PoSST, context,chapter string) map[ArrowPtr][]NodePtr` - return a map of groups of nodes formed as matroids to arrows of all types, classified by arrow type.

* `GetMatroidArrayBySSType(ctx PoSST) map[int][]NodePtr` - return a map of groups of nodes formed as matroids to arrows classified by STType.

* `GetMatroidHistogramByArrow(ctx PoSST) map[ArrowPtr]int` - Find the group (member frequency) sizes of the above groups by arrow.

* `GetMatroidHistogramBySSType(ctx PoSST) map[int]int` - Find the group (member frequency) sizes of the above groups by STType.

* `GetMatroidNodesByArrow(ctx PoSST) []ArrowMatroid` - Return the group members of a matroid by arrow type.

* `GetMatroidNodesBySTType(ctx PoSST) []STTypeMatroid` - Return the group members of a matroid by STType.











