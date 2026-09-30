package oto_test

import (
	"fmt"
	"strings"

	"github.com/yh2237/utauio/oto"
)

func ExampleScan() {
	diagnostics, err := oto.Scan(strings.NewReader("a.wav=あ,12.5,100,-20,30,10\nbroken\n"), "voice", func(entry oto.Entry) {
		fmt.Println(entry.Alias, entry.Offset, entry.Line)
	})
	fmt.Println(len(diagnostics), err)
	// Output:
	// あ 12.5 1
	// 1 <nil>
}
