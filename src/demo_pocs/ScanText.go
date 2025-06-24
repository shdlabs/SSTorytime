
//
// transform random text or book, suggesting arrow hints for N4LParser
//   Conclusion : this approach is misguided. Need something more authoritative 
//

package main

import (
	"strings"
	"os"
	"flag"
	"fmt"
	"regexp"

        SST "SSTorytime"
)

//**************************************************************

var CURRENT_FILE string
var VERBOSE bool
var LINE_NUM int

type Match struct {
	Arrow   string
	Before  string
	After   string
}

//**************************************************************
// BEGIN
//**************************************************************

func main() {

// load arrows

	args := Init()

	// input := "/home/mark/Laptop/Work/SST/data_samples/MobyDick.dat"

	for input := 0; input < len(args); input++ {

		NewFile(args[input])
		SST.FractionateTextFile(CURRENT_FILE)
		AnnotateFile()
	}
}

//**************************************************************

func AnnotateFile() {
	
	for n := SST.N_GRAM_MIN; n < SST.N_GRAM_MAX; n++ {
		for g := range SST.STM_NGRAM_RANK[n] {

			fmt.Println("ng",n,g)
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

	s = strings.Replace(s,"(","[",-1)
	s = strings.Replace(s,")","]",-1)

	m := regexp.MustCompile("<[^>]*>") 
	s = m.ReplaceAllString(s,":\n") 

	// Weird English abbrev
	s = strings.Replace(s,"Mr.","Mr",-1) 
	s = strings.Replace(s,"Ms.","Ms",-1) 
	s = strings.Replace(s,"Mrs.","Mrs",-1) 
	s = strings.Replace(s,"Dr.","Dr",-1)
	s = strings.Replace(s,"St.","St",-1) 

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

func ParaBySentence(paras []string) [][]string {
	
	var pbs [][]string
	
	for s := range paras {
		
		var para []string
		
		sentences := strings.Split(paras[s],"#")
		
		for s := range sentences {
			sentence := strings.TrimSpace(sentences[s])
			if sentence != ":" && sentence != "" {
				para = append(para,sentence)
			}
		}
		pbs = append(pbs,para)
	}

	return pbs
}

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



