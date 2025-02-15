
# Source code and design for the N4L knowledge representation

Some notes on the source code:

## `N4L.go` 

This program promises:

* To parse a configuration file in the current directory to define some reusable terms
* To parse a knowledge file of "notes" and convert into a graph
* To output the summary as text or upload data to a Postgres database

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

<pre> 
type NodeEventItemBlobs struct {

	// Power law n-gram frequencies

	N1grams map[string]ClassedNodePtr
	N1directory []NodeEventItem
	N1_top ClassedNodePtr

	N2grams map[string]ClassedNodePtr
	N2directory []NodeEventItem
	N2_top ClassedNodePtr

	N3grams map[string]ClassedNodePtr
	N3directory []NodeEventItem
	N3_top ClassedNodePtr

	// Use linear search on these exp fewer long strings

	LT128 []NodeEventItem
	LT128_top ClassedNodePtr
	LT1024 []NodeEventItem
	LT1024_top ClassedNodePtr
	GT1024 []NodeEventItem
	GT1024_top ClassedNodePtr
}
 
</pre>


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
index at each node) for immediate lookup.

<pre>

type NodeEventItem struct { // essentially the incidence matrix

	L int                 // length of name string
	S string              // name string itself

	Chap string           // section/chapter in which this was added
	SizeClass int         // the string class: N1-N3, LT128, etc
	NPtr NodeEventItemPtr // Pointer to self

	I [ST_TOP][]Link   // link incidence list, by arrow type
  	                   // NOTE: carefully how offsets represent negative SSTtypes
}

</pre>
* A pointer to such an object `NodeEventItemPtr` is classified for quickly reference
to which `swim lane` it belongs to.

<pre>
type NodeEventItemPtr struct {

	CPtr  ClassedNodePtr        // index of within name class lane
	Class int                   // Text size-class
}
</pre>

* The `ClassedNodePtr` is an alias for an integer pointer to pre-classified text in a lookup table.
We use this alias mainly to distinguish the kind of index, because several index
roles are in use.

## Finding nodes

We don't expect to have to use a function to get a pointer from a node text. Instead
we insist that string registration be transparently idempotent, so to get a pointer
for some text, simply register it with `IdempAddTextToNode()`, this updates
the node directory if necessary and returns a pointer, or it simply returns
a pointer if the record already existst. The converse (finding the text
pointed to by a text pointer) uses `GetTextFromPtr()`. 

## The long and short of arrows

Arrows in a graph represent relations with semantics. The semantics come from the
name and type of the arrow. When browsing data, we'd like links to be explanatory
with suitably long names. For jotting notes, this is cumbersome, so we'd like to
use short abbreviations. This is built into to the language as a requirement,
described in the configuration file. Thus, every link has both a long and a short name,
and every inverse link has the a long and a short inverse, because the user
doesn't want to have to think about semantic bureaucracy when making notes.

* `ARROW_SHORT_DIR` and `ARROW_LONG_DIR` are hashmaps for quickly finding an integer pointer
to an arrow.

## Arrow semantics

From a code perspective, the semantics are divided into four 
meta-types called `STtype`s (actually seven with the inverses). The types
determine the way the links will be searched.

<pre>
	NEAR = 0      // no inverse 0 = -0
	LEADSTO = 1   // +/-
	CONTAINS = 2  // +/-
	EXPRESS = 3   // +/-

	ST_OFFSET = EXPRESS // so that ST_OFFSET - EXPRESS == 0
	ST_TOP = ST_OFFSET + EXPRESS + 1
</pre>


Incoming and outgoing arrows are treated distinct. The order in which
they are entered is retained while parsing so that users see what they
have typed. However, when summing adjacencies for the graph we need to
be careful not to count arrows twice. We also need to be able to
search for both the incoming and outgoing arrows from a single
node. If we had simply created a table for all nodes-link triplets
(all_from,link,all_to), it would easy to search for links incoming and
and outgoing, but this is not very efficient for normal usage and would involve
a lot of redundant work for the usual use-cases.

This issues comes back to concern us when translating a graph model into SQL. SQL was designed
to work with normalized data models, which were (in turn) optimized for human entry in a random
access pattern. Here, on the
other hand, there is never any human entry into the database, and the APIs we provide here can
manage duplicated records easily (not least because data almost never change). We can make use
of the invariances of the data to separate data structures into different blocks, linked by pointers.

## Representing data structures in SQL

There are several issues around storing the data from N4L in a database. Apart from the
issue of speed and efficiency, we need to

* Represent lists of links, which are triplets of primary keys (a,b,c).
* Deal with unicode strings for storage and searching.

In Go(lang), links are represented as array slices, but databases do not typically support arrays
directly, and the normal thing to do would be to create another entity relation. This is possible, but
again would lead to adding redundant work during common searches. Call me old-fashioned, but I hate
waste sot he argument that modern CPUs will make it fast doesn't stop me from trying to optimize for the
actual processes involved in searching.

The key thing about saving data structured by N4L in SQL is that SQL
is not traditionally compatible with the datatypes and encapsulation
mechanisms used by Go (associative arrays or maps, etc). Although one
could simulate private arrays with a table restricted by a field
matching a primary key, this is wasteful and untidy. Its efficiency
becomes noticable as potentially the square of the number of nodes we
add.

So, traditional SQL doesn't have a way of encapsulating tables of a
particular type to primary key ranges that are private to a particular
parent node. Postgres, however, does have internal array types that
support this. Using this feature would make the database model non-portable,
but that was already foregone in the project proposal.

Ideally, we can hope that postgres arrays are efficient for lookup.
In the interests of science, I therefore commit to performing some
experiments and tests on these two alternatives.






