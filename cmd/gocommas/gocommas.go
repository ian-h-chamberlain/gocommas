package main

import (
	"flag"
	"fmt"
	"io"
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

	if flag.NArg() > 1 {
		flag.Usage()
		exitWithErrf(
			"unexpected number of arguments (%d). Only one input file is supported.",
			flag.NArg(),
		)
	}

	var (
		src   []byte
		fname string
		err   error
	)

	if flag.NArg() < 1 || flag.Arg(0) == "-" {
		if *writeFile {
			exitWithErrf("cannot use -w with standard input")
		}
		fname = "-"
		src, err = io.ReadAll(os.Stdin)
	} else {
		fname = flag.Arg(0)
		src, err = os.ReadFile(fname)
	}
	if err != nil {
		exitWithErrf("%v", err)
	}

	commasToAdd, err := fixer.FindMissingCommas(fname, src)
	if err != nil {
		exitWithErrf("%v", err)
		os.Exit(1)
	}

	for _, pos := range commasToAdd {
		fmt.Fprintf(os.Stderr, "Fixing missing comma: %s\n", pos)
	}

	fixedSrc := fixer.AddMissingCommas(src, commasToAdd)

	if *writeFile {
		stat, err := os.Stat(fname)
		if err != nil {
			exitWithErrf("%v", err)
		}

		err = os.WriteFile(fname, fixedSrc, stat.Mode().Perm())
		if err != nil {
			exitWithErrf("%v", err)
		}
	} else {
		fmt.Print(string(fixedSrc))
	}
}

func exitWithErrf(template string, a ...any) {
	fmt.Fprintf(os.Stderr, "error: "+template+"\n", a...)
	os.Exit(1)
}
