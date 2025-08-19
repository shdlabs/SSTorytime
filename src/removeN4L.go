//******************************************************************
//
// Try out neighbour search for all ST stypes together
//
// Prepare:
// cd examples
// ../src/N4L-db -u chinese.n4l
//
//******************************************************************

package main

import (
	"fmt"
	"strings"
        SST "SSTorytime"
)

var path [8][]string

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	chapter := "reminders"

	DeleteChapter(ctx,chapter)

	SST.Close(ctx)
}

//******************************************************************

func DeleteChapter(ctx SST.PoSST,chapter string) {

	qstr := fmt.Sprintf("select DeleteChapter('%s')",chapter)

	row,err := ctx.DB.Query(qstr)
	
	if err != nil {
		fmt.Println("Error running deletechapter function:",qstr,err)
	}

	row.Close()

}

//******************************************************************

func RemoveFromStringList(list,remove string) string {

	split := strings.Split(list,",")

	var retvar string

	for s := 0; s < len(split); s++ {

		lc := strings.ToLower(split[s])
		lt := strings.ToLower(remove)

// this won't work if there are accents
		if !strings.Contains(lc,lt) {
			retvar += split[s] + ","
		}
	}

	strings.Trim(retvar,",")
	fmt.Println("XX",retvar)
	return retvar
}

//******************************************************************

func UpdateDBNode(ctx SST.PoSST,nptr SST.NodePtr,edited string,list []SST.NodePtr) {

	node := SST.GetDBNodeByNodePtr(ctx,nptr)

	fmt.Println("\nEdit",node.Chap,"with",edited)

	for st := 0; st < SST.ST_TOP; st++ {
		for dst := range node.I[st] {
			fmt.Println("  delete lenk",node.I[st][dst].Dst)
		}
	}
}

//******************************************************************

func DeleteDBNodeArrowNode(ctx SST.PoSST,nptr SST.NodePtr) {

	fmt.Println("Remove NodeArrowNode",nptr)
}

//******************************************************************

func DeleteDBNode(ctx SST.PoSST,nptr SST.NodePtr) {

	fmt.Println("DELETE nptr",nptr)
}





