// prefixmapパッケージは復号済みprefix.mapの接頭辞・接尾辞を解析する。
package prefixmap

import (
	"fmt"
	"strings"
)

type Entry struct {
	Tone   string
	Prefix string
	Suffix string
	Line   int
}

type Diagnostic struct {
	Line    int
	Message string
}

// Scanは正常行を記載順に渡す。接辞の空白と空欄を保持し、不正行の診断を返す。
func Scan(text string, visit func(Entry)) ([]Diagnostic, error) {
	if visit == nil {
		return nil, fmt.Errorf("prefix.map visitor is required")
	}
	var diagnostics []Diagnostic
	lineNumber := 0
	for {
		line, rest, more := strings.Cut(text, "\n")
		text = rest
		lineNumber++
		line = strings.TrimSuffix(line, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, ";") {
			tone, tail, ok := strings.Cut(line, "\t")
			if !ok {
				diagnostics = append(diagnostics, Diagnostic{lineNumber, "invalid prefix.map line"})
			} else {
				prefix, suffix, _ := strings.Cut(tail, "\t")
				suffix, _, _ = strings.Cut(suffix, "\t")
				visit(Entry{Tone: strings.ToUpper(strings.TrimSpace(tone)), Prefix: prefix, Suffix: suffix, Line: lineNumber})
			}
		}
		if !more {
			break
		}
	}
	return diagnostics, nil
}
