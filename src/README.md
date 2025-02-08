
# Source code for the N4L knowledge representation

Some notes on the source code

## `N4L.go` 

This program promises:

-To parse a configuration file in the current directory to define some reusable terms
-To parse a knowledge file of "notes" and convert into a graph
-To output the summary as text or upload data to a Postgres database

The nodes (or vertices) of the graph are blobs of text, which might be
anything from a single word to a long passage pasted in from
somewhere, depending on what the user is trying to remember.  The
links (or edges) are semantic relationships that fall into four basic
meta-types (called STtypes for spacetime types). Users can define
directed links (and their inverses) with any names that are compatible with these types,
using the configuration file. After that, they can use simple abbreviations
in an easy way. The result is a graph structure, which is indexed by node text.

Since node texts are potentially long, we don't want to use them for any purpose
if we can avoid it. So we register every string in a directory or lookup table.
It's known, from previous research, that the frequency or probability of getting a
string of length *L* falls off as a power law with increasing *L*. Most strings are
short and the likelihood of seeing the same string more than once falls off rapidly (faster
than exponentially). So, knowing the kinds of operations we are going to need to do
on the data, we arbitrarily introduce six classes of text that are handled 
independently but transparently for the user: n-grams (words separated by spaces) of one, two,
and three words; short strings of less than 128 Unicode characters, longer strings less than 1024 characters, and everything else. These are turned into two dimensional coordinates (class,index)
given by their array keys and lengths. This text is kept in a number of swim lanes depending
on the process used to manage them:

`
 
type NodeEventItemBlobs struct {

	// Power law n-gram frequencies

	N1grams map[string]CTextPtr
	N1directory []NodeEventItem
	N1_top CTextPtr

	N2grams map[string]CTextPtr
	N2directory []NodeEventItem
	N2_top CTextPtr

	N3grams map[string]CTextPtr
	N3directory []NodeEventItem
	N3_top CTextPtr

	// Use linear search on these exp fewer long strings

	LT128 []NodeEventItem
	LT128_top CTextPtr
	LT1024 []NodeEventItem
	LT1024_top CTextPtr
	GT1024 []NodeEventItem
	GT1024_top CTextPtr
}
 
`


Now, searching can be done in a way that is appropriate for a string of unknown length.
For short strings, a hashmap/btree can be used for lookup. For long strings, we cans simply
do a linear search since there will not be many, since hashing a long string is costly.
The lookup tables associate an integer primary key index with a string. For long strings,
comparing the length of the string can quickly eliminate mismatches, without needing
to see the content. These details explain the apparent "over-thinking" of data representations.


## Knowledge graph and matrices

A graph is typically a sparse structure, i.e. the number of links is
much less than the square of the number of nodes. Graphs may be
represented by *Adjacency Matrices* and by *Incident Matrices*. Given
the directory lookup tables, which are simply allocated FIFO, it makes
sense to use an incident matrix representation to avoid an extra
lookup.

The main data structures are:

* `NODE_DIRECTORY NodeEventItemBlobs` - a list of nodes/events/items is stored in the structure
above, consisting of six arrays of type `[]NodeEventItemPtr`.

* `ARROW_DIRECTORY` - an array of arrow structures that can be searched linearly or using
pointer shortcuts kept by short and long name (`ARROW_SHORT_DIR` and `ARROW_LONG_DIR`).

* The `NodeEventItem` structure is the graph node, and a list of outgoing links with positive
STtypes. Incoming links have negative STtypes. Thus each node acts as a multiway switch (a local
index at each node) for immediate lookup. `

type NodeEventItem struct { // essentially the incidence matrix

	L int              // length of name string
	S string           // name string itself
	C int              // the string class: N1-N3, LT128, etc

	A [ST_TOP][]Link   // link incidence list, by arrow type
  	                   // NOTE: carefully how offsets represent negative SSTtypes
}

`
* A pointer to such an object `NodeEventItemPtr` is classified for quickly reference
to which `swim lane` it belongs to.
`
type NodeEventItemPtr struct {

	Ptr   CTextPtr              // index of within name class lane
	Class int                   // Text size-class
}
`

* The `CTextPtr` is an alias for an integer pointer to pre-classified text in a lookup table.
We use this alias mainly to distinguish the kind of index, because several index
roles are in use.

## The long and short of arrows

Arrows in a graph represent relations with semantics. The semantics come from the
name and type of the arrow. When browsing data, we'd like links to be explanatory
with suitably long names. For jotting notes, this is cumbersome, so we'd like to
use short abbreviations. This is built into to the language as a requirement,
described in the configuration file. Thus, every link has both a long and a short name,
and every inverse link has the a long and a short inverse, because the user
doesn't want to have to think about semantic bureaucracy when making notes.

From a code perspective, the semantics are divided into four 
meta-types called `STtype`s (actually seven with the inverses). The types
determine the way the links will be searched.

`
	NEAR = 0      // no inverse 0 = -0
	LEADSTO = 1   // +/-
	CONTAINS = 2  // +/-
	EXPRESS = 3   // +/-

	ST_OFFSET = EXPRESS // so that ST_OFFSET - EXPRESS == 0
	ST_TOP = ST_OFFSET + EXPRESS + 1
`
