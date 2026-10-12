package character

import (
	"reflect"
	"testing"
)

func TestParseText(t *testing.T) {
	got := ParseText("name=first\r\nunknown=value\n NAME = テスト=音源 \nname= \nnot a field\n")
	if got.Name != "テスト=音源" {
		t.Fatalf("name=%q", got.Name)
	}
	if got := ParseText("author=test\n"); got.Name != "" {
		t.Fatalf("missing name=%q", got.Name)
	}
}

func TestParseYAMLSubbanks(t *testing.T) {
	text := "default_phonemizer: 'ja-kana'\r\nsubbanks:\r\n" +
		"- color: \"\"\n  prefix: 'it''s'\n  suffix: \" C#4\" # comment\n  tone_ranges:\n  - B3-B4\n" +
		"- color: power\n  suffix: '# P'\n  tone_ranges: [C3-C5, Db4-F#4]\n  unknown: ignored\n"
	got, diagnostics := ParseYAML(text)
	want := YAML{DefaultPhonemizer: "ja-kana", Subbanks: []Subbank{
		{Prefix: "it's", Suffix: " C#4", ToneRanges: []ToneRange{{59, 71}}},
		{Color: "power", Suffix: "# P", ToneRanges: []ToneRange{{48, 72}, {61, 66}}},
	}}
	if !reflect.DeepEqual(got, want) || len(diagnostics) != 0 {
		t.Fatalf("got=%+v diagnostics=%+v", got, diagnostics)
	}
}

func TestParseYAMLDiagnosticsAndMissingRanges(t *testing.T) {
	got, diagnostics := ParseYAML("subbanks:\n- color: test\n  tone_ranges: [bad, C5-C3, C4-C5]\n- color: all\n  tone_ranges:\n  - not-a-tone\n- color: last\n")
	wantDiagnostics := []Diagnostic{
		{3, `invalid tone range "bad"`},
		{3, `invalid tone range "C5-C3"`},
		{6, `invalid tone range "not-a-tone"`},
	}
	if !reflect.DeepEqual(diagnostics, wantDiagnostics) {
		t.Fatalf("diagnostics=%+v", diagnostics)
	}
	want := []Subbank{
		{Color: "test", ToneRanges: []ToneRange{{60, 72}}},
		{Color: "all", ToneRanges: []ToneRange{{-1 << 30, 1<<30 - 1}}},
		{Color: "last", ToneRanges: []ToneRange{{-1 << 30, 1<<30 - 1}}},
	}
	if !reflect.DeepEqual(got.Subbanks, want) {
		t.Fatalf("subbanks=%+v", got.Subbanks)
	}
}

func TestParseYAMLCompatibilitySubset(t *testing.T) {
	got, diagnostics := ParseYAML("DEFAULT_PHONEMIZER: first # literal\ndefault_phonemizer: second\nsubbanks:\n- prefix: ignored\n- color: test\n  suffix: \"\\tC4\"\n  tone_ranges: [C4]\n")
	if got.DefaultPhonemizer != "first # literal" || len(got.Subbanks) != 1 || got.Subbanks[0].Suffix != "\tC4" {
		t.Fatalf("config=%+v", got)
	}
	if len(diagnostics) != 1 || diagnostics[0].Line != 7 {
		t.Fatalf("diagnostics=%+v", diagnostics)
	}
}

func FuzzParse(f *testing.F) {
	for _, text := range []string{"", "name=音源", "subbanks:\n- color: ''\n  tone_ranges: [C3-C5]\n", "subbanks:\n- color: test\n  tone_ranges: [C#-D4, bad]\n"} {
		f.Add(text)
	}
	f.Fuzz(func(t *testing.T, text string) {
		_ = ParseText(text)
		config, diagnostics := ParseYAML(text)
		for _, diagnostic := range diagnostics {
			if diagnostic.Line < 1 || diagnostic.Message == "" {
				t.Fatalf("diagnostic=%+v", diagnostic)
			}
		}
		for _, bank := range config.Subbanks {
			for _, span := range bank.ToneRanges {
				if span.Low > span.High {
					t.Fatalf("range=%+v", span)
				}
			}
		}
	})
}
