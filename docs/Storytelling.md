
## SSTorytelling

We pass on knowledge by telling stories or building narratives. Some of these
are rooted in facts, some are metaphors and plain fiction and we use them in different
ways in different contexts to pass on facts and hypothetical lessons.



## N4L, a note taking language

We want a simple free text format for entering data, without
specialized encodings, for the first phase of jotting down items of
information that we want to learn and know.  This language must support Unicode for multi-language
support.  Memory data may come in a variety of media formats: text,
images, audio etc. When we are searching, however, symbolic language
is a convenient interaction format. So, in a first instance, we can
use pattern matching transducers to convert multimedia formats into
``alt texts'' and categorize those in a knowledge reasoning structure.
This may not strictly be necessary in the long run, but it's a useful
place to begin: in particular we need the ability to identify
sub-parts of an image of audio segment in order to cross-reference it.

The basic syntactic format of the contextual note taking language takes the following form:

## The SST Graph

`N4L` results in a graph representation, called the Semantic Spacetime (SST) graph.
If you've heard about Knowledge Graphs such as those using the Resource Description Framework (RDF)
and its related Web Ontology Language (OWL) then you'll know something about the idea already.
SST Graphs are not RDF or OWL, indeed they reject those early principles as a flawed concept.
Still, the differences between the two are subtle.
SST graphs impose a different kind of discipline on knowledge representations than RDF/OWL:

* SST graphs do not have XML schemas, nor are there formal ontologies.
* SST graphs use implicit typing based on *relation* (link) names.
* Graph nodes are focal points for any kind of data, though some text is usual.
* Graph links between nodes, i.e. graph *relations*, are unidirectional and may have any name.
 However, te names must be classified into one of four possible meta-types
that describe spacetime semantics:
* * SIMILARITY - a degree of equivalence
* * LEADS TO - a causally ordered relationship
* * CONTAINS - is one thing a part of another?
* * PROPERTY - a descriptive orexpressive property of a node

*These classes make it easier and more meaningful to search the graph later,
because their meanings are aligned with the processes of searching. 
The main problem with ontology and RDF is that they encourage you to model
the world as a number of things of different types, rather than modelling what
processes those things are involved in, i.e. the things we are interested in.
If we make data searchable by design, we avoid gettting into trouble later.*


## Some examples

Here are some examples of knowledge graphs.

* Everyone knows about flow charts. These can be rather trivial, or very complicated. They are the basis for 
error and risk graphs too. In N4L, we might write:

<pre>

- flow chart

@question Do I have the key?

Start (next) Find Door (next) $question.1

 $question.1 (next if yes) Open Door (next) End

 $question.1 (next if no)  Get key (next) $question.1

</pre>

The picture looks like this:

![A Flow Chart is a knowledge representation](https://github.com/markburgess/SSTorytime/blob/main/docs/figs/flow.png 'Flow Charts Are Knowledge Graphs')




## More generally an API for interacting with the SST graph
