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
	"regexp"
)

//**************************************************************
// Global state
//**************************************************************

const (
	ROLE_EVENT = 1
	ROLE_RELATION = 2
	ROLE_SECTION = 3
	ROLE_CONTEXT = 4
	ROLE_CONTEXT_ADD = 5
	ROLE_CONTEXT_SUBTRACT = 6
	ROLE_BLANK_LINE = 7
	ROLE_LINE_ALIAS = 8
	ROLE_LOOKUP = 9

	ERR_NO_SUCH_FILE_FOUND = "No file found in the name "
	ERR_MISSING_EVENT = "Missing item? Dangling section, relation, or context"
	ERR_MISSING_SECTION = "Declarations outside a section or chapter"
	ERR_NO_SUCH_ALIAS = "No such alias or \" reference exists to fill in - aborting"
	ERR_MISSING_ITEM_SOMEWHERE = "Missing item somewhere"

	WARN_NOTE_TO_SELF = "WARNING: Found a note to self in the text"
	WARN_INADVISABLE_CONTEXT_EXPRESSION = "WARNING: Inadvisably complex/parenthetic context expression - simplify?"
)

var ( 
	LINE_NUM int = 1
	LINE_ITEM_CACHE = make(map[string][]string)
	LINE_RELN_CACHE = make(map[string][]string)
	LINE_ITEM_STATE int = ROLE_BLANK_LINE
	LINE_ALIAS string = ""
	LINE_ITEM_COUNTER int = 1
	LINE_RELN_COUNTER int = 0
	CONTEXT_STATE = make(map[string]bool)
	SECTION_STATE string
	VERBOSE bool = false
	CURRENT_FILE string
)

//**************************************************************

func main() {

	args := Init()

	//config := ReadFile("config.in")
	//ParseConfig(config)

	for input := 0; input < len(args); input++ {

		CURRENT_FILE = args[input]
		input := ReadFile(CURRENT_FILE)
		ParseN4L(input)
	}
}


//**************************************************************

func Init() []string {

	flag.Usage = usage
	verbosePtr := flag.Bool("v", false,"verbose")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		usage()
		os.Exit(1);
	}

	if *verbosePtr {
		VERBOSE = true
	}

	return args
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
		return "",pos

	case '/':
		if src[pos+1] == '/' {
			return "",pos
		}

	case '@':
		token,pos = ReadToLast(src,pos,' ')


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
		expression := ExtractContextExpression(token)
		Role("context reset:",expression)
		LINE_ITEM_STATE = ROLE_CONTEXT
		AssessGrammarCompletions(expression,LINE_ITEM_STATE)

	case '+':
		expression := ExtractContextExpression(token)
		Role("context augmentation:",expression)
		LINE_ITEM_STATE = ROLE_CONTEXT_ADD
		AssessGrammarCompletions(expression,LINE_ITEM_STATE)

	case '-':
		if token[1:2] == string(':') {
			expression := ExtractContextExpression(token)
			Role("context pruning:",expression)
			LINE_ITEM_STATE = ROLE_CONTEXT_SUBTRACT
			AssessGrammarCompletions(expression,LINE_ITEM_STATE)
		} else {
			section := strings.TrimSpace(token[1:])
			Role("notes section name:",section)
			LINE_ITEM_STATE = ROLE_SECTION
			AssessGrammarCompletions(section,LINE_ITEM_STATE)
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
		result := LookupAlias("PREV",LINE_ITEM_COUNTER)
		LINE_ITEM_CACHE["THIS"] = append(LINE_ITEM_CACHE["THIS"],result)
		AssessGrammarCompletions(result,LINE_ITEM_STATE)
		LINE_ITEM_COUNTER++

	case '@':
		Role("line-alias",token)
		LINE_ITEM_STATE = ROLE_LINE_ALIAS
		token  = strings.TrimSpace(token)
		LINE_ALIAS = token[1:]

	case '$':
		Role("variable-reference",token)
		actual := ResolveAliasedItem(token)
		Verbose(LINE_NUM,": -",token,"...resolved to",actual)

		AssessGrammarCompletions(actual,LINE_ITEM_STATE)
		LINE_ITEM_STATE = ROLE_LOOKUP

	default:
		Role("Event item:",token)

		LINE_ITEM_CACHE["THIS"] = append(LINE_ITEM_CACHE["THIS"],token)

		if LINE_ALIAS != "" {
			LINE_ITEM_CACHE[LINE_ALIAS] = append(LINE_ITEM_CACHE[LINE_ALIAS],token)
			//Verbose(LINE_ALIAS,"=",LINE_ITEM_CACHE[LINE_ALIAS],len(LINE_ITEM_CACHE[LINE_ALIAS]))
		}

		AssessGrammarCompletions(token,LINE_ITEM_STATE)

		LINE_ITEM_STATE = ROLE_EVENT
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

	for ; Collect(src,pos,stop,cpy); pos++ {

		// sanitize small case at start of item
		if stop == 'x' {
			src[pos] = unicode.ToLower(src[pos])
		}
		cpy = append(cpy,src[pos])
	}

	token := string(cpy)
	token = strings.TrimSpace(token)

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

	if pos > 0 && (src[pos-1] == stop || src[pos-1] == '\t') && src[pos] != stop {
		return true
	}

	return false
}

//**************************************************************

func UpdateLastLineCache() {

	if Dangler() {
		ParseError(ERR_MISSING_EVENT)
	}

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
	LINE_RELN_COUNTER = 0
	LINE_ALIAS = ""

	LINE_ITEM_STATE = ROLE_BLANK_LINE
}

//**************************************************************

func Dangler() bool {

	switch LINE_ITEM_STATE {

	case ROLE_EVENT:
		return false
	case ROLE_LOOKUP:
		return false
	case ROLE_BLANK_LINE:
		return false
	case ROLE_SECTION:
		return false
	case ROLE_CONTEXT:
		return false
	case ROLE_CONTEXT_ADD:
		return false
	case ROLE_CONTEXT_SUBTRACT:
		return false
	}

	return true
}

//**************************************************************

func ExtractContextExpression(token string) string {

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

func ResolveAliasedItem(token string) string {

	// split $alias.n into (alias string,n int)

	split := strings.Split(token[1:],".")
	name := strings.TrimSpace(split[0])

	var number int = 0
	fmt.Sscanf(split[1],"%d",&number)

	return LookupAlias(name,number)
}

//**************************************************************

func AssessGrammarCompletions(token string, prior_state int) {

 //using foreach in LINE_ITEM_CACHE["THIS"]  LINE_RELN_CACHE["THIS"]

	// completions
        // @alias --> set current in hash table @alias.$1 etc ... until \n (default alias = "")
	// :: :: ---> set current until next :: ::
	// ITEM1 ITEM2 -> install first item and keep ITEM2++
        // ITEM1 (reln) ITEM2 ITEM3 --> install whole relation and keep ITEM3++
        // ITEM1 (reln) ITEM2 (reln2) --> install whole relation and keep ITEM2++

	// item, context (special case of) item
	// (item1,context) (reln,context) (item2,context)

	if AllCaps(token) {
		ParseError(WARN_NOTE_TO_SELF+" ("+token+")")
		return
	}

	this_item := token

	switch prior_state {

	case ROLE_RELATION:

		CheckNonNegative(LINE_ITEM_COUNTER-2)
		last_item := LINE_ITEM_CACHE["THIS"][LINE_ITEM_COUNTER-2]
		last_reln := LINE_RELN_CACHE["THIS"][LINE_RELN_COUNTER-1]
		Verbose("New item/event:",this_item)
		Verbose("... Relation:",last_item,"--",last_reln,"->",this_item)
		CheckSection()

	case ROLE_CONTEXT:
		Box("Fresh context: ->",this_item)
		ContextEval(this_item,"=")

	case ROLE_CONTEXT_ADD:
		Verbose("Add to context:",this_item)
		ContextEval(this_item,"+")

	case ROLE_CONTEXT_SUBTRACT:
		Verbose("Remove from context:",this_item)
		ContextEval(this_item,"-")

	case ROLE_SECTION:
		Box("Set chapter/section: ->",this_item)
		SECTION_STATE = this_item

	default:
		Verbose("New item/event:",this_item)
		CheckSection()
	}
}

//**************************************************************
// Context logic
//**************************************************************

func ContextEval(s,op string) {

	expr := CleanExpression(s)

	or_parts := SplitWithParensIntact(expr,'|')

	if strings.Contains(s,"(") {
		ParseError(WARN_INADVISABLE_CONTEXT_EXPRESSION)
	}

	// +,-,= on CONTEXT_STATE

	switch op {
		
	case "=": 
		CONTEXT_STATE = make(map[string]bool)
		ModContext(or_parts,"+")
	default:
		ModContext(or_parts,op)
	}
}

// ***********************************************************************

func CleanExpression(s string) string {

	s = TrimParen(s)
	r1 := regexp.MustCompile("[|,]+") 
	s = r1.ReplaceAllString(s,"|") 
	r2 := regexp.MustCompile("[&]+") 
	s = r2.ReplaceAllString(s,".") 
	r3 := regexp.MustCompile("[.]+") 
	s = r3.ReplaceAllString(s,".") 

	return s
}

// ***********************************************************************

func SplitWithParensIntact(expr string,split_ch byte) []string {

	var token string = ""
	var set []string

	for c := 0; c < len(expr); c++ {

		switch expr[c] {

		case split_ch:
			set = append(set,token)
			token = ""

		case '(':
			subtoken,offset := Paren(expr,c)
			token += subtoken
			c = offset-1

		default:
			token += string(expr[c])
		}
	}

	if len(token) > 0 {
		set = append(set,token)
	}

	return set
} 

// ***********************************************************************

func Paren(s string, offset int) (string,int) {

	var level int = 0

	for c := offset; c < len(s); c++ {

		if s[c] == '(' {
			level++
			continue
		}

		if s[c] == ')' {
			level--
			if level == 0 {
				token := s[offset:c+1]
				return token, c+1
			}
		}
	}

	return "bad expression", -1
}


// ***********************************************************************

func TrimParen(s string) string {

	var level int = 0
	var trim = true

	if len(s) == 0 {
		return s
	}

	s = strings.TrimSpace(s)

	if s[0] != '(' {
		return s
	}

	for c := 0; c < len(s); c++ {

		if s[c] == '(' {
			level++
			continue
		}

		if level == 0 && c < len(s)-1 {
			trim = false
		}
		
		if s[c] == ')' {
			level--

			if level == 0 && c == len(s)-1 {
				
				var token string
				
				if trim {
					token = s[1:len(s)-1]
				} else {
					token = s
				}
				return token
			}
		}
	}
	
	return s
}

//**************************************************************

func ModContext(list []string,op string) {

	for or_frag := range list {

		frag := strings.TrimSpace(list[or_frag])

		if len(frag) == 0 {
			continue
		}

		switch op {
		case "+":
			CONTEXT_STATE[frag] = true

		case "-": // to remove, we also need to look at children
			for cand := range CONTEXT_STATE {
				and_parts := SplitWithParensIntact(cand,'.') 
		
				for part := range and_parts {

					if strings.Contains(and_parts[part],frag) {
						delete(CONTEXT_STATE,cand)
					}
				}
			}
		}

	}
}

//**************************************************************

func CheckNonNegative(i int) {

	if i < 0 {
		ParseError(ERR_MISSING_ITEM_SOMEWHERE)
		os.Exit(-1)
	}
}

//**************************************************************

func CheckSection() {

	if len(SECTION_STATE) == 0 {
		ParseError(ERR_MISSING_SECTION)
		os.Exit(-1)
	}
}

//**************************************************************

func AllCaps(s string) bool {

	if len(s) < 3 {
		return false
	}

	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

//**************************************************************
// Tools
//**************************************************************

func Role(role,item string) {

	//Verbose(LINE_NUM,":",role,item)

}

//**************************************************************

func ParseError(message string) {

	fmt.Print(LINE_NUM,":")
	fmt.Println("N4L",CURRENT_FILE,message,"at line", LINE_NUM)
}

//**************************************************************

func ReadTUF8File(filename string) []rune {
	
	content,err := ioutil.ReadFile(filename)
	
	if err != nil {
		ParseError(ERR_NO_SUCH_FILE_FOUND+filename)
		os.Exit(-1)
	}

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
	
	fmt.Fprintf(os.Stderr, "usage: go run N4L.go [-v] [file].dat\n")
	flag.PrintDefaults()
	os.Exit(2)
}

//**************************************************************

func Verbose(a ...interface{}) (n int, err error) {

	if VERBOSE {
		fmt.Print(LINE_NUM,":")
		n, err = fmt.Println(a...)
	}
	return
}

//**************************************************************

func Box(a ...interface{}) (n int, err error) {

	if VERBOSE {

		fmt.Println("------------------------------------")
		fmt.Print(LINE_NUM,":")
		fmt.Println(a...)
		fmt.Println("------------------------------------")
	}
	return
}