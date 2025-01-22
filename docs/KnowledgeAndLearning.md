
# Knowledge and learning

You might think you know how to acquire knowledge. You learn stuff.
You take a course, or read a book. There is much more to it than that.
If you don't do the exercises, or apply the knowledge hands-on, you
quickly find that very little of what you heard or read sticks.

There's an old joke: that Wikis are places where knowledge goes to die. 
Well meaning individuals may invest hours of work to write something down for
others. But no one forms an intimate relationship to what has been written,
it is not knowledge. It's just a graveyard of bits and bytes that means nothing
to anyone except the person who wrote it.

Knowledge is more than learning. You can learn a page of text by
heart and still have no idea how to use it. As long as it remains a
lump of data, in your head or simply on paper, in a computer, or and the back of your hand,
it's of little use to you. Knowledge comes from knowing things deeply---by
forming a relationship to material. You know something when you know it like a friend.

https://medium.com/@mark-burgess-oslo-mb/the-failure-of-knowledge-management-5d97bb748fc3

This is one reason why current AI language models that seem to `know things'
are in fact as clueless about their subject matter as you are the day after reading
their results.

## Limits on learning

Even simply learning by memorization is hard. There is a limit to the
scaling of learning. Machine learning is not something that will lead
to endless improvements, because the cost benefit of learning rapidly
becomes untenable. Adding resources to capture every last variable
takes too long and costs too much energy. It's most effective for
common short-term experiential learning, like pattern recognition in
the immediate environment of sensory pre-processing. After that, a
Monte Carlo approach to guesswork or ``prediction'' would be much more
energy efficient.

Of course, one might say that humans have developed vast learning
resources--far greater perhaps than could be accounted for by a strict
cost-benefit analysis.  Doesn't this indicate that learning still has
a benefit?  But evolution doesn't direct the course of change on an
actionable timescale, it merely expresses a lukewarm approval over
benefits still to be spat upon by competitive opportunism.  Adding new
data will at some point plateau in terms of semantics or meaningful
features.  What remains is quite battle tested over a wide range of
contexts but might be only very rarely useful.
What is rarely spoken about in machine learning is the importance of
forgetting data: of post selection to keep only the most useful
pattern memories.

We also need to expose contextual knowledge.  Brute force learning
also doesn't discriminate by context: there is a far greater number of
seemingly-infinite combinatoric contexts for relevant data. Each will
have its own learning plateau associated with its charcteristic scales
(particularly timescales). In machine learning, transformer
architectures have discovered this, but they are still based mainly on
recall.

Learning examines spacelike and timelike invariants.
\begin{itemize}
\item Learning things discovers mesoscopic invariants snapshot structures.
\item Learning stories discovers behavioural invariants (mesoscopic sequences).
\end{itemize}
Context: where, when

Think of learning a language. You can quickly the vocabulary but you
can't recall it on demand. It takes several years of constant use to
develop the recall methods.


In this project, we look at trying to separate the processes that
manage different scales of learning to establish efficient recall structures
after post-selection.

## SSTorytelling

To devise a knowledge management system, our aims are:

* To assist in overcoming human limitations, while respecting the reason for them.
* To devise a way of getting experiences and thoughts into a computer representation
that will be used actively and immediately. For this, we shall devise a language N4L
or Notes For Learning.

The structure of cognition lies in decoding input into a lasting invariant
representation, which is linked to a rich set of very familiar contexts,
 with reactive output and post-associative feedback that expands the integration of what you've
memorized.

The distinction between input representation and recall or information
``usability'' (the process of distilling something useful) is the
basic problem. It's the same process we face when attempting to learn
a foreign language: out mental model is spontaneous, but forcing it
through the bottleneck of language is hard, because we don't know all
the words and we may never have tried to say what we want to say
before.  We would like to ensure that there is a sufficient reservoir
of examples to draw on and use as templates.

We may try to `cram' words into memory by repetition, but in the heat of the moment
we're unable to recall any of it. That's because memory lookup is contextual,
and the context in which we experience it and the context in which we try to learn it
are totally different---and we don't know how to connect the two.

Part of our difficulty in making notes and organizing is also our lack
of an overview. It doesn't take very much information to fill a page
and exceed our visual resolution. Even if we could organize everything
on a single sheet, our ability to resolve and comprehend it is
limited. So what we need is a way of integrating disjointed fragmemts
of note-taking to create that integrated model.

## Tidying up is a learning strategy

If you do it once, it just hides information. If you keep doing to improve and fiddle with the
organization of things, you become intimately and cognitively familar with the placement and
usefulness of things. This is what our brains do. There has to be re-use activity, involving motor
functions to extract things from their tidy places.

## Memory strategy and organization

When we tidy, often we create subject boxes first. It's by putting something in the right box that
we believe we understand it.
For example, when writing this, I write a number of section titles and try to collect
everything related to it under each heading..

Reorganizing notes post-hoc is very time consuming, because it requires multiple passes and many decisions.
Unless we intentionally place fragments of experience into boxes initially and intentionally,
we lose the economic benefits of experience. Intent thus serves an economic function for a cognitive
agent.

But there is a problem with this. When we group things together, we made trade-offs and approximations
that we might not agree with later. For instance, is a a duckbilled platypus a mammal or a bird.
It's a warm blooded species that lays eggs. It flouts the boundaries of the largest boxes in biology
by belonging to two boxes or a box all by itself. When we rely og box logic, things always go wrong in the end.
This is the trouble with ontologies as we use them in technology today. The main benefit of doing this
is that it forces us to revisit the model over and over again and make connections. We then identify knowledge
with those "aha!" moments at which we independently had insights. Those are the moments we remember
and can recall what we learned.

Our philosophical affectations have led us to go too far, however, when we form entire world *ontologies*
and try to fit everything into a single common spanning tree of knowledge. This is an abuse of process,
because ontologies are only outgrowths of coordinate systems for indexing {\em contexts}. They belong
to a separate sensory language that represents experience, not the concepts derived from cumulative processing and metaphoric expansion.

The temptation to continue classifying and sub-classifying by inventing ad hoc discriminators is a 
powerful tendency that occasionally runs riot and appeals to the bureaucratic mind. Even though applying
discrimintators from general experience ad hoc seems like a smart approach, it easily
 leads to tunnel vision.




## N4L, a note taking language

### Minimal syntax

We would like a simple free text format for entering data, without
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

The basic syntatic format of the contextual note taking language takes the following form:
