

# An API for interacting with the SST graph

Once data have been entered into a SSToryline database, we want to be able to extract it again.
It's possible to create tools for this, but ultimately any set of tools will tend to limit the user.
A user's imagination should be the only limit. 

Many specialized graph databases offer graph languages, but they
expose an important problem with Domain Specific Languages, which is
that by trying to make simple things easy, they make less-simple
things hard. The most well known standard for data (Structured Query
Language, or SQL) is itself a Domain Specific Language with exactly
these problems. However, in Open Source Postgres there are plenty of
extensions that make it possible to overcome the limitations of SQL.

*This project uses Postgres because of that compromise between a well known
standard, and a battle-tested and extensible data platform.*

You will find examples of using Go(lang) code to write custom scripts
that interact with the database through the Go API [here](https://github.com/markburgess/SSTorytime/tree/main/src/demo_pocs).

## Basic queries from SQL

Using perfectly standard SQL, you can interrogate the database established by N4L or the low level API
functions.

### Tables

* To show the different tables:
<pre>
$ psql newdb

newdb=# \dt
              List of relations
 Schema |      Name      | Type  |   Owner    
--------+----------------+-------+------------
 public | arrowdirectory | table | sstoryline
 public | arrowinverses  | table | sstoryline
 public | node           | table | sstoryline
 public | nodearrownode  | table | sstoryline
(4 rows)

</pre>
* To query these, we look at the members:
<pre>
newdb=# \d node
                Table "public.node"
 Column |  Type   | Collation | Nullable | Default 
--------+---------+-----------+----------+---------
 nptr   | nodeptr |           |          | 
 l      | integer |           |          | 
 s      | text    |           |          | 
 chap   | text    |           |          | 
 im3    | link[]  |           |          | 
 im2    | link[]  |           |          | 
 im1    | link[]  |           |          | 
 in0    | link[]  |           |          | 
 il1    | link[]  |           |          | 
 ic2    | link[]  |           |          | 
 ie3    | link[]  |           |          | 
Indexes:
    "node_chan_l_s_idx" btree (((nptr).chan), l, s)

</pre>

### Nodes

Now try:
<pre>
newdb=# select S,chap from Node limit 10;
     s      |       chap       
------------+------------------
 please     | notes on chinese
 yes        | notes on chinese
 请          | notes on chinese
 qǐng       | notes on chinese
 thankyou   | notes on chinese
 Méiyǒu     | notes on chinese
 谢谢        | notes on chinese
 xièxiè     | notes on chinese
 是的        | notes on chinese
 请在这里等    | notes on chinese
(10 rows)

</pre>

* An alternative view of relations is provided by NodeArrowNode:
<pre>
newdb=# select *  from NodeArrowNode LIMIT 10;
 nfrom | sttype | arr | wgt |              ctx              |   nto   
-------+--------+-----+-----+-------------------------------+---------
 (1,0) |     -1 |  69 |   1 | {please,"thank you",thankyou} | (1,1)
 (1,1) |     -1 |  67 |   1 | {thankyou,please,"thank you"} | (1,2)
 (1,1) |      1 |  68 |   1 | {thankyou,please,"thank you"} | (1,0)
 (1,1) |      1 |  68 |   1 | {news,online}                 | (2,291)
 (1,2) |      1 |  66 |   1 | {thankyou,please,"thank you"} | (1,1)
 (1,3) |     -1 |  69 |   1 | {please,"thank you",thankyou} | (1,4)
 (1,4) |     -1 |  67 |   1 | {please,"thank you",thankyou} | (1,5)
 (1,4) |      1 |  68 |   1 | {please,"thank you",thankyou} | (1,3)
 (1,5) |      1 |  66 |   1 | {please,"thank you",thankyou} | (1,4)
 (1,6) |     -1 |  67 |   1 | {please,"thank you",thankyou} | (4,0)
(10 rows)

</pre>

Notice how nodes (`nfrom`,`nto`,`nptr? ) and arrows (`arr`) are represented by pointer references
that are integers. When working with the graph, we often don't need to know the names
of things, we can get away with deferring the lookup of the actual data until we find what we're
looking for. That information can be cached so as to minimize the data transferred over the wire.

<pre>
newdb=# select S from Node where NPtr=(1,5);
   s    
--------
 xièxiè
(1 row)

</pre>

## Links and Arrows

A link is a composite relation that involves an arrow (pointer), a context,
and a destination node. Links are anchored to their origin nodes in the `Node` table
in the six columns `im3`, `im2`, `im1`, `in0`, `il1`, `ic2`, `ie3`.  
To find the links of type `Leads to':
<pre>
newdb=# select Il1 from Node where NPtr=(1,5);
                                       il1                                        
----------------------------------------------------------------------------------
 {"(66,1,\"{ \"\"please\"\", \"\"thank you\"\", \"\"thankyou\"\" }\",\"(1,4)\")"}
(1 row)

</pre>

Arrows are defined for each arrow pointer in the arrow directory:

<pre>
newdb=# select * from arrowdirectory limit 10;
 staindex |         long         | short | arrptr 
----------+----------------------+-------+--------
        4 | leads to             | lt    |      0
        2 | arriving from        | af    |      1
        4 | forwards             | fwd   |      2
        2 | backwards            | bwd   |      3
        4 | affects              | aff   |      4
        2 | affected by          | baff  |      5
        4 | causes               | cf    |      6
        2 | is caused by         | cb    |      7
        4 | used for             | for   |      8
        2 | is a possible use of | use   |      9
(10 rows)

</pre>



