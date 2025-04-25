
package main

import (
	"fmt"

        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)
	nptr := SST.GetDBNodePtrMatchingName(ctx,"lamb","")
	const maxdepth = 4

	multicone := "{\n"

	for n := 0; n < len(nptr); n++ {

		cone,_ := SST.GetFwdPathsAsLinks(ctx,nptr[n],1,maxdepth)
		json := SST.JSONCone(ctx,cone,"",nil)

		const empty = 5

		if len(json) > empty {
			multicone += fmt.Sprintf("\"%v\" : %s ",nptr[n],json)
			if n < len(nptr)-1 {
				multicone += ",\n"
			}
		}
	}

	multicone += "\n}\n"

	fmt.Println(multicone)
	SST.Close(ctx)
}










