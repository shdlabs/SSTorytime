//******************************************************************
//
// single search string without complex options
//
// how shall we split a search query into parts to match against?
//
//******************************************************************

package main

import (
	"fmt"
	"os"
	"flag"

        SST "SSTorytime"
)

//******************************************************************

var TESTS = []string{"api for Orbits",	"get a visualization",	"notes in chinese",	"notes about chinese",	"integrate in math",	"arrows pe,ep, eh",	"paths from start to target",	"a1 to b6",	"hubs for",	"stories about (bjorvika)",	"sequences about ",	"notes context \"not only\", \"come in\"",	"containing / matching \"(),\"", "page 2 of notes on brains", "notes page 3 brain", "chinese kinds of meat", "images prince", "summary of chapter interference" , "showme greetings in norwegian" }

var KEYWORDS = []{ "note", "page", "about" ,"in", "story", "stories", "sequence", "context", "chapter", "image", "match", "contain", "kind", "summary", "show", "arrow", "link", "edge", "node", "vertex" }

//******************************************************************

func main() {

	args := GetArgs()

	fmt.Println(args)

	SST.MemoryInit()

	load_arrows := false
	ctx := SST.Open(load_arrows)

	for test := range TESTS {
		search := DecodeSearch(TESTS[test])
		Search(ctx,search)
	}

	SST.Close(ctx)
}

//**************************************************************

func Usage() {
	
	fmt.Printf("usage: ByYourCommand <search request>n")
	flag.PrintDefaults()

	os.Exit(2)
}

//**************************************************************

func GetArgs() []string {

	flag.Usage = Usage

	flag.Parse()
	return flag.Args()
}

//******************************************************************

func DecodeSearch(s string) SST.SearchParameters {

	var p SST.SearchParameters 

	fmt.Println(s)

	return p
}

//******************************************************************

func Search(ctx SST.PoSST, search SST.SearchParameters) {

	fmt.Println(search)
}









