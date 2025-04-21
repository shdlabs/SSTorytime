package main

import (
	"fmt"
	"sort"

        SST "SSTorytime"
)

func main() {

	var order []string

	ctx := SST.Open(false)

	chap := ""
	context := []string{}

	toc := SST.JSON_TableOfContents(ctx,chap,context)

	for keys := range toc {
		order = append(order,keys)
	}

	sort.Strings(order)

	for key := 0; key < len(order); key++ {
		fmt.Println(key,order[key],"=",toc[order[key]])
	}

	SST.Close(ctx)

}
