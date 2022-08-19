package main

import (
	"flag"
	"fmt"
	"go/scanner"
	"os"

	"github.com/ian-h-chamberlain/gocommas/fixer"
)

func main() {
	// TODO use proper CLI args instead, and some e2e testing or something
	// use testdata in unit tests instead of main
	writeFile := flag.Bool("w", false, "write to input file instead of stdout")
	helpRequested := flag.Bool("help", false, "display this help text")

	flag.Parse()

	if *helpRequested {
		flag.Usage()
		os.Exit(0)
	}
	if flag.NArg() != 1 {
		fmt.Fprintln(
			os.Stderr,
			"unexpected number of arguments. Only one input file is supported",
		)
		flag.Usage()
		os.Exit(1)
	}

	fname := flag.Arg(0)

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

	if *writeFile {
		stat, err := os.Stat(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		os.WriteFile(fname, fixedSrc, stat.Mode().Perm())
	} else {
		fmt.Print(string(fixedSrc))
	}
}
