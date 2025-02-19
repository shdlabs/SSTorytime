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

	if !CreateTable(db,"Person(name text, hasfriend text[], employs text[], primary key(name))") {
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

	//ReadNodes(db)

	fmt.Println("From a total space of",len(everyone))

	centre := "Mark"

	for radius := 1; radius < 7; radius++ {
		CalculatePureBall(db,centre,radius)
	}

	fmt.Println("----------------------------------------")

	for radius := 1; radius < 7; radius++ {
		CalculateHybridBall(db,centre,radius)
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

	qstr = fmt.Sprintf("INSERT INTO Person(name) VALUES ( '%s' ) RETURNING name",key)

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

func CalculatePureBall(db *sql.DB, centre string, radius int) {

	qstr := fmt.Sprintf("WITH RECURSIVE templist (name,friends,radius) " +
                "AS ( " +
                //  -- anchor member
                " SELECT name,unnest(hasfriend) friends,1 FROM entity WHERE name='%s' " +
                " UNION " +
                // -- recursive term
                " SELECT e.name,unnest(e.hasfriend),radius+1 " +
                " FROM entity e JOIN templist ON e.name = friends where radius < %d " +
                ")" +
                "SELECT DISTINCT friends FROM templist",centre,radius)

	row, err := db.Query(qstr)

	if err != nil {
		fmt.Println("Error executing query:",qstr,err)
	}

	var v string
	var ball []string

	for row.Next() {

		err = row.Scan(&v)

		if err != nil {
			fmt.Println("Error scanning row:",qstr,err)
		} else {
			ball = append(ball,v)
		}
	}

	fmt.Println("\nPURE Ball around ",centre,"radius",radius,": volume", len(ball))

	var cols int
	sort.Strings(ball)

	for r := range ball {

		if cols % 5 == 0 {
			fmt.Print("\n     ")
		}

		fmt.Printf(" %s,",ball[r])
		cols++
	}

	fmt.Println()
}

// **************************************************************************

func CalculateHybridBall(db *sql.DB, centre string, radius int) {

	// We use friends as a running list/set to be appended, and radius as a counter

	qstr := fmt.Sprintf("WITH RECURSIVE templist (name,friends,radius) " +
                "AS ( " +
                //  -- anchor member
                " SELECT name,unnest(array_cat(hasfriend,employs)),1 FROM entity WHERE name='%s' " +
                " UNION " +
                // -- recursive term
                " SELECT e.name,unnest(array_cat(e.hasfriend,e.employs)),radius+1 " +
                " FROM entity e JOIN templist t ON e.name = t.friends where radius < %d " +
                ")" +
                "SELECT DISTINCT friends FROM templist",centre,radius)

	row, err := db.Query(qstr)

	if err != nil {
		fmt.Println("Error executing query:",qstr,err)
	}

	var v string
	var ball []string

	for row.Next() {

		err = row.Scan(&v)

		if err != nil {
			fmt.Println("Error scanning row:",qstr,err)
		} else {
			ball = append(ball,v)
		}
	}

	fmt.Println("\nMIXED Ball around ",centre,"radius",radius,": volume", len(ball))

	var cols int
	sort.Strings(ball)

	for r := range ball {

		if cols % 5 == 0 {
			fmt.Print("\n     ")
		}

		fmt.Printf(" %s,",ball[r])
		cols++
	}
	fmt.Println()
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




