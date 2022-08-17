package main

import (
	"fmt"
	"go/scanner"
	"os"

	"github.com/ian-h-chamberlain/gocommas/fixer"
)

func main() {
	// TODO use CLI args instead, and some e2e testing or something
	// use testdata in unit tests instead of main

	fname := "/Users/ichamberlain/Documents/gocommas/fixer/testdata/missing_trailing.go"

	src, err := os.ReadFile(fname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	commasToAdd, err := fixer.FindMissingCommas(fname, src)
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		os.Exit(1)
	}

	for _, pos := range commasToAdd {
		fmt.Fprintf(os.Stderr, "Missing comma: %s\n", pos)
	}

	fixedSrc := fixer.AddMissingCommas(src, commasToAdd)
	fmt.Println(string(fixedSrc))
}
