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
	"regexp"
	"strings"

        SST "SSTorytime"
)

//******************************************************************

var TESTS = []string{ 
	"visual for ",	
	"visual of ",	
	"notes on",	
	"notes about",	
	"(1,1), (1,3), (4,4) (3,3)",
	"page 2 of notes on brains", 
	"notes page 3 brain", 
	"integrate in math",	
	"arrows pe,ep, eh",
	"arrows 1,-1",
	"paths to/from arrows pe,ep, eh",
	"paths from start to target",	
	"a1 to b6",
	"a1 to b6 arrows then",
	"forward cone for (bjorvika)",
	"backward cone for (bjorvika)",
	"sequences about ",	
	"stories about (bjorvika)",	
	"context \"not only\"", 
	"\"come in\"",	
	"containing / matching \"blub blub\"", 
	"chinese kinds of meat", 
	"images prince", 
	"summary of chapter interference",
	"showme greetings in norwegian",
        }

var KEYWORDS = []string{ 
	"note", "page","notes", 
	"visual","img","image",
	"story", "stories", 
	"sequence","story","stories",
	"context","not","used as", 
	"chapter","in","section",
	"node", "vertex","image","node","match","summary","show", 
	"arrow", "link", "edge", 
        }

var IGNORE = []string{"about", "for", "on", "of" }

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

func DecodeSearch(cmd string) SST.SearchParameters {

	// parentheses are reserved for unaccenting

	fmt.Println("\nDECODE",cmd)

	var param SST.SearchParameters 

	m := regexp.MustCompile("[ \t]+") 
	cmd = m.ReplaceAllString(cmd," ") 

	cmd = strings.TrimSpace(cmd)
	pts := SST.SplitPunctuationText(cmd)

	var parts [][]string
	var part []string

	for p := range pts {

// retain quoted strings
		subparts := strings.Split(pts[p]," ")

		for w := range subparts {
			if w > 0 && In(subparts[w],KEYWORDS) {
				parts = append(parts,part)
				part = nil
				part = append(part,subparts[w])
			} else if !In(subparts[w],IGNORE) {
				part = append(part,subparts[w])
			}
		}
	}

	parts=append(parts,part)

	for c := range parts {
		fmt.Println("CMD",parts[c])
	}

	return param
}

//******************************************************************

func In(s string,list []string) bool {

	for w := range list {
		if strings.Contains(s,list[w]) {
			return true
		}
	}
	return false
}

//******************************************************************

func Search(ctx SST.PoSST, search SST.SearchParameters) {


}









