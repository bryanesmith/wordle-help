package utils

import (
	"fmt"
	"sort"
	"strings"
)

// BuildRegex converts the validated guesses into a single regular expression string.
//
// The returned pattern is anchored (starts with '^' and ends with '$') so that it matches
// only full 5-letter words.
func BuildRegex(guesses []ParsedGuess) (string, error) {
	greens := [5]rune{}
	present := map[rune]struct{}{}
	for _, g := range guesses {
		for r := range g.MarkedSet {
			present[r] = struct{}{}
		}
	}

	globalAbsent := map[rune]struct{}{}
	posExclude := [5]map[rune]struct{}{}
	for i := 0; i < 5; i++ {
		posExclude[i] = map[rune]struct{}{}
	}

	for _, g := range guesses {
		for i := 0; i < 5; i++ {
			s := g.States[i]
			letter := g.Letters[i]
			switch s {
			case LetterStateGreen:
				if greens[i] != 0 && greens[i] != letter {
					return "", fmt.Errorf("%w: conflicting green letter at position %d", ErrContradiction, i+1)
				}
				greens[i] = letter
			case LetterStateYellow:
				// Yellow letters exist somewhere in the word, but cannot be in the position they were marked.
				posExclude[i][letter] = struct{}{}
			default:
				// A plain (unmarked) letter is often "gray", but if we also have evidence that the
				// letter exists elsewhere (yellow/green), then this tile only rules out the current
				// position (duplicate-letter case).
				if _, ok := present[letter]; ok {
					posExclude[i][letter] = struct{}{}
				} else {
					globalAbsent[letter] = struct{}{}
				}
			}
		}
	}

	parts := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		if greens[i] != 0 {
			parts = append(parts, string(greens[i]))
			continue
		}

		exclude := map[rune]struct{}{}
		for r := range globalAbsent {
			exclude[r] = struct{}{}
		}
		for r := range posExclude[i] {
			exclude[r] = struct{}{}
		}

		if len(exclude) == 0 {
			parts = append(parts, "[a-z]")
			continue
		}

		ex := make([]string, 0, len(exclude))
		for r := range exclude {
			ex = append(ex, string(r))
		}
		sort.Strings(ex)
		class := strings.Join(ex, "")
		parts = append(parts, "[^"+class+"]")
	}

	return "^" + strings.Join(parts, "") + "$", nil
}
