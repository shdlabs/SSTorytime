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

const (
	NEAR = 0
	LEADSTO = 1   // +/-
	CONTAINS = 2  // +/-
	EXPRESS = 3   // +/-

	ST_ZERO = EXPRESS
	ST_TOP = ST_ZERO + EXPRESS + 1

	N1GRAM = 1
	N2GRAM = 2
	N3GRAM = 3
	LT128 = 4
	LT1024 = 5
	GT1024 = 6
)

//**************************************************************

type Node struct { // essentially the incidence matrix

	L int                 // length of name string
	S string              // name string itself

	Chap string           // section/chapter in which this was added
	SizeClass int         // the string class: N1-N3, LT128, etc
	NPtr NodePtr // Pointer to self

	I [ST_TOP][]Link   // link incidence list, by arrow type
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

const NODE_TABLE = "CREATE TABLE IF NOT EXISTS Node " +
	"( " +
	"NPtr      int primary key," +
	"L         int,            " +
	"S         text,           " +
	"Chap      text,           " +
	"SizeClass int,            " +
	"I         Link[7][]       " +
	")"

const NODEPTR_TABLE = "CREATE TABLE IF NOT EXISTS NodePtr " +
	"( " +
	"CPtr int,               " +
	"Class int,              " +
	"primary key(Cptr,Class) " +
	")"

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

	if !CreateTable(ctx,NODE_TABLE) {
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

	qstr = fmt.Sprintf("INSERT INTO Node(L,S,Chap,SizeClass) VALUES (  ) RETURNING name")

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

func AppendDBLinkToNode(ctx PoSST, nodeptr NodePtr, lnk Link) bool {

	// Want to make this idempotent, because SQL is not (and not clause)

	qstr := fmt.Sprintf("update person set %s = array_append(%s,'%s') where name = '%s' and (%s is null or not '%s' = ANY(%s))")

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



