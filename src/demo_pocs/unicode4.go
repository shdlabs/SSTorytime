
package main

import (
    "fmt"

    "golang.org/x/text/runes"
    "golang.org/x/text/secure/precis"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)



func main() {

    trans := transform.Chain(
        norm.NFD,
        precis.UsernameCaseMapped.NewTransformer(),
        runes.Map(func(r rune) rune {
            switch r {
            case 'ą':
                return 'a'
            case 'ć':
                return 'c'
            case 'ę':
                return 'e'
            case 'ł':
                return 'l'
            case 'ń':
                return 'n'
            case 'ó':
                return 'o'
            case 'ś':
                return 's'
            case 'ż':
                return 'z'
            case 'ź':
                return 'z'
            }
            return r
        }),
        norm.NFC,
    )
    result, _, _ := transform.String(trans, "ŻóŁć")
    fmt.Println(result)
}



