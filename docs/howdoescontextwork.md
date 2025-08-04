
# How does context work?

Context is the "hard problem" of knowledge management. We use the context in which we experienced
something as a "lookup key" for memories. It's not like the primary keys we feed to a database, e.g.
something like a name, a phone or social security number, etc. Cognitive processes use sensory
inputs to encode memory. How does one feed that into a directory listing to get an answer? And how do we list
what was there?

The knowledge graph is a way of painting a picture of a scene, but we still need to find the right scene.
Modern recognition methods can match sensory inputs like vision and sound as well as writing now, but
they don't solve the problem of how to know which experience is the correct match given a generic
sensory input. There was that one time at band camp....

## An idealized approximation

The way we refer to contexts has to be simple and easy to document,
otherwise we won't capture it.  My experience with the CFEngine
cognitive agent taught some valuable lessons here. CFEngine
intuitively did several things right. What CFEngine did was to use
"smart sensors" to characterize its environment every time is woke up:
to ask, where am I, who am I? What am I supposed to be doing?

Using this approach in as simple a way as possible, let's define context
as a kind of *scope of experience*. Programmers understand scope as an
bounded environment in which certain variables are available and others
are hidden. It's like a chapter in a book, or a separate document.
So, we can model as conext in two parts for simplicity: 

<pre>
  ( scene   , thoughts and sensory impressions ... )
  ( chapter , environmental fragments   .... )
</pre>

Chapters or scenes are non-overlapping collections of information, events, etc, Christmas 2012.
Sensory fragments are potentially shared or overlapping aspects on an experience, e.g.
it was Monday, or the weather was hot. Other parts of our context come from our train of thought
at the time of the event in question. When we're trying to recall something, e.g. the word in
French for baguette, we want to label that knowledge with those  sensory cues: at the restaurant,
in the supermarche, etc.

* A chapter or scene is partitioned by the exterior physical world of space and time.
* A train of thought is partitioned by the interior virtual world of our imagined or remembered space and time.

## Labelling scenes

* We use the ` - chapter name` syntax to name a chapter in N4L.
* We use the `:: ... ::` syntax in N4L to label context.

When we record information, we list the characters in the scene as "nodes" and their relationships
and meanings through the "links" between them. But, we also collect several of those references under
a subheading of the chapter that lists the fragments of thought we would want to use to remember
that information: the where, when, what, how, and even why!

<pre>

 :: random thoughts ::

   We're using too much cloud CPU (example of) economic priorities

 :: Monday standup meeting, discussion about Easter holiday, urgent work ::

   Sally reported a problem with ssh keys (possibly about)  security
                                  "       (may cause bug)   can't access the dashboard


</pre>
In a more complex scenario, we might use context is a number of ways. What matters is that
it is intuitive to us, because it represents the way we think. Context is very personal.
What helps you to remember will not necessarily be what helps others to remember. This is
why notes and knowledge are not easy to share. Take this example:
<pre>

- cluedo: Forensic map of a Murder Most Horrid

 // invariants and dependencies define first, as we may need to refer to them

 #######################

 :: Dramatis personae ::

 scarlett (id) Miss Scarlett, The Woman in Red, New York socialite.
 plumb    (id) Professor Plumb, University of Oxford, Lincoln College.
 dibbly   (id) Florist working at the Summertown flower shop
          (also called) dibbly womble
 martin   (id) possible boyfriend, unknown

 doorman (id) Fabian Merryweather, former Army officer, 24 Summertown Road

 car (id)   Black Ford Cortina 1972, vehicle number 1234
  "  (role) the get away vehicle

 #######################

 :: locations, places :: 

 library
 Covent Garden Pub
 24 Summertown Road


 #######################

 :: key events ::

@party  party hosted by >scarlett in London flat, evening of the 23rd March
    "   (note) This was a month before >scarlett dated >martin

@scarlettarrives  entering the library on Monday at 10am
@doormanarrives   doorman arrived at library for work on Monday 7 am

@murder  Wednesday April 1st, around 11 am in the morning
</pre>

## Looking it up later

When we come to look up these facts, we use context to disambiguate.
What we remember may be limited, so we want to match the fragments we know:
We don't want to look up a road, something to do with summer. A person, scarlet or something like that.

<pre>
searchN4L  summer in context place
searchN4L  scarlet in the context of person
searchN4L  some car in chapter cluedo
</pre>

## Socializing knowledge

Wikis and databases are potentially places where knowledge goes to die. If we drop information
into a dark place where no one goes, it won't be found intentionally. 

To make knowledge effective, we need to generate the intent to find it and to use it.

The way individual knowledge becomes shared knowledge is through socializing. We have to talk
about it, repeat it, share it, and use it regularly. This is the opposite of thinking about
search keys. What this should tell us is that databases and wikis are not places for
shared knowledge. They are only usable by the particular group who happens to know the context
for looking up the information.

* We need to know that the information is in a datastore before trying to find it.
* We need to have a good idea about what we should look for, i.e. the right context.

Recall the project slogan:

"*It's not knowledge if you don't know it.*"