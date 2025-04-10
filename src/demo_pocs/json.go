
package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	for goes := 0; goes < 10; goes ++ {

		fmt.Println("\n\nEnter some text:")
		
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		
		SearchToJSON(ctx,text)
	}

	SST.Close(ctx)
}

//******************************************************************

func SearchToJSON(ctx SST.PoSST, text string) {

	text = strings.TrimSpace(text)

	var start_set []SST.NodePtr

	search_items := strings.Split(text," ")

	for w := range search_items {
		start_set = append(start_set,SST.GetDBNodePtrMatchingName(ctx,"",search_items[w])...)
	}

	for s := range start_set {
		r := SST.JSONNodeOrbit(ctx, start_set[s]) 
		fmt.Println(s,r)
	}
	
}










