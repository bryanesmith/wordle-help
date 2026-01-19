package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bryanesmith/wordle-help/utils"
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

	rated, err := utils.NextGuessRecommendations(guesses, utils.DefaultDictionaryPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	for _, r := range rated {
		fmt.Fprintf(os.Stdout, "%-10s (E_remaining = %.2f, E_eliminated = %.2f)\n", r.Guess, r.ERemaining, r.EEliminated)
	}

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
