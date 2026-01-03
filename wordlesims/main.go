package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bryanesmith/wordle-help/utils"
)

type simResult struct {
	StartingWord string
	Answer       string
	TotalGuesses int
}

func wordleHelpPath() string {
	return filepath.Join("bin", "wordle_help")
}

func runWordleHelp(guesses []string) (string, error) {
	args := make([]string, 0, len(guesses)*2)
	for _, g := range guesses {
		args = append(args, "-g", g)
	}

	cmdStr := append([]string{wordleHelpPath()}, args...)
	fmt.Fprintln(os.Stdout, strings.Join(cmdStr, " "))

	cmd := exec.Command(wordleHelpPath(), args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	line := strings.TrimSpace(string(out))
	if line == "" {
		return "", fmt.Errorf("wordle_help produced no output")
	}

	firstLineEnd := strings.IndexByte(line, '\n')
	if firstLineEnd >= 0 {
		line = line[:firstLineEnd]
	}

	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "", fmt.Errorf("failed to parse wordle_help output: %q", line)
	}

	return fields[0], nil
}

func run(argv []string) int {
	fs := flag.NewFlagSet("wordle-sims", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var answers []string
	fs.Func("a", "An answer to simulate (repeatable)", func(v string) error {
		answers = append(answers, v)
		return nil
	})
	fs.Func("answer", "An answer to simulate (repeatable)", func(v string) error {
		answers = append(answers, v)
		return nil
	})

	starting := fs.String("s", "", "Starting word")
	fs.StringVar(starting, "starting", "", "Starting word")

	help := fs.Bool("h", false, "Show help")
	fs.BoolVar(help, "help", false, "Show help")

	if err := fs.Parse(argv); err != nil {
		return 2
	}
	if *help {
		fs.Usage()
		return 0
	}

	if len(answers) == 0 {
		fmt.Fprintln(os.Stderr, "at least one -a/--answer must be provided")
		return 1
	}
	if *starting == "" {
		fmt.Fprintln(os.Stderr, "exactly one -s/--starting must be provided")
		return 1
	}

	results := make([]simResult, 0, len(answers))
	totalGuesses := 0

	for _, ans := range answers {
		answer := strings.ToLower(ans)
		currentGuess := strings.ToLower(*starting)
		guessEncodings := []string{}
		guessCount := 1

		for currentGuess != answer {
			result := utils.CheckGuess(currentGuess, answer)
			guessEncodings = append(guessEncodings, result)

			nextGuess, err := runWordleHelp(guessEncodings)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return 1
			}

			fmt.Fprintf(os.Stdout, "Guess: %s Result: %s Next Guess: %s (Answer: %s)\n", currentGuess, result, nextGuess, answer)
			guessCount++
			currentGuess = nextGuess
		}

		results = append(results, simResult{StartingWord: strings.ToLower(*starting), Answer: answer, TotalGuesses: guessCount})
		totalGuesses += guessCount
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "| Starting Word | Answer | Total Guesses |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	for _, r := range results {
		fmt.Fprintf(os.Stdout, "| %s | %s | %d |\n", r.StartingWord, r.Answer, r.TotalGuesses)
	}

	avg := float64(totalGuesses) / float64(len(results))
	fmt.Fprintf(os.Stdout, "\nAverage guesses: %.2f\n", avg)

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
