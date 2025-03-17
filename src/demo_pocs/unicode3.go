package main

import (
    "fmt"
    "unicode"

    "golang.org/x/text/secure/precis"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

// **************************************************************

func main() {

	// makes use of deprecated RemoveFunc
	loosecompare := precis.NewIdentifier(
		precis.AdditionalMapping(func() transform.Transformer {
			return transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
				return unicode.Is(unicode.Mn, r)
			}))
		}),
		precis.Norm(norm.NFC), // This is the default; be explicit though.
	)
	p, _ := loosecompare.String("žůžo")
	
	fmt.Println(p, loosecompare.Compare("žůžo", "zuzo"))
	// Prints "zuzo true"

}
