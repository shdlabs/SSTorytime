
# searchN4L

This is a tool for querying the database. This is redesigned to avoid having to use command line options.
The only command line option is `-v` for verbose output, which is helpful for debugging.

## Commands

The tool recognizes a number of words:

"on"
"for"
"about"

"note"
"page"

"path"
"seq"
"from"
"to"
"ctx"

"context"
"as"
"chapter"
"section"
"in"
"arrow"
"limit"
"depth"
"range"
"distance"


Using the pre-loaded examples, you can try:

## Search for nodes and their close neighbour orbits matching a name
<pre>
$ ./searchN4L Mark
------------------------------------------------------------------

0: supermarket
      -    (english has hanzi) - 超市
      -    (hanzi has pinyin) - chāoshì  .. food, shopping


1: I'm looking for snacks in the supermarket
      -    (english has hanzi) - 我在找超市的零食去
      -    (hanzi has pinyin) - wǒ zài zhǎo chāoshì de língshí qù  .. food, shopping


2: the supermarket is on the basement floor
      -    (english has hanzi) - 超市在地下一层
      -    (hanzi has pinyin) - chāoshì zài dìxià yī céng  .. configuration, containment, directions
      down, examples, location, orientation, position, up


3: uses a language that looks a lot like SQL but is markedly different - beware!
      -    (has name) - PLpgSQL
      -    (is a note or remark about) - stored procedures/functions in postgres

</pre>

## Searching when you can't type unicode accents

If you can only get English characters on your keyboard, you can still search for accented
words by placing parentheses around them "(...)":
<pre>
% ./searchN4L  "(fangzi)" chapter "chinese" 
------------------------------------------------------------------

0: fángzi
      -    (pinyin has hanzi) - 房子
           -    (hanzi has english) - house  .. at home, domestic


1: fángzǐ de fùjìn yǒu hěnduō piàoliang de huā
      -    (pinyin has hanzi) - 房子的附近有很多漂亮的花
           -    (hanzi has english) - there are many beautiful flowers near the house  .. area, environment
      neighbourhood, region


2: wǒ de chē zài fángzǐ pángbiān
      -    (pinyin has hanzi) - 我的车在房子旁边
           -    (hanzi has english) - my car is next to the house  .. configuration, directions,
     from, layout, position, toward
</pre>

## Searching by direct NodePtr references

If you know about the database internals, you can look up node pointers directly
as long as you quote the parentheses for the shell.
Notice how the indentation shows you the distance from the starting node.
<pre>
./searchN4L "(1,1)"
------------------------------------------------------------------

0: door
      -    (leads to) - passage
           -    (leads to) - target 1  .. connectivity, path example, physics
      -    (leads to) - road
           -    (english has hanzi) - 路  .. browsing, caution, walking
           -    (leads to) - target 2  .. connectivity, path example, physics
      -    (leads to) - river
           -    (english has hanzi) - 河  .. nature
           -    (english has hanzi) - 江  .. nature
           -    (leads to) - target 3  .. connectivity, path example, physics
      -    (comes from / arriving from) - start
           -    (english has hanzi) - 开始  .. common verbs, doing, look, see, using, wanting

</pre>

## Searching for related contexts

Context strings are clustered into groups. If you don't remember, you can search:
<pre>
% ./searchN4L context restaurant
  0. "buildings"
  1. "come"
  2. "come in"
  3. "cooking"
  4. "eating"
  5. "enter"
  6. "restaurant"
  7. "rooms"
  8. "transport"
  9. "vehicles"
</pre>


## Searching for arrows

You can look up arrow definitions too, by name, number, or spacetime type.
The output format is `arrowptr, sttype, long name`:
<pre>
$ ./searchN4L arrow ph,pe
192. (3) ph -> pinyin has hanzi
190. (3) pe -> pinyin has english

$ ./searchN4L arrow 125
125. (-2) during -> happened during

$ ./searchN4L arrow -2
  9. (-2) in -> is in
 11. (-2) is an emphatic proto-concept in -> is emph in
 13. (-2) is mentioned in -> ismentin
 89. (-2) part -> is component of
 91. (-2) inset -> is part of set
 93. (-2) subby -> is subsumed by
 95. (-2) isencl -> is enclosed by
 97. (-2) swby -> is swallowed by
 99. (-2) pt -> is part of
101. (-2) wordin -> is a word used in

</pre>

## Searching for paths

You can search for paths from one location to another:
<pre>
 ./searchN4L from start to "target 1"
------------------------------------------------------------------

     - story path:  start  -(leads to)->  door  -(leads to)->  passage  -(debug)->  target 1
     -  [ Link STTypes: -(+leads to)->  -(+leads to)->  -(+leads to)-> . ]
</pre>
The default path length limtis to 5 hops. There might be longer paths, so you can add a depth
to force a larger search:

<pre>
$ ./searchN4L paths from a7 to i6 depth 16
</pre>
or simply
<pre>
$ ./searchN4L a7 to i6 depth 16
------------------------------------------------------------------

     - story path:  maze_a7  -(forwards)->  maze_b7  -(forwards)->  maze_b6  -(forwards)->  maze_c6
      -(forwards)->  maze_c5  -(forwards)->  maze_b5  -(forwards)->  maze_b4
      -(forwards)->  maze_a4  -(forwards)->  maze_a3  -(forwards)->  maze_b3
      -(forwards)->  maze_c3  -(forwards)->  maze_d3  -(forwards)->  maze_d2
      -(forwards)->  maze_e2  -(forwards)->  maze_e3  -(debug)->  maze_f3
      -(debug)->  maze_f4  -(debug)->  maze_e4  -(debug)->  maze_e5
      -(debug)->  maze_f5  -(debug)->  maze_f6  -(debug)->  maze_g6
      -(debug)->  maze_g5  -(debug)->  maze_g4  -(debug)->  maze_h4
      -(debug)->  maze_h5  -(debug)->  maze_h6  -(debug)->  maze_i6
     -  [ Link STTypes: -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)->  -(+leads to)-> . ]

</pre>

## Searching in note form

Sometimes you want to see your full notes, the way you ordered them:
<pre>
$ ./searchN4L notes brain

---------------------------------------------

Title: neuroscience brain
Context: oscillations waves 
---------------------------------------------


alpha waves (has frequency) 5-15 Hz 
alpha waves (is characterized by) very relaxed, light or passive attention 
beta waves (has frequency) 12-35 Hz 
beta waves (is characterized by) medium attention,anxiety dominant, active, external attention 
gamma waves (has frequency) 32-100 Hz 
gamma waves (note/remark) 40 Hz of special interest 
gamma waves (is characterized by) concentration 
gamma waves (occurs in) premotor cortex 
gamma waves (occurs in) parietal cortex 
gamma waves (occurs in) temporal cortex 
gamma waves (occurs in) frontal cortex 

</pre>
