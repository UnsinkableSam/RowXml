package parser

import (
	"strings"
)

type token struct {
	typ   string   // "P","T","A","F"
	parts []string // fields after split and trimmed
	raw   string   // original trimmed raw line
}

func lexLines(input string) []token {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(input, "\n")
	var toks []token
	for _, l := range lines {
		line := strings.TrimSpace(l)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		toks = append(toks, token{
			typ:   parts[0],
			parts: parts,
			raw:   line,
		})
	}
	return toks
}

