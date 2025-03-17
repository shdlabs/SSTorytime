

package main

import (
    "fmt"
    "unicode"

    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

// *****************************************************************

func main() {

        // This function is deprecated
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	result, _, _ := transform.String(t, "žůžo  ÅåøØææÆÆ Tā pǎole shànglái")

	fmt.Println(result)

	// r = transform.NewReader(r, t) 
	// read as before ...


	// wc := norm.NFC.Writer(w)
	// defer wc.Close()
	// write as before...


	// New replacement function

	input := []byte(`tschüß; до свидания`)

	b := make([]byte, len(input))

	t = transform.RemoveFunc(unicode.IsSpace)
	n, _, _ := t.Transform(b, input, true)
	fmt.Println(string(b[:n]))

	t = transform.RemoveFunc(func(r rune) bool {
		return !unicode.Is(unicode.Latin, r)
	})
	n, _, _ = t.Transform(b, input, true)
	fmt.Println(string(b[:n]))

	n, _, _ = t.Transform(b, norm.NFD.Bytes(input), true)
	fmt.Println(string(b[:n]))


}

// *****************************************************************

func isMn(r rune) bool {

	// Example derived from: https://go.dev/blog/normalization

	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}





