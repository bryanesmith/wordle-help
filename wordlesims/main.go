package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/bryanesmith/wordle-help/utils"
)

var ErrWordleHelpNoOutput = errors.New("wordle_help produced no output")
var ErrWordleHelpTooManyGuesses = errors.New("no more than six -g/--guess values may be provided")

type simResult struct {
	StartingWord        string
	Answer              string
	TotalGuesses        int
	TotalGuessesDisplay string
	Succeeded           bool
}

func isValidDictWord(word string) bool {
	if len([]rune(word)) != 5 {
		return false
	}
	for _, r := range word {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func pickRandomAnswers(dictPath string, n int) ([]string, error) {
	f, err := os.Open(dictPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	words := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		w := strings.TrimSpace(scanner.Text())
		if w == "" {
			continue
		}
		w = strings.ToLower(w)
		if !isValidDictWord(w) {
			continue
		}
		words = append(words, w)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(words) == 0 {
		return nil, fmt.Errorf("no valid words found in %s", dictPath)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

	if n > len(words) {
		n = len(words)
	}

	return words[:n], nil
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
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if strings.Contains(msg, ErrWordleHelpTooManyGuesses.Error()) {
			return "", ErrWordleHelpTooManyGuesses
		}
		return "", err
	}

	line := strings.TrimSpace(string(out))
	if line == "" {
		return "", ErrWordleHelpNoOutput
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

	if *starting == "" {
		fmt.Fprintln(os.Stderr, "exactly one -s/--starting must be provided")
		return 1
	}

	if len(answers) == 0 {
		picked, err := pickRandomAnswers("/usr/share/dict/words", 50)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		answers = picked
	}

	results := make([]simResult, 0, len(answers))
	totalGuesses := 0
	successfulSims := 0

	for _, ans := range answers {
		answer := strings.ToLower(ans)
		currentGuess := strings.ToLower(*starting)
		guessEncodings := []string{}
		guessCount := 1
		succeeded := true
		failedDisplay := ""

		for currentGuess != answer {
			result := utils.CheckGuess(currentGuess, answer)
			guessEncodings = append(guessEncodings, result)

			nextGuess, err := runWordleHelp(guessEncodings)
			if err != nil {
				if errors.Is(err, ErrWordleHelpNoOutput) {
					fmt.Fprintf(os.Stderr, "Skipping simulation for answer %s: %s\n", answer, err.Error())
					succeeded = false
					failedDisplay = "-"
					break
				}
				if errors.Is(err, ErrWordleHelpTooManyGuesses) {
					fmt.Fprintf(os.Stderr, "Skipping simulation for answer %s: %s\n", answer, err.Error())
					succeeded = false
					failedDisplay = "fail"
					break
				}
				fmt.Fprintln(os.Stderr, err.Error())
				return 1
			}

			fmt.Fprintf(os.Stdout, "Guess: %s Result: %s Next Guess: %s (Answer: %s)\n", currentGuess, result, nextGuess, answer)
			guessCount++
			currentGuess = nextGuess
		}

		if !succeeded {
			results = append(results, simResult{StartingWord: strings.ToLower(*starting), Answer: answer, TotalGuesses: 0, TotalGuessesDisplay: failedDisplay, Succeeded: false})
			continue
		}

		results = append(results, simResult{StartingWord: strings.ToLower(*starting), Answer: answer, TotalGuesses: guessCount, TotalGuessesDisplay: fmt.Sprintf("%d", guessCount), Succeeded: true})
		totalGuesses += guessCount
		successfulSims++
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "| Starting Word | Answer | Total Guesses |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	for _, r := range results {
		if !r.Succeeded {
			fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n", r.StartingWord, r.Answer, r.TotalGuessesDisplay)
			continue
		}
		fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n", r.StartingWord, r.Answer, r.TotalGuessesDisplay)
	}

	if successfulSims == 0 {
		fmt.Fprintln(os.Stdout, "\nAverage guesses: -")
		return 0
	}

	avg := float64(totalGuesses) / float64(successfulSims)
	fmt.Fprintf(os.Stdout, "\nAverage guesses: %.2f\n", avg)

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
