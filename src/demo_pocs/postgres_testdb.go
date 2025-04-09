package main

import (
        SST "SSTorytime"
)


func main() {

	ctx := SST.Open(false)

	SST.Close(ctx)
}

