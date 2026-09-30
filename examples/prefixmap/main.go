package main

import (
	"fmt"
	"github.com/yh2237/utauio/prefixmap"
)

func main() {
	diagnostics, err := prefixmap.Scan("C4\t\t C4\n", func(e prefixmap.Entry) { fmt.Printf("%s prefix=%q suffix=%q\n", e.Tone, e.Prefix, e.Suffix) })
	fmt.Println(len(diagnostics), err)
}
