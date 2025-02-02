
# The approach

*If you want to know the deep thought behind an apparently ad hoc bit of
software, please take a look at [the project's own research website](http://markburgess.org/spacetime.html).*

Knowledge management is like tending a garden. You plant, you have to
tend, you organize, and you have to pull weeds. Eventually, you learn
your way around the garden, where to pick flowers and vegetables. The
bigger the garden, the more effort it takes to reach this level of
familiarity.

When we're familar with knowledge, we communicate by telling stories
about it.  A story is much more than idle gossip about unfamiliar
hearsay: it has to make sense, it has to be grounded in intimate
detail. Condensing knowledge into *intentional* stories (at different
levels of detail for different purposes) helps us to pass on
knowledge.  Some stories are about solving problems, others are merely
descriptive lessons learned.

Gardens aside, one could say that generating familiarity and
maintaining knowledge is like scaling the activity in a city. Time is
driven by parallel events arriving. The more people, the faster the
pace or pulse. We have to let this evolve, like the weather. You can't
easy capture the weather, but with a little ingenuity you can put up a
sail to ride the wind, and learn yourself to know its currents.

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

