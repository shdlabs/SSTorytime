<!--
 SSTorytime - a ChiTek-i project by Mark Burgess

 Semantic Spacetime Story graph database library over postgresql (SSTorytime)
 This is an NLnet sponsored project, See https://nlnet.nl/project/SmartSemanticDataLookup/

-->

# SSTorytime

* This is a work in progress for 2025, as part of an [NLnet project](https://nlnet.nl/project/SmartSemanticDataLookup/) *

This project aims to make knowledge capture, querying, and dissemination easy for humans
for general use. AI can only capture knowledge from humans, so even if we
want to use AI, we'd better get the knowledge representations right.

The roots of this project go back almost 20 years for me, when I was working in configuration
management (with the CFEngine project) and realized that the main problem there was not
fixing machinery but rather understanding the monster you've created! Knowledge Management
was built into CFEngine 3, but later removed again when `the market wasn't ready'. Over those
20 years, I've studied and learned how to approach the problem in better ways. I've implemented
the concepts using a variaety of technologies, including my own. In this latest version, I'm
combining those lessons to make a version that builds on standard Postgres.

Knowledge capture requires tools for collecting factual notes, data
relationships, and structures for representing and organizing them, so
that they can be found easily. Many mistakes have been made around
this in the past, trying to force discipline onto people at the wrong
moment and neglecting to do so when it matters. As a lifelong teacher
of ideas, I've studied this problem and think we can do better than
what we have today.

* [The Mission and Approach](docs/approach.md)
* [Basics of Knowledge Engineering](docs/KnowledgeAndLearning.md)
* [SSTorytelling](docs/Storytelling.md)
* [N4L - Notes For Learning/Loading](docs/N4L.md)

One of the goals of this project is to better understand what we call "reasoning".
One used to think of reasoning, philosophically, as logical argumentation. As computers
entered society we replaced this idea with actual first order logic. But, if you ask
a teacher (and if we've learned anything from the Artificial Intelligence journey)
then we realize that the way humans arrive at conclusions has a more complicated
relationship to logic. We first decide emotionally, narrowly or expansively depending
on our context, and then we try to formulate a "logical" story to support that.
This is why we strive to study the role of stories in learning and understanding for this project.

## The tools

The tool-set consistent of several components, starting with:

* [N4L](docs/N4L.md) - a Unicode text based note taking language for jotting down notes in a way
        that can be parsed and loaded into a semantic database. 
        N = note, 4 = 4 semantic type, L = language
        N4L = notes for loading


