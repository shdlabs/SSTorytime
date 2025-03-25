
//
// Simplest text based set-overlap match test
//

package main

import (
	"fmt"

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)


	// Show me the nodes in this context

	qstr := "select AllNCPathsAsLinks('(1,116)','chinese','{}','any',-1)"

	row,err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("FAILED \n",qstr,err)
	}

	var whole string
	var retval [][]SST.Link

	for row.Next() {		
		err = row.Scan(&whole)
		retval = SST.ParseLinkPath(whole)
	}

	row.Close()

	fmt.Println("GOT EMPTY MATCH",retval)


	// Show me the nodes in this context

	qstr = "select AllNCPathsAsLinks('(1,116)','chinese','{trivia}','any',-1)"

	row,err = ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("FAILED \n",qstr,err)
	}

	for row.Next() {		
		err = row.Scan(&whole)
		retval = SST.ParseLinkPath(whole)
	}

	row.Close()

	fmt.Println("GOT CONTEXT MATCH",retval)

	SST.Close(ctx)
}
