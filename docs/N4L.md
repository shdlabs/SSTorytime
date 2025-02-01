
# N4L - notes for loading utility

N4L is a simple language parser for keeping notes and scanning into a
structured format for uploading into a database or for use with other tools.

The purpose of using a simple yet semi-formal language as a starting
point is to avoid the "information model trap" that befalls many data
representations, i.e. forcing users to put everything into a pre-approved,
like filling out a form. Without any structure, it's only guesswork to
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
A (relation) B                   # Relationship
A (relation) B (relation) C      # Chain relationship
" (relation) D                   # Continuination of chain from previous single item
$1 (relation) D                  # Continuination of chain from previous first item
$2 (relation) E                  # Continuation from second previous

@myalias                            # alias this line for easy reference
$myalias.1                          # a reference to the aliased line for easy reference

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
Literal parentheses can be quoted

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

## How references work

References are written in parentheses. They are designed to be highly
abbreviated for note taking. As they are written, the examples above
look like RDF (Resource Description Framework) triplets. However, they
are much more powerful than the ad hoc references in RDF.  In order
for references to be used for reasoning and effective semantic search,
they need to be declared with more properties. Declarations are made in the configuration file. 

Each relationship needs to be classified as one of four types depending
on how it is to be interpreted. This might be tricky in the beginning, so you
can stick to some predefined relation.

Reserved topics and their aliases include the four spacetime meta-semantic types:
* leadsto    / affects, causes
* contains   / contains
* properties / express
* similarity / near, alike

These four classes of association behave like placing items around
each other in a mind-map on paper. Things that belong close together
because they are aliases for one another are *similar*.  If one thing
leads to another, e.g. because it causes it or because it precedes it
in a history then we use *leadsto*. Some items are parts of other items,
so we use *contains*. Finally, something that's purely descriptive
or is expressed by an item, e.g. "is blue" or 

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
made independently of all the data where you use it. This is the usefulness
of a language interface.