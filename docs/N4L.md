
# N4L - notes for loading utility

*Notes for loading*<br>
*Narrative for logic*<br>
*Network for learning*<br>
*Nourishment for life*

N4L is a simple language parser for keeping notes and scanning into a
structured format for uploading into a database or for use with other tools.

The purpose of using a simple yet semi-formal language as a starting
point is to avoid the "information model trap" that befalls many data
representations, i.e. forcing users to put everything into a pre-approved model,
like filling out a rigid form. This makes it hard to back out of decisions
and change our minds. It makes modelling fragile and fraught with risk.

Without any structure, it's only guesswork to
understand intent. N4L is a compromise that allows you to use any kind of
familiar editor to write notes in pure text (Unicode).

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
By adding the sequence directive to a context
<pre>

 +:: _sequence_ , poem ::   // starting sequence mode

 Mary had a little lamb         (note) Had means possessed not gave birth to
 Whose fleece was white as snow
 And everywhere that Mary went

 The lamb was sure to go        (note) SatNav invented later

 -:: _sequence_ ::          // ending sequence mode

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
* 1 **leadsto    / affects, causes** one thing follows from the other (sequences)
* 2 **contains   / contains** something is a part of something else (boxes in boxes)
* 3 **properties / express** something just has a name or an attribute (descriptive)
* 4 **similarity / near, alike** something is close to something else (proximity,closeness)

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


