package textdecode

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
)

func TestDecodeSupportedEncodings(t *testing.T) {
	text := "あ.wav=あ,0,100,-20,30,10\r\n"
	shift, err := japanese.ShiftJIS.NewEncoder().Bytes([]byte(text))
	if err != nil {
		t.Fatal(err)
	}
	little, big := []byte{0xff, 0xfe}, []byte{0xfe, 0xff}
	for _, unit := range utf16.Encode([]rune(text)) {
		little = binary.LittleEndian.AppendUint16(little, unit)
		big = binary.BigEndian.AppendUint16(big, unit)
	}
	for _, tc := range []struct {
		name, encoding string
		data           []byte
	}{{"UTF8", "UTF-8", []byte(text)}, {"UTF8 BOM", "UTF-8", []byte("\ufeff" + text)},
		{"UTF16 LE", "UTF-16LE", little}, {"UTF16 BE", "UTF-16BE", big}, {"ShiftJIS", "Shift_JIS", shift}} {
		t.Run(tc.name, func(t *testing.T) {
			got, encoding, err := Decode(tc.data)
			if err != nil || got != text || encoding != tc.encoding {
				t.Fatalf("Decode = %q, %q, %v", got, encoding, err)
			}
		})
	}
	for _, data := range [][]byte{{0xff, 0xfe, 0}, {0xfe, 0xff, 0}} {
		if _, _, err := Decode(data); err == nil {
			t.Fatal("odd UTF-16 accepted")
		}
	}
}

func TestDecodeUTF8BOMBeforeShiftJIS(t *testing.T) {
	text := "あ.wav=あ,0,100,-20,30,10\r\n"
	shift, err := japanese.ShiftJIS.NewEncoder().Bytes([]byte(text))
	if err != nil {
		t.Fatal(err)
	}
	data := append([]byte{0xef, 0xbb, 0xbf}, shift...)
	got, encoding, err := Decode(data)
	if err != nil || got != text || encoding != "Shift_JIS" {
		t.Fatalf("Decode = %q, %q, %v", got, encoding, err)
	}
}

var decodeAllocationSink string

func TestDecodeUTF8Allocations(t *testing.T) {
	base := bytes.Repeat([]byte("原音.wav=あ,0,0,-100,10,5\n"), 8)
	if len(base) != 232 {
		t.Fatalf("fixture length = %d, want 232", len(base))
	}
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"UTF8", base},
		{"UTF8 BOM", append([]byte{0xef, 0xbb, 0xbf}, base...)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if text, encoding, err := Decode(tc.data); err != nil || encoding != "UTF-8" || text == "" {
				t.Fatalf("Decode = %q, %q, %v", text, encoding, err)
			}
			allocs := testing.AllocsPerRun(100, func() {
				text, _, _ := Decode(tc.data)
				decodeAllocationSink = text
			})
			if allocs != 1 {
				t.Fatalf("Decode allocated %.1f times per call, want 1", allocs)
			}
		})
	}
}

func FuzzDecode(f *testing.F) {
	for _, data := range [][]byte{{}, []byte("あ.wav"), {0xff, 0xfe}, {0xfe, 0xff, 0}, {0x82, 0xa0},
		[]byte("\ufeffあ.wav"), {0xff, 0xfe, 0x42, 0x30}, {0xfe, 0xff, 0x30, 0x42}} {
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 65536 {
			t.Skip()
		}
		text, encoding, err := Decode(data)
		if err == nil && (!utf8.ValidString(text) || encoding == "") {
			t.Fatal("successful decode produced invalid UTF-8 or no encoding")
		}
	})
}
