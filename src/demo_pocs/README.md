
# Writing the SQL queries

An important part of implementing graph structures is figuring out the algorithms
expressed as SQL for parsing the data structres.

Ordinary search results from a `select *...` are called "rows".
They are instances of tables, just as variables in Go are instances of data types.
We can call array elements with a row
of table instance "columns".

## Basics

Consider the table:
<pre>
create table Entity
   (
   name text, 
   hasfriend text[], 
   employs text[], 
   primary key(name)
   );
</pre>
Which can be viewed from the shell tool psql:
<pre>
mark% psql newdb

newdb=# \d Entity
                Table "public.entity"
  Column   |  Type  | Collation | Nullable | Default 
-----------+--------+-----------+----------+---------
 name      | text   |           |          | 
 hasfriend | text[] |           |          | 
 employs   | text[] |           |          | 

</pre>
We can add test data and select single rows:
<pre>
newdb=# select * from Entity where name='Mark';
 name |                                     hasfriend                                     | employs 
------+-----------------------------------------------------------------------------------+---------
 Mark | {Silvy,Mandy,Brent,Zhao,Doug,Tore,Joyce,Mike,Carol,Ali,Matt,Bjørn,Tamar,Kat,Hans} | 
(1 row)
</pre>
As we see, the member `hasfriend` has columns.

This association of rows and columns with variables in Go is somewhat confusing,
as both directions would typically be though of as arrays.

## Comparing SQL and Go

A SQL `select * from TABLE` statement is roughly analogous to a `for variable := range table` statement is Go.

Where the two arrays can contain lists of associates. Starting from a node with name="mark"
we aggregate friends and employees as one into a a temporary list called templist:


The `unnest()` is postgres function splits up an array of columns into rows, something like
this pseudo-code:
<pre>
func unnest(array []table) {

   for col := range array {
       print array[col]
   }
}
</pre>

## Nested loops, recursion, and function-like behaviour

So-called recursive queries are something like nested for loops in postres SQL.
Consider the following example used in `postgres_ball.go`
<pre>
WITH RECURSIVE templist (name,friend,radius)
AS (
    SELECT name,unnest(hasfriend), 1 FROM entity WHERE name='Mark'
    UNION
     SELECT e.name,unnest(e.hasfriend),radius+1 FROM entity e JOIN templist ON e.name = friend where radius < 2
)
SELECT DISTINCT friend FROM templist;
</pre>
The syntax is idiosyncratic, yet corresponds to something like this:

<pre>
-------------------------------------------------------------------
1. WITH RECURSIVE templist (name,friends,radius)
2. AS (
3.     -- anchor member FOR (name,friends,radius) :=
4.     SELECT name,unnest(hasfriend), 1 FROM entity WHERE name='Mark'
5.     UNION
6.    -- recursive term
7.     SELECT e.name,unnest(e.hasfriend),radius+1 FROM entity e JOIN templist ON e.name = friends where radius < 2
8. )
9. SELECT DISTINCT friends FROM templist;
-------------------------------------------------------------------
</pre>

The recursive statement defines a re-entrant quasi-function object, with formal parameters
called name, x, and radius. It works something like a for loop

* We are searching a table indexed by name so effectively `Entity` is sort of an array of structs,
where the `name` member doubles as the array index. Think of how Go represents linked lists as array slices.

* `templist` is a TABLE of `(name,friend,radius)` whose types are defined in the initializer


More fully, since ` Entity` is a many element array of rows,
in struct form, `templist` is a triplet consisting of a name, a list, and an integer.

<pre>
templist.name = Entity['Mark'].name   // name is the index variable/primary key so this looks redundant
templist.{friends} = Entity['Mark'].{hasfriend}
templist.radius = 1

for tmp := range templist.{friends}
   {
   for e = range Entity[]
       {
       if Entity[e].name = templist.{friends}[tmp]
          {
          join_delta = (Entity[e].name,Entity[e].{hasfriend}[Entity[e].name],radius+1)
          templist += join_delta
          }
       }  
    }

 for t = distinct range templist
    {
    print templist[t].{friends}
    }

</pre>
In SQL language,
* The intial values are set by the "anchor line" (line 4)
`SELECT name,unnest(hasfriend), 1 FROM entity WHERE name='Mark'`
Notice that this is the line that contains the "boundary conditions" or invariant constant 'Mark' and unit value 1.

* The templist is not the actually the name of a function, but its resulting output formal parameter, i.e. a selection of rows that is appended by a join on each iteration.

## Remark - SQL aliasing (missing AS)

In a select statement (although it's perhaps not the first thing one learns), one can define an alias for the
items and columns by suffixing the search part with alias identifiers, which can then
be used like struct variables to refer to column members  e.g. in line 7
<pre>
SELECT e.name,unnest(e.hasfriend),radius+1 FROM Entity e JOIN templist ON e.name = friends where radius < 2
</pre>
Notice how "e" is used after the table-name or datatype Entity as the current instance of that value
in order to distinguish it from friends, whose source is the iteration algorithm.

This is essentially the keyword AS in SQL that's implicit. One is free to omit it though that leads to some confusion.

## Coaxing SQL to do what it doesn't want to

As Mick Jagger taught us, you caint always get wha' yooo wan'. This is so with domain specific
languages that are always designed to make limited cases easy. The downfall is that they can laso make
difficult cases hard. SQL is a language that has been grown unwillingly to handle things that were never intended,
but we're stuck with the imperfect results.

It would be nice to prune a list of layers in a graph explosion, from some starting point, so that each user
only appeared once. This could save much transfer time for large datasets. However, the limitations of recursive
or iterative evaluation make this impossible.
<pre>

WITH RECURSIVE cone (name,member,past,depth)
AS (
    SELECT name,unnest(hasfriend), Array['Mark']::text[], 1 FROM entity WHERE name='Mark'
    UNION
   SELECT e.name,unnest(e.hasfriend),e.name||past,depth+1 FROM entity e JOIN cone ON e.name = member where (depth < 7 and not member = ANY(past))
)
SELECT member,depth,past FROM cone order by depth ;

</pre>
This query results in something like this:
<pre>
       member       | depth |                      past                       
--------------------+-------+-------------------------------------------------
 Silvy              |     1 | {Mark}
 Mandy              |     1 | {Mark}
 Brent              |     1 | {Mark}
 Zhao               |     2 | {Mark,Mandy}
 Doug               |     2 | {Mark,Mandy}
 Tore               |     2 | {Mark,Mandy}
 Joyce              |     2 | {Mark,Mandy}
 Mike               |     2 | {Mark,Mandy}
 Carol              |     2 | {Mark,Mandy}
 Ali                |     2 | {Mark,Mandy}
 Matt               |     2 | {Mark,Mandy}
 Bjørn              |     2 | {Mark,Mandy}
 Tamar              |     2 | {Mark,Mandy}
 Kat                |     2 | {Mark,Mandy}
 Hans               |     2 | {Mark,Mandy}
 Mark               |     3 | {Mark,Mandy,Mike}
 Jane1              |     3 | {Mark,Mandy,Mike}
 Jane2              |     3 | {Mark,Mandy,Mike}
 Jan                |     3 | {Mark,Mandy,Mike}
 Alfie              |     3 | {Mark,Mandy,Mike}
 Jungi              |     3 | {Mark,Mandy,Mike}
 Peter              |     3 | {Mark,Mandy,Mike}
 Paul               |     3 | {Mark,Mandy,Mike}
 Adam               |     4 | {Mark,Mandy,Mike,Jan}
 Jane1              |     4 | {Mark,Mandy,Mike,Jan}
 Jane               |     4 | {Mark,Mandy,Mike,Jan}
 Company of Friends |     5 | {Mark,Mandy,Mike,Jan,Adam}
 Paul               |     5 | {Mark,Mandy,Mike,Jan,Adam}
 Matt               |     5 | {Mark,Mandy,Mike,Jan,Adam}
 Billie             |     5 | {Mark,Mandy,Mike,Jan,Adam}
 Chirpy Cheep Cheep |     5 | {Mark,Mandy,Mike,Jan,Adam}
 Taylor Swallow     |     5 | {Mark,Mandy,Mike,Jan,Adam}
 Matt               |     6 | {Mark,Mandy,Mike,Jan,Adam,"Company of Friends"}
 Jane1              |     6 | {Mark,Mandy,Mike,Jan,Adam,"Company of Friends"}
(34 rows)
</pre>
What we'd like is for the left column to avoid looping around.