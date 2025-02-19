//******************************************************************
//
// Demo of accessing postgres with custom data structures and arrays
// ( for graphs )
//
//******************************************************************

package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sort"

	_ "github.com/lib/pq"

)

//******************************************************************

const (
	host     = "localhost"
	port     = 5432
	user     = "sstoryline"
	password = "sst_1234"
	dbname   = "newdb"
)

//******************************************************************

const NODE_TABLE = "Node as ("

        "key int," +	
	"L int," +
	"S text," +
	"Chap text," + 

	"primary key(S)," +

        "I_mexp Link[]," +
	"I_mcon Link[]," +
	"I_mfol Link[]," +
	"I_near Link[]," +
	"I_pfol Link[]," +
	"I_pcon Link[]," +
	"I_pexp Link[]," +
	")"

type NodeEventItem struct { // essentially the incidence matrix

	L int                 // length of name string
	S string              // name string itself

	Chap string           // section/chapter in which this was added
	SizeClass int         // the string class: N1-N3, LT128, etc
	NPtr NodeEventItemPtr // Pointer to self

	I [ST_TOP][]Link   // link incidence list, by arrow type
  	                   // NOTE: carefully how offsets represent negative SSTtypes
}

//**************************************************************

const LINK_TABLE = "Link as ("
        "weight real," + 
        "arrow int," +
	"dest int" +
	"primary key(dest,arrow)" +
")"

type Link struct {  // A link is a type of arrow, with context
                    // and maybe with a weightfor package math
	Arr ArrowPtr         // type of arrow, presorted
	Wgt float64          // numerical weight of this link
	Ctx []string         // context for this pathway
	Dst NodeEventItemPtr // adjacent event/item/node
}

type LinkPtr int

//**************************************************************

const ARROW_TABLE = "ArrowDir as ("
        "sttype int," + 
        "long text," +
        "short text," +
        "arrowptr int," +
        "primary key(short,arrowptr)," +
")"

type ArrowDirectory struct {

	STtype  int
	Long    string
	Short   string
	Ptr     ArrowPtr
}

type ArrowPtr int // ArrowDirectory index

 // all fwd arrow types have a simple int representation > 0
 // all bwd/inverse arrow readings have the negative int for fwd
 // Hashed by long and short names

//******************************************************************

func main() {

        connStr := "user="+user+" dbname="+dbname+" password="+password+" sslmode=disable"

        db, err := sql.Open("postgres", connStr)

	if err != nil {
	   	fmt.Println("Error connecting to the database: ", err)
		os.Exit(-1)
	}
	
	defer db.Close()
	
	err = db.Ping()
	
	if err != nil {
		fmt.Println("Error pinging the database: ", err)
		os.Exit(-1)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	if !CreateTable(db,NODE_TABLE) {
	   os.Exit(-1)
	}

	if !CreateTable(db,LINK_TABLE) {
	   os.Exit(-1)
	}

	if !CreateTable(db,NODE_TABLE) {
	   os.Exit(-1)
	}

	if !CreateTable(db,ARROW_TABLE) {
	   os.Exit(-1)
	}

}

// **************************************************************************

func CreateTable(db *sql.DB,defn string) bool {

        fmt.Println("Create table from type...")
	
	_,err := db.Query("CREATE TABLE IF NOT EXISTS "+defn)
	
	if err != nil {
		s := fmt.Sprintln("Failed to create a table of type PGLink ",err)
		
		if strings.Contains(s,"already exists") {
			return true
		} else {
			fmt.Println("Y",s)
			return false
		}
	}

	return true
}

// **************************************************************************

func CreateNode(db *sql.DB, key string) bool {

	var qstr string

	qstr = fmt.Sprintf("INSERT INTO Node(name) VALUES ( '%s' ) RETURNING name",key)

	_,err := db.Query(qstr)

	if err != nil {
		s := fmt.Sprint("Failed to insert",key,err)
		
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

func AppendLink(db *sql.DB, arrow,name,fr string) bool {

	// Want to make this idempotent, because SQL is not (and not clause)

	qstr := fmt.Sprintf("update person set %s = array_append(%s,'%s') where name = '%s' and (%s is null or not '%s' = ANY(%s))",arrow,arrow,fr,name,arrow,fr,arrow)

	_,err := db.Query(qstr)

	if err != nil {
		fmt.Println("Failed to append",err)
	       return false
	}

	return true
}

// **************************************************************************

func ReadNodes(db *sql.DB) {

	var node string

	rows, err := db.Query("SELECT name,hasfriend FROM Person")

	if err != nil {
		fmt.Println("Error executing query: ", err)
	}

	defer rows.Close()

	for rows.Next() {

                // pq can't handle postgres arrays, so we have to
	    	var whole_array string
		
		rows.Scan(&node,&whole_array)

                list := ParseLinkArray(whole_array)
		
		fmt.Println(" -- Person",node,"claims friends",list)
	}
}

// **************************************************************************

func GetLinksFromNode(db *sql.DB, key string) []string {

	qstr := fmt.Sprintf("select hasfriend from Person where name='%s'",key)

	row, err := db.Query(qstr)

	if err != nil {
		fmt.Println("Error executing query:",qstr,err)
	}

	var whole_array string

	for row.Next() {

		err = row.Scan(&whole_array)

		if err != nil {
			fmt.Println("Error scanning row:",qstr,err)
		}
	}

	return ParseLinkArray(whole_array)
}



// **************************************************************************
// Tools
// **************************************************************************

func ParseLinkArray(whole_array string) []string {

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




