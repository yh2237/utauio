// characterパッケージは復号済みcharacter.txtとcharacter.yamlの音源メタデータを解析する。
package character

import (
	"fmt"
	"strconv"
	"strings"
)

type Text struct {
	Name string
}

type YAML struct {
	DefaultPhonemizer string
	Subbanks          []Subbank
}

type Subbank struct {
	Color      string
	Prefix     string
	Suffix     string
	ToneRanges []ToneRange
}

// ToneRangeは両端を含むMIDI音階番号の範囲を表す。
type ToneRange struct {
	Low  int
	High int
}

type Diagnostic struct {
	Line    int
	Message string
}

// ParseTextは名前を読む。空の名前は無視し、重複時は最後の非空値を採用する。
func ParseText(text string) Text {
	var result Text
	for line := range strings.SplitSeq(text, "\n") {
		key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
		if ok && strings.EqualFold(strings.TrimSpace(key), "name") {
			if name := strings.TrimSpace(value); name != "" {
				result.Name = name
			}
		}
	}
	return result
}

// ParseYAMLはUtauTTSから継承した限定的なYAML表記を読む。
// 完全なYAMLパーサーではなく、不正な音域は診断に回して後続の解析を続ける。
func ParseYAML(text string) (YAML, []Diagnostic) {
	var result YAML
	for line := range strings.SplitSeq(text, "\n") {
		key, value, ok := splitField(strings.TrimSpace(line))
		if ok && strings.EqualFold(key, "default_phonemizer") {
			result.DefaultPhonemizer = scalar(value)
			break
		}
	}
	var (
		current     *Subbank
		inList      bool
		inRanges    bool
		diagnostics []Diagnostic
	)
	finish := func() {
		if current == nil {
			return
		}
		if len(current.ToneRanges) == 0 {
			current.ToneRanges = []ToneRange{{Low: -1 << 30, High: 1<<30 - 1}}
		}
		result.Subbanks = append(result.Subbanks, *current)
		current = nil
	}
	lineNumber := 0
	for raw := range strings.SplitSeq(text, "\n") {
		lineNumber++
		line := strings.TrimSpace(stripComment(strings.TrimSuffix(raw, "\r")))
		if line == "" {
			continue
		}
		if line == "subbanks:" {
			finish()
			inList, inRanges = true, false
			continue
		}
		if !inList {
			continue
		}
		if strings.HasPrefix(line, "-") {
			rest := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if strings.HasPrefix(rest, "color:") {
				finish()
				current = &Subbank{}
				inRanges = false
				assignField(current, rest)
				continue
			}
			if current != nil && inRanges {
				toneRange, err := parseToneRange(scalar(rest))
				if err != nil {
					diagnostics = append(diagnostics, Diagnostic{Line: lineNumber, Message: err.Error()})
				} else {
					current.ToneRanges = append(current.ToneRanges, toneRange)
				}
				continue
			}
		}
		if current == nil {
			continue
		}
		key, value, ok := splitField(line)
		if !ok {
			continue
		}
		if key == "tone_ranges" {
			inRanges = true
			for _, item := range inlineList(value) {
				toneRange, err := parseToneRange(item)
				if err != nil {
					diagnostics = append(diagnostics, Diagnostic{Line: lineNumber, Message: err.Error()})
				} else {
					current.ToneRanges = append(current.ToneRanges, toneRange)
				}
			}
			continue
		}
		assignField(current, line)
	}
	finish()
	return result, diagnostics
}

func assignField(subbank *Subbank, line string) {
	key, value, ok := splitField(line)
	if !ok {
		return
	}
	value = scalar(value)
	switch key {
	case "color":
		subbank.Color = value
	case "prefix":
		subbank.Prefix = value
	case "suffix":
		subbank.Suffix = value
	}
}

func splitField(line string) (key, value string, ok bool) {
	key, value, ok = strings.Cut(line, ":")
	key = strings.TrimSpace(key)
	if !ok || key == "" {
		return "", "", false
	}
	return key, strings.TrimSpace(value), true
}

func scalar(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		if decoded, err := strconv.Unquote(value); err == nil {
			return decoded
		}
	}
	if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
		return strings.ReplaceAll(value[1:len(value)-1], "''", "'")
	}
	return value
}

func inlineList(value string) []string {
	value = strings.TrimSpace(value)
	if len(value) < 2 || value[0] != '[' || value[len(value)-1] != ']' {
		return nil
	}
	var result []string
	for item := range strings.SplitSeq(value[1:len(value)-1], ",") {
		if item = scalar(item); item != "" {
			result = append(result, item)
		}
	}
	return result
}

func parseToneRange(value string) (ToneRange, error) {
	value = scalar(value)
	lowText, highText, ok := strings.Cut(value, "-")
	if !ok {
		return ToneRange{}, fmt.Errorf("invalid tone range %q", value)
	}
	low, lowOK := toneNumber(lowText)
	high, highOK := toneNumber(highText)
	if !lowOK || !highOK || low > high {
		return ToneRange{}, fmt.Errorf("invalid tone range %q", value)
	}
	return ToneRange{Low: low, High: high}, nil
}

func toneNumber(tone string) (int, bool) {
	tone = strings.ToUpper(strings.TrimSpace(tone))
	if len(tone) < 2 {
		return 0, false
	}
	semitones := map[byte]int{'C': 0, 'D': 2, 'E': 4, 'F': 5, 'G': 7, 'A': 9, 'B': 11}
	semitone, ok := semitones[tone[0]]
	if !ok {
		return 0, false
	}
	index := 1
	if tone[index] == '#' || tone[index] == 'B' {
		if tone[index] == '#' {
			semitone++
		} else {
			semitone--
		}
		index++
	}
	octave, err := strconv.Atoi(tone[index:])
	if err != nil {
		return 0, false
	}
	return (octave+1)*12 + semitone, true
}

func stripComment(value string) string {
	quoted := byte(0)
	for index := 0; index < len(value); index++ {
		switch value[index] {
		case '\'', '"':
			if quoted == 0 {
				quoted = value[index]
			} else if quoted == value[index] {
				quoted = 0
			}
		case '#':
			if quoted == 0 && (index == 0 || value[index-1] == ' ' || value[index-1] == '\t') {
				return value[:index]
			}
		}
	}
	return value
}
