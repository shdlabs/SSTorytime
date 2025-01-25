
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

@myalias                            # alias this line for easy reference
$myalias.1                          # a reference to the aliased line for easy reference

"paragraph =specialword paragraph paragraph paragraph paragraph
 paragraph paragraph paragraph paragraph paragraph
  paragraph paragraph =specialword paragraph paragraph paragraph
paragraph paragraph paragraph paragraphparagraph"

[=,*,..]A                        # implicit relation marker

</pre>
Here A,B,C,D,E stand for unicode strings
Parentheses are reserved symbols. Literal parentheses can be quoted

## Example

Assocations have explanatory power, so we want to take advantage of that.
On the other hand, we don't want to type a lot when making notes, so
it's sensible to make extensive use of abbreviations.

<pre>

::food::

  meat    (english for pinyin) ròu
          (english for hanzi)  肉
  chicken (english for pinyin) jīròu
          (english for hanzi)  鸡肉 

  lamb (english for pinyin) yángròu (pinyin for hanzi) 羊肉
  beef (english for pinyin) niúròu  (pinyin for hanzi) 牛肉
  milk (english for pinyin) niúnǎi  (pinyin for hanzi) 牛奶

  # more realistic with abbreviations ...

菜 (hp) Cài (pe) vegetable 
meat (eh) 肉 (hp) Ròu
beef  (eh) 牛肉  (hp) Niúròu
lamb  (eh) 羊肉  (hp) Yángròu
chicken (eh)  鸡肉 (hp)  Jīròu
pork  (eh) 猪肉  (hp) Zhūròu
soup (eh)  汤 (hp) Tāng
sugar (eh)  糖 (hp) Táng
porridge  (eh) 粥 (hp) zhōu

:: phrases, in the hotel ::

@robot I'm waiting for some food from the robot (eh) 我在等机器人送来的食物 (hp) Wǒ zài děng jīqìrén sòng lái de shíwù

:: technology ::

jīqìrén (pe) robot (example) $robot.1

</pre>

Notice how the different ends of the implicit arrows in (pe) from
pinyin to English effectively define the type of thing they are
attached to. So we don't need to define the ontology for things
because it can be allowed to emergy from the already implicit choices
we've made about the kinds of relationship things can have.

Semantic reasoning can make use of the fuzziness of associative types
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