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
	"strconv"
	"sort"
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

	WORD_MISTAKE_LEN = 3 // a string shorter than this is probably a mistake

	ERR_NO_SUCH_FILE_FOUND = "No file found in the name "
	ERR_MISSING_EVENT = "Missing item? Dangling section, relation, or context"
	ERR_MISSING_SECTION = "Declarations outside a section or chapter"
	ERR_NO_SUCH_ALIAS = "No such alias or \" reference exists to fill in - aborting"
	ERR_NO_SUCH_ARROW = "No such arrow has been declared in the configuration: "
	ERR_MISSING_ITEM_SOMEWHERE = "Missing item somewhere"
	ERR_MISSING_ITEM_RELN = "Missing item or double relation"
	ERR_MISMATCH_QUOTE = "Apparent missing or mismatch in ', \" or ( )"
	ERR_ILLEGAL_CONFIGURATION = "Error in configuration, no such section"
	ERR_BAD_LABEL_OR_REF = "Badly formed label or reference (@label becomes $label.n) in "
	WARN_NOTE_TO_SELF = "WARNING: Found a note to self in the text"
	WARN_INADVISABLE_CONTEXT_EXPRESSION = "WARNING: Inadvisably complex/parenthetic context expression - simplify?"
	ERR_ILLEGAL_QUOTED_STRING_OR_REF = "WARNING: Something wrong, bad quoted string or mistaken back reference. Close any space after a quote..."
	ERR_ANNOTATION_BAD = "Annotation marker should be short mark of non-space, non-alphanumeric character "
	ERR_BAD_ABBRV = "abbreviation out of place"
	ERR_BAD_ALIAS_REFERENCE = "Alias references start from $name.1"
	ERR_ANNOTATION_MISSING = "Missing non-alphnumeric annotation marker or stray relation"
	ERR_ANNOTATION_REDEFINE = "Redefinition of annotation character"
	ERR_SIMILAR_NO_SIGN = "Arrows for similarity do not have signs, they are directionless"
	ERR_ARROW_SELFLOOP = "Arrow's origin points to itself"
	ERR_NEGATIVE_WEIGHT = "Arrow relation has a negative weight, which is disallowed. Use a NOT relation if you want to signify inhibition: "
	ERR_TOO_MANY_WEIGHTS = "More than one weight value in the arrow relation "
        ERR_STRAY_PAREN="Stray ) in an event/item - illegal character"
	ERR_MISSING_LINE_LABEL_IN_REFERENCE="Missing a line label in reference, should be int he form $label.n"
	ERR_NON_WORD_WHITE="Non word (whitespace) character after an annotation: "
	ERR_SHORT_WORD="Short word, probably a mistake: "
	ERR_ARR_REDEFINITION="Redefinition of arrow "
)

//**************************************************************

var ( 
	LINE_NUM int = 1
	LINE_ITEM_CACHE = make(map[string][]string)  // contains current and labelled line elements
	LINE_ITEM_REFS []NodePtr                     // contains current line integer references
	LINE_RELN_CACHE = make(map[string][]Link)
	LINE_ITEM_STATE int = ROLE_BLANK_LINE
	LINE_ALIAS string = ""
	LINE_ITEM_COUNTER int = 1
	LINE_RELN_COUNTER int = 0

	FWD_ARROW string
	BWD_ARROW string
	FWD_INDEX ArrowPtr
	BWD_INDEX ArrowPtr
	ANNOTATION = make(map[string]string)
	INVERSE_ARROWS = make(map[ArrowPtr]ArrowPtr)

	CONTEXT_STATE = make(map[string]bool)
	SECTION_STATE string

	SEQUENCE_MODE bool = false
	SEQUENCE_RELN string = "then" 
	LAST_IN_SEQUENCE string = ""

	// Flags

	VERBOSE bool = false
	DIAGNOSTIC bool = false
	UPLOAD bool = false
	SUMMARIZE bool = false
	CREATE_ADJACENCY bool = false
	ADJ_LIST string

	CURRENT_FILE string
	TEST_DIAG_FILE string

	RELN_BY_SST [4][]ArrowPtr // From an EventItemNode
	SST_NAMES[4] string
)

//**************************************************************
// DATA structures for input
//**************************************************************

// See the notes in README

//**************************************************************

const (
	NEAR = 0
	LEADSTO = 1   // +/-
	CONTAINS = 2  // +/-
	EXPRESS = 3   // +/-

	ST_ZERO = EXPRESS // so that ST_ZERO - EXPRESS == 0
	ST_TOP = ST_ZERO + EXPRESS + 1

	N1GRAM = 1
	N2GRAM = 2
	N3GRAM = 3
	LT128 = 4
	LT1024 = 5
	GT1024 = 6
)

//**************************************************************

type Node struct { // essentially the incidence matrix

	L int                 // length of name string
	S string              // name string itself

	Chap string           // section/chapter in which this was added
	SizeClass int         // the string class: N1-N3, LT128, etc
	NPtr NodePtr          // Pointer to self

	I [ST_TOP][]Link   // link incidence list, by arrow type
  	                   // NOTE: carefully how offsets represent negative SSTtypes
}

//**************************************************************

type NodePtr struct {

	CPtr  ClassedNodePtr // index of within name class lane
	Class int            // Text size-class
}

type ClassedNodePtr int  // Internal pointer type of size-classified text

//**************************************************************

type RCtype struct {
	Row NodePtr
	Col NodePtr
}

//**************************************************************

type Link struct {  // A link is a type of arrow, with context
                    // and maybe with a weightfor package math
	Arr ArrowPtr         // type of arrow, presorted
	Wgt float64          // numerical weight of this link
	Ctx []string         // context for this pathway
	Dst NodePtr // adjacent event/item/node
}

//**************************************************************

type ArrowDirectory struct {

	STAindex int
	Long    string
	Short   string
	Ptr     ArrowPtr
}

type ArrowPtr int // ArrowDirectory index

 // all fwd arrow types have a simple int representation > 0
 // all bwd/inverse arrow readings have the negative int for fwd
 // Hashed by long and short names

//**************************************************************

type NodeDirectory struct {

	// Power law n-gram frequencies

	N1grams map[string]ClassedNodePtr
	N1directory []Node
	N1_top ClassedNodePtr

	N2grams map[string]ClassedNodePtr
	N2directory []Node
	N2_top ClassedNodePtr

	N3grams map[string]ClassedNodePtr
	N3directory []Node
	N3_top ClassedNodePtr

	// Use linear search on these exp fewer long strings

	LT128 []Node
	LT128_top ClassedNodePtr
	LT1024 []Node
	LT1024_top ClassedNodePtr
	GT1024 []Node
	GT1024_top ClassedNodePtr
}

//**************************************************************
// Lookup tables
//**************************************************************

var ( 
	ARROW_DIRECTORY []ArrowDirectory
	ARROW_SHORT_DIR = make(map[string]ArrowPtr) // Look up short name int referene
	ARROW_LONG_DIR = make(map[string]ArrowPtr)  // Look up long name int referene
	ARROW_DIRECTORY_TOP ArrowPtr = 0

	NODE_DIRECTORY NodeDirectory  // Internal histo-representations
	NO_NODE_PTR NodePtr           // see Init()
)

//**************************************************************

func main() {

	args := Init()

	NewFile("N4Lconfig.in")
	config := ReadFile(CURRENT_FILE)

	AddMandatory()
	ParseConfig(config)

	//SummarizeAndTestConfig()

	for input := 0; input < len(args); input++ {
		NewFile(args[input])
		input := ReadFile(CURRENT_FILE)
		ParseN4L(input)
	}

	if SUMMARIZE {
		SummarizeGraph()
	}

	if CREATE_ADJACENCY {
		dim, key, d_adj, u_adj := CreateAdjacencyMatrix(ADJ_LIST)
		PrintMatrix("directed adjacency sub-matrix",dim,key,d_adj)
		PrintMatrix("undirected adjacency sub-matrix",dim,key,u_adj)
		evc := ComputeEVC(dim,u_adj)
		PrintNZVector("Eigenvector centrality (EVC) score for symmetrized graph",dim,key,evc)
	}

	if UPLOAD {
	}
}

//**************************************************************

func Init() []string {

	flag.Usage = Usage
	verbosePtr := flag.Bool("v", false,"verbose")
	diagPtr := flag.Bool("d", false,"diagnostic mode")
	uploadPtr := flag.Bool("u", false,"upload")
	incidencePtr := flag.Bool("s", false,"summary (node,links...)")
	adjacencyPtr := flag.String("adj", "none", "a quoted, comma-separated list of short link names")

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		Usage()
		os.Exit(1);
	}

	if *verbosePtr {
		VERBOSE = true
	}

	if *diagPtr {
		VERBOSE = true
		DIAGNOSTIC = true
	}

	if *uploadPtr {
		UPLOAD = true
	}
	if *incidencePtr {
		SUMMARIZE = true
	}

	if *adjacencyPtr != "none" {
		CREATE_ADJACENCY = true
		ADJ_LIST = *adjacencyPtr
	}

	NO_NODE_PTR.Class = 0
	NO_NODE_PTR.CPtr =  -1

	NODE_DIRECTORY.N1grams = make(map[string]ClassedNodePtr)
	NODE_DIRECTORY.N2grams = make(map[string]ClassedNodePtr)
	NODE_DIRECTORY.N3grams = make(map[string]ClassedNodePtr)

	SST_NAMES[NEAR] = "Near"
	SST_NAMES[LEADSTO] = "LeadsTo"
	SST_NAMES[CONTAINS] = "Contains"
	SST_NAMES[EXPRESS] = "Express"

	return args
}

//**************************************************************

func NewFile(filename string) {

	CURRENT_FILE = filename
	TEST_DIAG_FILE = DiagnosticName(filename)

	Box("Parsing new file",filename)

	LINE_ITEM_STATE = ROLE_BLANK_LINE
	LINE_NUM = 1
	LINE_ITEM_CACHE = make(map[string][]string)
	LINE_RELN_CACHE = make(map[string][]Link)
	LINE_ITEM_REFS = nil
	LINE_ITEM_COUNTER = 1
	LINE_RELN_COUNTER = 0
	LINE_ALIAS = ""
	LAST_IN_SEQUENCE = ""
	FWD_ARROW = ""
	BWD_ARROW = ""
	SECTION_STATE = ""
	CONTEXT_STATE = make(map[string]bool)
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
				CheckArrow(reln,BWD_ARROW)
				BWD_INDEX = InsertArrowDirectory(SECTION_STATE,reln,BWD_ARROW,"-")
				InsertInverseArrowDirectory(FWD_INDEX,BWD_INDEX)
			} else if LINE_ITEM_STATE == HAVE_PLUS {
				CheckArrow(reln,FWD_ARROW)
				FWD_INDEX = InsertArrowDirectory(SECTION_STATE,reln,FWD_ARROW,"+")
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
				index := InsertArrowDirectory(SECTION_STATE,reln,BWD_ARROW,"both")
				InsertInverseArrowDirectory(index,index)

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

func InsertArrowDirectory(sec,alias,name,pm string) ArrowPtr {

	// Insert an arrow into the forward/backward indices

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
		newarrow.STAindex = ST_ZERO + LEADSTO * sign
	case "contains":
		newarrow.STAindex = ST_ZERO + CONTAINS * sign
	case "properties":
		newarrow.STAindex = ST_ZERO + EXPRESS * sign
	case "similarity":
		newarrow.STAindex = ST_ZERO + NEAR
	}

	newarrow.Long = name
	newarrow.Short = alias
	newarrow.Ptr = ARROW_DIRECTORY_TOP

	ARROW_DIRECTORY = append(ARROW_DIRECTORY,newarrow)
	ARROW_SHORT_DIR[alias] = ARROW_DIRECTORY_TOP
	ARROW_LONG_DIR[name] = ARROW_DIRECTORY_TOP
	ARROW_DIRECTORY_TOP++

	return ARROW_DIRECTORY_TOP-1
}

//**************************************************************

func InsertInverseArrowDirectory(fwd,bwd ArrowPtr) {

	// Lookup inverse by long name, only need this in search presentation

	INVERSE_ARROWS[fwd] = bwd
	INVERSE_ARROWS[bwd] = fwd
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
	fmt.Println("\nTEXT\n\n",NODE_DIRECTORY)
}

//**************************************************************

func SummarizeGraph() {

	Box("SUMMARIZE GRAPH.....\n")

	var count_nodes int = 0
	var count_links [4]int
	var total int

	for class := N1GRAM; class <= GT1024; class++ {
		switch class {
		case N1GRAM:
			for n := range NODE_DIRECTORY.N1directory {
				fmt.Println(n,"\t",NODE_DIRECTORY.N1directory[n].S)
				count_nodes++
				for sttype := range NODE_DIRECTORY.N1directory[n].I {
					for lnk := range NODE_DIRECTORY.N1directory[n].I[sttype] {
						count_links[FlatSTIndex(sttype)]++
						PrintLink(NODE_DIRECTORY.N1directory[n].I[sttype][lnk])
					}
				}
				fmt.Println()
			}
		case N2GRAM:
			for n := range NODE_DIRECTORY.N2directory {
				fmt.Println(n,"\t",NODE_DIRECTORY.N2directory[n].S)
				count_nodes++
				for sttype := range NODE_DIRECTORY.N2directory[n].I {
					for lnk := range NODE_DIRECTORY.N2directory[n].I[sttype] {
						count_links[FlatSTIndex(sttype)]++
						PrintLink(NODE_DIRECTORY.N2directory[n].I[sttype][lnk])
					}

				}
				fmt.Println()
			}
		case N3GRAM:
			for n := range NODE_DIRECTORY.N3directory {
				fmt.Println(n,"\t",NODE_DIRECTORY.N3directory[n].S)
				count_nodes++
				for sttype := range NODE_DIRECTORY.N3directory[n].I {
					for lnk := range NODE_DIRECTORY.N3directory[n].I[sttype] {
						count_links[FlatSTIndex(sttype)]++
						PrintLink(NODE_DIRECTORY.N3directory[n].I[sttype][lnk])
					}
				}
				fmt.Println()
			}
		case LT128:
			for n := range NODE_DIRECTORY.LT128 {
				fmt.Println(n,"\t",NODE_DIRECTORY.LT128[n].S)
				count_nodes++
				for sttype := range NODE_DIRECTORY.LT128[n].I {
					for lnk := range NODE_DIRECTORY.LT128[n].I[sttype] {
						count_links[FlatSTIndex(sttype)]++
						PrintLink(NODE_DIRECTORY.LT128[n].I[sttype][lnk])
					}
				}
				fmt.Println()
			}
		case LT1024:
			for n := range NODE_DIRECTORY.LT1024 {
				fmt.Println(n,"\t",NODE_DIRECTORY.LT1024[n].S)
				count_nodes++
				for sttype := range NODE_DIRECTORY.LT1024[n].I {
					for lnk := range NODE_DIRECTORY.LT1024[n].I[sttype] {
						count_links[FlatSTIndex(sttype)]++
						PrintLink(NODE_DIRECTORY.LT1024[n].I[sttype][lnk])
					}
				}
				fmt.Println()
			}

		case GT1024:
			for n := range NODE_DIRECTORY.GT1024 {
				fmt.Println(n,"\t",NODE_DIRECTORY.GT1024[n].S)
				count_nodes++
				for sttype := range NODE_DIRECTORY.GT1024[n].I {
					for lnk := range NODE_DIRECTORY.GT1024[n].I[sttype] {
						count_links[FlatSTIndex(sttype)]++
						PrintLink(NODE_DIRECTORY.GT1024[n].I[sttype][lnk])
					}
				}
				fmt.Println()
			}
		}
	}

	fmt.Println("-------------------------------------")
	fmt.Println("Incidence summary of raw declarations")
	fmt.Println("-------------------------------------")

	fmt.Println("Total nodes",count_nodes)

	for st := 0; st < 4; st++ {
		total += count_links[st]
		fmt.Println("Total directed links of type",SST_NAMES[st],count_links[st])
	}

	complete := count_nodes * (count_nodes-1)
	fmt.Println("Total links",total,"sparseness (fraction of completeness)",float64(total)/float64(complete))
}

//**************************************************************

func CreateAdjacencyMatrix(searchlist string) (int,[]NodePtr,[][]float64,[][]float64) {

	search_list := ValidateLinkArgs(searchlist)

	// the matrix is dim x dim

	filtered_node_list,path_weights := AssembleInvolvedNodes(search_list)

	dim := len(filtered_node_list)

	for f := 0; f < len(filtered_node_list); f++ {
		Verbose("    - row/col key [",f,"/",dim,"]",GetNodeTxtFromPtr(filtered_node_list[f]))
	}

	// Debugging mainly
	//for f := range path_weights {
	//	Verbose("    - path weight",path_weights[f],"from",GetNodeTxtFromPtr(f.Row),"to",GetNodeTxtFromPtr(f.Col))
	//}

	var subadj_matrix [][]float64 = make([][]float64,dim)
	var symadj_matrix [][]float64 = make([][]float64,dim)

	for row := 0; row < dim; row++ {
		subadj_matrix [row] = make([]float64,dim)
		symadj_matrix [row] = make([]float64,dim)
	}

	for row := 0; row < dim; row++ {
		for col := 0; col < dim; col++ {

			var rc, rcT RCtype
			rc.Row = filtered_node_list[row]
			rc.Col = filtered_node_list[col]

			rcT.Row = filtered_node_list[col]
			rcT.Col = filtered_node_list[row]

			subadj_matrix[row][col] = path_weights[rc]

			symadj_matrix[row][col] = path_weights[rc] + path_weights[rcT]
			symadj_matrix[col][row] = path_weights[rc] + path_weights[rcT]
		}
	}

	return dim, filtered_node_list, subadj_matrix, symadj_matrix
}

//**************************************************************

func PrintMatrix(name string, dim int, key []NodePtr, matrix [][]float64) {


	s := fmt.Sprintln("\n",name,"...\n")
	Verbose(s)

	for row := 0; row < dim; row++ {
		
		s = fmt.Sprintf("%20.15s ..\r\t\t\t(",GetNodeTxtFromPtr(key[row]))
		
		for col := 0; col < dim; col++ {
			
			const screenwidth = 12
			
			if col > screenwidth {
				s += fmt.Sprint("\t...")
				break
			} else {
				s += fmt.Sprintf("  %4.1f",matrix[row][col])
			}
			
		}
		s += fmt.Sprint(")")
		Verbose(s)
	}
}

//**************************************************************

func PrintNZVector(name string, dim int, key []NodePtr, vector[]float64) {

	s := fmt.Sprintln("\n",name,"...\n")
	Verbose(s)

	type KV struct {
		Key string
		Value float64
	}

	var vec []KV = make([]KV,dim)

	for row := 0; row < dim; row++ {
		vec[row].Key = GetNodeTxtFromPtr(key[row])
		vec[row].Value = vector[row]
	}

	sort.SliceStable(vec, func(i, j int) bool {
		return vec[i].Value > vec[j].Value
	})

	for row := 0; row < dim; row++ {
		if vec[row].Value > 0.1 {
			s = fmt.Sprintf("ordered by EVC:  (%4.1f)  ",vec[row].Value)
			s += fmt.Sprintf("%-80.79s",vec[row].Key)
			Verbose(s)
		}
	}
}

//**************************************************************

func ComputeEVC(dim int,adj [][]float64) []float64 {

	v := MakeInitVector(dim,1.0)
	vlast := v

	const several = 6

	for i := 0; i < several; i++ {

		v = MatrixOpVector(dim,adj,vlast)
		maxval := GetVecMax(v)
		v = NormalizeVec(v,maxval)

		if CompareVec(v,vlast) < 0.1 {
			break
		}
		vlast = v
	}
	return v
}

//**************************************************************

func MakeInitVector(dim int, init_value float64) []float64 {

	var v = make([]float64,dim)

	for r := 0; r < dim; r++ {
		v[r] = init_value
	}

	return v
}

//**************************************************************

func MatrixOpVector(dim int,m [][]float64, v []float64) []float64 {

	var vp = make([]float64,dim)

	for r := 0; r < dim; r++ {
		for c := 0; c < dim; c++ {
			vp[r] += m[r][c] * v[c]
		}
	}
	return vp
}

//**************************************************************

func GetVecMax(v []float64) float64 {

	var max float64 = -1

	for r := range v {
		if v[r] > max {
			max = v[r]
		}
	}

	return max
}

//**************************************************************

func NormalizeVec(v []float64, div float64) []float64 {

	for r := range v {
		v[r] = v[r] / div
	}

	return v
}

//**************************************************************

func CompareVec(v1,v2 []float64) float64 {

	var max float64 = -1

	for r := range v1 {
		diff := v1[r]-v2[r]

		if diff < 0 {
			diff = -diff
		}

		if diff > max {
			max = diff
		}
	}

	return max
}

//**************************************************************

func FlatSTIndex(stindex int) int {

	// Return positive STtype from STAindex

	p_sttype := stindex - ST_ZERO

	if p_sttype < 0 {
		p_sttype = -p_sttype
	}

	return p_sttype
}

//**************************************************************

func PrintLink(l Link) {

	to := GetNodeTxtFromPtr(l.Dst)
	arrow := ARROW_DIRECTORY[l.Arr]
	Verbose("\t ... --(",arrow.Long,",",l.Wgt,")->",to,l.Ctx," \t . . .",PrintSTAIndex(arrow.STAindex))
}

//**************************************************************

func ValidateLinkArgs(s string) []ArrowPtr {

	list := strings.Split(s,",")
	var search_list []ArrowPtr

	if s == "" || s == "all" {
		return nil
	}

	for i := range list {
		v,ok := ARROW_SHORT_DIR[list[i]]

		if ok {
			typ := ARROW_DIRECTORY[v].STAindex-ST_ZERO
			if typ < 0 {
				typ = -typ
			}

			name := ARROW_DIRECTORY[v].Long
			ptr := ARROW_DIRECTORY[v].Ptr

			fmt.Println(" - including search pathway STtype",SST_NAMES[typ],"->",name)
			search_list = append(search_list,ptr)

			if typ != NEAR {
				inverse := INVERSE_ARROWS[ptr]
				fmt.Println("   including inverse meaning",ARROW_DIRECTORY[inverse].Long)
				search_list = append(search_list,inverse)
			}
		} else {
			fmt.Println("\nThere is no link abbreviation called ",list[i])
			os.Exit(-1)
		}
	}

	return search_list
}

//**************************************************************

func AssembleInvolvedNodes(search_list []ArrowPtr) ([]NodePtr,map[RCtype]float64) {

	var node_list []NodePtr
	var weights = make(map[RCtype]float64)

	for class := N1GRAM; class <= GT1024; class++ {

		switch class {
		case N1GRAM:
			for n := range NODE_DIRECTORY.N1directory {
				node_list = SearchIncidentRowClass(NODE_DIRECTORY.N1directory[n],search_list,node_list,weights)
			}
		case N2GRAM:
			for n := range NODE_DIRECTORY.N2directory {
				node_list = SearchIncidentRowClass(NODE_DIRECTORY.N2directory[n],search_list,node_list,weights)
			}
		case N3GRAM:
			for n := range NODE_DIRECTORY.N3directory {
				node_list = SearchIncidentRowClass(NODE_DIRECTORY.N3directory[n],search_list,node_list,weights)
			}
		case LT128:
			for n := range NODE_DIRECTORY.LT128 {
				node_list = SearchIncidentRowClass(NODE_DIRECTORY.LT128[n],search_list,node_list,weights)
			}
		case LT1024:
			for n := range NODE_DIRECTORY.LT1024 {
				node_list = SearchIncidentRowClass(NODE_DIRECTORY.LT1024[n],search_list,node_list,weights)
			}
		case GT1024:
			for n := range NODE_DIRECTORY.GT1024 {
				node_list = SearchIncidentRowClass(NODE_DIRECTORY.GT1024[n],search_list,node_list,weights)
			}
		}
	}

	return node_list,weights
}

//**************************************************************

func SearchIncidentRowClass(node Node, searcharrows []ArrowPtr,node_list []NodePtr,ret_weights map[RCtype]float64) []NodePtr {

	var row_nodes = make(map[NodePtr]bool)
	var ret_nodes []NodePtr

        var rc,cr RCtype

	rc.Row = node.NPtr // transposes
        cr.Col = node.NPtr

	// flip backward facing arrows
	const inverse_flip_arrow = ST_ZERO

        // Only sum over outgoing (+) links
	
	for sttype := ST_ZERO; sttype < len(node.I); sttype++ {
		
		for lnk := range node.I[sttype] {
			arrowptr := node.I[sttype][lnk].Arr
			
			if len(searcharrows) == 0 {
				match := node.I[sttype][lnk]
				row_nodes[match.Dst] = true
				rc.Col = match.Dst
				cr.Row = match.Dst

				if sttype < inverse_flip_arrow {
					ret_weights[cr] += match.Wgt  // flip arrow
				} else {
					ret_weights[rc] += match.Wgt
				}
			} else {
				for l := range searcharrows {
					if arrowptr == searcharrows[l] {
						match := node.I[sttype][lnk]
						row_nodes[match.Dst] = true
						rc.Col = match.Dst
						cr.Row = match.Dst
						if sttype < inverse_flip_arrow {
							ret_weights[cr] += match.Wgt  // flip arrow
						} else {
							ret_weights[rc] += match.Wgt
						}
						
					}
				}
			}
		}
	}
	
	if len(row_nodes) > 0 {
		row_nodes[node.NPtr] = true // Add the parent if it has children
	}

	for nptr := range node_list {
		row_nodes[node_list[nptr]] = true
	}
	
	// Merge idempotently
	
	for nptr := range row_nodes {
		ret_nodes = append(ret_nodes,nptr)
	}

	return ret_nodes
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

func AddMandatory() {

	//   + then the next is (then) - previous (prior)

	arr := InsertArrowDirectory("leadsto",SEQUENCE_RELN,SEQUENCE_RELN,"+")
	inv := InsertArrowDirectory("leadsto","prev","follows on from","-")
	InsertInverseArrowDirectory(arr,inv)

	// for rendering from the database in a web browser

	arr = InsertArrowDirectory("properties","url","has URL","+")
        inv = InsertArrowDirectory("properties","isurl","is a URL for","-")
	InsertInverseArrowDirectory(arr,inv)

	arr = InsertArrowDirectory("properties","img","has image","+")
        inv = InsertArrowDirectory("properties","isimg","is an image for","-")
	InsertInverseArrowDirectory(arr,inv)

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

		if IsQuote(quote) && IsBackReference(src,pos) {
			token = "\""
			pos++
		} else {
			if pos+2 < len(src) && IsWhiteSpace(src[pos+1],src[pos+2]) {
				ParseError(ERR_ILLEGAL_QUOTED_STRING_OR_REF)
				os.Exit(-1)
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
		link := GetLinkArrowByName(token)
		LINE_ITEM_STATE = ROLE_RELATION
		LINE_RELN_CACHE["THIS"] = append(LINE_RELN_CACHE["THIS"],link)
		LINE_RELN_COUNTER++

	case '"': // prior reference
		result := LookupAlias("PREV",LINE_ITEM_COUNTER)
		LINE_ITEM_CACHE["THIS"] = append(LINE_ITEM_CACHE["THIS"],result)
		StoreAlias(result)
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
		PVerbose("fyi, line reference",token,"resolved to",actual)
		AssessGrammarCompletions(actual,LINE_ITEM_STATE)
		LINE_ITEM_STATE = ROLE_LOOKUP
		LINE_ITEM_COUNTER++

	default:
		LINE_ITEM_CACHE["THIS"] = append(LINE_ITEM_CACHE["THIS"],token)
		StoreAlias(token)
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
		last_iptr := LINE_ITEM_REFS[LINE_ITEM_COUNTER-2]
		this_iptr := HandleNode(this_item)
		IdempAddLink(last_item,last_iptr,last_reln,this_item,this_iptr)
		CheckSection()

	case ROLE_CONTEXT:
		Box("Reset context: ->",this_item)
		ContextEval(this_item,"=")
		CheckSection()

	case ROLE_CONTEXT_ADD:
		PVerbose("Add to context:",this_item)
		ContextEval(this_item,"+")
		CheckSection()

	case ROLE_CONTEXT_SUBTRACT:
		PVerbose("Remove from context:",this_item)
		ContextEval(this_item,"-")
		CheckSection()

	case ROLE_SECTION:
		Box("Set chapter/section: ->",this_item)
		SECTION_STATE = this_item

	default:
		CheckSection()

		if AllCaps(token) {
			ParseError(WARN_NOTE_TO_SELF+" ("+token+")")
			return
		}

		HandleNode(this_item)
		LinkUpStorySequence(this_item)
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

func StoreAlias(name string) {

	if LINE_ALIAS != "" {
		PVerbose("-- Storing alias",LINE_ITEM_CACHE[LINE_ALIAS],name,"as",LINE_ALIAS)
		LINE_ITEM_CACHE[LINE_ALIAS] = append(LINE_ITEM_CACHE[LINE_ALIAS],name)
	}
}

//**************************************************************

func IdempAddLink(from string, frptr NodePtr, link Link,to string, toptr NodePtr) {

	// Add a link index cache pointer directly to a from node

	if from == to {
		ParseError(ERR_ARROW_SELFLOOP)
		os.Exit(-1)
	}

	if link.Wgt != 1 {
		PVerbose("... Relation:",from,"--(",ARROW_DIRECTORY[link.Arr].Long,",",link.Wgt,")->",to,link.Ctx)
	} else {
		PVerbose("... Relation:",from,"--",ARROW_DIRECTORY[link.Arr].Long,"->",to,link.Ctx)
	}

	AppendLinkToNode(frptr,link,toptr)

	// Double up the reverse definition for easy indexing of both in/out arrows
	// But be careful not the make the graph undirected by mistake

	invlink := GetLinkArrowByName(ARROW_DIRECTORY[INVERSE_ARROWS[link.Arr]].Short)

	AppendLinkToNode(toptr,invlink,frptr)

}

//**************************************************************

func HandleNode(annotated string) NodePtr {

	clean_ptr,clean_version := IdempAddNode(annotated)

	PVerbose("Event/item/node:",clean_version,"in chapter",SECTION_STATE)

	LINE_ITEM_REFS = append(LINE_ITEM_REFS,clean_ptr)
	
	if len(clean_version) != len(annotated) {
		AddBackAnnotations(clean_version,clean_ptr,annotated)
	}

	return clean_ptr
}

//**************************************************************

func IdempAddNode(s string) (NodePtr,string) {

	clean_version := StripAnnotations(s)

	l := len(s)
	c := ClassifyString(s,l)

	var new_nodetext Node
	new_nodetext.S = clean_version
	new_nodetext.L = l
	new_nodetext.Chap = SECTION_STATE
	new_nodetext.NPtr.Class = c

	iptr := AppendTextToDirectory(new_nodetext)

	return iptr,clean_version
}

//**************************************************************

func GetNodeTxtFromPtr(frptr NodePtr) string {

	class := frptr.Class
	index := frptr.CPtr

	var node Node

	switch class {
	case N1GRAM:
		node = NODE_DIRECTORY.N1directory[index]
	case N2GRAM:
		node = NODE_DIRECTORY.N2directory[index]
	case N3GRAM:
		node = NODE_DIRECTORY.N3directory[index]
	case LT128:
		node = NODE_DIRECTORY.LT128[index]
	case LT1024:
		node = NODE_DIRECTORY.LT1024[index]
	case GT1024:
		node = NODE_DIRECTORY.GT1024[index]
	}

	return node.S
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

func AppendTextToDirectory(event Node) NodePtr {

	var cnode_slot ClassedNodePtr = -1
	var ok bool = false
	var node_alloc_ptr NodePtr

	switch event.SizeClass {
	case N1GRAM:
		cnode_slot,ok = NODE_DIRECTORY.N1grams[event.S]
	case N2GRAM:
		cnode_slot,ok = NODE_DIRECTORY.N2grams[event.S]
	case N3GRAM:
		cnode_slot,ok = NODE_DIRECTORY.N3grams[event.S]
	case LT128:
		cnode_slot,ok = LinearFindText(NODE_DIRECTORY.LT128,event)
	case LT1024:
		cnode_slot,ok = LinearFindText(NODE_DIRECTORY.LT1024,event)
	case GT1024:
		cnode_slot,ok = LinearFindText(NODE_DIRECTORY.GT1024,event)
	}

	node_alloc_ptr.Class = event.SizeClass

	if ok {
		node_alloc_ptr.CPtr = cnode_slot
		return node_alloc_ptr
	}

	switch event.SizeClass {
	case N1GRAM:
		cnode_slot = NODE_DIRECTORY.N1_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.N1directory = append(NODE_DIRECTORY.N1directory,event)
		NODE_DIRECTORY.N1grams[event.S] = cnode_slot
		NODE_DIRECTORY.N1_top++ 
		return node_alloc_ptr
	case N2GRAM:
		cnode_slot = NODE_DIRECTORY.N2_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.N2directory = append(NODE_DIRECTORY.N2directory,event)
		NODE_DIRECTORY.N2grams[event.S] = cnode_slot
		NODE_DIRECTORY.N2_top++
		return node_alloc_ptr
	case N3GRAM:
		cnode_slot = NODE_DIRECTORY.N3_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.N3directory = append(NODE_DIRECTORY.N3directory,event)
		NODE_DIRECTORY.N3grams[event.S] = cnode_slot
		NODE_DIRECTORY.N3_top++
		return node_alloc_ptr
	case LT128:
		cnode_slot = NODE_DIRECTORY.LT128_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.LT128 = append(NODE_DIRECTORY.LT128,event)
		NODE_DIRECTORY.LT128_top++
		return node_alloc_ptr
	case LT1024:
		cnode_slot = NODE_DIRECTORY.LT1024_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.LT1024 = append(NODE_DIRECTORY.LT1024,event)
		NODE_DIRECTORY.LT1024_top++
		return node_alloc_ptr
	case GT1024:
		cnode_slot = NODE_DIRECTORY.GT1024_top
		node_alloc_ptr.CPtr = cnode_slot
		event.NPtr = node_alloc_ptr
		NODE_DIRECTORY.GT1024 = append(NODE_DIRECTORY.GT1024,event)
		NODE_DIRECTORY.GT1024_top++
		return node_alloc_ptr
	}

	return NO_NODE_PTR
}

//**************************************************************

func AppendLinkToNode(frptr NodePtr,link Link,toptr NodePtr) {

	frclass := frptr.Class
	frm := frptr.CPtr
	sttype := ARROW_DIRECTORY[link.Arr].STAindex

	link.Dst = toptr // fill in the last part of the reference

	switch frclass {

	case N1GRAM:
		NODE_DIRECTORY.N1directory[frm].I[sttype] = append(NODE_DIRECTORY.N1directory[frm].I[sttype],link)
	case N2GRAM:
		NODE_DIRECTORY.N2directory[frm].I[sttype] = append(NODE_DIRECTORY.N2directory[frm].I[sttype],link)
	case N3GRAM:
		NODE_DIRECTORY.N3directory[frm].I[sttype] = append(NODE_DIRECTORY.N3directory[frm].I[sttype],link)
	case LT128:
		NODE_DIRECTORY.LT128[frm].I[sttype] = append(NODE_DIRECTORY.LT128[frm].I[sttype],link)
	case LT1024:
		NODE_DIRECTORY.LT1024[frm].I[sttype] = append(NODE_DIRECTORY.LT1024[frm].I[sttype],link)
	case GT1024:
		NODE_DIRECTORY.GT1024[frm].I[sttype] = append(NODE_DIRECTORY.GT1024[frm].I[sttype],link)
	}
}

//**************************************************************

func LinearFindText(in []Node,event Node) (ClassedNodePtr,bool) {

	for i := 0; i < len(in); i++ {

		if event.L != in[i].L {
			continue
		}

		if in[i].S == event.S {
			return ClassedNodePtr(i),true
		}
	}

	return -1,false
}

//**************************************************************
// Scan text input
//**************************************************************

func ReadFile(filename string) []rune {

	text := ReadTUF8File(filename)

	// clean unicode nonsense

	for r := range text {
		switch text[r] {
		case NON_ASCII_LQUOTE,NON_ASCII_RQUOTE:
			text[r] = '"'
		}
	}

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
			PVerbose("\nStart sequence mode for items")
			SEQUENCE_MODE = true
			LAST_IN_SEQUENCE = ""

		case '-':
			PVerbose("End sequence mode for items\n")
			SEQUENCE_MODE = false
		}
	}

}

//**************************************************************

func LinkUpStorySequence(this string) {

	// Join together a sequence of nodes using default "(then)"

	if SEQUENCE_MODE && this != LAST_IN_SEQUENCE {

		if LINE_ITEM_COUNTER == 1 && LAST_IN_SEQUENCE != "" {
			
			PVerbose("* ... Sequence addition: ",LAST_IN_SEQUENCE,"-(",SEQUENCE_RELN,")->",this,"\n")
			
			last_iptr,_ := IdempAddNode(LAST_IN_SEQUENCE)
			this_iptr,_ := IdempAddNode(this)
			link := GetLinkArrowByName("(then)")
			AppendLinkToNode(last_iptr,link,this_iptr)

			invlink := GetLinkArrowByName(ARROW_DIRECTORY[INVERSE_ARROWS[link.Arr]].Short)
			AppendLinkToNode(this_iptr,invlink,last_iptr)
		}
		
		LAST_IN_SEQUENCE = this
	}
}

//**************************************************************

func CheckArrow(alias,name string) {

	prev,ok := ARROW_SHORT_DIR[alias]
	if ok {
		ParseError(ERR_ARR_REDEFINITION+"\""+alias+"\" previous short name: "+ARROW_DIRECTORY[prev].Short)
		os.Exit(-1)
	}
	
	prev,ok = ARROW_LONG_DIR[name]
	if ok {
		ParseError(ERR_ARR_REDEFINITION+"\""+name+"\" previous long name: "+ARROW_DIRECTORY[prev].Long)
		os.Exit(-1)
	}
}

//**************************************************************

func GetLinkArrowByName(token string) Link {

	// Return a preregistered link/arrow ptr bythe name of a link

	var reln []string
	var weight float64 = 1
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
				weight = v
				weightcount++
			} else {
				ctx = append(ctx,reln[i])
			}
		}
	}

	// First check if this is an alias/short name

	ptr, ok := ARROW_SHORT_DIR[name]

	// If not, then check longname

	if !ok {
		ptr, ok = ARROW_LONG_DIR[name]

		if !ok {
			ParseError(ERR_NO_SUCH_ARROW+name)
			os.Exit(-1)
		}
	}

	var link Link

	link.Arr = ptr
	link.Wgt = weight
	link.Ctx = GetContext(ctx)
	return link
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

	if len(token) < 1 {
		return ""
	}

	split := strings.Split(token[1:],".")

	if len(split) < 2 {
		ParseError(ERR_MISSING_LINE_LABEL_IN_REFERENCE)
		os.Exit(-1)
	}

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

func StripAnnotations(fulltext string) string {

	var protected bool = false
	var deloused []rune
	var preserve_unicode = []rune(fulltext)

	for r := 0; r < len(preserve_unicode); r++ {

		if preserve_unicode[r] == '"' {
			protected = !protected
		}

		if !protected {
			skip,symb := EmbeddedSymbol(preserve_unicode,r)
			if skip > 0 {
				r += skip-1
				if unicode.IsSpace(preserve_unicode[r]) {
					ParseError(ERR_NON_WORD_WHITE+symb)
				}
				continue
			}
		}

		deloused = append(deloused,preserve_unicode[r])
	}

	return string(deloused)
}

//**************************************************************

func AddBackAnnotations(cleantext string,cleanptr NodePtr,annotated string) {

	var protected bool = false

	reminder := fmt.Sprintf("%.30s...",cleantext)
	PVerbose("\n        Adding annotations from \""+reminder+"\"")

	for r := 0; r < len(annotated); r++ {

		if annotated[r] == '"' {
			protected = !protected
		} else {
			if !protected {
				skip,symb := EmbeddedSymbol([]rune(annotated),r)
				if skip > 0 {
					link := GetLinkArrowByName(ANNOTATION[symb])
					this_item := ExtractWord(annotated,r)
					this_iptr,_ := IdempAddNode(this_item)
					IdempAddLink(reminder,cleanptr,link,this_item,this_iptr)
					r += skip-1
					continue
				}
			}
		}
	}
}

//**************************************************************

func EmbeddedSymbol(fulltext []rune,offset int) (int,string) {

	for an := range ANNOTATION {

		// Careful of unicode, convert to runes
		uni := []rune(an)
		match := true

		for r := 0; r < len(uni) && r+offset < len(fulltext); r++ {
			if uni[r] != fulltext[offset+r] {
				match = false
			}
		} 

		if match {
			return len(an),an
		}
	}

	return 0,"UNKNOWN SYMBOL"
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

	if len(s) <= WORD_MISTAKE_LEN {
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

func Usage() {
	
	fmt.Printf("usage: N4L [-v] [-u] [-s] [file].dat\n")
	flag.PrintDefaults()
	os.Exit(2)
}

//**************************************************************

func Verbose(a ...interface{}) {

	line := fmt.Sprintln(a...)
	
	if DIAGNOSTIC {
		AppendStringToFile(TEST_DIAG_FILE,line)
	}


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

func DiagnosticName(filename string) string {

	return "test_output/"+filename+"_test_log"

}

//**************************************************************

func Diag(a ...interface{}) {

	// Log diagnostic output for self-diagnostic tests

	if DIAGNOSTIC {
		s := fmt.Sprintln(a...)
		prefix := fmt.Sprint(LINE_NUM,":")
		AppendStringToFile(TEST_DIAG_FILE,prefix+s)
	}
}

//**************************************************************

func PrintSTAIndex(st int) string {

	st = st - ST_ZERO
	var ty string

	switch st {
	case -EXPRESS:
		ty = "-(expressed by)"
	case -CONTAINS:
		ty = "-(part of)"
	case -LEADSTO:
		ty = "-(arriving from)"
	case NEAR:
		ty = "(close to)"
	case LEADSTO:
		ty = "+(leading to)"
	case CONTAINS:
		ty = "+(containing)"
	case EXPRESS:
		ty = "+(expressing)"
	default:
		ty = "unknown relation!"
	}

	const green = "\x1b[36m"
	const endgreen = "\x1b[0m"

	return green + ty + endgreen
}

//**************************************************************

func AppendStringToFile(name string, s string) {

	// strip out \r that mess up the file format but are useful for term

	san := strings.Replace(s,"\r","",-1)

	f, err := os.OpenFile(name,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Couldn't open for write/append to",name,err)
		return
	}

	_, err = f.WriteString(san)

	if err != nil {
		fmt.Println("Couldn't write/append to",name,err)
	}

	f.Close()
}

