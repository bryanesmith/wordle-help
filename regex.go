package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func BuildRegex(guesses []ParsedGuess) (string, error) {
	greens := [5]rune{}
	globalExclude := map[rune]struct{}{}

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
				// Yellow letters are known to exist somewhere in the word, so we do not exclude them.
			default:
				globalExclude[letter] = struct{}{}
			}
		}
	}

	parts := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		if greens[i] != 0 {
			parts = append(parts, regexp.QuoteMeta(string(greens[i])))
			continue
		}

		excluded := map[rune]struct{}{}
		for r := range globalExclude {
			excluded[r] = struct{}{}
		}

		ex := make([]string, 0, len(excluded))
		for r := range excluded {
			ex = append(ex, string(r))
		}
		sort.Strings(ex)

		class := strings.Join(ex, "")
		class = regexp.QuoteMeta(class)
		parts = append(parts, "[^"+class+"]")
	}

	return "^" + strings.Join(parts, "") + "$", nil
}
