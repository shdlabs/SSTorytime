//******************************************************************
//
// Demo of accessing postgres with custom data structures and arrays
// ( for graphs )
//
//******************************************************************

package main

import (
	"fmt"
	"sort"

        SST "SSTorytime"
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

        ctx := SST.Open()

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
		SST.CreateNode(ctx, entity)
		everyone[entity]++
		for fr := range friends[entity] {
			SST.CreateNode(ctx, friends[entity][fr])
			everyone[friends[entity][fr]]++
			SST.AppendLink(ctx,"hasfriend",entity,friends[entity][fr])
		}
	}

	for entity := range employees {
		SST.CreateNode(ctx, entity)
		everyone[entity]++
		for fr := range employees[entity] {
			SST.CreateNode(ctx, employees[entity][fr])
			everyone[employees[entity][fr]]++
			SST.AppendLink(ctx,"employs",entity,employees[entity][fr])
		}
	}

	//ReadNodes(db)

	fmt.Println("From a total space of",len(everyone))

	centre := "Mark"

	for radius := 1; radius < 7; radius++ {
		CalculatePureBall(ctx,centre,radius)
	}

	SST.Close(ctx)
}

// **************************************************************************

func CalculatePureBall(ctx SST.PoSST, centre string, radius int) {

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

	row, err := ctx.DB.Query(qstr)

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






