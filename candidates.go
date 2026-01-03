package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var ErrDictOpen = errors.New("failed to open dictionary")

type RatedGuess struct {
	Guess       string
	ERemaining  float64
	EEliminated float64
}

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

// wordleFeedbackPattern returns the Wordle feedback pattern you would observe if you guessed
// `guess` when the true answer is `answer`, including duplicate-letter rules.
func wordleFeedbackPattern(guess string, answer string) string {
	g := []rune(guess)
	a := []rune(answer)

	pattern := make([]byte, 5)
	remaining := map[rune]int{}

	// First pass: mark greens and count remaining (unmatched) letters in the answer.
	for i := 0; i < 5; i++ {
		if g[i] == a[i] {
			pattern[i] = '2'
			continue
		}
		remaining[a[i]]++
	}

	// Second pass: for non-green positions, mark yellow only if there is remaining supply
	// of that letter in the answer; otherwise mark gray.
	for i := 0; i < 5; i++ {
		if pattern[i] == '2' {
			continue
		}
		if remaining[g[i]] > 0 {
			pattern[i] = '1'
			remaining[g[i]]--
			continue
		}
		pattern[i] = '0'
	}

	return string(pattern)
}

func SortCandidates(candidates []string) []RatedGuess {
	if len(candidates) == 0 {
		return nil
	}

	N := len(candidates)
	rated := make([]RatedGuess, 0, N)

	for _, g := range candidates {
		bucketCounts := map[string]int{}
		for _, a := range candidates {
			p := wordleFeedbackPattern(g, a)
			bucketCounts[p]++
		}

		sumSquares := 0
		for _, c := range bucketCounts {
			sumSquares += c * c
		}

		eRemaining := float64(sumSquares) / float64(N)
		eEliminated := float64(N) - eRemaining
		rated = append(rated, RatedGuess{Guess: g, ERemaining: eRemaining, EEliminated: eEliminated})
	}

	sort.Slice(rated, func(i, j int) bool {
		if rated[i].ERemaining == rated[j].ERemaining {
			return rated[i].Guess < rated[j].Guess
		}
		return rated[i].ERemaining < rated[j].ERemaining
	})

	return rated
}
