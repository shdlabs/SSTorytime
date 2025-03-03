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

	// Node these are "plpgsql" language, NOT pure SQL. They are very different.
	
	// Insert a node structure, also an anchor for and containing link arrays
	
	qstr := "CREATE OR REPLACE FUNCTION IdempInsertNode(iLi INT, iszchani INT, icptri INT, iSi TEXT, ichapi TEXT)\n" + 
		"RETURNS TABLE (    \n" +
			"    ret_cptr INTEGER," +
			"    ret_channel INTEGER" +
			") AS $fn$ " +
			"DECLARE \n" +
			"BEGIN\n" +
			"  IF NOT EXISTS (SELECT (NPtr).Chan,(NPtr).CPtr FROM Node WHERE s = iSi) THEN\n" +
			"     INSERT INTO Node (Nptr.Cptr,L,S,chap,Nptr.Chan) VALUES (icptri, iLi, iSi, ichapi, iszchani);" +
			"  end if;\n" +
			"      return query SELECT (NPtr).Chan,(NPtr).CPtr FROM Node WHERE s = iSi;\n" +
			"END ;\n" +
			"$fn$ LANGUAGE plpgsql;";

		_, err := ctx.DB.Query(qstr)
		
		if err != nil {
			fmt.Println("Error defining postgres function:",qstr,err)
		}

	// Construct an empty link pointing nowhere as a starting node

	qstr = "CREATE OR REPLACE FUNCTION GetSingletonAsLinkArray(start NodePtr)\n"+
		"RETURNS Link[] AS $nums$\n"+
		"DECLARE \n"+
		"    level Link[] := Array[] :: Link[];\n"+
		"    lnk Link := (0,0.0,Array[]::text[],(0,0));\n"+
		"BEGIN\n"+
		" lnk.Dst = start;\n"+
		" level = array_append(level,lnk);\n"+
		"RETURN level; \n"+
		"END ;\n"+
		"$nums$ LANGUAGE plpgsql;"

	_, err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	// Construct an empty link pointing nowhere as a starting node

	qstr = "CREATE OR REPLACE FUNCTION GetSingletonAsLink(start NodePtr)\n"+
		"RETURNS Link AS $nums$\n"+
		"DECLARE \n"+
		"    lnk Link := (0,0.0,Array[]::text[],(0,0));\n"+
		"BEGIN\n"+
		" lnk.Dst = start;\n"+
		"RETURN lnk; \n"+
		"END ;\n"+
		"$nums$ LANGUAGE plpgsql;"

	_, err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	// Construct search by sttype. since table names are static we need a case statement

	qstr = "CREATE OR REPLACE FUNCTION GetNeighboursByType(start NodePtr, sttype int)\n"+
		"RETURNS Link[] AS $nums$\n"+
		"DECLARE \n"+
		"    fwdlinks Link[] := Array[] :: Link[];\n"+
		"    lnk Link := (0,0.0,Array[]::text[],(0,0));\n"+
		"BEGIN\n"+
		"   CASE sttype \n"
	
	for st := -EXPRESS; st <= EXPRESS; st++ {
		qstr += fmt.Sprintf("WHEN %d THEN\n"+
			"     SELECT %s INTO fwdlinks FROM Node WHERE Nptr=start;\n",st,STTypeDBChannel(st));
	}
	qstr += "ELSE RAISE EXCEPTION 'No such sttype %', sttype;\n" +
		"END CASE;\n" +
		"    RETURN fwdlinks; \n" +
		"END ;\n" +
		"$nums$ LANGUAGE plpgsql;\n"

	_, err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	// Get the nearest neighbours as NPtr, with respect to each of the four STtype

	qstr = fmt.Sprintf("CREATE OR REPLACE FUNCTION GetFwdNodes(start NodePtr,exclude NodePtr[],sttype int)\n"+
		"RETURNS NodePtr[] AS $nums$\n" +
		"DECLARE \n" +
		"    neighbours NodePtr[];\n" +
		"    fwdlinks Link[];\n" +
		"    lnk Link;\n" +
		"BEGIN\n" +

		"    fwdlinks =GetNeighboursByType(start,sttype);\n"+

		"    IF fwdlinks IS NULL THEN\n" +
		"        RETURN '{}';\n" +
		"    END IF;\n" +

		"    neighbours := ARRAY[]::NodePtr[];\n" +

		"    FOREACH lnk IN ARRAY fwdlinks\n" +
		"    LOOP\n"+
		"      IF exclude is not NULL AND NOT lnk.dst=ANY(exclude) THEN\n" +
		"         neighbours := array_append(neighbours, lnk.dst);\n" +
		"      END IF; \n" +
		"    END LOOP;\n" +

		"    RETURN neighbours; \n" +
		"END ;\n" +
		"$nums$ LANGUAGE plpgsql;\n")

	_, err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	qstr = fmt.Sprintf("CREATE OR REPLACE FUNCTION GetFwdLinks(start NodePtr,exclude NodePtr[],sttype int)\n"+
		"RETURNS Link[] AS $nums$\n" +
		"DECLARE \n" +
		"    neighbours Link[];\n" +
		"    fwdlinks Link[];\n" +
		"    lnk Link;\n" +
		"BEGIN\n" +

		"    fwdlinks =GetNeighboursByType(start,sttype);\n"+

		"    IF fwdlinks IS NULL THEN\n" +
		"        RETURN '{}';\n" +
		"    END IF;\n" +
		"    neighbours := ARRAY[]::Link[];\n" +
		"    FOREACH lnk IN ARRAY fwdlinks\n" +
		"    LOOP\n"+
		"      IF exclude is not NULL AND NOT lnk.dst=ANY(exclude) THEN\n" +
		"         neighbours := array_append(neighbours, lnk);\n" +
		"      END IF; \n" +
		"    END LOOP;\n" +
		"    RETURN neighbours; \n" +
		"END ;\n" +
		"$nums$ LANGUAGE plpgsql;\n")
	
	_, err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}
	
	// Get the forward cone / half-ball as NPtr

	qstr = "CREATE OR REPLACE FUNCTION FwdConeAsNodes(start NodePtr,sttype INT, maxdepth INT)\n"+
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
		"      partlevel := GetFwdNodes(neigh,exclude,sttype);\n" +
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
	
	
          /* e.g. select unnest(fwdconeaslinks) from FwdConeAsLinks('(4,1)',1,4);
                           unnest                           
             ------------------------------------------------------------
              (0,0,{},"(4,1)")
              (77,0.34,"{ ""fairy castles"", ""angel air"" }","(4,2)")
              (77,0.34,"{ ""fairy castles"", ""angel air"" }","(4,3)")
              (77,0.34,"{ ""steamy hot tubs"" }","(4,5)")
              (77,0.34,"{ ""fairy castles"", ""angel air"" }","(4,4)")
              (77,0.34,"{ ""steamy hot tubs"", ""lady gaga"" }","(4,6)")
             (6 rows)

          */

	qstr = "CREATE OR REPLACE FUNCTION FwdConeAsLinks(start NodePtr,sttype INT,maxdepth INT)\n"+
		"RETURNS Link[] AS $nums$\n" +
		"DECLARE \n" +
		"    nextlevel Link[];\n" +
		"    partlevel Link[];\n" +
		"    level Link[] = ARRAY[]::Link[];\n" +
		"    exclude NodePtr[] = ARRAY['(0,0)']::NodePtr[];\n" +
		"    cone Link[];\n" +
		"    neigh Link;\n" +
		"    frn Link;\n" +
		"    counter int := 0;\n" +

		"BEGIN\n" +

		"level := GetSingletonAsLinkArray(start);\n"+

		"LOOP\n" +
		"  EXIT WHEN counter = maxdepth+1;\n" +

		"  IF level IS NULL THEN\n" +
		"     RETURN cone;\n" +
		"  END IF;\n" +

		"  nextlevel := ARRAY[]::Link[];\n" +

		"  FOREACH frn IN ARRAY level "+
		"  LOOP \n"+
		"     nextlevel = array_append(nextlevel,frn);\n" +
		"  END LOOP;\n" +

		"  IF nextlevel IS NULL THEN\n" +
		"     RETURN cone;\n" +
		"  END IF;\n" +

		"  FOREACH neigh IN ARRAY nextlevel LOOP \n"+
		"    IF NOT neigh.Dst = ANY(exclude) THEN\n" +
		"      cone = array_append(cone,neigh);\n" +
		"      exclude := array_append(exclude,neigh.Dst);\n" +
		"      partlevel := GetFwdLinks(neigh.Dst,exclude,sttype);\n" +
		"    END IF;" +
		"    IF partlevel IS NOT NULL THEN\n" +
		"         level = array_cat(level,partlevel);\n"+
		"    END IF;\n" +
		"  END LOOP;\n" +

		"  counter = counter + 1;\n" +
		"END LOOP;\n" +

		"RETURN cone; \n" +
		"END ;\n" +
		"$nums$ LANGUAGE plpgsql;\n"

	_, err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	// Orthogonal (depth first) paths from origin spreading out

	qstr = "CREATE OR REPLACE FUNCTION FwdPathsAsLinks(start NodePtr,sttype INT,maxdepth INT)\n"+
		"RETURNS Text AS $nums$\n" +
		"DECLARE\n" +
		"   hop Text;\n" +
		"   path Text;\n"+
		"   summary_path Text[];\n"+
		"   exclude NodePtr[] = ARRAY['(0,0)']::NodePtr[];\n" +
		"   ret_paths Text;\n" +
		"   startlnk Link;"+

		"BEGIN\n" +

		"startlnk := GetSingletonAsLink(start);\n"+
		"path := Format('%s',startlnk::Text);\n"+
		"ret_paths := SumFwdPaths(startlnk,path,sttype,1,maxdepth,exclude);" +

		"RETURN ret_paths; \n" +
		"END ;\n" +
		"$nums$ LANGUAGE plpgsql;\n"

        // select FwdPathsAsLinks('(4,1)',1,3)

	_, err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}

	// Return end of path branches as aggregated text summaries

	qstr = "CREATE OR REPLACE FUNCTION SumFwdPaths(start Link,path TEXT, sttype INT,depth int, maxdepth INT,exclude NodePtr[])\n"+
		"RETURNS Text AS $nums$\n" +
		"DECLARE \n" + 
		"    fwdlinks Link[];\n" +
		"    empty Link[] = ARRAY[]::Link[];\n" +
		"    lnk Link;\n" +
		"    fwd Link;\n" +
		"    ret_paths Text;\n" +
		"    appendix Text;\n" +
		"    tot_path Text;\n"+
		"BEGIN\n" +

		"IF depth = maxdepth THEN\n"+
//		"  RAISE NOTICE 'Xend path %',path;"+
		"  ret_paths := Format('%s\n%s',ret_paths,path);\n"+
		"  RETURN ret_paths;\n"+
		"END IF;\n"+

		"fwdlinks := GetFwdLinks(start.Dst,exclude,sttype);\n" +

		"FOREACH lnk IN ARRAY fwdlinks LOOP \n" +
		"   IF NOT lnk.Dst = ANY(exclude) THEN\n"+
		"      exclude = array_append(exclude,lnk.Dst);\n" +
		"      IF lnk IS NULL THEN" +
		"         ret_paths := Format('%s\n%s',ret_paths,path);\n"+
//		"         RAISE NOTICE 'Yend path %',tot_path;"+
		"      ELSE"+
		"         tot_path := Format('%s;%s',path,lnk::Text);\n"+
		"         appendix := SumFwdPaths(lnk,tot_path,sttype,depth+1,maxdepth,exclude);\n" +
		"         IF appendix IS NOT NULL THEN\n"+
		"            ret_paths := Format('%s\n%s',ret_paths,appendix);\n"+
		"         END IF;"+
		"      END IF;"+
		"   END IF;"+
		"END LOOP;"+


		"RETURN ret_paths; \n" +
		"END ;\n" +
		"$nums$ LANGUAGE plpgsql;\n"

	// select SumFwdPaths('(4,1)',1,1,3);

	_, err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error defining postgres function:",qstr,err)
	}
}

// **************************************************************************
// Retrieve
// **************************************************************************

func GetFwdConeAsNodes(ctx PoSST, start NodePtr, sttype,depth int) []NodePtr {

	qstr := fmt.Sprintf("select unnest(fwdconeasnodes) from FwdConeAsNodes('(%d,%d)',%d,%d);",start.Class,start.CPtr,sttype,depth)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY to FwdConeAsNodes Failed",err)
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

func GetFwdConeAsLinks(ctx PoSST, start NodePtr, sttype,depth int) []Link {

	qstr := fmt.Sprintf("select unnest(fwdconeaslinks) from FwdConeAsLinks('(%d,%d)',%d,%d);",start.Class,start.CPtr,sttype,depth)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY to FwdConeAsLinkss Failed",err)
	}

	var whole string
	var retval []Link

	for row.Next() {		
		err = row.Scan(&whole)
		l := ParseSQLLinkString(whole)
		retval = append(retval,l)
	}

	return retval
}

// **************************************************************************

func GetFwdPathsAsLinks(ctx PoSST, start NodePtr, sttype,depth int) [][]Link {

	qstr := fmt.Sprintf("select FwdPathsAsLinks from FwdPathsAsLinks('(%d,%d)',%d,%d);",start.Class,start.CPtr,sttype,depth)

	row, err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("QUERY to FwdPathsAsLinkss Failed",err)
	}

	var whole string
	var retval [][]Link

	for row.Next() {		
		err = row.Scan(&whole)
		retval = ParseLinkPath(whole)
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

// **************************************************************************

func ParseSQLLinkString(s string) Link {

        // e.g. (77,0.34,"{ ""fairy castles"", ""angel air"" }","(4,2)")
	// This feels dangerous. Is postgres consistent here?

      	var l Link

    	s = strings.Replace(s,"(","",-1)
    	s = strings.Replace(s,")","",-1)
	s = strings.Replace(s,"\"\"",";",-1)
	s = strings.Replace(s,"\"","",-1)
	
        items := strings.Split(s,",")

	fmt.Sscanf(items[0],"%d",&l.Arr)
	fmt.Sscanf(items[1],"%f",&l.Wgt)

	var array []string

	for i := 2; i <= len(items)-3; i++ {
		items[i] = strings.Replace(items[i],"{","",-1)
		items[i] = strings.Replace(items[i],"}","",-1)
		items[i] = strings.Replace(items[i],";","",-1)
		items[i] = strings.TrimSpace(items[i])
		array = append(array,items[i])
	}

	l.Ctx = array

	fmt.Sscanf(items[len(items)-2],"%d",&l.Dst.Class)
	fmt.Sscanf(items[len(items)-1],"%d",&l.Dst.CPtr)

	return l
}

//**************************************************************

func ParseLinkPath(s string) [][]Link {

	// Each path will start on a new line, with comma sep Link encodings

	var array = make([][]Link,1)
	var index int = 0

	lines := strings.Split(s,"\n")

	for line := range lines {

		if len(lines[line]) > 0 {

			links := strings.Split(lines[line],";")

			if len(links) < 2 {
				continue
			}

			array = append(array,make([]Link,1))

			for l := 0; l < len(links); l++ {

				array[index] = append(array[index],ParseSQLLinkString(links[l]))
			}

			index++
		}
	}

	if index < 1 {
		return nil
	}
	
	return array
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





