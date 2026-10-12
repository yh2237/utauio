package main

import (
	"fmt"

	"github.com/yh2237/utauio/character"
)

func main() {
	info := character.ParseText("name=Example voicebank\n")
	config, diagnostics := character.ParseYAML("default_phonemizer: ja-kana\nsubbanks:\n- color: power\n  suffix: ' P'\n  tone_ranges: [C3-C5]\n")
	fmt.Printf("name=%s phonemizer=%s\n", info.Name, config.DefaultPhonemizer)
	for _, bank := range config.Subbanks {
		fmt.Printf("color=%q prefix=%q suffix=%q ranges=%v\n", bank.Color, bank.Prefix, bank.Suffix, bank.ToneRanges)
	}
	for _, diagnostic := range diagnostics {
		fmt.Printf("line %d: %s\n", diagnostic.Line, diagnostic.Message)
	}
}
