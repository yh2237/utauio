package prefixmap

import (
	"reflect"
	"testing"
)

func TestScanPreservesAffixesOrderAndDiagnostics(t *testing.T) {
	text := ";comment\r\n c4 \t\t C4\textra\r\nC4\t前\t後\r\nG4\t\r\nbad\r\n\tX\tY\n"
	var entries []Entry
	diagnostics, err := Scan(text, func(e Entry) { entries = append(entries, e) })
	want := []Entry{{"C4", "", " C4", 2}, {"C4", "前", "後", 3}, {"G4", "", "", 4}, {"", "X", "Y", 6}}
	if err != nil || !reflect.DeepEqual(entries, want) || len(diagnostics) != 1 || diagnostics[0].Line != 5 {
		t.Fatalf("%+v %+v %v", entries, diagnostics, err)
	}
	if _, err := Scan("", nil); err == nil {
		t.Fatal("nil visitor accepted")
	}
}

func FuzzScan(f *testing.F) {
	f.Add("C4\t\t_C4\n")
	f.Add("bad\n\tX\tY\r\n")
	f.Fuzz(func(t *testing.T, text string) {
		if len(text) > 65536 {
			t.Skip()
		}
		last := 0
		_, err := Scan(text, func(e Entry) {
			if e.Line <= last {
				t.Fatal("entry order changed")
			}
			last = e.Line
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}
