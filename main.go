package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

type multiFlag []string

func (m *multiFlag) String() string {
	return fmt.Sprint([]string(*m))
}

func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func run(argv []string) int {
	fs := flag.NewFlagSet("wordle-help", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var guesses multiFlag
	fs.Var(&guesses, "g", "A guess, e.g. slate or [p][r](o)ps")
	fs.Var(&guesses, "guess", "A guess, e.g. slate or [p][r](o)ps")

	help := fs.Bool("h", false, "Show help")
	fs.BoolVar(help, "help", false, "Show help")

	if err := fs.Parse(argv); err != nil {
		return 2
	}
	if *help {
		fs.Usage()
		return 0
	}

	parsed, err := ValidateGuesses([]string(guesses))
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
