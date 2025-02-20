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

	if !CreateTable(db,"Entity(name text, hasfriend text[], employs text[], primary key(name))") {
	   os.Exit(-1)
	}

	var friends = make(map[string][]string)
	var employees = make(map[string][]string)
	var everyone= make(map[string]int)

	friends["Mark"] = []string{ "Silvy","Mandy","Brent"}
	friends["Mike"] = []string{"Mark","Jane1","Jane2","Jan","Alfie","Jungi","Peter","Paul"}
	friends["Jan"] = []string{"Adam","Jane1","Jane"}
	friends["Adam"] = []string{"Company of Friends","Paul","Matt","Billie","Chirpy Cheep Cheep","Taylor Swallow"}
	friends["Mandy"] = []string{"Zhao","Doug","Tore","Joyce","Mike","Carol","Ali","Matt","Bjørn","Tamar","Kat","Hans"}
	friends["Company of Friends"] = []string{"Matt","Jane1"}

	employees["Company of Friends"] = []string{"Robo1","Robo2","Bot1","Bot2","Bot3","Bot4","Rob1Bot21"}

	for entity := range friends {
		CreateNode(db, entity)
		everyone[entity]++
		for fr := range friends[entity] {
			CreateNode(db, friends[entity][fr])
			everyone[friends[entity][fr]]++
			AppendLink(db,"hasfriend",entity,friends[entity][fr])
		}
	}

	for entity := range employees {
		CreateNode(db, entity)
		everyone[entity]++
		for fr := range employees[entity] {
			CreateNode(db, employees[entity][fr])
			everyone[employees[entity][fr]]++
			AppendLink(db,"employs",entity,employees[entity][fr])
		}
	}

	centre := "Mark"

	for radius := 1; radius < 10; radius++ {
		GetFutureCone(db,centre,radius)
	}

}

// **************************************************************************

func CreateTable(db *sql.DB,defn string) bool {

        fmt.Println("Create table from type...")
	
	_,err := db.Query("CREATE TABLE IF NOT EXISTS "+defn)
	
	if err != nil {
		s := fmt.Sprintln("Failed to create a table ",err)
		
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

	qstr = fmt.Sprintf("INSERT INTO Entity(name) VALUES ( '%s' ) RETURNING name",key)

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

	qstr := fmt.Sprintf("update Entity set %s = array_append(%s,'%s') where name = '%s' and (%s is null or not '%s' = ANY(%s))",arrow,arrow,fr,name,arrow,fr,arrow)

	_,err := db.Query(qstr)

	if err != nil {
		fmt.Println("Failed to append",err)
	       return false
	}

	return true
}

// **************************************************************************

func GetLinksFromNode(db *sql.DB, key string) []string {

	qstr := fmt.Sprintf("select hasfriend from Entity where name='%s'",key)

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

func GetFutureCone(db *sql.DB, centre string, radius int) {

	fmt.Println("--- Future cone by layers --- from ",centre,"depth",radius)

	row, err := db.Query("drop table output") // No error check, if output exists next will fail

	/* We group SQL commands by starting BEGIN ... COMMIT */

	qstr := fmt.Sprintf("BEGIN;" +
		"WITH RECURSIVE FutureCone(name,member,depth) "+
		"AS ( " +
		"SELECT name,unnest(hasfriend), 1 FROM entity WHERE name='%s' "+
		"UNION "+
		"SELECT e.name,unnest(e.hasfriend),depth+1 FROM entity e JOIN FutureCone ON e.name = member where depth < %d"+
		") "+
		"SELECT member,depth into output FROM FutureCone where not member='%s' order by depth; "+
		"select * from output; "+
		"commit;",centre,radius,centre)

	row, err = db.Query(qstr)

	if err != nil {
		fmt.Println("Error executing query:",qstr,err)
	}

	const maxdepth = 10

	var v string
	var l int
	var cone = make(map[int][]string,1)

	cone[0] = append(cone[0],centre)

	for row.Next() {

		err = row.Scan(&v,&l)

		if err != nil {
			fmt.Println("Error scanning row:",qstr,err)
		} else {
			_,ok := cone[l]

			if !ok {
				cone[l] = make([]string,1)
			}

			if !Already(v,cone) {
				cone[l] = append(cone[l],v)
			}
		}
	}

	for l := 0; l < len(cone); l++ {
		fmt.Println(l,cone[l])
	}

	fmt.Println()
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




