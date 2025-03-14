
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


### A flow chart example

Everyone knows about flow charts. These can be rather trivial, or very complicated. They are the basis for finite state machines (FSM), as well as error and risk graphs too. In N4L, we might write:

<pre>

- flow chart

@question Do I have the key?

Start (next) Find Door (next) $question.1

 $question.1 (next if yes) Open Door (next) End

 $question.1 (next if no)  Get key (next) $question.1

</pre>

The picture looks like this:

![A Flow Chart is a knowledge representation](https://github.com/markburgess/SSTorytime/blob/main/docs/figs/flow.png 'Flow Charts Are Knowledge Graphs')

In this case, we defined the arrows in the `N4Lconfig.in` file.

<pre>
- leadsto

 	# Define arrow causal directions ... left to right

        + is followed by (next) - is preceded by (prev)    
        + then the next is (then) - previous (prior)

        // Flow charts / FSMs etc

	+ next if yes (ifyes) - is a positive outcome of (bifyes)
	+ next if no (if no)  - is anegitive outcome of (bifno)

</pre>

### Dial M for Murder

Associations of clues and bits of information form forensic trails that are ideally suited
to graph representations. You can imagine a crime solving team entering all their evidence into
a graph and searching it for possible connections using inferences along the way. This is more
powerful than simply applying logical rules to an ontology, because logic can never tell you any more
than you explicitly stated in the beginning. Using `fuzzy' inferences, on the other hand, we can
perform lateral reasoning just like humans do. The goal is to be able to tell a plausible story
about something.

The details of the graph below are not yet defined, but you can imagine that they lead to
an organization of thought something like the picture below.

![A study or murder](https://github.com/markburgess/SSTorytime/blob/main/docs/figs/knowledge.png 'The large scale structure of a well-formed knowledge graph organizes knowledge into regions that lead from one to the other.')



## More generally an API for interacting with the SST graph

Once data have been entered into a SSToryline database, we want to be able to extract it again.