package main

import (
	"flag"
	"fmt"
	"os"

	"RowXml.com/internal/logging"
	"RowXml.com/internal/parser"
	"RowXml.com/internal/validator"
	"RowXml.com/internal/xmlrender"
)

func main() {
	var fromStdin bool
	flag.BoolVar(&fromStdin, "stdin", false, "read input from stdin (default: first arg)")
	flag.Parse()

	logger := logging.NewStdLogger()
	var raw string
	var err error

	if fromStdin {
		// read all from stdin
		b, readErr := os.ReadFile("/dev/stdin")
		if readErr != nil {
			logger.Error("failed to read stdin", "err", readErr)
			fmt.Fprintln(os.Stderr, readErr)
			os.Exit(1)
		}
		raw = string(b)
	} else {
		if flag.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "Usage: rowxml [--stdin] \"INPUT_STRING\"")
			os.Exit(2)
		}
		raw = flag.Arg(0)
	}

	people, err := parser.Parse(raw)
	if err != nil {
		logger.Error("parse failed", "err", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := validator.ValidatePeople(&people); err != nil {
		logger.Error("validation failed", "err", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	out, err := xmlrender.RenderXML(people)
	if err != nil {
		logger.Error("render failed", "err", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Print(out)
}

