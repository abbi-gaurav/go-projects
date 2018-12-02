package main

import (
	"unicode/utf8"
	"fmt"
)

func main() {
	s := "π† some funky chars in start"
	var buf [utf8.UTFMax] byte

	for i, r := range s {
		rl := utf8.RuneLen(r)
		si := i + rl

		copy(buf[:], s[i:si])

		fmt.Printf("%2d: %q; codepoint: %#6x; encoded bytes: %#v\n", i,r,r, buf[:rl])
	}
}
