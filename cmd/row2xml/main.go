package main

import (
    "flag"
	"fmt"

)

func main() {
    input := flag.String("input", "", "input row file")
	flag.Parse()

	if *input == "" {
		fmt.Println("input file is required")
		return
	}

	fmt.Println("input:", *input)


}
