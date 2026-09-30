package presamp

import (
	"reflect"
	"testing"
)

func TestParseClassesReplacementsAndSelectedEndings(t *testing.T) {
	text := "; comment\r\n[vowel]\r\na=a= あ,い,,う =100=extra\r\nb=b=あ=100\r\n[consonant]\r\nk=か,き=0\r\n[replace]\r\nx=1=2\r\n[ENDTYPE2]\r\nb\r\n[ENDTYPE1]\r\na\r\n[ENDFLAG]\r\n1"
	got := Parse(text)
	if !reflect.DeepEqual(got.Vowels, map[string]string{"あ": "b", "い": "a", "う": "a"}) || !reflect.DeepEqual(got.Consonants, map[string]string{"か": "k", "き": "k"}) || got.Replacements["x"] != "1=2" || !reflect.DeepEqual(got.Endings, []string{"a"}) {
		t.Fatalf("unexpected configuration: %+v", got)
	}
}

func TestParseMalformedRowsUnknownSectionsAndEndFlags(t *testing.T) {
	got := Parse("[VOWEL]\nbad\n =a=x\na=a=,x,,y,\n[REPLACE]\n=x\ny=\nz=  x  \n[UNKNOWN]\na=b\n[ENDTYPE]\na\n[ENDTYPEbad]\nb\n[ENDTYPE0]\nc\n[ENDFLAG]\ninvalid\n")
	if len(got.Vowels) != 2 || len(got.Replacements) != 1 || got.Replacements["z"] != "x" || !reflect.DeepEqual(got.Endings, []string{"a", "b", "c"}) {
		t.Fatal(got)
	}
	if Parse("").Vowels == nil {
		t.Fatal("empty result maps must be initialized")
	}
	negative := Parse("[ENDTYPE2]\nb\n[ENDTYPE1]\na\n[ENDFLAG]\n-1\n")
	if !reflect.DeepEqual(negative.Endings, []string{"a", "b"}) {
		t.Fatal("ending order changed", negative.Endings)
	}
}

func FuzzParse(f *testing.F) {
	f.Add("[VOWEL]\na=a=あ,い=100\n")
	f.Add("[ENDTYPE999999999999999999999]\na\n[ENDFLAG]\n-1\n")
	f.Fuzz(func(t *testing.T, text string) {
		if len(text) > 65536 {
			t.Skip()
		}
		got := Parse(text)
		for _, classes := range []map[string]string{got.Vowels, got.Consonants, got.Replacements} {
			if classes == nil {
				t.Fatal("uninitialized map")
			}
			for key, value := range classes {
				if key == "" || value == "" {
					t.Fatal("empty key/value accepted")
				}
			}
		}
	})
}
