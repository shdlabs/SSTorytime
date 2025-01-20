<!--
 SSTorytime - a ChiTek-i project by Mark Burgess

 Semantic Spacetime Story graph database library over postgresql (SSTorytime)
 This is an NLnet sponsored project, See https://nlnet.nl/project/SmartSemanticDataLookup/

-->

# SSTorytime

Semantic knowledge representations have not evolved since the Semantic
Web was proposed during the 1990s. Modern graph databases offer new
possibilities for knowledge representation, but the methods are poorly
developed and require the use of specialized query languages and
clumsy outdated formats. This project aims to make semantic maps easy
for general use, using standard SQL databases and modern lightweight
data formats. A user workflow starts from a simple note-taking
language, then ingesting into a database using a graph model based on
the causal semantic spacetime model, to the use of a simple web
application for supporting graph searches and data presentation. The
aim is to make a generally useful library for incorporating into other
applications, or running as a standalone notebook service.

 The project's own research website: http://markburgess.org/spacetime.html

## Notes and thoughts

Thoughts that we have, things that we learn are easily lost and
forgotten.  Writing down knowledge (intructions, experiences, etc) is
useless if no one reads what is written regularly. This is why most
wikis and documentations fail. Writing notes for occasional use once a
year, or only in case of emergency is pointless unless it's rehearsed.
Wikipedia succeeds because it's many to many interaction, which keeps
the pulse of interaction alive because someone is always looking at
the information. Even Wikipedia information is not knowledge unless
it's in someone's head. To the causal browser, it is merely
information, even hearsay.

Scaling knowledge is like scaling the activity in a city. Time is driven by parallel events arriving. The more people, the faster the pace or pulse.


   THE STATUS OF THIS PROJECT IS CURRENTLY R&D. DETAILS MAY CHANGE.


The toolset consistent of several components.

N4L - a unicode text based note taking language for jotting down notes in a way
      that can be parsed and loaded into a semantic database. 
      N = note, 4 = 4 semantic type, L = language
      N4L = notes for loading

s