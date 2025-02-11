
# N4L - Notes for Learning
## A simple knowledge management language

*Notes for learning*<br>
*Narrative for loading*<br>
*Network for logical inference*<br>
*Nourishment for life*

N4L is an intentionally simple language, for keeping notes and
scanning into a structured format for uploading into a database or for
use with other tools.  The language is designed to encourage you to
think about how you express and structure notes. That, in turn,
encourages you to revisit, tidy and organise the notes again and again, while
being able to quickly turn them into a searchable graphical database, from which
and can reason through stories.

The purpose of using a simple yet semi-formal language as a starting
point is to avoid the "information model trap" that befalls many data
representations, i.e. forcing users to put everything into a pre-approved model,
like filling out a rigid form. This makes it hard to back out of decisions
and change our minds. It makes modelling fragile and fraught with risk.

Without any structure, it's only guesswork to
understand intent. N4L is a compromise that allows you to use any kind of
familiar editor to write notes in pure text (Unicode).

## Command line tool

The N4L tool ingests a file of "notes" written in a simple language
and turns it into a machine representation in the form of a "Semantic Spacetime" graph
 (a form of knowledge graph). This format is only tangentially related to the
 usual Resource Description Framework (RDF)), so we shall not use of
 refer to RDF in what follows, except to occasionally clarify the distinction. 
The command options currently include:
<pre>
usage: go run N4L.go [-v] [-u] [-s] [file].dat
  -adj string
        a quoted, comma-separated list of short link names (default "none")
  -i    incidence summary (node,links...)
  -u    upload
  -v    verbose
</pre>
For example, to parse and validate a file of notes, one can simply type:
<pre>
$ N4L chinese.in
$ N4L chinese.in Mary.in kubernetes.in
</pre>
Any errors will be flagged for correction. Using verbose mode gives extensive
commentary on the file, line by line:
<pre>
$ N4L -v chinese.in
</pre>
The final goal will normally be to upload the contents of the file to a database:
<pre>
$ N4L -u chinese.in
</pre>
However, before that, there are several operations than can be performed more efficiently
just from the command line for many data sets. This is because most knowledge input
is quite small in size, and quick feedback is very useful for ironing out flaws
and improving your source note material.

We can look at the subset of notes that are related by
a certain kind of relation, using abbreviated labels for relations. 
For example, to look for items linked by relation "(pe)" (which stands
for Pinyin to Hanzi translation) in a file of Chinese language, we could write:
<pre>
$ N4L -v -i -adj="pe" chinese.in
</pre>
We can add other kinds of relation too to expand the set:
<pre>
$ N4L -v -i -adj="pe,he" chinese.in
</pre>
This extracts a sub-graph from the total graph. It can be quite effective,
because most knowledge graphs are only sparsely linked (which is why logical
searches tend to yield nothing of interest).

In verbose mode, the standard output shows a summary of the text (events or items, etc)
and an excerpt of the adjacency matrix.
<pre>

$ N4L -v -i -adj="" Mary.in


------------------------------------
SUMMARIZE GRAPH.....

------------------------------------

0        Mary's mum

1        Nursery rhyme

0        SatNav invented later

0        Mary had a little lamb
         ... --( example of , 1 )-> Nursery rhyme [cutting edge high brow poem]
         ... --( written by , 1 )-> Mary's mum [poem cutting edge high brow _sequence_]
         ... --( then the next is , 1 )-> Whose fleece was white as snow [poem cutting edge high brow _sequence_]
         ... --( note/remark , 1 )-> Had means possessed not gave birth to [_sequence_ poem cutting edge high brow]

1        Had means possessed not gave birth to

2        Whose fleece was white as snow
         ... --( then the next is , 1 )-> And everywhere that Mary went [poem cutting edge high brow _sequence_]

3        And everywhere that Mary went
         ... --( then the next is , 1 )-> The lamb was sure to go [cutting edge high brow _sequence_ poem]

4        The lamb was sure to go
         ... --( note/remark , 1 )-> SatNav invented later [cutting edge high brow _sequence_ poem]

-------------------------------------
Incidence summary of raw declarations
-------------------------------------
Total nodes 8
Total directed links of type Near 0
Total directed links of type LeadsTo 4
Total directed links of type Contains 1
Total directed links of type Express 2
Total links 7 sparseness (fraction of completeness) 0.125
    - row/col key [ 0 / 8 ] Had means possessed not gave birth to
    - row/col key [ 1 / 8 ] SatNav invented later
    - row/col key [ 2 / 8 ] The lamb was sure to go
    - row/col key [ 3 / 8 ] Mary had a little lamb
    - row/col key [ 4 / 8 ] Whose fleece was white as snow
    - row/col key [ 5 / 8 ] Nursery rhyme
    - row/col key [ 6 / 8 ] And everywhere that Mary went
    - row/col key [ 7 / 8 ] Mary's mum

 directed adjacency sub-matrix ...

     Had means posse .. (   0.0   0.0   0.0   0.0   0.0   0.0   0.0   0.0)
     SatNav invented .. (   0.0   0.0   0.0   0.0   0.0   0.0   0.0   0.0)
     The lamb was su .. (   0.0   1.0   0.0   0.0   0.0   0.0   0.0   0.0)
     Mary had a litt .. (   1.0   0.0   0.0   0.0   1.0   1.0   0.0   1.0)
     Whose fleece wa .. (   0.0   0.0   0.0   0.0   0.0   0.0   1.0   0.0)
       Nursery rhyme .. (   0.0   0.0   0.0   0.0   0.0   0.0   0.0   0.0)
     And everywhere  .. (   0.0   0.0   1.0   0.0   0.0   0.0   0.0   0.0)
          Mary's mum .. (   0.0   0.0   0.0   0.0   0.0   0.0   0.0   0.0)

 undirected adjacency sub-matrix ...

     Had means posse .. (   0.0   0.0   0.0   1.0   0.0   0.0   0.0   0.0)
     SatNav invented .. (   0.0   0.0   1.0   0.0   0.0   0.0   0.0   0.0)
     The lamb was su .. (   0.0   1.0   0.0   0.0   0.0   0.0   1.0   0.0)
     Mary had a litt .. (   1.0   0.0   0.0   0.0   1.0   1.0   0.0   1.0)
     Whose fleece wa .. (   0.0   0.0   0.0   1.0   0.0   0.0   1.0   0.0)
       Nursery rhyme .. (   0.0   0.0   0.0   1.0   0.0   0.0   0.0   0.0)
     And everywhere  .. (   0.0   0.0   1.0   0.0   1.0   0.0   0.0   0.0)
          Mary's mum .. (   0.0   0.0   0.0   1.0   0.0   0.0   0.0   0.0)

 Eigenvector centrality score for symmetrized graph ...

     Had means posse .. (   0.7)
     SatNav invented .. (   0.2)
     The lamb was su .. (   0.4)
     Mary had a litt .. (   0.9)
     Whose fleece wa .. (   1.0)
       Nursery rhyme .. (   0.7)
     And everywhere  .. (   0.5)
          Mary's mum .. (   0.7)

</pre>
A useful ranking of nodes (known as EVC, or Eigenvector Centrality, which is something like Google's PageRank) 
can be calculated from the weighted graph matrix (see below). The higher the score number, the more
interconnected or "important" a term of text is, e.g.
<pre>
$ ../src/N4L -v -i -adj="" chinese.in

  ...

 Eigenvector centrality score for symmetrized graph ...

            Fángjiān .. (   0.3)
             jiàoshì .. (   0.8)
              Kètáng .. (   0.2)
                   教室 .(   0.2)
     place/area/dist .. (   0.1)
                   表现 .(   0.7)
            Biǎoxiàn .. (  0.6)
                   课堂 .(   0.8)
         performance .. (   0.5)
                   房间 .(   0.3)
                   地方 .(   0.2)
                   表演 .(   0.1)
           classroom .. (   1.0)
                room .. (   0.2)
              Dìfāng .. (   0.3)

</pre>

## Language syntax

<pre>

#  a comment for the rest of the line
// also a comment for the rest of the line

-section/chapter                 # declare section/chapter as the subject

: list, context, words :         # context (persistent) set
::  list, context, words ::      # any number of :: is ok

+:: extend-list, context, words :: # extend the existing context set
-:: delete, words :                # prune the existing context set

A                                # Item
Any text not including a "("     # Item
"A..."                           # Quoted item
'also "quoted" item'             # Useful if item contains double quotes
A (relation) B                   # Relationship
A (relation) B (relation) C      # Chain relationship
" (relation) D                   # Continuination of chain from previous single item
$1 (relation) D                  # Continuination of chain from previous first item
$2 (relation) E                  # Continuation from second previous

@myalias                         # alias this line for easy reference
$myalias.1                       # a reference to the aliased line for easy reference

NOTE TO SELF ALLCAPS             # picked up as a "to do" item, not actual knowledge

"paragraph =specialword paragraph paragraph paragraph paragraph
 paragraph paragraph paragraph paragraph paragraph
  paragraph paragraph =specialword *paragraph paragraph paragraph
paragraph paragraph paragraph paragraphparagraph"

where [=,*,..]A                        # implicit relation marker

</pre>
Here A,B,C,D,E stand for unicode strings. Reserved symbols:
<pre>
(), +, -, @, $, and # 
</pre>
Literal parentheses can be quoted. There should be no whitespace after the initial quote
of a quoted string.

## Sequence mode ##

Sometimes it's useful to link items together into a chain or sequence.
By adding the sequence directive to a context. From the example of the Mary had a little lamb above:
<pre>

$ more Mary.in

-poetry

 :: cutting edge, high brow ::

 +:: _sequence_ , poem ::   // starting sequence mode

@title Mary had a little lamb  (note) Had means possessed not gave birth to
              "                (written by) Mary's mum

       Whose fleece was white as snow
       And everywhere that Mary went

       // no need to be contiguous

       The lamb was sure to go        (note) SatNav invented later

 -:: _sequence_ ::          // ending sequence mode

 $title.1 (example of) Nursery rhyme

</pre>
This results is a sequence of lines linked by `then' arrows, until the context is removed.
<pre>
Mary had a little lamb (then) Whose fleece was white as snow (then) ...
</pre>
Then is a pre-defined and effectively reserved association.

* Only the first items on a line are linked. 
* Only new items are linked, so the use of a " or variable reference will not trigger a new item.

## Example

Assocations have explanatory power, so we want to take advantage of that.
On the other hand, we don't want to type a lot when making notes, so
it's sensible to make extensive use of abbreviations.

<pre>
-chinese notes

::food::

  meat    (is english for the pinyin) ròu
   "      (is english for the chinese or hanzi)  肉

  # more realistic with abbreviations ...

 菜 (hp) Cài (pe) vegetable 
 meat (eh) 肉 (hp) Ròu
 beef  (eh) 牛肉  (hp) Niúròu
 lamb  (eh) 羊肉  (hp) Yángròu
 chicken (eh)  鸡肉 (hp)  Jīròu

:: phrases, in the hotel ::

@robot I'm waiting for some food from the robot (eh) 我在等机器人送来的食物 (hp) Wǒ zài děng jīqìrén sòng lái de shíwù

:: technology ::

jīqìrén (pe) robot (example) $robot.1

</pre>

Notice how the implicit "arrows" in relations like 
<pre>(is english for the pinyin)</pre> or its short form
<pre>(pe)</pre> effectively define the `types' of thing they are 
attached to at either end. So we don't need to define the ontology for things
because it emerges automatically from the names
we've given to relationships.

Semantic reasoning can make use of both the precision and the fuzziness of associative types
when reasoning. This is a powerful feature that enables automated
inference with lateral thinking, just as humans do. Languages that use
logic to define ontologies are greatly over-constrained and make
reasoning precise but trivial, because they can only retrieve exactly
what you typed into the model.

## How relationships work

A piece of text can be thought of as an item or an event.
Relationships between items or events are written inside parentheses, as in the
examples above. They are designed to be highly
abbreviated for note taking. 

As written, the examples above look a bit like any old RDF (Resource
Description Framework) triplets. However, with the underlying
assumptions of the language that we'll lay out below, they are much
more powerful than the ad hoc references in RDF, because RDF
relationships are just strings without any semantics.

In order for references to be used for reasoning (and effective
semantic search), they need some basic properties. The simplest thing
we can do is to classify each relationship as though it were a special
case of one of four basic types, depending on how you want to
interpret it. This might be tricky in the beginning, so you can stick
to some predefined relation.

It turns out that every relationship basically falls into one of
four basic types that help you to imagine sketching the items on a map.
Here are the four types:
* 0 **similarity / near, alike** something is close to something else (proximity,closeness)
* 1 **leadsto    / affects, causes** one thing follows from the other (sequences)
* 2 **contains   / contains** something is a part of something else (boxes in boxes)
* 3 **properties / express** something just has a name or an attribute (descriptive)

For example:
* 0 - A **(sounds like)** B, or B **(sounds like)** A
* 1 - A **(causes)** B , or B **(is caused by)** A
* 2 - A **(is the boss of)** B, or B **(has boss)** A
* 3 - A **(has a degree in)** B, B **(is a qualification of)** A

*(Technical note the use of integers allows us to use signs for orientation.
Similarity is directionless 0 = -0; for the others there is a difference between
positive and negative inverses.)*

These four classes of association can be literal or metaphorical (all language
is an outgrowth of [metaphors for space and time](https://www.amazon.com/Smart-Spacetime-information-challenges-process/dp/B091F18Q8K/ref=tmm_hrd_swatch_0)).
behave like placing items around
each other in a mind-map on paper. Things that belong close together
because they are aliases for one another are *similar*.  If one thing
leads to another, e.g. because it causes it or because it precedes it
in a history then we use *leadsto*. Some items are parts of other items,
so we use *contains*. Finally, something that's purely descriptive
or is expressed by an item, e.g. "is blue" or 

Some authors who write about semantic networks have suggested that the
way to think about arrows and nodes is as "nouns" (things) and "verbs"
(actions). This is a simple idea, but it's not quite right. The catch lies
in the way language semantics rely almost entirely on metaphors to express
ideas. We frequently speak of "nouning verbs" and "verbing nouns", e.g.
in Silicon Valley speak:
<pre>
 The company's spend is ...   (vs)    I need to spend .. an expenditure
 I have a big ask ...         (vs)    I need to ask you .. a question

 I question your use of language ... with a question
 I expensed by trip ... as an expense
</pre>
Spend is a verb (expenditure or budget are nouns. Ask is a verb, question is
a noun, but we now use both for both!
We see that language is used and abused in fluid ways, so we need more
discipline in thinking about what the functions of terms are.


## Examples and pitfalls in modelling


When we say that A follows B, this may apply to things or actions.
* Space travel came after aircraft. 
* Shopping is done after work.
* Hammering is done after assembly.
Order applies to both processes and objects.


We could imagine a supply-chain worker noting:
<pre>
 delivery 123 (damaged) 2 boxes
</pre>
It's a fair thing to write in a moment of unexpected pressure. But which of the
four relations is this? That's the same as asking: what could we use this note
for later? The problem with it is that it's ambiguous.

The left hand side "delivery 123" is clear enough. It represents some shipment
and we could embellish this description like this
<pre>
 delivery 123 (contains) shoes
     "        (came from) Italy
     "        (received by) shift crew 12
</pre>
and so on. So no problem here. The relation "damaged" becomes an issue however
because it's referring to the condition or state of the delivery. 
A more flexible approach would be to rewrite this as
<pre>
 delivery 123 (condition) 2 boxes damaged
</pre>
because now
* condition is a generic and reusable relation, which is a propery attribute (type 3) of the delivery
* "2 boxes damaged" is an event that can be explained easily
For instance, now we can explain the event further:
<pre>
  2 boxes damaged (condition) water damage
         "        (contains) red stiletto box 1445
         "        (contains) black stiletto box 1446
</pre>

### The "is a" fallacy

During the OO-movement to sanctify Object Orientation as a software modelling approach, many
superifical ideas were proposed. Object Orientation rubber stamps
the idea that objects, i.e. "things" are the most important entities in a model, 
leaving *processes* asking: what am I 
then? (The answer was usually that processes should be thought of as methods that affect
objects, which is extremely limiting.)
Classification of objects into types was the goal of OO, because this is a way to simply
map ideas into first order logic, and that makes programming easy to understand.
Alas, squeezing processes into this isn't always easy.
The answer commonly associated with this was to use the "is a" or "is an instance of" relation
as the way of thinking about things.
<pre>
Object X is an instance of a class Square
A Square is a special case (inheriting) the class of Rectangle
etc.
</pre>
The trouble with this idea is that it attempts to assert an *static* or *invariant* truth
about the role of something (the square). But squares, indeed any properties or
roles, are typically context dependent. We use the same concept in different ways.
<pre>
In DIY: A hammer is a tool.
In music: A hammer is a musical instrument
In DIY: a drill is a tool for making holes.
In operations, a drill is a practice episode.
</pre>
If we insist of having different types for each of these cases (a type polymorphic approach),
we push the responsibilty of the technology back onto the person using it. Technology
is supposed to work for humans, not the other way around.

The example above of damaged delivery  is a good example of how this becomes
problematic. Suppose we introduce an object for a delivery, is that
"Delivery" or "Shoes"? Should we have a separate object for "Damaged delivery" or is
damage an attribute of the object. What could it mean? how would we explain it?

The virtue of a semantic language is that we never have to shoe-horn
(no pun intended) an idea into a rigid box, as we do when we try to
lock down data types. This is an affectation of logical reasoning,
but logic is highly restrictive (on purpose, as a matter of design).
That makes it precise, but also extremely fragile to variability.

### Belonging

Some relationships can be tricky to fathom. The semantics of ownership,
for example, are not completely unambiguous. Suppose you want to say

<pre>
The bracelet "belongs to" Martin 
</pre> 

Is the bracelet a property of Martin or a part of him?  As an object,
we might choose to make this a part the "extended space of
martin". There is no right answer. You can choose what works for you.
The difference between the two is how they are searched.  If we
interpret the bracelet as "a part of" Martin then we can also say that
the bracelet contains a diamond and thus the diamond is also a part of
Martin, because "part of" is a transitive relationship. But if we say
that the bracelet is just something that characterizes him, it's not
clear that that is transitive because a bracelet may be characterized
by being golden but this does not imply that Martin is golden!

You might make the wrong choices about things initially, but it's easy to
change your decision because the definition of the relationship is
made independently of all the data where you use it. You'll figure out
the bugs in your wordings as you go, and it's precisely this reworking
that is learning.

The usefulness
of a language interface becomes clear now. It's much easier to edit your notes than to maintain
a database.


