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
// DATA structures
//**************************************************************

/* 1. We want every reusable text to be addressible by a pointer/index
      in a simple lookup table. This includes
        - item/event nodes
        - arrows referenced by full text or a short alias in N4L
        - context expressions

   2. We want relations/arrows to be (Arrow ArrowPtr,To EventItemPtr,Context CtxPtr)
      that are attached to a (From EventItemPtr) record as an array of lists,
      indexed by SST types 0-3

   3. For arrows, short aliases are often used, but longer names may be used too
      A Directory maps names to ArrowPtr
*/

//**************************************************************

type EventItemPtr int   // EventItem index
type ArrowPtr int       // ArrowDirectory index
type CtxPtr int
type NodeTextPtr int

//**************************************************************

const (
	SSTtypes = 4
	LEADSTO = 1
	CONTAINS = 2
	EXPRESS = 3
	NEAR = 4

	N1GRAM = 1
	N2GRAM = 2
	N3GRAM = 3
	LT128 = 4
	LT1024 = 5
	GT1024 = 6
)

//**************************************************************

type EventItem struct {

	Name NodeTextPtr               // use a reference to avoid loading a large object
	Links [SSTtypes][]ArrowPtr
}

//**************************************************************

type Text struct {
	L  int
	S  string
	C int
}

//**************************************************************

type ArrowRelation struct {       // used only to fill in output text

	LFrom    string
	LTo      string
	SFrom    string
	STo      string
	SSTType  int
}

//**************************************************************

type ArrowDirectory struct {

	STtype  int
	Long    string
	Short   string
	Ptr     ArrowPtr
}

 // all fwd arrow types have a simple int representation > 0
 // all bwd/inverse arrow readings have the negative int for fwd
 // Hashed by long and short names

//**************************************************************

type NodeTextBlobs struct {

	// Power law n-gram frequencies

	N1grams map[string]NodeTextPtr
	N1directory []Text
	N1_top NodeTextPtr

	N2grams map[string]NodeTextPtr
	N2directory []Text
	N2_top NodeTextPtr

	N3grams map[string]NodeTextPtr
	N3directory []Text
	N3_top NodeTextPtr

	// Use linear search on these exp fewer long strings

	LT128 []Text
	LT128_top NodeTextPtr
	LT1024 []Text
	LT1024_top NodeTextPtr
	GT1024 []Text
	GT1024_top NodeTextPtr
}

//**************************************************************

var ( 
	ARROW_DIRECTORY []ArrowDirectory
	ARROW_SHORT_DIR = make(map[string]ArrowPtr) // Look up short name int referene
	ARROW_LONG_DIR = make(map[string]ArrowPtr)  // Look up long name int referene
	ARROW_DIRECTORY_TOP ArrowPtr = 0

	TEXT_DIRECTORY NodeTextBlobs
)

//**************************************************************
// Global parsing state
//**************************************************************

const (
	ALPHATEXT = 'x'

        HAVE_PLUS = 11
        HAVE_MINUS = 22
	ROLE_ABBR = 33

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
	ERR_NO_SUCH_ARROW = "No such arrow has been declared in the configuration: "
	ERR_MISSING_ITEM_SOMEWHERE = "Missing item somewhere"
	ERR_MISSING_ITEM_RELN = "Missing item or double relation"
	ERR_ILLEGAL_CONFIGURATION = "Error in configuration, no such section"
	ERR_BAD_LABEL_OR_REF = "Badly formed label or reference (@label becomes $label.n) in "
	WARN_NOTE_TO_SELF = "WARNING: Found a note to self in the text"
	WARN_INADVISABLE_CONTEXT_EXPRESSION = "WARNING: Inadvisably complex/parenthetic context expression - simplify?"
	WARN_ILLEGAL_QUOTED_STRING_OR_REF = "WARNING: Something wrong, bad quoted string or mistaken back reference"
	ERR_ANNOTATION_TOO_LONG = "Annotation marker should be a single non-alphnumeric character "
	ERR_BAD_ABBRV = "abbreviation out of place"
	ERR_BAD_ALIAS_REFERENCE = "Alias references start from $name.1"
	ERR_ANNOTATION_MISSING = "Missing non-alphnumeric annotation marker or stray relation"
	ERR_ANNOTATION_REDEFINE = "Redefinition of annotation character"
	ERR_SIMILAR_NO_SIGN = "Arrows for similarity do not have signs, they are directionless"
	ERR_ARROW_SELFLOOP = "Arrow's origin points to itself"
)

var ( 
	LINE_NUM int = 1
	LINE_ITEM_CACHE = make(map[string][]string)
	LINE_RELN_CACHE = make(map[string][]string)
	LINE_ITEM_STATE int = ROLE_BLANK_LINE
	LINE_ALIAS string = ""
	LINE_ITEM_COUNTER int = 1
	LINE_RELN_COUNTER int = 0

	FWD_ARROW string
	BWD_ARROW string
	ANNOTATION = make(map[string]string)

	CONTEXT_STATE = make(map[string]bool)
	SECTION_STATE string

	SEQUENCE_MODE bool = false
	SEQUENCE_RELN string = "then" 
	LAST_IN_SEQUENCE string = ""

	VERBOSE bool = false
	CURRENT_FILE string
	TEST_DIAG_FILE string

	RELN_BY_SST [4][]ArrowPtr // From an EventItemNode
)

//**************************************************************

func main() {

	args := Init()

	NewFile("N4Lconfig.in")
	config := ReadFile(CURRENT_FILE)
	ParseConfig(config)

	//SummarizeAndTestConfig()

	for input := 0; input < len(args); input++ {

		NewFile(args[input])
		input := ReadFile(CURRENT_FILE)
		ParseN4L(input)
	}

	fmt.Println("\nTEXT\n\n",TEXT_DIRECTORY)
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

	TEXT_DIRECTORY.N1grams = make(map[string]NodeTextPtr)
	TEXT_DIRECTORY.N2grams = make(map[string]NodeTextPtr)
	TEXT_DIRECTORY.N3grams = make(map[string]NodeTextPtr)

	return args
}

//**************************************************************

func NewFile(filename string) {

	CURRENT_FILE = filename
	TEST_DIAG_FILE = DiagnosticName(filename)

	Box("Parsing new file",filename)

	LINE_ITEM_STATE = ROLE_BLANK_LINE
	LINE_NUM = 1
	LINE_ITEM_CACHE["THIS"] = nil
	LINE_RELN_CACHE["THIS"] = nil
	LINE_ITEM_COUNTER = 1
	LINE_RELN_COUNTER = 0
	LINE_ALIAS = ""
	FWD_ARROW = ""
	BWD_ARROW = ""
}

//**************************************************************
// N4L configuration
//**************************************************************

func ParseConfig(src []rune) {

	var token string

	for pos := 0; pos < len(src); {

		pos = SkipWhiteSpace(src,pos)
		token,pos = GetConfigToken(src,pos)

		ClassifyConfigRole(token)
	}
}

//**************************************************************

func GetConfigToken(src []rune, pos int) (string,int) {

	// Handle concatenation of words/lines and separation of types

	var token string

	if pos >= len(src) {
		return "", pos
	}

	switch (src[pos]) {

	case '+':
		token,pos = ReadToLast(src,pos,ALPHATEXT)

	case '-':
		token,pos = ReadToLast(src,pos,ALPHATEXT)

	case '(':
		token,pos = ReadToLast(src,pos,')')  // alias

	case '#':
		return "",pos

	case '/':
		if src[pos+1] == '/' {
			return "",pos
		}

	default: // similarity
		token,pos = ReadToLast(src,pos,ALPHATEXT)

	}

	return token, pos
}

//**************************************************************

func ClassifyConfigRole(token string) {

	if len(token) == 0 {
		return
	}

	if token[0] == '-' && LINE_ITEM_STATE == ROLE_BLANK_LINE {
		SECTION_STATE = strings.TrimSpace(token[1:])
		Box("Configuration of",SECTION_STATE)
		LINE_ITEM_STATE = ROLE_SECTION
		return
	}

	switch SECTION_STATE {

	case "leadsto","contains","properties":

		switch token[0] {

		case '+':
			FWD_ARROW = strings.TrimSpace(token[1:])
			LINE_ITEM_STATE = HAVE_PLUS
			Diag("fwd arrow in",SECTION_STATE, token)
			
		case '-':
			BWD_ARROW = strings.TrimSpace(token[1:])
			LINE_ITEM_STATE = HAVE_MINUS
			Diag("bwd arrow in",SECTION_STATE, token)

		case '(':
			reln := token[1:len(token)-1]
			reln = strings.TrimSpace(reln)

			if LINE_ITEM_STATE == HAVE_MINUS {
				InsertArrowDirectory(SECTION_STATE,reln,BWD_ARROW,"-")
			} else if LINE_ITEM_STATE == HAVE_PLUS {
				InsertArrowDirectory(SECTION_STATE,reln,FWD_ARROW,"+")
			} else {
				ParseError(ERR_BAD_ABBRV)
				os.Exit(-1)
			}
		}

	case "similarity":

		switch token[0] {

		case '(':
			reln := token[1:len(token)-1]
			reln = strings.TrimSpace(reln)

			if LINE_ITEM_STATE == HAVE_MINUS {
				InsertArrowDirectory(SECTION_STATE,reln,BWD_ARROW,"both")
			} else {
				PVerbose(SECTION_STATE,"abbreviation out of place")
			}

		case '+','-':
			ParseError(ERR_SIMILAR_NO_SIGN)
			os.Exit(-1)

		default:
			similarity := strings.TrimSpace(token)
			FWD_ARROW = similarity
			BWD_ARROW = similarity
			LINE_ITEM_STATE = HAVE_MINUS
		}

	case "annotations":

		switch token[0] {

		case '(':
			if LINE_ITEM_STATE != HAVE_PLUS {
				ParseError(ERR_ANNOTATION_MISSING)
			}

			FWD_ARROW = StripParen(token)
			PVerbose("Annotation marker",LAST_IN_SEQUENCE,"defined as arrow:",FWD_ARROW)

			value,defined := ANNOTATION[LAST_IN_SEQUENCE]

			if defined && value != FWD_ARROW {
				ParseError(ERR_ANNOTATION_REDEFINE)
				os.Exit(-1)
			}

			ANNOTATION[LAST_IN_SEQUENCE] = FWD_ARROW
			LINE_ITEM_STATE = ROLE_BLANK_LINE
			
		default:
			if (len(token) > 1 || unicode.IsLetter(rune(token[0]))) {
				ParseError(ERR_ANNOTATION_TOO_LONG)
			}
			Diag("Markup character defined in",SECTION_STATE, token)
			LINE_ITEM_STATE = HAVE_PLUS
			LAST_IN_SEQUENCE = token

		}

	default:
		ParseError(ERR_ILLEGAL_CONFIGURATION+" "+SECTION_STATE)
		os.Exit(-1)
	}
}

//**************************************************************

func InsertArrowDirectory(sec,alias,name,pm string) {

	PVerbose("In",sec,"short name",alias,"for",name,", direction",pm)

	var sign int

	switch pm {
	case "+":
		sign = 1
	case "-":
		sign = -1
	}

	var newarrow ArrowDirectory

	switch sec {
	case "leadsto":
		newarrow.STtype = LEADSTO * sign
	case "contains":
		newarrow.STtype = CONTAINS * sign
	case "properties":
		newarrow.STtype = EXPRESS * sign
	case "similarity":
		newarrow.STtype = NEAR
	}

	newarrow.Long = name
	newarrow.Short = alias
	newarrow.Ptr = ARROW_DIRECTORY_TOP

	ARROW_DIRECTORY = append(ARROW_DIRECTORY,newarrow)
	ARROW_SHORT_DIR[alias] = ARROW_DIRECTORY_TOP
	ARROW_LONG_DIR[name] = ARROW_DIRECTORY_TOP
	ARROW_DIRECTORY_TOP++
}

//**************************************************************

func SummarizeAndTestConfig() {

	Box("Raw Summary")
	fmt.Println("..\n")
	fmt.Println("ANNOTATION MARKS", ANNOTATION)
	fmt.Println("..\n")
	fmt.Println("DIRECTORY", ARROW_DIRECTORY)
	fmt.Println("..\n")
	fmt.Println("SHORT",ARROW_SHORT_DIR)
	fmt.Println("..\n")
	fmt.Println("LONG",ARROW_LONG_DIR)
	
	fmt.Println("..\n")
	relation := "pe"
	p := ARROW_SHORT_DIR[relation]
	all := ARROW_DIRECTORY[p]

	fmt.Println("Short form",relation,"with ptr",p,"means",all)

	relation = "pinyin for english"
	p = ARROW_LONG_DIR[relation]
	all = ARROW_DIRECTORY[p]

	fmt.Println("Long form",relation,"with ptr",p,"means",all)
}

//**************************************************************
// N4L language
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

func SkipWhiteSpace(src []rune, pos int) int {

	for ; pos < len(src) && IsWhiteSpace(src[pos],src[pos]); pos++ {

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

	if pos >= len(src) {
		return "", pos
	}

	switch (src[pos]) {

	case '+':  // could be +:: 

		switch (src[pos+1]) {

		case ':':
			token,pos = ReadToLast(src,pos,':')
		default:
			token,pos = ReadToLast(src,pos,ALPHATEXT)
		}

	case '-':  // could -:: or -section

		switch (src[pos+1]) {

		case ':':
			token,pos = ReadToLast(src,pos,':')
		default:
			token,pos = ReadToLast(src,pos,ALPHATEXT)
		}

	case ':':
		token,pos = ReadToLast(src,pos,':')

	case '(':
		token,pos = ReadToLast(src,pos,')')

        case '"','\'':
		quote := src[pos]
		if quote == '"' && IsBackReference(src,pos) {
			token = "\""
			pos++
		} else {
			if pos+2 < len(src) && IsWhiteSpace(src[pos+1],src[pos+2]) {
				ParseError(WARN_ILLEGAL_QUOTED_STRING_OR_REF)
			}
			token,pos = ReadToLast(src,pos,quote)
			strip := strings.Split(token,string(quote))
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
		token,pos = ReadToLast(src,pos,ALPHATEXT)

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
		CheckSequenceMode(expression,'+')
		LINE_ITEM_STATE = ROLE_CONTEXT
		AssessGrammarCompletions(expression,LINE_ITEM_STATE)

	case '+':
		expression := ExtractContextExpression(token)
		CheckSequenceMode(expression,'+')
		LINE_ITEM_STATE = ROLE_CONTEXT_ADD
		AssessGrammarCompletions(expression,LINE_ITEM_STATE)

	case '-':
		if token[1:2] == string(':') {
			expression := ExtractContextExpression(token)
			CheckSequenceMode(expression,'-')
			LINE_ITEM_STATE = ROLE_CONTEXT_SUBTRACT
			AssessGrammarCompletions(expression,LINE_ITEM_STATE)
		} else {
			section := strings.TrimSpace(token[1:])
			LINE_ITEM_STATE = ROLE_SECTION
			AssessGrammarCompletions(section,LINE_ITEM_STATE)
		}

		// No quotes here in a string, we need to allow quoting in excerpts.

	case '(':
		if LINE_ITEM_STATE == ROLE_RELATION {
			ParseError(ERR_MISSING_ITEM_RELN)
			os.Exit(-1)
		}
		reln,_ := FindAssociation(token)
		LINE_ITEM_STATE = ROLE_RELATION
		LINE_RELN_CACHE["THIS"] = append(LINE_RELN_CACHE["THIS"],reln)
		LINE_RELN_COUNTER++

	case '"': // prior reference
		result := LookupAlias("PREV",LINE_ITEM_COUNTER)
		LINE_ITEM_CACHE["THIS"] = append(LINE_ITEM_CACHE["THIS"],result)
		AssessGrammarCompletions(result,LINE_ITEM_STATE)
		LINE_ITEM_STATE = ROLE_EVENT
		LINE_ITEM_COUNTER++

	case '@':
		LINE_ITEM_STATE = ROLE_LINE_ALIAS
		token  = strings.TrimSpace(token)
		LINE_ALIAS = token[1:]
		CheckLineAlias(LINE_ALIAS)

	case '$':
		CheckLineAlias(token[1:])
		actual := ResolveAliasedItem(token)
		LINE_ITEM_CACHE["THIS"] = append(LINE_ITEM_CACHE["THIS"],actual)
		Verbose("fyi, line reference",token,"resolved to",actual)
		AssessGrammarCompletions(actual,LINE_ITEM_STATE)
		LINE_ITEM_STATE = ROLE_LOOKUP
		LINE_ITEM_COUNTER++

	default:
		LINE_ITEM_CACHE["THIS"] = append(LINE_ITEM_CACHE["THIS"],token)

		if LINE_ALIAS != "" {
			LINE_ITEM_CACHE[LINE_ALIAS] = append(LINE_ITEM_CACHE[LINE_ALIAS],token)
		}

		AssessGrammarCompletions(token,LINE_ITEM_STATE)

		LINE_ITEM_STATE = ROLE_EVENT
		LINE_ITEM_COUNTER++
	}
}

//**************************************************************

func AssessGrammarCompletions(token string, prior_state int) {

	if len(token) == 0 {
		return
	}

	this_item := token

	switch prior_state {

	case ROLE_RELATION:

		CheckNonNegative(LINE_ITEM_COUNTER-2)
		last_item := LINE_ITEM_CACHE["THIS"][LINE_ITEM_COUNTER-2]
		last_reln := LINE_RELN_CACHE["THIS"][LINE_RELN_COUNTER-1]
		IdempAddText(this_item)
		IdempAddArrow(last_item,last_reln,this_item)
		CheckSection()

	case ROLE_CONTEXT:
		Box("Reset context: ->",this_item)
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
		if AllCaps(token) {
			ParseError(WARN_NOTE_TO_SELF+" ("+token+")")
			return
		}

		IdempAddText(this_item)
		CheckSection()
		LinkSequence(token)

	}
}

//**************************************************************

func CheckLineAlias(s string) {

	if strings.Contains(s,"$") || strings.Contains(s,"@") {
		ParseError(ERR_BAD_LABEL_OR_REF+s)
		os.Exit(-1)
	}

}

//**************************************************************

func IdempAddArrow(from,arrowptr,to string) {

	if from == to {
		ParseError(ERR_ARROW_SELFLOOP)
		os.Exit(-1)
	}

	Verbose("... Relation:",from,"--",arrowptr,"->",to)

// NEW VARIABLE TO HOLD list of links

// CONTEXT..
}

//**************************************************************

func IdempAddText(s string) NodeTextPtr {

	Verbose("Event/item:",s)

	l := len(s)
	c := ClassifyString(s,l)

	var newtext Text
	newtext.S = s
	newtext.L = l
	newtext.C = c

	return AppendTextToDirectory(newtext)
}

//**************************************************************

func ClassifyString(s string,l int) int {

	var spaces int = 0

	for i := 0; i < l; i++ {

		if s[i] == ' ' {
			spaces++
		}

		if spaces > 2 {
			break
		}
	}

	// Text usage tends to fall into a number of different roles, with a power law
        // frequency of occurrence in a text, so let's classify in order of likely usage
	// for small and many, we use a hashmap/btree

	switch spaces {
	case 0:
		return N1GRAM
	case 1:
		return N2GRAM
	case 2:
		return N3GRAM
	}

	// For longer strings, a linear search is probably fine here
        // (once it gets into a database, it's someone else's problem)

	if l < 128 {
		return LT128
	}

	if l < 1024 {
		return LT1024
	}

	return GT1024

}

//**************************************************************

func AppendTextToDirectory(txt Text) NodeTextPtr {

	var ptr NodeTextPtr = -1
	var ok bool = false

	switch txt.C {
	case N1GRAM:
		ptr,ok = TEXT_DIRECTORY.N1grams[txt.S]
	case N2GRAM:
		ptr,ok = TEXT_DIRECTORY.N2grams[txt.S]
	case N3GRAM:
		ptr,ok = TEXT_DIRECTORY.N3grams[txt.S]
	case LT128:
		ptr,ok = LinearFindText(TEXT_DIRECTORY.LT128,txt)
	case LT1024:
		ptr,ok = LinearFindText(TEXT_DIRECTORY.LT1024,txt)
	case GT1024:
		ptr,ok = LinearFindText(TEXT_DIRECTORY.GT1024,txt)

	}

	if ok {
		return ptr
	}

	switch txt.C {
	case N1GRAM:
		TEXT_DIRECTORY.N1directory = append(TEXT_DIRECTORY.N1directory,txt)
		ptr = TEXT_DIRECTORY.N1_top
		TEXT_DIRECTORY.N1grams[txt.S] = ptr
		TEXT_DIRECTORY.N1_top++
		return ptr
	case N2GRAM:
		TEXT_DIRECTORY.N2directory = append(TEXT_DIRECTORY.N2directory,txt)
		ptr = TEXT_DIRECTORY.N2_top
		TEXT_DIRECTORY.N2grams[txt.S] = ptr
		TEXT_DIRECTORY.N2_top++
		return ptr
	case N3GRAM:
		TEXT_DIRECTORY.N3directory = append(TEXT_DIRECTORY.N3directory,txt)
		ptr = TEXT_DIRECTORY.N3_top
		TEXT_DIRECTORY.N3grams[txt.S] = ptr
		TEXT_DIRECTORY.N3_top++
		return ptr
	case LT128:
		TEXT_DIRECTORY.LT128 = append(TEXT_DIRECTORY.LT128,txt)
		TEXT_DIRECTORY.LT128_top++
		return TEXT_DIRECTORY.LT128_top-1
	case LT1024:
		TEXT_DIRECTORY.LT1024 = append(TEXT_DIRECTORY.LT1024,txt)
		TEXT_DIRECTORY.LT1024_top++
		return TEXT_DIRECTORY.LT1024_top-1
	case GT1024:
		TEXT_DIRECTORY.GT1024 = append(TEXT_DIRECTORY.GT1024,txt)
		TEXT_DIRECTORY.GT1024_top++
		return TEXT_DIRECTORY.GT1024_top-1
	}

	return -1
}

//**************************************************************

func LinearFindText(in []Text,txt Text) (NodeTextPtr,bool) {

	for i := 0; i < len(in); i++ {

		if txt.L != in[i].L {
			continue
		}

		if in[i].S == txt.S {
			return NodeTextPtr(i),true
		}
	}

	return -1,false
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
		if stop == ALPHATEXT {
			//src[pos] = unicode.ToLower(src[pos])
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

	// Quoted strings are tricky

	if stop == '"' || stop == '\'' {
		var is_end bool

		if pos+1 >= len(src) {
			is_end= true
		} else {
			is_end = IsWhiteSpace(src[pos],src[pos+1])
		}

		if src[pos-1] == stop && is_end {
			return false
		} else {
			return true
		}

	}

	if pos >= len(src) || src[pos] == '\n' {
		return false
	}

	if stop == ALPHATEXT {
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
		if stop != '"' {
			return true
		}
	}

	if pos > 0 && src[pos-1] == stop && src[pos] != stop {
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

func IsWhiteSpace(r,rn rune) bool {

	return (unicode.IsSpace(r) || r == '#' || r == '/' && rn == '/')
}

//**************************************************************

func IsBackReference(src []rune,pos int) bool {

	// Any non-whitespace before \n or ( means it's not a back reference

	for pos++; pos < len(src); pos++ {

		if src[pos] == '(' || src[pos] == '\n' || src[pos] == '#' {
			return true
		} else {
			if !unicode.IsSpace(src[pos]) {
				return false
			}
		}
	}

	return false
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
	case HAVE_MINUS:
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

func CheckSequenceMode(context string, mode rune) {

	if (strings.Contains(context,"_sequence_")) {

		switch mode {
		case '+':
			Verbose("\nStart sequence mode for items")
			SEQUENCE_MODE = true
		case '-':
			Verbose("End sequence mode for items\n")
			SEQUENCE_MODE = false
			LAST_IN_SEQUENCE = ""
		}
	}

}

//**************************************************************

func LinkSequence(this string) {

	if SEQUENCE_MODE {
		if LINE_ITEM_COUNTER == 1 && LAST_IN_SEQUENCE != "" {
			Verbose("... Append sequence:",SEQUENCE_RELN,"->",this)
		}
		LAST_IN_SEQUENCE = this
	}
}

//**************************************************************

func FindAssociation(token string) (string,[]string) {

	var fullname string
	var reln []string

	name := token[1:len(token)-1]
	name = strings.TrimSpace(name)

	if strings.Contains(name,",") {
		reln = strings.Split(name,",")
		name = reln[0]
	}

	ptr, ok := ARROW_SHORT_DIR[name]

	if !ok {
		ptr, ok = ARROW_LONG_DIR[name]

		if !ok {
			ParseError(ERR_NO_SUCH_ARROW+name)
			os.Exit(-1)
		}
	}

	fullname = ARROW_DIRECTORY[ptr].Long

	return fullname, reln
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

	if number < 1 {
		ParseError(ERR_BAD_ALIAS_REFERENCE)
		os.Exit(-1)
	}

	return LookupAlias(name,number)
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
		if !unicode.IsUpper(r) && unicode.IsLetter(r) || unicode.IsNumber(r) {
			return false
		}
	}

	return true
}

//**************************************************************

func StripParen(token string) string {

	token =	strings.TrimSpace(token[1:])

	if token[0] == '(' {
		token =	strings.TrimSpace(token[1:])
	}

	if token[len(token)-1] == ')' {
		token =	token[:len(token)-1]
	}

	return token
}

//**************************************************************
// Tools
//**************************************************************

func ParseError(message string) {

	const red = "\033[31;1;1m"
	const endred = "\033[0m"

	fmt.Print("\n",LINE_NUM,":",red)
	fmt.Println("N4L",CURRENT_FILE,message,"at line", LINE_NUM,endred)
	Diag("N4L",CURRENT_FILE,message,"at line", LINE_NUM)
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

func Verbose(a ...interface{}) {

	if VERBOSE {
		fmt.Print(LINE_NUM,":\t")
		fmt.Println(a...)
	}

	Diag(a...)
}

//**************************************************************

func PVerbose(a ...interface{}) {

	const green = "\x1b[36m"
	const endgreen = "\x1b[0m"

	if VERBOSE {
		fmt.Print(LINE_NUM,":\t",green)
		fmt.Println(a...)
		fmt.Print(endgreen)
	}

	Diag(a...)
}

//**************************************************************

func Box(a ...interface{}) {

	if VERBOSE {

		fmt.Println("\n------------------------------------")
		//fmt.Print(LINE_NUM,":")
		fmt.Println(a...)
		fmt.Println("------------------------------------\n")
	}

	Diag(a...)
}

//**************************************************************

func DiagnosticName(filename string) string {

	return filename+"_test_log"

}

//**************************************************************

func Diag(a ...interface{}) {

	// Log diagnostic output for self-diagnostic tests

	prefix := fmt.Sprint(LINE_NUM,":")
	s := fmt.Sprintln(a...)
	AppendStringToFile(TEST_DIAG_FILE,prefix+s)
}

//**************************************************************

func AppendStringToFile(name string, s string) {

	f, err := os.OpenFile(name,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Couldn't open for write/append to",name,err)
		return
	}

	_, err = f.WriteString(s)

	if err != nil {
		fmt.Println("Couldn't write/append to",name,err)
	}

	f.Close()
}

