//******************************************************************
//
/* 
Demo of accessing postgres with multiple arrays (2d)
 Array are really indexed lists that are one dimensional, so
 extensive use of multdimensional arrays is not going to be efficient 
 because all the work is pushed into a string encoding. In that
case, since we only need a fiuxed number of these lists, it's
better to name them differently and avoid useless overhead
*/

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

	CreateTable(db,"drop table marray")

	const table = "CREATE TABLE IF NOT EXISTS marray" +
		"( " +
		"name text primary key,   " +
		//"myarray text[7][]      " + don't do this...
		"myarray_0 text[],      " +
		"myarray_1 text[],      " +
		"myarray_2 text[],      " +
		"myarray_3 text[]       " +
		")"

	if !CreateTable(db,table) {
	   os.Exit(-1)
	}


	friends := []string{ "Silvy","Mandy","Brent"}

	CreateNode(db, "testnode")

	for sttype := 4; sttype >= 0; sttype-- {

		for l := 0; l < len(friends) && l <= sttype; l++ {

			dst := fmt.Sprintf("%s_%d",friends[l],sttype)

			AppendLink(db,"myarray",sttype,"testnode",dst)
		}

		fmt.Println("Setting type: ",sttype, "...notice how intermediate blank cols block later values until all have data")
		ReadNodes(db)
	}
}

// **************************************************************************

func CreateTable(db *sql.DB,defn string) bool {

        fmt.Println("Create table from type...")
	
	_,err := db.Query(defn)
	
	if err != nil {
		s := fmt.Sprintln("Failed to create a table (%s)",defn,err)
		
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

	qstr = fmt.Sprintf("INSERT INTO MARRAY(name) VALUES ( '%s' ) RETURNING name",key)

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

func AppendLink(db *sql.DB, arrow string,sttype int,anchor string,destination string) bool {

	// Want to make this idempotent, because SQL is not (and not clause)

	name := fmt.Sprintf("%s_%d",arrow,sttype)

	qstr := fmt.Sprintf("update marray set %s = array_append(%s,'%s') where name = '%s' and (%s is null or not '%s' = ANY(%s))",name,name,destination,anchor,name,destination,name)

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

	rows, err := db.Query("SELECT name,myarray_0,myarray_1,myarray_2,myarray_3 FROM marray")

	if err != nil {
		fmt.Println("Error executing query: ", err)
	}

	defer rows.Close()

	for rows.Next() {

                // pq can't handle postgres arrays, so we have to
	    	var arr0,arr1,arr2,arr3 string
		
		rows.Scan(&node,&arr0,&arr1,&arr2,&arr3)

		// NOTE! If arr0 is nil, arr1++ will not be read!!!

		fmt.Println(" -- record name",node,"\n - col 0:",arr1,"\n - col 1:",arr2,"\n - col 2:",arr3,"\n - col 3:",arr3)
	}
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




