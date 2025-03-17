
package main

import (
	"unicode"
	"fmt"
	
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)


// ******************************************************

func main() {

	s := "žůžo  ÅåøØææÆÆ Tā pǎole shànglái"

	fmt.Println(Normalize(s))
}

// ******************************************************

func Normalize(s string) (string, error) {

    t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

    result, _, err := transform.String(t, s)
    if err != nil {
        return "", err
    }

    return result, nil
}




