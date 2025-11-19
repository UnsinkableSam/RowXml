package main

import (
	"flag"
	"fmt"
	"os"
	"io"

	"RowXml.com/internal/logging"
	"RowXml.com/internal/parser"
	"RowXml.com/internal/validator"
	"RowXml.com/internal/xmlrender"
)

func main() {
	fileInput := flag.String("file", "", "Path to input file")
	stringInput := flag.String("string", "", "Raw input string")
	flag.Parse()

	var data string
	var err error

	logger := logging.NewStdLogger()

	switch {
	case *fileInput != "":
		bytes, err := os.ReadFile(*fileInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		data = string(bytes)

	case *stringInput != "":
		data = *stringInput

	default:
		// If nothing provided, read from stdin
		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		data = string(stdinBytes)
	}

	people, err := parser.Parse(data)
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
