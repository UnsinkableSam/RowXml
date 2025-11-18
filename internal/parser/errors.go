package parser

import "fmt"

// ParseError is a structured parse error with a line number and message.
type ParseError struct {
	Line int
	Msg  string
	Raw  string
}

func (e ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d: %s", e.Line, e.Msg)
	}
	return e.Msg
}

func newParseError(line int, msg string) ParseError {
	return ParseError{Line: line, Msg: msg}
}

