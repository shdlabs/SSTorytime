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

type ConfigParser struct 
{
stream_position  int
mode             int
topic            string   // in case we mix together several
context_set      []string
item_set         []string
relation_set     []string
}

//**************************************************************

type NoteParser struct 
{
pos              int
state            int
alias            string
topic            string   // in case we mix together several
context_set      []string

item_set         []string  // these reset on each line
relation_set     []string  // < item_set
}

//**************************************************************
// Globals
//**************************************************************

var LINE_ITEM_CACHE []string
var LINE_RELN_CACHE []string
var LINE_ITEM_COUNTER int = 1

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

	var p NoteParser
	var token string

	//fmt.Println(string(src))

	for p.pos = 0; p.pos < len(src); {

		p.pos = SkipWhiteSpace(src,p.pos)
		token,p.pos = GetToken(src,p.pos)

		ClassifyTokenRole(token,p.state)

	}
}

//**************************************************************
// Parsing objects
//**************************************************************

func SkipWhiteSpace(src []rune, pos int) int {
	
	for ; pos < len(src) && (unicode.IsSpace(src[pos]) || src[pos] == '#' || src[pos] == '/') ; pos++ {

		if src[pos] == '\n' {
			UpdateLastLineCache() 
		}
		
		if src[pos] == '#' || (src[pos] == '/' && src[pos+1] == '/') {
			
			for ; pos < len(src) && src[pos] != '\n'; pos++ {
			}
			
			UpdateLastLineCache() 
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
			token = strings.TrimSpace(token)
			token = strings.Trim(token,"\"")
		}

	case '#':
		token,pos = ReadToLast(src,pos,'\n')
		UpdateLastLineCache() 

	case '/':
		if src[pos+1] == '/' {
			token,pos = ReadToLast(src,pos,'\n')
			UpdateLastLineCache() 
		}

	default: // a text item that could end with any of the above
		token,pos = ReadToLast(src,pos,'x')
	}

	return token, pos

}

//**************************************************************

func ClassifyTokenRole(token string,state int) {

	if len(token) == 0 {
		return
	}

	switch token[0] {

	case ':':
		expression := ContextExpression(token)
		fmt.Println("context reset:",expression)

	case '+':
		expression := ContextExpression(token)
		fmt.Println("context augmentation:",expression)

	case '-':
		if token[1:2] == string(':') {
			expression := ContextExpression(token)
			fmt.Println("context pruning:",expression)
		} else {
			section := strings.TrimSpace(token[1:])
			fmt.Println("notes section name:",section)
		}

		// No quotes here in a string, we need to allow quoting in excerpts.

	case '(':
		reln := FindAssociation(token)
		fmt.Println("Relationship:",reln)

	case '"':
		fmt.Println("prior-reference: $n",LINE_ITEM_COUNTER)
		LINE_ITEM_COUNTER++

	default:

		fmt.Println("Node item:",token)
		LINE_ITEM_COUNTER++
	}

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

	if stop == 'x' {
		
		collect = IsGeneralString(src,pos)

	} else {
		// a ::: cluster is special, we don't care how many

		if stop != ':' && stop != '"' { 
			return !LastSpecialChar(src,pos,stop)
		} else {
			var groups int = 0

			for r := 1; r < len(cpy); r++ {

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
		UpdateLastLineCache() 
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
		UpdateLastLineCache() 
		return true
	}

	if src[pos-1] == stop && src[pos] != stop {
		return true
	}

	return false
}

//**************************************************************

func UpdateLastLineCache() {

	// reset $n variables

	LINE_ITEM_CACHE = nil
	LINE_RELN_CACHE = nil
	LINE_ITEM_COUNTER = 1

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
// Tools
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
