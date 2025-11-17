package logging

import (
	"log"
)

// Logger is a tiny interface so we can swap implementations in tests.
type Logger interface {
	Info(msg string, kv ...interface{})
	Error(msg string, kv ...interface{})
}

// NewStdLogger returns a minimal stdlib-backed logger implementation suitable
// for CLI tools and unit tests.
func NewStdLogger() Logger { return stdLogger{} }

type stdLogger struct{}

func (stdLogger) Info(msg string, kv ...interface{}) {
	log.Println("INFO:", msg, kv)
}

func (stdLogger) Error(msg string, kv ...interface{}) {
	log.Println("ERROR:", msg, kv)
}
