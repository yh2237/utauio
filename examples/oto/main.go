package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yh2237/utauio/oto"
	"github.com/yh2237/utauio/textdecode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./examples/oto path/to/oto.ini")
		os.Exit(2)
	}
	if err := run(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	text, encoding, err := textdecode.Decode(data)
	if err != nil {
		return err
	}
	count := 0
	diagnostics, err := oto.Scan(strings.NewReader(text), filepath.Dir(path), func(entry oto.Entry) { count++ })
	if err != nil {
		return err
	}
	fmt.Printf("encoding=%s entries=%d diagnostics=%d\n", encoding, count, len(diagnostics))
	for _, diagnostic := range diagnostics {
		fmt.Printf("line %d: %s\n", diagnostic.Line, diagnostic.Message)
	}
	return nil
}
