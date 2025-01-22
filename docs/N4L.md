
# N4L - notes for loading utility


The purpose of modelling a language for the information as a starting
point is to avoid the "information model trap" that befalls many data
representations. If we can first determine the necessary and
sufficient roles for semantic expression without assuming a normal form
that later becomes a liability, then we start from the correct user perspective.
We also understand what data structures will be needed on a pragmatic level.

N4L parse text files that contain a semmi-formal language for easily taking notes and commenting
on their meaning. The idea is to jot down examples and nuggets of meaning.

## Language syntax

<pre>

-section                         # topic or reserved-topic

A                                # Item
"A"                              # Quoted item
A (relation) B                   # Relationship
A (relation) B (relation) C      # Chain relationship
" (relation) D                   # Continuination of chain from previous single item
$1 (relation) D                  # Continuination of chain from previous first item
$2 (relation) E                  # Continuation from second previous
: list, context, words :         # context (persistent) set
::  list, context, words ::
+: extend-list, context, words : # extend context set
-: delete, words :               # prune context set

@name                            # alias this line for easy reference
@name.$1                         # alias column in a line for easy reference

"paragraph =specialword paragraph paragraph paragraph paragraph
 paragraph paragraph paragraph paragraph paragraph
  paragraph paragraph =specialword paragraph paragraph paragraph
paragraph paragraph paragraph paragraphparagraph"

[=,*,..]A                        # implicit relation marker

</pre>

Here A,B,C,D,E are unicode strings

parentheses are reserved symbols. Literal parentheses can be quoted


Reserved topics and their aliases include the four spacetime meta-semantic types:
* leadsto    / affects, causes
* contains   / contains
* properties / express
* similarity / near, alike

## Running state vector

The interpretation of the language has items, relationships, and context.
Within a stream of Unicode runes:

<pre>
# parser state

type Parser struct 
{
stream_position  int
context_set      []string
item_set         []string
relation_set     []string
}

</pre>
<pre>
# abbreviation lookup table

type Alias map[string]string

</pre>
<pre>

# relation lookup

type LinkType map[string]string

</pre>
<pre>
# relation inverse and type table

type Association struct {

type int
fwd  string
bwd  string    # currently undecied how to represesnt negative patterns NOT, !, exceptions
}

</pre>



## TODO

Implement aliasing and inferences of graph structures