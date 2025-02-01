<!--
 SSTorytime - a ChiTek-i project by Mark Burgess

 Semantic Spacetime Story graph database library over postgresql (SSTorytime)
 This is an NLnet sponsored project, See https://nlnet.nl/project/SmartSemanticDataLookup/

-->

# SSTorytime

* This is a work in progress, planned for 2025 *

This project aims to make knowledge capture easy for humans
for general use. AI can only capture knowledge from humans, so even if we
want to use AI, we'd better get the knowledge representations right.

Knowledge capture requires tools for collecting notes, and structures
for representing and organizing them, so that they can be found
easily. Many mistakes have been made around this in the past, trying
to force discipline onto people at the wrong moment and neglecting to
do so when it matters. As a lifelong teacher of ideas, I've studied
this problem and think we can do better than what we have today.

## The approach

*If you want to know the deep thought behind an apparently ad hoc bit of
software, please take a look at [the project's own research website](http://markburgess.org/spacetime.html).*

Maintaining knowledge is like scaling the activity in a city. Time is
driven by parallel events arriving. The more people, the faster the
pace or pulse. We have to let this evolve, like the weather. You can't easy capture
the weather, but with a little ingenuity you can put up a sail to ride the wind,
and learn yourself to know its currents.

Memory has a horizon.  Things we learn, thoughts that we have, are
easily lost and forgotten.  Writing stuff down (intructions,
experiences, etc) seems like a good idea, but it's useless if no one
reads what is written regularly. It doesn't matter how you try to
remember something. If you don't use the memory, revisit it
frequently, it will decay to nothing.  This is why most Wikis and
documentation efforts fail. Writing notes for occasional use once a
year, or only in case of emergency is pointless unless it's used and
rehearsed often.  Wikipedia succeeds on average for a population
because it's many to many interaction, which keeps the pulse of
interaction alive: someone is always looking at the information. But
that doen't mean that you, an individual, will learn from what is
written there. Even Wikipedia information is not knowledge unless it's
in someone's head. To the causal browser, it is merely information,
even hearsay. Knowledge for the population is different from knowledge
for a single person. A library doesn't make you smartif you don't read
the books.

Technology tries to help, but it can also get us stuck doing the wrong
thing.  Semantic knowledge representations have not evolved since the
Semantic Web was proposed during the 1990s, at a time when the
technology was primitive and the technologists weren't themselves
experienced in how to use it.

Recently, graph databases have seemed to offer new possibilities for
knowledge representation, but the methods for using them have been
poorly developed and require the use of specialized query languages
and clumsy outdated formats.  In this project, I'm shifting the focus
away from technology to to content. We can use standard SQL databases
and modern lightweight data formats so embrace familiarity. The aim is
not to pickle a a static knowledge graph like an RDF structure, but
rather to establish a context dependent switching network of thoughts
that form stories. There's some overlap with what Large Language Models
do, but they take away control from the user. Here we give it back.

A user workflow starts from a simple note-taking language, then by
ingesting it into a database using a graph model based on the causal
semantic spacetime model, to the use of a simple web application for
supporting graph searches and data presentation. The aim is to make a
generally useful library for incorporating into other applications, or
running as a standalone notebook service.


## The tools

The toolset consistent of several components.

* N4L - a unicode text based note taking language for jotting down notes in a way
        that can be parsed and loaded into a semantic database. 
        N = note, 4 = 4 semantic type, L = language
        N4L = notes for loading


