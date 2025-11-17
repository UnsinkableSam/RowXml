package cli

import (
	"fmt"

	"RowXml.com/internal/parser"
	"RowXml.com/internal/logging"
)

// Runner orchestrates parsing, validation and rendering. Kept minimal so it's
// trivial to test with unit tests / table driven tests in future.
type Runner struct {
	logger logging.Logger
}

func NewRunner(l logging.Logger) *Runner { return &Runner{logger: l} }

// Run accepts the raw input string, returns pretty XML or an error. The
// function purposefully returns actionable errors (wrapped) for callers.
func (r *Runner) Run(input string) (string, error) {
	r.logger.Info("starting run")
	people, err := parser.Parse(input)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	xml, err := parser.ToXML(people)
	if err != nil {
		return "", fmt.Errorf("render error: %w", err)
	}

	r.logger.Info("run completed")
	return xml, nil
}
