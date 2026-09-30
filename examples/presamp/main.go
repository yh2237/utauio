package main

import (
	"fmt"
	"github.com/yh2237/utauio/presamp"
)

func main() {
	config := presamp.Parse("[VOWEL]\na=a=あ,い=100\n[ENDTYPE1]\n%v% R\n")
	fmt.Println(config.Vowels["あ"], config.Vowels["い"], config.Endings)
}
