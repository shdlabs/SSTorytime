
# searchN4L

This is an experimental tool for querying the database. The details
are likely to change in the near future as the software is tested in use.

e.g. try the examples
<pre>
$ cd examples
$ ../src/N4L-db -wipe -u chinese*n4l doors.n4l Mary.n4l brains.n4l

$ ./searchN4L  -chapter multi start 
--------------------------------------------------
Looking for relevant nodes by start
--------------------------------------------------
Search separately by start,..
XXX select NPtr from Node where S LIKE '%start%' AND chap LIKE '%multi%'
found 1 possible relevant nodes:

#1 (search start => start)
--------------------------------------------
 +leads to: 'start'     in chapter notes on chinese,multi slit interference
 +leads to: 'door'      in chapter multi slit interference
 +leads to: 'port'      in chapter multi slit interference
 +leads to: 'hole'      in chapter multi slit interference
 +leads to: 'gate'      in chapter multi slit interference
 +leads to: 'passage'   in chapter multi slit interference
 +leads to: 'road'      in chapter notes on chinese,multi slit interference
 +leads to: 'river'     in chapter notes on chinese,multi slit interference
 +leads to: 'tram'      in chapter multi slit interference
 +leads to: 'bike'      in chapter multi slit interference
 +leads to: 'target 1'  in chapter multi slit interference
 +leads to: 'target 2'  in chapter multi slit interference
 +leads to: 'target 3'  in chapter multi slit interference

  Story:1: start  -(leads to)->  door  -(leads to)->  passage...
  Story:2: start  -(leads to)->  door  -(leads to)->  road...
  Story:3: start  -(leads to)->  door  -(leads to)->  river...
  Story:4: start  -(leads to)->  port  -(leads to)->  river...
  Story:5: start  -(leads to)->  port  -(leads to)->  tram...
  Story:6: start  -(leads to)->  hole  -(leads to)->  tram...
  Story:7: start  -(leads to)->  gate  -(leads to)->  tram...
  Story:8: start  -(leads to)->  gate  -(leads to)->  bike...

 -comes from: 'start'   in chapter notes on chinese,multi slit interference
 -comes from: '开始'    in chapter notes on chinese
 -comes from: 'Kāishǐ'  in chapter notes on chinese

  Story:1: start  -(english for hanzi)->  开始  -(hanzi for pinyin)->  Kāishǐ...


(No relevant matroid patterns matching by arrow)

Check for story paths of length 3
No stories

</pre>

And

<pre>
 ./searchN4L -chapter chinese tiger
--------------------------------------------------
Looking for relevant nodes by tiger
--------------------------------------------------
Search separately by tiger,..
XXX select NPtr from Node where S LIKE '%tiger%' AND chap LIKE '%chinese%'
found 1 possible relevant nodes:

#1 (search tiger => two tigers, two tigers)
--------------------------------------------
 -comes from: 'two tigers, two tigers'  in chapter notes on chinese
 -comes from: '两只老虎, 两只老虎'      in chapter notes on chinese
 -comes from: 'Liǎng zhī lǎohǔ, liǎng zhī lǎohǔ'        in chapter notes on chinese

  Story:1: two tigers, two tigers  -(english for hanzi)->  两只老虎, 两只老虎  -(hanzi for pinyin)->  Liǎng zhī lǎohǔ, liǎng zhī lǎohǔ...


(No relevant matroid patterns matching by arrow)

Check for story paths of length 3
No stories
</pre>