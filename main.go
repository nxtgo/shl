package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
)

type Rule struct {
	Pattern string   `json:"pattern"`
	Color   string   `json:"color"`
	Nested  []string `json:"nested,omitempty"`
	Capture int      `json:"capture,omitempty"`
}

var colors = map[string]string{
	"black": "30", "red": "31", "green": "32", "yellow": "33",
	"blue": "34", "magenta": "35", "cyan": "36", "white": "37",
	"brightBlack": "90", "brightRed": "91", "brightGreen": "92",
	"brightYellow": "93", "brightBlue": "94", "brightMagenta": "95",
	"brightCyan": "96", "brightWhite": "97",
}

func colorCode(s string) string {
	if code, ok := colors[s]; ok {
		return code
	}
	if len(s) == 7 && s[0] == '#' {
		var r, g, b int
		fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
		return fmt.Sprintf("38;2;%d;%d;%d", r, g, b)
	}
	var r, g, b int
	if n, _ := fmt.Sscanf(s, "rgb(%d,%d,%d)", &r, &g, &b); n == 3 {
		return fmt.Sprintf("38;2;%d;%d;%d", r, g, b)
	}
	return s
}

type match struct {
	start, end int
	rule       int
	nested     []match
	capStart   int
	capEnd     int
}

func findMatches(text string, res []*regexp.Regexp, rules []Rule, ruleIndices []int) []match {
	var matches []match

	for _, idx := range ruleIndices {
		re := res[idx]
		if rules[idx].Capture > 0 {
			allMatches := re.FindAllStringSubmatchIndex(text, -1)
			for _, loc := range allMatches {
				capIdx := rules[idx].Capture * 2
				if capIdx+1 < len(loc) && loc[capIdx] >= 0 && loc[capIdx+1] >= 0 {
					matches = append(matches, match{
						start:    loc[0],
						end:      loc[1],
						rule:     idx,
						capStart: loc[capIdx],
						capEnd:   loc[capIdx+1],
					})
				}
			}
		} else {
			for _, loc := range re.FindAllStringIndex(text, -1) {
				if loc[0] < loc[1] {
					matches = append(matches, match{
						start:    loc[0],
						end:      loc[1],
						rule:     idx,
						capStart: loc[0],
						capEnd:   loc[1],
					})
				}
			}
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].start != matches[j].start {
			return matches[i].start < matches[j].start
		}
		return matches[i].rule < matches[j].rule
	})

	used := make([]bool, len(text))
	var kept []match
	for _, m := range matches {
		ok := true
		for j := m.start; j < m.end; j++ {
			if used[j] {
				ok = false
				break
			}
		}
		if ok {
			for j := m.start; j < m.end; j++ {
				used[j] = true
			}

			if len(rules[m.rule].Nested) > 0 {
				nestedIndices := make([]int, 0)
				for _, pattern := range rules[m.rule].Nested {
					for i, r := range rules {
						if r.Pattern == pattern {
							nestedIndices = append(nestedIndices, i)
							break
						}
					}
				}
				if len(nestedIndices) > 0 {
					matchedText := text[m.capStart:m.capEnd]
					m.nested = findMatches(matchedText, res, rules, nestedIndices)
				}
			}

			kept = append(kept, m)
		}
	}

	return kept
}

func printMatches(line string, matches []match, rules []Rule, parentColor string) {
	pos := 0
	for _, m := range matches {
		if m.start > pos {
			if parentColor != "" {
				fmt.Printf("\033[%sm", parentColor)
			}
			fmt.Print(line[pos:m.start])
			if parentColor != "" {
				fmt.Print("\033[0m")
			}
		}

		if m.capStart > m.start {
			if parentColor != "" {
				fmt.Printf("\033[%sm", parentColor)
			}
			fmt.Print(line[m.start:m.capStart])
			if parentColor != "" {
				fmt.Print("\033[0m")
			}
		}

		code := colorCode(rules[m.rule].Color)
		fmt.Printf("\033[%sm", code)

		if len(m.nested) == 0 {
			fmt.Print(line[m.capStart:m.capEnd])
		} else {
			printMatches(line[m.capStart:m.capEnd], m.nested, rules, code)
		}

		fmt.Print("\033[0m")

		if m.capEnd < m.end {
			if parentColor != "" {
				fmt.Printf("\033[%sm", parentColor)
			}
			fmt.Print(line[m.capEnd:m.end])
			if parentColor != "" {
				fmt.Print("\033[0m")
			}
		}

		pos = m.end
	}
	if pos < len(line) {
		if parentColor != "" {
			fmt.Printf("\033[%sm", parentColor)
		}
		fmt.Print(line[pos:])
		if parentColor != "" {
			fmt.Print("\033[0m")
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: program <rules.json> [input_file]")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var rules []Rule
	if json.Unmarshal(data, &rules) != nil || len(rules) == 0 {
		fmt.Fprintln(os.Stderr, "invalid rules file")
		os.Exit(1)
	}

	res := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		res[i] = regexp.MustCompile(r.Pattern)
	}

	in := os.Stdin
	if len(os.Args) > 2 {
		in, err = os.Open(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer in.Close()
	}

	s := bufio.NewScanner(in)
	lineNum := 1
	for s.Scan() {
		line := s.Text()
		fmt.Printf("\033[90m%4d\033[0m ", lineNum)
		lineNum++

		allIndices := make([]int, len(rules))
		for i := range rules {
			allIndices[i] = i
		}
		matches := findMatches(line, res, rules, allIndices)

		printMatches(line, matches, rules, "")
		fmt.Println()
	}
}
