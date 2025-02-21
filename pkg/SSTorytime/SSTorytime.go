
package SSTorytime

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	//"sort"

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

type PoSST struct {

   DB *sql.DB
}

//******************************************************************

func Open() PoSST {

	var ctx PoSST
	var err error

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

func Close(ctx PoSST) {
	ctx.DB.Close()
}

// **************************************************************************

func CreateTable(ctx PoSST,defn string) bool {

        fmt.Println("Create table from type...")
	
	_,err := ctx.DB.Query("CREATE TABLE IF NOT EXISTS "+defn)
	
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

func CreateNode(ctx PoSST, key string) bool {

	var qstr string

	qstr = fmt.Sprintf("INSERT INTO Person(name) VALUES ( '%s' ) RETURNING name",key)

	_,err := ctx.DB.Query(qstr)

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

func AppendLink(ctx PoSST, arrow,name,fr string) bool {

	// Want to make this idempotent, because SQL is not (and not clause)

	qstr := fmt.Sprintf("update person set %s = array_append(%s,'%s') where name = '%s' and (%s is null or not '%s' = ANY(%s))",arrow,arrow,fr,name,arrow,fr,arrow)

	_,err := ctx.DB.Query(qstr)

	if err != nil {
		fmt.Println("Failed to append",err)
	       return false
	}

	return true
}

// **************************************************************************

func ReadNodes(ctx PoSST) {

	var node string

	rows, err := ctx.DB.Query("SELECT name,hasfriend FROM Person")

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

func GetLinksFromNode(ctx PoSST, key string) []string {

	qstr := fmt.Sprintf("select hasfriend from Person where name='%s'",key)

	row, err := ctx.DB.Query(qstr)

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



