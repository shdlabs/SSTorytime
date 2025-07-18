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
	"head used as chinese stuff",
	"head context neuro",
	"leg in chapter bodyparts",
	"foot in bodyparts2",
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
	"paths from arrows pe,ep, eh",
	"paths from start to target",	
	"paths to target3",	
	"a2 to b5",
	"to a5",
	"from start",
	"from (1,6)",
	"a1 to b6 arrows then",
        }


//******************************************************************

func main() {

	args := GetArgs()

	fmt.Println(args)

	SST.MemoryInit()

	load_arrows := false
	ctx := SST.Open(load_arrows)

	for test := range TESTS {
		fmt.Println("......................")
		search := DecodeSearchField(TESTS[test])
		Search(ctx,search,TESTS[test])
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

type SearchParameters struct {

	Name     []string
	From     []string
	To       []string
	Chapter  string
	Context  []string
	Arrows   []string
	PageNr   int
	Sequence bool
}

const (

	CMD_NOTE = "note"
	CMD_PAGE = "page"
	CMD_PATH = "path"
	CMD_STORY = "story"
	CMD_SEQ = "sequence"
	CMD_FROM = "from"
	CMD_TO = "to"
	CMD_CTX = "ctx"
	CMD_CONTEXT = "context"
	CMD_AS = "as"
	CMD_CHAPTER = "chapter"
	CMD_SECTION = "section"
	CMD_IN = "in"
	CMD_ARROW = "arrows"
)

//******************************************************************

func DecodeSearchField(cmd string) SearchParameters {

	var keywords = []string{ 
		CMD_NOTE, CMD_PATH,
		CMD_PATH,CMD_FROM,CMD_TO,CMD_STORY,
		CMD_SEQ,
		CMD_CONTEXT,CMD_CTX,CMD_AS,
		CMD_CHAPTER,CMD_IN,CMD_SECTION,
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
				// special case for TO with implicit FROM, and USED AS

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

func FillInParameters(cmd_parts [][]string) SearchParameters {

	var param SearchParameters 

	for c := 0; c < len(cmd_parts); c++ {

		lenp := len(cmd_parts[c])

		for p := range cmd_parts[c] {

			switch cmd_parts[c][p] {

			case CMD_CHAPTER, CMD_IN:
				if lenp > p+1 {
					param.Chapter = cmd_parts[c][p+1]
					break
				}

			case CMD_NOTE:
				if lenp > p+1 {
					param.Chapter = cmd_parts[c][p+1]
					param.PageNr = 1
					break
				}

			case CMD_PAGE:
				// = GetIntParam(cmd_parts[c][p])
				if lenp > p+1 {
					var no int = 1
					fmt.Sscanf(cmd_parts[c][p+1],"%d",&no)
					param.PageNr = no
					break
				}

			case CMD_ARROW:
				if lenp > p+1 {
					for pp := p+1; pp < lenp; pp++ {
						p++
						ult := strings.Split(cmd_parts[c][pp],",")
						for u := range ult {
							param.Arrows = append(param.Arrows,ult[u])
						}
					}
					break
				}

			case CMD_CONTEXT, CMD_CTX,CMD_AS:
				if lenp > p+1 {
					for pp := p+1; pp < lenp; pp++ {
						p++
						ult := strings.Split(cmd_parts[c][pp],",")
						for u := range ult {
							param.Context = append(param.Context,ult[u])
						}
					}
					break
				}

			case CMD_FROM:
				if lenp > p+1 {
					for pp := p+1; pp < lenp; pp++ {
						p++
						ult := strings.Split(cmd_parts[c][pp],",")
						for u := range ult {
							param.From = append(param.From,ult[u])
						}
					}
					break
				}

			case CMD_TO:
				if p > 0 && lenp > p+1 {

					if param.From == nil {
						param.From = append(param.From,cmd_parts[c][p-1])
					}

					for pp := p+1; pp < lenp; pp++ {
						p++
						ult := strings.Split(cmd_parts[c][pp],",")
						for u := range ult {
							param.To = append(param.To,ult[u])
						}
					}
					break
				}

				if lenp > p+1 {
					for pp := p+1; pp < lenp; pp++ {
						p++
						ult := strings.Split(cmd_parts[c][pp],",")
						for u := range ult {
							param.To = append(param.To,ult[u])
						}
					}
					break
				}



			case CMD_PATH,CMD_STORY,CMD_SEQ:
				param.Sequence = true

			default:
				//param.Name = append(param.Name,cmd_parts[c][p])

				if lenp > p+1 {

					if cmd_parts[c][p+1] == CMD_TO {
						continue
					}

					for pp := p; pp < lenp; pp++ {
						p++
						ult := strings.Split(cmd_parts[c][pp],",")
						for u := range ult {
							param.Name = append(param.Name,ult[u])
						}
					}
					break
				}
			}
			break
		}
	}

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
	var block_quote bool = false

	quotes := strings.Count(s,"\"")

	if quotes % 2 != 0 {
		fmt.Println("Unpaired quotes in search",s,quotes)
	}

	cmd := []rune(s)

	for r := 0; r < len(cmd); r++ {

		switch cmd[r] {

		case ' ':
			if !block_quote {
				items = append(items,string(upto))
				upto = nil
				continue
			}
			break

		case '"':
			if block_quote {
				items = append(items,string(upto))
				upto = nil
			}
			block_quote = !block_quote
			continue
		}

		upto = append(upto,cmd[r])
	}

	items = append(items,string(upto))
	return items
}

//******************************************************************

func Search(ctx SST.PoSST, search SearchParameters,line string) {

	fmt.Println("FROM",line)
	fmt.Println(" - name:",search.Name)
	fmt.Println(" - from:",search.From)
	fmt.Println(" - to:",search.To)
	fmt.Println(" - chap:",search.Chapter)
	fmt.Println(" - context:",search.Context)
	fmt.Println(" - arrows:",search.Arrows)
	fmt.Println(" - pagenr:",search.PageNr)
	fmt.Println(" - seq:",search.Sequence)
}









