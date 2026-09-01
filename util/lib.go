package util

import (
	"flag"
	"os"
)

func If[T any](cond bool, vtrue, vfalse T) T {
    if cond {
        return vtrue
    }
    return vfalse
}

type Options struct {
	Verbose *bool
	Repl bool
	Files []string
}

func Args() Options {
	flag.CommandLine.Name()
	opt := Options{
		Verbose: flag.Bool("verbose", false, "Enable verbose logging"),
		Repl: false,
		Files: []string{},
	}

	flag.BoolVar(&opt.Repl, "repl", false, "Run as in REPL mode")

	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		flag.Usage();
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		if !opt.Repl {
			print("No files provided, running in REPL")
		}
		opt.Repl = true
	} else {
		opt.Files = flag.Args()
	}

	return opt
}