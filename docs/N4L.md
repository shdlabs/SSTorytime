
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

" (relation) D                   # Continuination of chain

"2 (relation) E

: list, context, words :         # context (persistent) set

::  list, context, words ::

+: extend-list, context, words : # extend context set

-: delete, words :               # prune context set

paragraph =specialword paragraph paragraph paragraph paragraph
 paragraph paragraph paragraph paragraph paragraph
  paragraph paragraph =specialword paragraph paragraph paragraph
paragraph paragraph paragraph paragraphparagraph

=A                               # implicit example of/occurs in relation

</pre>

Here A,B,C,D,E are unicode strings

parentheses are reserved symbols. Literal parentheses can be quoted


Reserved topics and their aliases include:
* arrows     / follows
* members    / contains
* properties / express
* similarity / near

## Running state vector

The interpretation of the language has items, relationships, and context.
Within a stream of Unicode runes:

<pre>
type Parser struct 
{
stream_position  int
context_set      []string
item_set         []string
relation_set     []string
}

type Alias map[string]string

type Relation struct {

type int
fwd  string
bwd  string
}

</pre>