package main

import (
	"flag"
	"fmt"
	"os"

	"RowXml.com/internal/logging"
	"RowXml.com/internal/cli"
)


func main() {
	flag.Parse()


	logger := logging.NewStdLogger()
	runner := cli.NewRunner(logger)

	// The CLI takes one argument: the raw input string.
	if flag.NArg() < 1 {
		logger.Error("missing input string")
		fmt.Fprintln(os.Stderr, "Usage: parse \"INPUT_STRING\"")
		os.Exit(2)
	}

	input := flag.Arg(0)

	xml, err := runner.Run(input)
	if err != nil {
		logger.Error("failed to run parser", "err", err)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println(xml)
}
