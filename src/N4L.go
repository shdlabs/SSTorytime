//
// N4LParser
//

package main

import (
	"strings"
	"os"
	"io/ioutil"
	"flag"
	"fmt"
	"unicode/utf8"
	"unicode"
)

//**************************************************************
// Globals
//**************************************************************

const (
	ROLE_EVENT = 1
	ROLE_RELATION = 2
	ROLE_SECTION = 3
	ROLE_CONTEXT = 4
	ROLE_BLANK_LINE = 5
	ROLE_LINE_ALIAS = 6

	ERR_MISSING_EVENT = "Missing item? Dangling section, relation, or context"
	ERR_NO_SUCH_ALIAS = "No such alias or \" reference exists to fill in - aborting"

)

var LINE_NUM int = 1
var LINE_ITEM_CACHE = make(map[string][]string)
var LINE_RELN_CACHE = make(map[string][]string)
var LINE_ITEM_STATE int = ROLE_BLANK_LINE
var LINE_ALIAS string = ""
var LINE_ITEM_COUNTER int = 1
var LINE_RELN_COUNTER int = 1

//**************************************************************

func main() {

	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	
	if len(args) < 1 {
		usage()
		os.Exit(1);
	}

	//config := ReadFile("config.in")
	//ParseConfig(config)

	input := ReadFile(args[0])
	ParseN4L(input)
}


//**************************************************************

func ParseN4L(src []rune) {

	var token string

	for pos := 0; pos < len(src); {

		pos = SkipWhiteSpace(src,pos)
		token,pos = GetToken(src,pos)

		ClassifyTokenRole(token)

	}

	if Dangler() {
		ParseError(ERR_MISSING_EVENT)
	}
}

//**************************************************************
// Parsing objects
//**************************************************************

func SkipWhiteSpace(src []rune, pos int) int {
	
	for ; pos < len(src) && (unicode.IsSpace(src[pos]) || src[pos] == '#' || src[pos] == '/') ; pos++ {

		if src[pos] == '\n' {
			UpdateLastLineCache() 
		} else {
		
			if src[pos] == '#' || (src[pos] == '/' && src[pos+1] == '/') {
				
				for ; pos < len(src) && src[pos] != '\n'; pos++ {
				}

				UpdateLastLineCache() 
			}
		}
	}

	return pos
}

//**************************************************************

func GetToken(src []rune, pos int) (string,int) {

	// Handle concatenation of words/lines and separation of types

	var token string

	if pos == len(src) {
		return "", pos
	}

	switch (src[pos]) {

	case '+':  // could be +:: 

		switch (src[pos+1]) {

		case ':':
			token,pos = ReadToLast(src,pos,':')
		default:
			token,pos = ReadToLast(src,pos,'x')
		}

	case '-':  // could -:: or -section

		switch (src[pos+1]) {

		case ':':
			token,pos = ReadToLast(src,pos,':')
		default:
			token,pos = ReadToLast(src,pos,'x')
		}

	case ':':
		token,pos = ReadToLast(src,pos,':')

	case '(':
		token,pos = ReadToLast(src,pos,')')

        case '"': // only a quoted string must end with the same followed by one of above
		if unicode.IsSpace(src[pos+1]) {
			token = "\""
			pos++
		} else {
			token,pos = ReadToLast(src,pos,'"')
			strip := strings.Split(token,"\"")
			token = strip[1]
		}

	case '#':
		token,pos = ReadToLast(src,pos,'\n')

	case '@':
		token,pos = ReadToLast(src,pos,' ')

	case '/':
		if src[pos+1] == '/' {
			token,pos = ReadToLast(src,pos,'\n')
		}

	default: // a text item that could end with any of the above
		token,pos = ReadToLast(src,pos,'x')

	}

	return token, pos

}

//**************************************************************

func ClassifyTokenRole(token string) {

	if len(token) == 0 {
		return
	}

	switch token[0] {

	case ':':
		expression := ContextExpression(token)
		Role("context reset:",expression)
		LINE_ITEM_STATE = ROLE_CONTEXT

	case '+':
		expression := ContextExpression(token)
		Role("context augmentation:",expression)
		LINE_ITEM_STATE = ROLE_CONTEXT

	case '-':
		if token[1:2] == string(':') {
			expression := ContextExpression(token)
			Role("context pruning:",expression)
			LINE_ITEM_STATE = ROLE_CONTEXT
		} else {
			section := strings.TrimSpace(token[1:])
			Role("notes section name:",section)
			LINE_ITEM_STATE = ROLE_SECTION
		}

		// No quotes here in a string, we need to allow quoting in excerpts.

	case '(':
		reln := FindAssociation(token)
		Role("Relationship:",reln)
		LINE_ITEM_STATE = ROLE_RELATION
		LINE_RELN_CACHE["THIS"] = append(LINE_RELN_CACHE["THIS"],token)

		LINE_RELN_COUNTER++

	case '"':
		Role("prior-reference",LookupAlias("PREV",LINE_ITEM_COUNTER))
		LINE_ITEM_STATE = ROLE_EVENT
		LINE_ITEM_CACHE["THIS"] = append(LINE_ITEM_CACHE["THIS"],LookupAlias("PREV",LINE_ITEM_COUNTER))

		LINE_ITEM_COUNTER++

	case '@':
		Role("line-alias",token)
		LINE_ITEM_STATE = ROLE_LINE_ALIAS
		LINE_ALIAS = token[1:]

	case '$':
		Role("variable-reference",token)
		actual := HandleAliasedItem(token)
		fmt.Println("...resolved",actual)

	default:
		Role("Event item:",token)
		LINE_ITEM_STATE = ROLE_EVENT

		// need to check if we have () embedded between items or missing....

		LINE_ITEM_CACHE["THIS"] = append(LINE_ITEM_CACHE["THIS"],token)
		if LINE_ALIAS != "" {
			LINE_ITEM_CACHE[LINE_ALIAS] = append(LINE_ITEM_CACHE["THIS"],token)
		}
		LINE_ITEM_COUNTER++
	}

	//To do, store classified parts for grammar rules

	AssessGrammarCompletions()

}

//**************************************************************
// Scan text input
//**************************************************************

func ReadFile(filename string) []rune {

	return ReadTUF8File(filename)
}


//**************************************************************

func ReadToLast(src []rune,pos int, stop rune) (string,int) {

	var cpy []rune

	for ; pos > 0 && Collect(src,pos,stop,cpy); pos++ {

		cpy = append(cpy,src[pos])
	}

	token := string(cpy)
	
	return token,pos
}

//**************************************************************

func Collect(src []rune,pos int, stop rune,cpy []rune) bool {

	var collect bool = true

	if src[pos] == '\n' {
		return false
	}

	if stop == 'x' {
		collect = IsGeneralString(src,pos)
	} else {
		// a ::: cluster is special, we don't care how many

		if stop != ':' && stop != '"' { 
			return !LastSpecialChar(src,pos,stop)
		} else {
			var groups int = 0

			for r := 1; r < len(cpy)-1; r++ {

				if cpy[r] != ':' && cpy[r-1] == ':' {
					groups++
				}

				if cpy[r] != '"' && cpy[r-1] == '"' {
					groups++
				}
			}

			if groups > 1 {
				collect = !LastSpecialChar(src,pos,stop)
			}
		} 
	}

	return collect
}

//**************************************************************

func IsGeneralString(src []rune,pos int) bool {

	switch src[pos] {

	case '(':
		return false
	case '#':
		return false
	case '\n':
		return false

	case '/':
		if src[pos+1] == '/' {
			return false
		}
	}

	return true
}

//**************************************************************

func LastSpecialChar(src []rune,pos int, stop rune) bool {

	if src[pos] == '\n' {
		return true
	}

	// tabs are divisors, don't use them!

	if (src[pos-1] == stop || src[pos-1] == '\t') && src[pos] != stop {
		return true
	}

	return false
}

//**************************************************************

func UpdateLastLineCache() {

	if Dangler() {
		ParseError(ERR_MISSING_EVENT)
	}

// check if len(itemcache) = len(relcache)+1 --- something wrong

	LINE_NUM++

	// If this line was not blank, overwrite previous settings and reset

	if LINE_ITEM_STATE != ROLE_BLANK_LINE {
		
		if LINE_ITEM_CACHE["THIS"] != nil {
			LINE_ITEM_CACHE["PREV"] = LINE_ITEM_CACHE["THIS"]
		}
		if LINE_RELN_CACHE["THIS"] != nil {
			LINE_RELN_CACHE["PREV"] = LINE_RELN_CACHE["THIS"]
		}

	} 

	LINE_ITEM_CACHE["THIS"] = nil
	LINE_RELN_CACHE["THIS"] = nil
	LINE_ITEM_COUNTER = 1
	LINE_ALIAS = ""

	LINE_ITEM_STATE = ROLE_BLANK_LINE
}

//**************************************************************

func Dangler() bool {

	switch LINE_ITEM_STATE {

	case ROLE_EVENT:
		return false
	case ROLE_BLANK_LINE:
		return false
	case ROLE_SECTION:
		return false
	case ROLE_CONTEXT:
		return false
	}

	return true
}

//**************************************************************

func ContextExpression(token string) string {

	var expression string

	s := strings.Split(token, ":")

	for i := 1; i < len(s); i++ {
		if len(s[i]) > 1 {
			expression = strings.TrimSpace(s[i])
			break
		}
	}
	
	return expression
}

//**************************************************************

func FindAssociation(token string) string {

	name := token[1:len(token)-1]

	// lookup in the alias table

	return strings.TrimSpace(name)
}

//**************************************************************

func LookupAlias(alias string, counter int) string {
	
	value,ok := LINE_ITEM_CACHE[alias]

	if !ok || counter > len(value) {
		ParseError(ERR_NO_SUCH_ALIAS)
		os.Exit(1)
	}
	
	return LINE_ITEM_CACHE[alias][counter-1]

}

//**************************************************************

func HandleAliasedItem(token string) string {

 // construct lookup from $alias.n or $n
 // if define

	return ""
}

//**************************************************************

func AssessGrammarCompletions() {

 //using foreach in LINE_ITEM_CACHE["THIS"]  LINE_RELN_CACHE["THIS"]

	// completions
        // @alias --> set current in hash table @alias.$1 etc ... until \n (default alias = "")
	// :: :: ---> set current until next :: ::
	// ITEM1 ITEM2 -> install first item and keep ITEM2++
        // ITEM1 (reln) ITEM2 ITEM3 --> install whole relation and keep ITEM3++
        // ITEM1 (reln) ITEM2 (reln2) --> install whole relation and keep ITEM2++

	// item, context (special case of) item
	// (item1,context) (reln,context) (item2,context)

}

//**************************************************************
// Tools
//**************************************************************

func Role(role,item string) {

	fmt.Println(LINE_NUM,":",role,item)

}

//**************************************************************

func ParseError(message string) {

	fmt.Println("N4L",message,"at line", LINE_NUM)
}

//**************************************************************

func ReadTUF8File(filename string) []rune {
	
	content, _ := ioutil.ReadFile(filename)
	
	var unicode []rune
	
	for i, w := 0, 0; i < len(content); i += w {
                runeValue, width := utf8.DecodeRuneInString(string(content)[i:])
                w = width

		unicode = append(unicode,runeValue)
	}
	
	return unicode
}

//**************************************************************

func usage() {
	
	fmt.Fprintf(os.Stderr, "usage: go run N4L.go [file].dat\n")
	flag.PrintDefaults()
	os.Exit(2)
}
