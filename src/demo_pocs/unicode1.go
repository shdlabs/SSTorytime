

package main

import (
    "fmt"
    "unicode"

    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

// *****************************************************************

func main() {

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	result, _, _ := transform.String(t, "žůžo  ÅåøØææÆÆ Tā pǎole shànglái")

	fmt.Println(result)

	// r = transform.NewReader(r, t) 
	// read as before ...


	// wc := norm.NFC.Writer(w)
	// defer wc.Close()
	// write as before...

}

// *****************************************************************

func isMn(r rune) bool {

	// Example derived from: https://go.dev/blog/normalization

	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}





