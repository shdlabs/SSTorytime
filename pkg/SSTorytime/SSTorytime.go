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
	I_NEAR = "In0"
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
	Dst NodePtr // adjacent event/item/node
}


const LINK_TYPE = "CREATE TYPE Link AS  " +
	"(                    " +
	"ArrowPtr int,        " +
	"Wgt      real,       " +
	"Ctx      text[],     " +
	"Dst      int         " +
	")"

const NODE_DEF = "" +
	"( " +
	"NPtr      int primary key," +
	"L         int,            " +
	"S         text,           " +
	"Chap      text,           " +
	"SizeClass int,            " +
	I_MEXPR+"  Link[],         " +
	I_MCONT+"  Link[],         " +
	I_MLEAD+"  Link[],         " +
	I_NEAR +"  Link[],         " +
	I_PLEAD+"  Link[],         " +
	I_PCONT+"  Link[],         " +
	I_PEXPR+"  Link[]          " +
	")"

const N1_TABLE = "CREATE TABLE IF NOT EXISTS N1GRAM " + NODE_DEF
const N2_TABLE = "CREATE TABLE IF NOT EXISTS N2GRAM " + NODE_DEF
const N3_TABLE = "CREATE TABLE IF NOT EXISTS N3GRAM " + NODE_DEF
const N4_TABLE = "CREATE TABLE IF NOT EXISTS LT128 " + NODE_DEF
const N5_TABLE = "CREATE TABLE IF NOT EXISTS LT1024 " + NODE_DEF
const N6_TABLE = "CREATE TABLE IF NOT EXISTS GT1024 " + NODE_DEF

const NODEPTR_TABLE = "CREATE TABLE IF NOT EXISTS NodePtr " +
	"( " +
	"CPtr int,               " +
	"Class int,              " +
	"primary key(Cptr,Class) " +
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

	return ctx
}

// **************************************************************************

func Configure(ctx PoSST) {

	if !CreateType(ctx,LINK_TYPE) {
		os.Exit(-1)
	}

	// There is no separate link table, as links are an array under nodes
	// There is an adjacency table, however

	if !CreateTable(ctx,N1_TABLE) {
		os.Exit(-1)
	}

	if !CreateTable(ctx,N2_TABLE) {
		os.Exit(-1)
	}

	if !CreateTable(ctx,N3_TABLE) {
		os.Exit(-1)
	}

	if !CreateTable(ctx,N4_TABLE) {
		os.Exit(-1)
	}

	if !CreateTable(ctx,N5_TABLE) {
		os.Exit(-1)
	}

	if !CreateTable(ctx,N6_TABLE) {
		os.Exit(-1)
	}

	if !CreateTable(ctx,NODEPTR_TABLE) {
		os.Exit(-1)
	}

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
			fmt.Println("X",s)
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

func CreateDBNode(ctx PoSST, n Node) bool {

	var qstr string

	qstr = fmt.Sprintf("INSERT INTO Node(Nptr,L,S,Chap,SizeClass) VALUES (%d,%d,'%s','%s',%d)",n.Nptr,n.L,n.S,n.Chap,n.SizeClass)

	_,err := ctx.DB.Query(qstr)

	if err != nil {
		s := fmt.Sprint("Failed to insert",err)
		
		if strings.Contains(s,"duplicate key") {
			return true
		} else {
			fmt.Println(s,"\n",qstr,err)
			return false
		}
	}
	
	return true
}

// **************************************************************************

func AppendDBLinkToNode(ctx PoSST, nodeptr NodePtr, lnk Link, sttype int) bool {

	// Want to make this idempotent, because SQL is not (and not clause)

	if sttype < 0 || sttype > ST_ZERO+EXPRESS {
		fmt.Println(ERR_ST_OUT_OF_BOUNDS,sttype)
		os.Exit(-1)
	}

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
	}

	qstr := fmt.Sprintf("update person set %s = array_append(%s,'%s') where Nptr = '%d' and (%s is null or not '%s' = ANY(%s))")

	_,err := ctx.DB.Query(qstr)

	if err != nil {
		fmt.Println("Failed to append",err)
	       return false
	}

	return true
}


// **************************************************************************
// Tools
// **************************************************************************

func ParseArrayString(whole_array string) []string {

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



