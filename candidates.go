package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var ErrDictOpen = errors.New("failed to open dictionary")

func buildYellowConstraints(guesses []ParsedGuess) (map[rune]struct{}, [5]map[rune]struct{}) {
	requiredYellow := map[rune]struct{}{}
	forbiddenYellowByPos := [5]map[rune]struct{}{}
	for i := 0; i < 5; i++ {
		forbiddenYellowByPos[i] = map[rune]struct{}{}
	}

	for _, g := range guesses {
		for i := 0; i < 5; i++ {
			for r := range g.YellowByPos[i] {
				requiredYellow[r] = struct{}{}
				forbiddenYellowByPos[i][r] = struct{}{}
			}
		}
	}

	return requiredYellow, forbiddenYellowByPos
}

func isValidGuessWord(word string) bool {
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

// Candidates reads the given dictionary file and returns all unique, lowercase, five-letter
// words that match the provided regular expression.
func Candidates(re *regexp.Regexp, guesses []ParsedGuess, dictPath string) ([]string, error) {
	f, err := os.Open(dictPath)
	if err != nil {
		return nil, ErrDictOpen
	}
	defer f.Close()

	out := []string{}
	seen := map[string]struct{}{}
	requiredYellow, forbiddenYellowByPos := buildYellowConstraints(guesses)

	scanner := bufio.NewScanner(f)
scan:
	for scanner.Scan() {
		// Ignore blank lines.
		w := strings.TrimSpace(scanner.Text())
		if w == "" {
			continue scan
		}

		// Wordle is case-insensitive; normalize words.
		w = strings.ToLower(w)

		// Only consider valid 5-letter words comprised entirely of letters.
		if !isValidGuessWord(w) {
			continue scan
		}

		// Fast-path: regex pre-filtering for greens/grays.
		if !re.MatchString(w) {
			continue scan
		}

		// Yellow letters cannot be in the position they were marked yellow.
		wordRunes := []rune(w)
		for i := 0; i < 5; i++ {
			if _, ok := forbiddenYellowByPos[i][wordRunes[i]]; ok {
				continue scan
			}
		}

		// Yellow letters must still exist somewhere in the candidate.
		for r := range requiredYellow {
			if !strings.ContainsRune(w, r) {
				continue scan
			}
		}

		if _, ok := seen[w]; ok {
			continue scan
		}
		seen[w] = struct{}{}
		out = append(out, w)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
