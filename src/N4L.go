
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
)

//**************************************************************

func main() {

	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	
	if len(args) < 1 {
		usage()
		os.Exit(1);
	}

	ReadFile(args[0])
}

//**************************************************************
// Parsing objects
//**************************************************************

func GetItem(input string, pos int) (string,int) {

	return item, newpos
}

//**************************************************************

func GetRelation(input string, pos int) {


}

//**************************************************************
// Scan text input
//**************************************************************

func ReadFile(filename string) {

	proto_text := ReadTUF8File(filename)

	ParseN4L(proto_text)
}

//**************************************************************

func ParseN4L(s []rune) {

	for r := 0; r < len(s); r++ {

		fmt.Print(string(s[r]))

	}

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
