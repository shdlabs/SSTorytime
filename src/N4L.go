//
// N4LParser
//

package main

import (
	//"strings"
	"os"
	//"io"
	"io/ioutil"
	"flag"
	"fmt"
	"unicode/utf8"
	"unicode"
)

//**************************************************************

// Notes states

const (
	start = iota
	section
	notes
	context
)

// Config states

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

type NoteParser struct 
{
pos              int
mode             int
topic            string   // in case we mix together several
context_set      []string
item_set         []string
relation_set     []string
}

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

		fmt.Println("GOT",token)
	}

/*	switch (p.mode) {

	case init: 
	case section:
	case notes:
	case context:

	}*/

}

//**************************************************************
// Parsing objects
//**************************************************************

func SkipWhiteSpace(src []rune, pos int) int {

	for ; unicode.IsSpace(src[pos]) || src[pos] == '#' || src[pos] == '/' ; pos++ {

		if src[pos] == '#' || (src[pos] == '/' && src[pos+1] == '/') {
			for ; src[pos] != '\n'; pos++ {
			}
		}
	}

	return pos
}

//**************************************************************

func GetToken(src []rune, pos int) (string,int) {

	// Handle concatenation of words/lines and separation of types

	var token string

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
		token,pos = ReadToLast(src,pos,'"')

	default: // a text item that could end with any of the above
		token,pos = ReadToLast(src,pos,'x')
	}

	return token, pos

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

	if stop == 'x' {

		for ; pos > 0 && !(unicode.IsLetter(src[pos-1]) && !unicode.IsLetter(src[pos])); pos++ {
			cpy = append(cpy,src[pos])
		}

	} else {
		var passed_lhs bool = false

		for ; pos > 0 && !(passed_lhs && src[pos-1] == stop && src[pos] != stop); pos++ {

			// In the case of ::, we need to count a contiguous block

			if src[pos] != ':' && src[pos+1] != ':' {
				passed_lhs = true
			}

			cpy = append(cpy,src[pos])
		}

	}

	token := string(cpy)

	return token,pos
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
