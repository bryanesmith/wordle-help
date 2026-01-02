package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

func run(argv []string) int {
	fs := flag.NewFlagSet("wordle-help", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var guesses []string
	fs.Func("g", "A guess, e.g. slate or [p][r](o)ps", func(v string) error {
		guesses = append(guesses, v)
		return nil
	})
	fs.Func("guess", "A guess, e.g. slate or [p][r](o)ps", func(v string) error {
		guesses = append(guesses, v)
		return nil
	})

	help := fs.Bool("h", false, "Show help")
	fs.BoolVar(help, "help", false, "Show help")

	if err := fs.Parse(argv); err != nil {
		return 2
	}
	if *help {
		fs.Usage()
		return 0
	}

	parsed, err := ValidateGuesses(guesses)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	pattern, err := BuildRegex(parsed)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	cands, err := Candidates(re, "/usr/share/dict/words")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	for _, w := range cands {
		fmt.Fprintln(os.Stdout, w)
	}

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
