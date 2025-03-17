
//
// Simplest text based set-overlap match test
//

package main

import (
	"fmt"

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

	load_arrows := false
	ctx := SST.Open(load_arrows)

	qstr := "SELECT S from Node where unaccent(S) LIKE '%xue%'"

	fmt.Println("TRY",qstr)

	row,err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("FAILED \n",qstr,err)
	}

	var whole string

	for row.Next() {		
		err = row.Scan(&whole)
		fmt.Println("GOT",whole)
	}

	row.Close()

	SST.Close(ctx)
}
