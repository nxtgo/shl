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
	Pattern string
	Color   string
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "wrong usage CUCK")
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
	for s.Scan() {
		line := s.Text()
		used := make([]bool, len(line))

		type match struct {
			start, end, rule int
		}
		var matches []match
		for i, re := range res {
			for _, loc := range re.FindAllStringIndex(line, -1) {
				if loc[0] < loc[1] {
					matches = append(matches, match{loc[0], loc[1], i})
				}
			}
		}

		sort.Slice(matches, func(i, j int) bool {
			if matches[i].start != matches[j].start {
				return matches[i].start < matches[j].start
			}
			return matches[i].rule < matches[j].rule
		})

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
				kept = append(kept, m)
			}
		}

		pos := 0
		for _, m := range kept {
			if m.start > pos {
				fmt.Print(line[pos:m.start])
			}
			code := colorCode(rules[m.rule].Color)
			fmt.Printf("\033[%sm%s\033[0m", code, line[m.start:m.end])
			pos = m.end
		}
		if pos < len(line) {
			fmt.Print(line[pos:])
		}
		fmt.Println()
	}
}
