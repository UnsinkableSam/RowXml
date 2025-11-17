package main

import (
	"flag"
	"fmt"
	"encoding/xml"

	"RowXml.com/internal/parser"

)

func main() {
	input := flag.String("input", "", "String to convert to XML")

	flag.Parse()

	if *input == "" {
		fmt.Println("Usage: mycli -msg \"Hello\"")
		return
	}

	people, err := parser.Parse(*input)
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	output, err := xml.MarshalIndent(people, "", "  ")
	if err != nil {
		fmt.Println("xml error:", err)
		return
	}

	fmt.Println(string(output))
}
