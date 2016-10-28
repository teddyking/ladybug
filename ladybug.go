package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewParser(nil, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
	}
}
