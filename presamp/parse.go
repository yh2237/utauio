// presampパッケージは復号済みpresamp.iniの発音設定を解析する。
package presamp

import (
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	Vowels       map[string]string
	Consonants   map[string]string
	Replacements map[string]string
	Endings      []string
}

func Parse(text string) Config {
	result := Config{Vowels: map[string]string{}, Consonants: map[string]string{}, Replacements: map[string]string{}}
	section := ""
	endingTypes := map[int][]string{}
	endingFlag := 0
	for {
		raw, rest, more := strings.Cut(text, "\n")
		text = rest
		line := strings.TrimSpace(strings.TrimSuffix(raw, "\r"))
		if line != "" && !strings.HasPrefix(line, ";") && !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				section = strings.ToUpper(strings.TrimSpace(line[1 : len(line)-1]))
			} else {
				switch {
				case section == "VOWEL":
					parseClass(line, 2, result.Vowels)
				case section == "CONSONANT":
					parseClass(line, 1, result.Consonants)
				case section == "REPLACE":
					key, value, ok := strings.Cut(line, "=")
					if ok && strings.TrimSpace(key) != "" && strings.TrimSpace(value) != "" {
						result.Replacements[strings.TrimSpace(key)] = strings.TrimSpace(value)
					}
				case strings.HasPrefix(section, "ENDTYPE"):
					index := endingType(section)
					endingTypes[index] = append(endingTypes[index], line)
				case section == "ENDFLAG":
					endingFlag, _ = strconv.Atoi(line)
				}
			}
		}
		if !more {
			break
		}
	}
	indices := make([]int, 0, len(endingTypes))
	for index := range endingTypes {
		indices = append(indices, index)
	}
	sort.Ints(indices)
	for _, index := range indices {
		if endingFlag == 0 || endingFlag&(1<<(index-1)) != 0 {
			result.Endings = append(result.Endings, endingTypes[index]...)
		}
	}
	return result
}

func endingType(section string) int {
	if index, err := strconv.Atoi(strings.TrimPrefix(section, "ENDTYPE")); err == nil && index > 0 {
		return index
	}
	return 1
}

func parseClass(line string, aliasesIndex int, destination map[string]string) {
	class, tail, ok := strings.Cut(line, "=")
	if !ok {
		return
	}
	class = strings.TrimSpace(class)
	if class == "" {
		return
	}
	for index := 1; index < aliasesIndex; index++ {
		_, tail, ok = strings.Cut(tail, "=")
		if !ok {
			return
		}
	}
	aliases, _, _ := strings.Cut(tail, "=")
	for {
		alias, rest, more := strings.Cut(aliases, ",")
		if alias = strings.TrimSpace(alias); alias != "" {
			destination[alias] = class
		}
		if !more {
			break
		}
		aliases = rest
	}
}
