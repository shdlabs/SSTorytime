//
// Command-line tool for scanning documents and extracting sentences
// that are high in "intentionality" or potential knowledge significance
//

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shdlabs/SSTorytime/services/text2n4l"
)

var TARGET_PERCENT float64 = 50.0

//**************************************************************
// BEGIN
//**************************************************************

func main() {
	input := getArgs()

	err := text2n4l.ProcessFile(input, TARGET_PERCENT)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}
}

//**************************************************************

func getArgs() string {
	flag.Usage = usage

	limitPtr := flag.Float64("%", 50, "approximate percentage of file to skim (overestimates for small values)")

	flag.Parse()
	args := flag.Args()

	TARGET_PERCENT = *limitPtr

	if len(args) != 1 {
		fmt.Println("Missing pure text filename to scan")
		os.Exit(-2)
	}

	return args[0]
}

//**************************************************************

func usage() {
	fmt.Println("usage: Text2N4L [-% percent] filename\n")
	flag.PrintDefaults()
	os.Exit(2)
}
