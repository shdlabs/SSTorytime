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
	"head used as chinese",
	"head context neuro",
	"leg in chapter bodyparts",
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
	"a2 to b5",
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
	"summary chapter interference",
	"showme greetings in norwegian",
        }


//******************************************************************

func main() {

	args := GetArgs()

	fmt.Println(args)

	SST.MemoryInit()

	load_arrows := false
	ctx := SST.Open(load_arrows)

	for test := range TESTS {
		search := DecodeSearchField(TESTS[test])
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

const (
	CMD_NOTE = "note"
	CMD_TO = "to"
	CMD_CTX = "ctx"
	CMD_CONTEXT = "context"
	CMD_CHAPTER = "chapter"
	CMD_SECTION = "section"
	CMD_IN = "in"
	CMD_ARROW = "arrow"
	CMD_USEDAS = "as"
)

//******************************************************************

func DecodeSearchField(cmd string) SST.SearchParameters {

	var keywords = []string{ 
		CMD_NOTE, "page",
		"visual","img","image",
		"path",CMD_TO,"story", "stories", 
		"sequence","story","stories",
		CMD_CONTEXT,CMD_CTX,CMD_USEDAS,
		CMD_CHAPTER,CMD_IN,CMD_SECTION,
		"node", "vertex","image","node","match","summary","show", 
		CMD_ARROW,
        }
	
	var ignore = []string{"about", "for", "on", "of", "used" }
	
	// parentheses are reserved for unaccenting

	m := regexp.MustCompile("[ \t]+") 
	cmd = m.ReplaceAllString(cmd," ") 

	cmd = strings.TrimSpace(cmd)
	pts := SST.SplitPunctuationText(cmd)

	var parts [][]string
	var part []string

	for p := range pts {

		subparts := SplitQuotes(pts[p])

		for w := range subparts {
			if w > 0 && In(subparts[w],keywords) {
				if strings.HasPrefix(subparts[w],"to") {
					part = append(part,subparts[w])
				} else {
					parts = append(parts,part)
					part = nil
					part = append(part,subparts[w])
				}
			} else if !In(subparts[w],ignore) {
				part = append(part,subparts[w])
			}
		}
	}

	parts = append(parts,part)

	// command is now segmented

	param := FillInParameters(parts)

	return param
}

//******************************************************************

func FillInParameters(cmd_parts [][]string) SST.SearchParameters {

	var param SST.SearchParameters 

	for c := range cmd_parts {
		for p := range cmd_parts[c] {

			fmt.Println("dealing with CMD -->",cmd_parts[c][p],"of",cmd_parts)

			switch cmd_parts[c][p] {
				
			case CMD_CHAPTER:
				//			param.Chapter
			}
		}
	}

	var nptr SST.NodePtr
	param.NPtr = append(param.NPtr,nptr)

	param.Context = append(param.Context,)
	param.Arrows  = append(param.Arrows,)

	return param
}

//******************************************************************

func In(s string,list []string) bool {

	for w := range list {
		if strings.HasPrefix(s,list[w]) {
			return true
		}
	}
	return false
}

//******************************************************************

func SplitQuotes(s string) []string {

	var items []string
	var upto []rune
	var blocked bool = false

	quotes := strings.Count(s,"\"")

	if quotes % 2 != 0 {
		fmt.Println("Unpaired quotes in search",s,quotes)
	}

	cmd := []rune(s)

	for r := 0; r < len(cmd); r++ {

		switch cmd[r] {
		case ' ':
			if !blocked {
				items = append(items,string(upto))
				upto = nil
				continue
			}
			break

		case '"':
			if blocked {
				items = append(items,string(upto))
				upto = nil
			}
			blocked = !blocked
			continue
		}

		upto = append(upto,cmd[r])
	}

	items = append(items,string(upto))
	return items
}

//******************************************************************

func Search(ctx SST.PoSST, search SST.SearchParameters) {


}









