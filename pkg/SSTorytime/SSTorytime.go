//**************************************************************
//
// An interface for postgres for graph analytics and semantics
//
//**************************************************************

package SSTorytime

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	//"sort"

	_ "github.com/lib/pq"

)

//**************************************************************
// Errors
//**************************************************************

const (
	ERR_ST_OUT_OF_BOUNDS="Link STtype is out of bounds - "
	ERR_ILLEGAL_LINK_CLASS="ILLEGAL LINK CLASS"
)

//**************************************************************

const (
	NEAR = 0
	LEADSTO = 1   // +/-
	CONTAINS = 2  // +/-
	EXPRESS = 3   // +/-

	ST_ZERO = EXPRESS
	ST_TOP = ST_ZERO + EXPRESS + 1

	// For the SQL table, as 2d arrays not good

	I_MEXPR = "Im3"
	I_MCONT = "Im2"
	I_MLEAD = "Im1"
	I_NEAR  = "In0"
	I_PLEAD = "Il1"
	I_PCONT = "Ic2"
	I_PEXPR = "Ie3"

	// For separating text types

	N1GRAM = 1
	N2GRAM = 2
	N3GRAM = 3
	LT128 = 4
	LT1024 = 5
	GT1024 = 6
)

//**************************************************************

type Node struct { // essentially the incidence matrix

	L         int     // length of text string
	S         string  // text string itself

	Chap      string  // section/chapter name in which this was added
	SizeClass int     // the string class: N1-N3, LT128, etc for separating types

	NPtr      NodePtr // Pointer to self index

	I [ST_TOP][]Link  // link incidence list, by arrow type
  	                  // NOTE: carefully how offsets represent negative SSTtypes
}

//**************************************************************

type Link struct {  // A link is a type of arrow, with context
                    // and maybe with a weightfor package math
	Arr ArrowPtr         // type of arrow, presorted
	Wgt float64          // numerical weight of this link
	Ctx []string         // context for this pathway
	Dst NodePtr          // adjacent event/item/node
}

const NODEPTR_TYPE = "CREATE TYPE NodePtr AS  " +
	"(                    " +
	"Chan     int,        " +
	"CPtr     int         " +
	")"

const LINK_TYPE = "CREATE TYPE Link AS  " +
	"(                    " +
	"Arr      int,        " +
	"Wgt      real,       " +
	"Ctx      text,       " +
	"Dst      NodePtr     " +
	")"

const NODE_TABLE = "CREATE TABLE IF NOT EXISTS Node " +
	"( " +
	"NPtr      NodePtr,         " +
	"L         int,            " +
	"S         text,           " +
	"Chap      text,           " +
	I_MEXPR+"  Link[],         " + // Im3
	I_MCONT+"  Link[],         " + // Im2
	I_MLEAD+"  Link[],         " + // Im1
	I_NEAR +"  Link[],         " + // In0
	I_PLEAD+"  Link[],         " + // Il1
	I_PCONT+"  Link[],         " + // Ic2
	I_PEXPR+"  Link[]          " + // Ie3
	")"

const LINK_TABLE = "CREATE TABLE IF NOT EXISTS NodeLinkNode " +
	"( " +
	"NFrom    NodePtr, " +
	"Lnk      Link,    " +
	"NTo      NodePtr " +
	")"

//**************************************************************

type NodeDirectory struct {

	// Power law n-gram frequencies

	N1grams map[string]ClassedNodePtr
	N1directory []Node
	N1_top ClassedNodePtr

	N2grams map[string]ClassedNodePtr
	N2directory []Node
	N2_top ClassedNodePtr

	N3grams map[string]ClassedNodePtr
	N3directory []Node
	N3_top ClassedNodePtr

	// Use linear search on these exp fewer long strings

	LT128 []Node
	LT128_top ClassedNodePtr
	LT1024 []Node
	LT1024_top ClassedNodePtr
	GT1024 []Node
	GT1024_top ClassedNodePtr
}

// This is better represented as separate tables in SQL, one for each class


//**************************************************************

type NodePtr struct {

	CPtr  ClassedNodePtr // index of within name class lane
	Class int            // Text size-class
}

type ClassedNodePtr int  // Internal pointer type of size-classified text

//**************************************************************

type ArrowDirectory struct {

	STtype  int
	Long    string
	Short   string
	Ptr     ArrowPtr
}

type ArrowPtr int // ArrowDirectory index


const ARROW_DIRECTORY_TABLE = "CREATE TABLE IF NOT EXISTS Arrow_Directory " +
	"(    " +
	"STtype int,             " +
	"Long text,              " +
	"Short text,             " +
	"ArrPtr int primary key  " +
	")"

const ARROW_INVERSES_TABLE = "CREATE TABLE IF NOT EXISTS Arrow_Inverses " +
	"(    " +
	"Plus int,  " +
	"Minus int,  " +
	"Primary Key(Plus,Minus)," +
	")"

//******************************************************************
// LIBRARY
//******************************************************************

type PoSST struct {

   DB *sql.DB
}

//******************************************************************

func Open() PoSST {

	var ctx PoSST
	var err error

	const (
		host     = "localhost"
		port     = 5432
		user     = "sstoryline"
		password = "sst_1234"
		dbname   = "newdb"
	)

        connStr := "user="+user+" dbname="+dbname+" password="+password+" sslmode=disable"

        ctx.DB, err = sql.Open("postgres", connStr)

	if err != nil {
	   	fmt.Println("Error connecting to the database: ", err)
		os.Exit(-1)
	}
	
	err = ctx.DB.Ping()
	
	if err != nil {
		fmt.Println("Error pinging the database: ", err)
		os.Exit(-1)
	}

	Configure(ctx)

	return ctx
}

// **************************************************************************

func Configure(ctx PoSST) {

	// Tmp reset

	ctx.DB.Query("drop table Node")
	ctx.DB.Query("drop table NodeLinkNode")
	ctx.DB.Query("drop type Link")
	ctx.DB.Query("drop type NodePtr")

	if !CreateType(ctx,NODEPTR_TYPE) {
		fmt.Println("Unable to create type as, ",NODEPTR_TYPE)
		os.Exit(-1)
	}

	if !CreateType(ctx,LINK_TYPE) {
		fmt.Println("Unable to create type as, ",LINK_TYPE)
		os.Exit(-1)
	}

	if !CreateTable(ctx,NODE_TABLE) {
		fmt.Println("Unable to create table as, ",NODE_TABLE)
		os.Exit(-1)
	}

	if !CreateTable(ctx,LINK_TABLE) {
		fmt.Println("Unable to create table as, ",LINK_TABLE)
		os.Exit(-1)
	}

	DefineStoredFunctions(ctx)
}

// **************************************************************************

func Close(ctx PoSST) {
	ctx.DB.Close()
}

// **************************************************************************

func CreateType(ctx PoSST, defn string) bool {

	_,err := ctx.DB.Query(defn)

	if err != nil {
		s := fmt.Sprintln("Failed to create datatype PGLink ",err)
		
		if strings.Contains(s,"already exists") {
			return true
		} else {
			return false
		}
	}

	return true
}

// **************************************************************************

func CreateTable(ctx PoSST,defn string) bool {

	_,err := ctx.DB.Query(defn)
	
	if err != nil {
		s := fmt.Sprintln("Failed to create a table %.10 ...",defn,err)
		
		if strings.Contains(s,"already exists") {
			return true
		} else {
			return false
		}
	}

	return true
}

// **************************************************************************
// Store
// **************************************************************************

func CreateDBNode(ctx PoSST, n Node) Node {

	var qstr string

	// No need to trust the values

        n.L,n.SizeClass = StorageClass(n.S)
	
	cptr := n.NPtr.CPtr

	qstr = fmt.Sprintf("SELECT IdempInsertNode(%d,%d,%d,'%s','%s')",n.L,n.SizeClass,cptr,n.S,n.Chap)
	
	row,err := ctx.DB.Query(qstr)
	
	if err != nil {
		s := fmt.Sprint("Failed to insert",err)
		
		if strings.Contains(s,"duplicate key") {
		} else {
			fmt.Println(s,"FAILED \n",qstr,err)
		}
	}

	var whole string
	var cl,ch int

	for row.Next() {		
		err = row.Scan(&whole)
		fmt.Sscanf(whole,"(%d,%d)",&cl,&ch)
	}

	n.NPtr.Class = cl
	n.NPtr.CPtr = ClassedNodePtr(ch)

	return n
}

// **************************************************************************

func AppendDBLinkToNode(ctx PoSST, n1ptr NodePtr, lnk Link, n2ptr NodePtr, sttype int) bool {

	// Want to make this idempotent, because SQL is not (and not clause)

	if sttype < 0 || sttype > ST_ZERO+EXPRESS {
		fmt.Println(ERR_ST_OUT_OF_BOUNDS,sttype)
		os.Exit(-1)
	}

	lnk.Dst.Class = n2ptr.Class
	lnk.Dst.CPtr = n2ptr.CPtr

	//                       Arr,Wgt,Ctx,  Dst
	linkval := fmt.Sprintf("(%d, %f, %s, (%d,%d)::NodePtr)",lnk.Arr,lnk.Wgt,FormatSQLArray(lnk.Ctx),lnk.Dst.Class,lnk.Dst.CPtr)

	literal := fmt.Sprintf("%s::Link",linkval)

	link_table := STTypeDBChannel(sttype)

	qstr := fmt.Sprintf("UPDATE NODE set %s=array_append(%s,%s) where (NPtr).CPtr = '%d' and (NPtr).Chan = '%d' and (%s is null or not %s = ANY(%s))",
		link_table,
		link_table,
		literal,
		n1ptr.CPtr,
		n1ptr.Class,
		link_table,
		literal,
		link_table)

	_,err := ctx.DB.Query(qstr)

	if err != nil {
		fmt.Println("Failed to append",err,qstr)
	       return false
	}

	return true
}


// **************************************************************************

func DefineStoredFunctions(ctx PoSST) {
	
	// Insert node
	
	qstr := "CREATE OR REPLACE FUNCTION IdempInsertNode(iLi INT, iszchani INT, icptri INT, iSi TEXT, ichapi TEXT)" + 
		"RETURNS TABLE (    " +
			"    ret_cptr INTEGER," +
			"    ret_channel INTEGER" +
			") AS $fn$ " +
			"DECLARE " +
			"BEGIN" +
			"  IF NOT EXISTS (SELECT (NPtr).Chan,(NPtr).CPtr FROM Node WHERE s = iSi) THEN" +
			"     INSERT INTO Node (Nptr.Cptr,L,S,chap,Nptr.Chan) VALUES (icptri, iLi, iSi, ichapi, iszchani);" +
			"  end if;" +
			"      return query SELECT (NPtr).Chan,(NPtr).CPtr FROM Node WHERE s = iSi;" +
			"END ;" +
			"$fn$ LANGUAGE plpgsql;";

		_, err := ctx.DB.Query(qstr)
		
		if err != nil {
			fmt.Println("Error defining postgres function:",qstr,err)
		}

	// Search along a vector for the future cone
	// e.g. select * from GetFwdNeigh('(4,2)',ARRAY[ '(4,6)' ]::NodePtr[]);

	//for sttype := -EXPRESS; sttype <= EXPRESS; sttype++ {

	for sttype := LEADSTO; sttype <= LEADSTO; sttype++ {

	stname := STTypeDBChannel(sttype)
 
		// Get the nearest neighbours as NPtr

		qstr := fmt.Sprintf("CREATE OR REPLACE FUNCTION GetFwdNeigh(start NodePtr,exclude NodePtr[] )\n"+
			"RETURNS NodePtr[] AS $nums$\n" +
			"DECLARE \n" +
			"    neighbours NodePtr[];\n" +
			"    fwdlinks Link[];\n" +
			"    lnk Link;\n" +
			"BEGIN\n" +
			//   Note, you can't select into in the stored function language
			"    SELECT %s INTO fwdlinks FROM Node WHERE Nptr=start;\n"+
			"    IF fwdlinks IS NULL THEN\n" +
			"        RETURN '{}';\n" +
			"    END IF;\n" +
			"    neighbours := ARRAY[]::Link[];\n" +
			"    FOREACH lnk IN ARRAY fwdlinks\n" +
			"    LOOP\n"+
			"      IF exclude is not NULL AND NOT lnk.dst=ANY(exclude) THEN\n" +
			"         neighbours := array_append(neighbours, lnk.dst);\n" +
			"      END IF; \n" +
			"    END LOOP;\n" +
			"    RETURN neighbours; \n" +
			"END ;\n" +
			"$nums$ LANGUAGE plpgsql;\n",stname)

		_, err := ctx.DB.Query(qstr)
		
		if err != nil {
			fmt.Println("Error defining postgres function:",qstr,err)
		}
	}
	
	// Get the forward cone / half-ball as NPtr

	qstr =  "CREATE OR REPLACE FUNCTION GetFwdCone(start NodePtr,maxdepth INTEGER)\n"+
		"RETURNS NodePtr[] AS $nums$\n" +
		"DECLARE \n" +
		"    nextlevel NodePtr[];\n" +
		"    partlevel NodePtr[];\n" +
		"    level NodePtr[] = ARRAY[start]::NodePtr[];\n" +
		"    exclude NodePtr[] = ARRAY['(0,0)']::NodePtr[];\n" +
		"    cone NodePtr[];\n" +
		"    neigh NodePtr;\n" +
		"    frn NodePtr;\n" +
		"    counter int := 0;\n" +

		"BEGIN\n" +

		"LOOP\n" +
		"  EXIT WHEN counter = maxdepth+1;\n" +

		"  IF level IS NULL THEN\n" +
		"     RETURN cone;\n" +
		"  END IF;\n" +

		"  nextlevel := ARRAY[]::NodePtr[];\n" +

		"  FOREACH frn IN ARRAY level "+
		"  LOOP \n"+
		"     nextlevel = array_append(nextlevel,frn);\n" +
		"  END LOOP;\n" +

		"  IF nextlevel IS NULL THEN\n" +
		"     RETURN cone;\n" +
		"  END IF;\n" +

		"  FOREACH neigh IN ARRAY nextlevel LOOP \n"+
		"    IF NOT neigh = ANY(exclude) THEN\n" +
		"      cone = array_append(cone,neigh);\n" +
		"      exclude := array_append(exclude,neigh);\n" +
		"      partlevel := GetFwdNeigh(neigh,exclude);\n" +
		"    END IF;" +
		"    IF partlevel IS NOT NULL THEN\n" +
		"         level = array_cat(level,partlevel);\n"+
		"    END IF;\n" +
		"  END LOOP;\n" +

		// Next, continue, foreach
		"  counter = counter + 1;\n" +
		"END LOOP;\n" +

		"RETURN cone; \n" +
		"END ;\n" +
		"$nums$ LANGUAGE plpgsql;\n"

		_, err = ctx.DB.Query(qstr)
		
		if err != nil {
			fmt.Println("Error defining postgres function:",qstr,err)
		}

		// Drop Functions later?
}

// **************************************************************************
// Retrieve
// **************************************************************************

func GetFwdCone(ctx PoSST, start NodePtr, sttype,depth int) []NodePtr {

	qstr := fmt.Sprintf("select unnest(getfwdcone) from GetFwdCone('(%d,%d)',%d);",start.Class,start.CPtr,depth)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY to GetFwdCone Failed",err)
	}

	var whole string
	var n NodePtr
	var retval []NodePtr

	for row.Next() {		
		err = row.Scan(&whole)
		fmt.Sscanf(whole,"(%d,%d)",&n.Class,&n.CPtr)
		retval = append(retval,n)
	}

	return retval
}

// **************************************************************************
// Tools
// **************************************************************************

func ParseSQLArrayString(whole_array string) []string {

   // array as {"(1,2,3)","(4,5,6)"}

      	var l []string

    	whole_array = strings.Replace(whole_array,"{","",-1)
    	whole_array = strings.Replace(whole_array,"}","",-1)
	whole_array = strings.Replace(whole_array,"\",\"",";",-1)
	whole_array = strings.Replace(whole_array,"\"","",-1)
	
        items := strings.Split(whole_array,";")

	for i := range items {
	    var v string
	    s := strings.TrimSpace(items[i])
	    fmt.Sscanf(s,"%s",&v)
	    l = append(l,v)
	    }

	return l
}

// **************************************************************************

func FormatSQLArray(array []string) string {

        if len(array) == 0 {
	   return ""
        }

	var ret string = "'{ "
	
	for i := 0; i < len(array); i++ {
		ret += fmt.Sprintf("\"%s\"",array[i])
	    if i < len(array)-1 {
	    ret += ", "
	    }
        }

	ret += " }' "

	return ret
}

//**************************************************************

func StorageClass(s string) (int,int) {
	
	var spaces int = 0
	
	var l = len(s)
	
	for i := 0; i < l; i++ {
		
		if s[i] == ' ' {
			spaces++
		}
		
		if spaces > 2 {
			break
		}
	}
	
	// Text usage tends to fall into a number of different roles, with a power law
	// frequency of occurrence in a text, so let's classify in order of likely usage
	// for small and many, we use a hashmap/btree
	
	switch spaces {
	case 0:
		return l,N1GRAM
	case 1:
		return l,N2GRAM
	case 2:
		return l,N3GRAM
	}
	
	// For longer strings, a linear search is probably fine here
        // (once it gets into a database, it's someone else's problem)
	
	if l < 128 {
		return l,LT128
	}
	
	if l < 1024 {
		return l,LT1024
	}
	
	return l,GT1024
}

// **************************************************************************

func STtype(st int) string {

	st = st - ST_ZERO
	var ty string

	switch st {
	case -EXPRESS:
		ty = "-(expressed by)"
	case -CONTAINS:
		ty = "-(part of)"
	case -LEADSTO:
		ty = "-(arriving from)"
	case NEAR:
		ty = "(close to)"
	case LEADSTO:
		ty = "+(leading to)"
	case CONTAINS:
		ty = "+(containing)"
	case EXPRESS:
		ty = "+(expressing)"
	default:
		ty = "unknown relation!"
	}

	const green = "\x1b[36m"
	const endgreen = "\x1b[0m"

	return green + ty + endgreen
}

// **************************************************************************

func STTypeDBChannel(sttype int) string {

	var link_channel string
	switch sttype {

	case NEAR:
		link_channel = I_NEAR
	case LEADSTO:
		link_channel = I_PLEAD
	case CONTAINS:
		link_channel = I_PCONT
	case EXPRESS:
		link_channel = I_PEXPR
	case -LEADSTO:
		link_channel = I_MLEAD
	case -CONTAINS:
		link_channel = I_MCONT
	case -EXPRESS:
		link_channel = I_MEXPR
	default:
		fmt.Println(ERR_ILLEGAL_LINK_CLASS)
		os.Exit(-1)
	}

	return link_channel
}

// **************************************************************************

func Already (s string, cone map[int][]string) bool {

	for l := range cone {
		for n := 0; n < len(cone[l]); n++ {
			if s == cone[l][n] {
				return true
			}
		}
	}

	return false
}





