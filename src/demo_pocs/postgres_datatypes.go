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

type Link struct {
     Weight float64
     Reln int
     To   int
}

//******************************************************************

type NodeEventItem struct {
	Key   string
	Value float64
	Links []Link
}

//******************************************************************

func main() {

	// db, err := sql.Open("postgres", "postgres://sstoryline:sst_1234@localhost:5432/sst?sslmode=disable")

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

	if !CreateType(db,"PGLink AS (weight real, arrow int,dest int)") {
	   os.Exit(-1)
	}

	if !CreateTable(db,"NodeEventItem(key text, value int, links PGLink [], primary key (key))") {
	   os.Exit(-1)
	}

	var lnk = Link{Weight: 1.2, Reln: 6, To: 5}
	
	var links []Link

	links = append(links,lnk)
	links = append(links,Link{Weight: 99, Reln: 9, To: 9999})

	if !CreateNodeEventItem(db, "ninetynine", 999,links) {
	   os.Exit(-1)
	}

	CreateNodeEventItem(db, "two", 9999,nil)

	AppendToLinks(db,"ninetynine",Link{Weight: 1.0, Reln: 2, To: 5456})
	AppendToLinks(db,"ninetynine",Link{Weight: 1.0, Reln: 2, To: 5456})
	AppendToLinks(db,"ninetynine",Link{Weight: 1.0, Reln: 2, To: 5456})


	records, err := ReadNodeEventItems(db)

	if err != nil {
		fmt.Println("Error reading records: ", err)
	}

	fmt.Println("All records:")

	for _, r := range records {
		fmt.Println("Key:", r.Key, "Value:", r.Value,r.Links)
	}

        somelinks := GetLinksFromNode(db,"ninetynine")

	for l := range somelinks {
		fmt.Println("  - Node has arrow",somelinks[l])
	}
}

// **************************************************************************

func CreateType(db *sql.DB, defn string) bool {

        fmt.Println("Create type ...")
	
	_,err := db.Query("CREATE TYPE "+defn)

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

func CreateNodeEventItem(db *sql.DB, vkey string, vvalue int,array []Link) bool {

	var qstr string

	if !IdempotentArray(array) {
		fmt.Println("The proposed links contain duplicates",array)
		os.Exit(-1)
	}

	if array != nil {
		varray := FormatLinkArray(array)
		qstr = fmt.Sprintf("INSERT INTO NodeEventItem(key,value,links) VALUES ( '%s', '%d', %s ) RETURNING key",vkey,vvalue,varray)
	} else {
		qstr = fmt.Sprintf("INSERT INTO NodeEventItem(key,value) VALUES ( '%s', '%d' ) RETURNING key",vkey,vvalue)
	}

	_,err := db.Query(qstr)

	if err != nil {
		s := fmt.Sprint("Failed to insert",vkey,vvalue,err)
		
		if strings.Contains(s,"duplicate key") {
			return true
		} else {
			fmt.Println(s,"\n",qstr)
			return false
		}
	}
	
	return true
}

// **************************************************************************

func AppendToLinks(db *sql.DB, key string, l Link) bool {

	// Want to make this idempotent, because SQL is not (and not clause)

	qstr := fmt.Sprintf("update NodeEventItem set links = array_append(links, '(%f, %d, %d)' ) where key = '%s'and not '(%f, %d, %d)'::PGLink = ANY(links)",l.Weight,l.Reln,l.To,key,l.Weight,l.Reln,l.To)

	_,err := db.Query(qstr)

	if err != nil {
		fmt.Println("Failed to append",err)
	       return false
	}
	return true
}

// **************************************************************************

func ReadNodeEventItems(db *sql.DB) ([]NodeEventItem, error) {

	var node NodeEventItem

	rows, err := db.Query("SELECT key,value,links FROM NodeEventItem")

	if err != nil {
		fmt.Println("Error executing query: ", err)
	}

	defer rows.Close()

	var records []NodeEventItem

	for rows.Next() {

                // pq can't handle postgres arrays, so we have to
	    	var whole_array string
		
		err := rows.Scan(&node.Key,&node.Value,&whole_array)

                node.Links = ParseLinkArray(whole_array)
		
		if err != nil {
		   //fmt.Println("Error: ", err)
		} else {
		  records = append(records, node)
		}
	}

	return records, nil
}

// **************************************************************************

func GetLinksFromNode(db *sql.DB, key string) []Link {

	qstr := fmt.Sprintf("select links from NodeEventItem where key='%s'",key)

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

func ParseLinkArray(whole_array string) []Link {

   // array as {"(1,2,3)","(4,5,6)"}

      	var l []Link

    	whole_array = strings.Replace(whole_array,"{","",-1)
    	whole_array = strings.Replace(whole_array,"}","",-1)
	whole_array = strings.Replace(whole_array,"\",\"",";",-1)
	whole_array = strings.Replace(whole_array,"\"","",-1)
	
        items := strings.Split(whole_array,";")

	for i := range items {
	    var lnk Link
	    s := strings.TrimSpace(items[i])
	    fmt.Sscanf(s,"(%f,%d,%d)",&lnk.Weight,&lnk.Reln,&lnk.To)
	    l = append(l,lnk)
	    }

	return l
}

// **************************************************************************

func FormatLinkArray(array []Link) string {

	// ARRAY ['(1,2,3)' :: Link , '(4,5,6)' :: Link ]  ;

        if len(array) == 0 {
	   return ""
        }

	var ret string = "ARRAY ["
	
	for i := 0; i < len(array); i++ {
	    ret += fmt.Sprintf("'(%f,%d,%d)' :: PGLink ",array[i].Weight,array[i].Reln,array[i].To)
	    if i < len(array)-1 {
	    ret += ", "
	    }
        }

	ret += "]"

	return ret
}

// **************************************************************************

func IdempotentArray(array []Link) bool {

	var check = make(map[string]int)

	for i := 0; i < len(array); i++ {
		s := fmt.Sprintf("(%f,%d,%d)",array[i].Weight,array[i].Reln,array[i].To)
		check[s]++
	}

	for i := range check {
		if check[i] > 1 {
			return false
		}
	}
	return true
}






