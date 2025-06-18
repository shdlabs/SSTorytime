
//
// transform random ngrams to graph
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
	//"sort"
	"strconv"
	"math"

        SST "SSTorytime"
)

//**************************************************************
// Parsing state variables
//**************************************************************

const (
	ALPHATEXT = 'x'
	NON_ASCII_LQUOTE = '“'
	NON_ASCII_RQUOTE = '”'

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

	N_GRAM_MAX = 5

)

//**************************************************************

var ( 
	LINE_NUM int = 1
	WORD_MISTAKE_LEN = 3 // a string shorter than this is probably a mistake

	// Flags

	VERBOSE bool = false
	CURRENT_FILE string
	LINE_ITEM_CACHE = make(map[string][]string)  // contains current and labelled line elements
	LINE_ITEM_REFS []SST.NodePtr                     // contains current line integer references
	LINE_RELN_CACHE = make(map[string][]SST.Link)
	LINE_ITEM_STATE int = ROLE_BLANK_LINE
	LINE_ALIAS string = ""
	LINE_ITEM_COUNTER int = 1
	LINE_RELN_COUNTER int = 0

	FWD_ARROW string
	BWD_ARROW string
	FWD_INDEX SST.ArrowPtr
	BWD_INDEX SST.ArrowPtr
	ANNOTATION = make(map[string]string)
	SECTION_STATE string
	CONTEXT_STATE = make(map[string]bool)

	ERR_MISSING_EVENT = "Missing item? Dangling section, relation, or context"
	ERR_ILLEGAL_CONFIGURATION = "Error in configuration, no such section"
	ERR_ARR_REDEFINITION="Redefinition of arrow "
	ERR_NEGATIVE_WEIGHT = "Arrow relation has a negative weight, which is disallowed. Use a NOT relation if you want to signify inhibition: "
	ERR_TOO_MANY_WEIGHTS = "More than one weight value in the arrow relation "

	ERR_BAD_ABBRV = "abbreviation out of place"
	ERR_SIMILAR_NO_SIGN = "Arrows for similarity do not have signs, they are directionless"
	ERR_ANNOTATION_MISSING = "Missing non-alphnumeric annotation marker or stray relation"
	ERR_ANNOTATION_REDEFINE = "Redefinition of annotation character"
	ERR_ANNOTATION_BAD = "Annotation marker should be short mark of non-space, non-alphanumeric character "
	ERR_ILLEGAL_ANNOT_CHAR="Cannot use +/- reserved tokens for annotation"
	ERR_MISMATCH_QUOTE = "Apparent missing or mismatch in ', \" or ( )"
        ERR_STRAY_PAREN="Stray ) in an event/item - illegal character"
	ERR_SHORT_WORD="Short word, probably a mistake: "
	ERR_NO_SUCH_FILE_FOUND = "No file found in the name "

	LAST_IN_SEQUENCE string = ""

)

//**************************************************************

// Promise bindings in English. This domain knowledge saves us a lot of training analysis

var FORBIDDEN_ENDING = []string{"but", "and", "the", "or", "a", "an", "its", "it's", "their", "your", "my", "of", "as", "are", "is", "be", "with", "using", "that", "who", "to" ,"no", "because","at","but","yes","no","yeah","yay", "in", "which", "what","as","he","she","they"}

var FORBIDDEN_STARTER = []string{"and","or","of","the","it","because","in","that","these","those","is","are","was","were","but","yes","no","yeah","yay","also"}

// ****************************************************************************

var EXCLUSIONS []string
var LEG_WINDOW int = 100  // sentences per leg
var STM_NGRAM_RANK [N_GRAM_MAX]map[string]float64

// ***************************************************************************

type Match struct {
	Arrow   string
	Before  string
	After   string
}

//**************************************************************
// BEGIN
//**************************************************************

func main() {

	args := Init()

	NewFile("N4Lconfig.in")
	config := ReadFile(CURRENT_FILE)

	ParseConfig(config)

	for i := 1; i < N_GRAM_MAX; i++ {

		STM_NGRAM_RANK[i] = make(map[string]float64)
	}

	for input := 0; input < len(args); input++ {

		NewFile(args[input])

		//text := string(ReadFile(CURRENT_FILE)) - why does this not terminate?

		rawtext,err := ioutil.ReadFile(CURRENT_FILE)

		if err != nil {
			fmt.Println("Can't read",CURRENT_FILE)
			os.Exit(-1)
		}

		text := CleanText(string(rawtext))
		paras := strings.Split(text,">>")

		pbs,count := ParaBySentence(paras)

		scale := float64(count) / float64(LEG_WINDOW)

		var ltm_every_ngram_occurrence [N_GRAM_MAX]map[string][]int
		for i := 1; i < N_GRAM_MAX; i++ {
			ltm_every_ngram_occurrence[i] = make(map[string][]int)
		} 

		for p := range pbs {
			for s := range pbs[p] {
				
				score := FractionateThenRankSentence(s,pbs[p][s],scale,ltm_every_ngram_occurrence)
				fmt.Println(" - ",score,pbs[p][s])

			}
		}
	}


	for n := range STM_NGRAM_RANK {
		fmt.Println("n --------------- ",n)
		for w := range STM_NGRAM_RANK[n] {
			fmt.Println(n,w)
		}
	}
}

//**************************************************************

func Analyze(s string) []Match {

	var matches []Match

	for arr := range SST.ARROW_DIRECTORY {

		arrow := SST.ARROW_DIRECTORY[arr].Long

		if arrow == "is not" {
			continue
		}

		arrow = strings.Replace(arrow,"is","",-1)
		arrow = strings.Replace(arrow,"has","",-1)

		var match Match
		pos := strings.Index(s,arrow)

		if pos >= 0 {
			match.Arrow = SST.ARROW_DIRECTORY[arr].Long
			match.After = s[pos:]
			match.Before = s[:pos]
			matches = append(matches,match)
			continue
		}

		if len(SST.ARROW_DIRECTORY[arr].Short) > 3 {
			pos = strings.Index(s,SST.ARROW_DIRECTORY[arr].Short)

			if pos >= 0 {
				match.Arrow = SST.ARROW_DIRECTORY[arr].Long
				match.After = s[pos:]
				match.Before = s[:pos]
				matches = append(matches,match)
			}
		}
	}

	return matches
}

//**************************************************************

func CleanText(s string) string {

	// Start by stripping HTML / XML tags before para-split
	// if they haven't been removed already

	m := regexp.MustCompile("<[^>]*>") 
	s = m.ReplaceAllString(s,":\n") 

	// Weird English abbrev
	s = strings.Replace(s,"Mr.","Mr",-1) 
	s = strings.Replace(s,"Ms.","Ms",-1) 
	s = strings.Replace(s,"Mrs.","Mrs",-1) 
	s = strings.Replace(s,"Dr.","Dr",-1)
	s = strings.Replace(s,"St.","St",-1) 
	s = strings.Replace(s,"[","",-1) 
	s = strings.Replace(s,"]","",-1) 

	// Encode end of sentence markers with a # for later splitting

	m = regexp.MustCompile("[\n][\n]")
	s = m.ReplaceAllString(s,">>\n")

	m = regexp.MustCompile("[?!.]+[ \n]")
	s = m.ReplaceAllString(s,"$0#")

	m = regexp.MustCompile("[\n]+")
	s = m.ReplaceAllString(s," ")

	return s
}

//**************************************************************

func ParaBySentence(paras []string) ([][]string,int) {
	
	var pbs [][]string
	var count = 0
	
	for s := range paras {
		
		var para []string
		
		sentences := strings.Split(paras[s],"#")
		
		for s := range sentences {
			sentence := strings.TrimSpace(sentences[s])
			if sentence != ":" && sentence != "" {
				para = append(para,sentence)
				count++
			}
		}
		pbs = append(pbs,para)
	}

	return pbs,count
}

//**************************************************************

func FractionateThenRankSentence(s_idx int, sentence string, scale float64,ltm_every_ngram_occurrence [N_GRAM_MAX]map[string][]int) float64 {

	// A round robin cyclic buffer for taking fragments and extracting
	// n-ngrams of 1,2,3,4,5,6 words separateed by whitespace, passing

	var rrbuffer [N_GRAM_MAX][]string
	var sentence_meaning_rank float64 = 0
	var rank float64
	
	// split sentence on any residual punctuation here, because punctuation cannot be in the middle
	// of an n-gram by definition of punctuation's promises, and we are not interested in word groups
	// that unintentionally straddle punctuation markers, since they are false signals
	
	re := regexp.MustCompile("[,.;:!?]")
	sentence_frags := re.Split(sentence, -1)
	
	for f := range sentence_frags {
		
		// For one sentence, break it up into codons and sum their importances
		
		clean_sentence := strings.Split(string(sentence_frags[f])," ")
		
		for word := range clean_sentence {
			
			// This will be too strong in general - ligatures and foreign languages etc
			
			m := regexp.MustCompile("[/()?!]*") 
			cleanjunk := m.ReplaceAllString(clean_sentence[word],"") 
			cleanword := strings.Trim(strings.ToLower(cleanjunk)," ")
			
			if len(cleanword) == 0 {
				continue
			}
			
			// Shift all the rolling longitudinal Ngram rr-buffers by one word
			
			rank, rrbuffer = NextWordAndUpdateLTMNgrams(s_idx,cleanword, rrbuffer,scale,ltm_every_ngram_occurrence)
			sentence_meaning_rank += rank
		}
	}
	
	return sentence_meaning_rank
}

//**************************************************************

func NextWordAndUpdateLTMNgrams(s_idx int, word string, rrbuffer [N_GRAM_MAX][]string,scale float64,ltm_every_ngram_occurrence [N_GRAM_MAX]map[string][]int) (float64,[N_GRAM_MAX][]string) {

	// Word by word, we form a superposition of scores from n-grams of different lengths
	// as a simple sum. This means lower lengths will dominate as there are more of them
	// so we define intentionality proportional to the length also as compensation

	var rank float64 = 0

	for n := 2; n < N_GRAM_MAX; n++ {
		
		// Pop from round-robin

		if (len(rrbuffer[n]) > n-1) {
			rrbuffer[n] = rrbuffer[n][1:n]
		}
		
		// Push new to maintain length

		rrbuffer[n] = append(rrbuffer[n],word)

		// Assemble the key, only if complete cluster
		
		if (len(rrbuffer[n]) > n-1) {
			
			var key string
			
			for j := 0; j < n; j++ {
				key = key + rrbuffer[n][j]
				if j < n-1 {
					key = key + " "
				}
			}

			if ExcludedByBindings(rrbuffer[n][0],rrbuffer[n][n-1]) {

				continue
			}

			STM_NGRAM_RANK[n][key]++
			rank += Intentionality(n,key,scale)

			ltm_every_ngram_occurrence[n][key] = append(ltm_every_ngram_occurrence[n][key],s_idx)

		}
	}

	STM_NGRAM_RANK[1][word]++
	rank += Intentionality(1,word,scale)

	ltm_every_ngram_occurrence[1][word] = append(ltm_every_ngram_occurrence[1][word],s_idx)

	return rank, rrbuffer
}

//**************************************************************
// Heuristics
//**************************************************************

func ExcludedByBindings(firstword,lastword string) bool {

	// A standalone fragment can't start/end with these words, because they
	// Promise to bind to something else...
	// Rather than looking for semantics, look at spacetime promises only - words that bind strongly
	// to a prior or posterior word.

	if (len(firstword) == 1) || len(lastword) == 1 {
		return true
	}

	for s := range FORBIDDEN_ENDING {
		if lastword == FORBIDDEN_ENDING[s] {
			return true
		}
	}
	
	for s := range FORBIDDEN_STARTER {
		if firstword == FORBIDDEN_STARTER[s] {
			return true
		}
	}

	return false 
}

//**************************************************************

func Intentionality(n int, s string, scale float64) float64 {

	// Compute the effective intent of a string s at a position count
	// within a document of many sentences. The weighting due to
	// inband learning uses an exponential deprecation based on
	// SST scales (see "leg" meaning).

	occurrences := STM_NGRAM_RANK[n][s]
	work := float64(len(s))

	if occurrences < 3 {
		return 0
	}

	if work < 5 {
		return 0
	}

	// lambda should have a cutoff for insignificant words, like "a" , "of", etc that occur most often

	lambda := occurrences / float64(LEG_WINDOW)

	// This constant is tuned to give words a growing importance up to a limit
	// or peak occurrences, then downgrade

	// Things that are repeated too often are not important
	// but length indicates purposeful intent

	meaning := lambda * work / (1.0 + math.Exp(lambda-scale))

return meaning
}

//**************************************************************
// IMPORTED: N4L configuration
//**************************************************************

func Init() []string {

	flag.Usage = Usage
	verbosePtr := flag.Bool("v", false,"verbose")

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		Usage()
		os.Exit(1);
	}

	if *verbosePtr {
		VERBOSE = true
	}

	SST.MemoryInit()

	return args
}

//**************************************************************

func NewFile(filename string) {

	CURRENT_FILE = filename

	Box("Parsing new file",filename)

	LINE_NUM = 1
}

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
			
		case '-':
			BWD_ARROW = strings.TrimSpace(token[1:])
			LINE_ITEM_STATE = HAVE_MINUS

		case '(':
			reln := token[1:len(token)-1]
			reln = strings.TrimSpace(reln)

			if LINE_ITEM_STATE == HAVE_MINUS {
				CheckArrow(reln,BWD_ARROW)
				BWD_INDEX = SST.InsertArrowDirectory(SECTION_STATE,reln,BWD_ARROW,"-")
				SST.InsertInverseArrowDirectory(FWD_INDEX,BWD_INDEX)
				PVerbose("In",SECTION_STATE,"short name",reln,"for",BWD_ARROW,", direction","-")
			} else if LINE_ITEM_STATE == HAVE_PLUS {
				CheckArrow(reln,FWD_ARROW)
				FWD_INDEX = SST.InsertArrowDirectory(SECTION_STATE,reln,FWD_ARROW,"+")
				PVerbose("In",SECTION_STATE,"short name",reln,"for",FWD_ARROW,", direction","+")
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
				index := SST.InsertArrowDirectory(SECTION_STATE,reln,BWD_ARROW,"both")
				SST.InsertInverseArrowDirectory(index,index)
				PVerbose("In",SECTION_STATE,reln,"for",BWD_ARROW,", direction","both")
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

			for r := range token {
				if unicode.IsLetter(rune(token[r])) {
					ParseError(ERR_ANNOTATION_BAD)
				}
			}

			if token[0] == '+' || token[0] == '-' {
				ParseError(ERR_ILLEGAL_ANNOT_CHAR)
				os.Exit(-1)
			}

			LINE_ITEM_STATE = HAVE_PLUS
			LAST_IN_SEQUENCE = token

		}

	default:
		ParseError(ERR_ILLEGAL_CONFIGURATION+" "+SECTION_STATE)
		os.Exit(-1)
	}
}

//**************************************************************

func CheckArrow(alias,name string) {

	prev,ok := SST.ARROW_SHORT_DIR[alias]
	if ok {
		ParseError(ERR_ARR_REDEFINITION+"\""+alias+"\" previous short name: "+SST.ARROW_DIRECTORY[prev].Short)
		os.Exit(-1)
	}
	
	prev,ok = SST.ARROW_LONG_DIR[name]
	if ok {
		ParseError(ERR_ARR_REDEFINITION+"\""+name+"\" previous long name: "+SST.ARROW_DIRECTORY[prev].Long)
		os.Exit(-1)
	}
}

//**************************************************************

func GetLinkArrowByName(token string) SST.Link {

	// Return a preregistered link/arrow ptr bythe name of a link

	var reln []string
	var weight float32 = 1
	var weightcount int
	var ctx []string
	var name string

	if token[0] == '(' {
		name = token[1:len(token)-1]
	} else {
		name = token
	}

	name = strings.TrimSpace(name)

	if strings.Contains(name,",") {
		reln = strings.Split(name,",")
		name = reln[0]

		// look at any comma separated notes after the arrow name
		for i := 1; i < len(reln); i++ {

			v, err := strconv.ParseFloat(reln[i], 64)

			if err == nil {
				if weight < 0 {
					ParseError(ERR_NEGATIVE_WEIGHT+token)
					os.Exit(-1)
				}
				if weightcount > 1 {
					ParseError(ERR_TOO_MANY_WEIGHTS+token)
					os.Exit(-1)
				}
				weight = float32(v)
				weightcount++
			} else {
				ctx = append(ctx,reln[i])
			}
		}
	}

	// First check if this is an alias/short name

	ptr, ok := SST.ARROW_SHORT_DIR[name]
	
	// If not, then check longname
	
	if !ok {
		ptr, ok = SST.ARROW_LONG_DIR[name]
		
		if !ok {
			ParseError(SST.ERR_NO_SUCH_ARROW+"("+name+")")
			os.Exit(-1)
		}
	}

	var link SST.Link

	link.Arr = ptr
	link.Wgt = weight
	link.Ctx = GetContext(ctx)
	return link
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
	LINE_ITEM_REFS = nil
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
	case HAVE_MINUS:
		return false
	}

	return true
}

//**************************************************************

func GetContext(ctx []string) []string {

	var merge = make(map[string]bool)
	var clist []string

	for c := range CONTEXT_STATE {
		merge[c] = true
	}


	for c := range ctx {
		merge[ctx[c]] = true
	}

	for c := range merge {
		clist = append(clist,c)
	}

	return clist
}

//**************************************************************
// Scan text input
//**************************************************************

func ReadFile(filename string) []rune {

	text := ReadTUF8File(filename)
	return text
}


//**************************************************************

func ReadToLast(src []rune,pos int, stop rune) (string,int) {

	var cpy []rune

	var starting_at = LINE_NUM

	for ; Collect(src,pos,stop,cpy) && pos < len(src); pos++ {
		cpy = append(cpy,src[pos])
	}

	if IsQuote(stop) && src[pos-1] != stop {
		e := fmt.Sprintf("%s starting at line %d (found token %s)",ERR_MISMATCH_QUOTE,starting_at,string(cpy))
		ParseError(e)
		os.Exit(-1)
	}

	token := string(cpy)

	token = strings.TrimSpace(token)

	return token,pos
}

//**************************************************************

func Collect(src []rune,pos int, stop rune,cpy []rune) bool {

	var collect bool = true

	// Quoted strings are tricky

	if IsQuote(stop) {
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

		if stop != ':' && !IsQuote(stop) { 
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

        case ')':
	        ParseError(ERR_STRAY_PAREN)
		os.Exit(-1)
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

func IsQuote(r rune) bool {

	switch r {
	case '"','\'',NON_ASCII_LQUOTE,NON_ASCII_RQUOTE:
		return true
	}

	return false
}

//**************************************************************

func LastSpecialChar(src []rune,pos int, stop rune) bool {

	if src[pos] == '\n' {
		if stop != '"' {
			return true
		}
	}

	// Special case, but still don't understand why?!

	if src[pos] == '@' {
		return false
	}

	if pos > 0 && src[pos-1] == stop && src[pos] != stop {
		return true
	}

	return false
}

//**************************************************************

func IsWhiteSpace(r,rn rune) bool {

	return (unicode.IsSpace(r) || r == '#' || r == '/' && rn == '/')
}

//**************************************************************

func ExtractWord(fulltext string,offset int) string {

	var protected bool = false
	var word string

	for r := offset+1; r < len(fulltext); r++ {

		if fulltext[r] == '"' {
			protected = !protected
		}

		if !protected && !unicode.IsLetter(rune(fulltext[r])) {
			word = strings.Trim(strings.TrimSpace(word),"\" ")
			return word
		}

		word += string(fulltext[r])
	}
	
	word = strings.Trim(strings.TrimSpace(word),"\" ")
	
	if len(word) <= WORD_MISTAKE_LEN {
		ParseError(ERR_SHORT_WORD+word)
	}
	
	return word
}

// **************************************************************************

func ParseError(message string) {

	const red = "\033[31;1;1m"
	const endred = "\033[0m"

	fmt.Print("\n",LINE_NUM,":",red)
	fmt.Println("ScanText",CURRENT_FILE,message,"at line", LINE_NUM,endred)

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

func Usage() {
	
	fmt.Printf("usage: ScanText [-v] [file].dat\n")
	flag.PrintDefaults()
	os.Exit(2)
}

//**************************************************************

func Verbose(a ...interface{}) {

	line := fmt.Sprintln(a...)
	
	if VERBOSE {
		fmt.Print(line)
	}
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
}

//**************************************************************

func Box(a ...interface{}) {

	if VERBOSE {

		fmt.Println("\n------------------------------------")
		fmt.Println(a...)
		fmt.Println("------------------------------------\n")
	}
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



