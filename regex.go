package main

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

	// Build the exclusion character class once, since it does not vary by position.
	ex := make([]string, 0, len(globalExclude))
	for r := range globalExclude {
		ex = append(ex, string(r))
	}
	sort.Strings(ex)
	class := strings.Join(ex, "")

	parts := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		if greens[i] != 0 {
			parts = append(parts, string(greens[i]))
			continue
		}

		if class == "" {
			parts = append(parts, "[a-z]")
			continue
		}

		parts = append(parts, "[^"+class+"]")
	}

	return "^" + strings.Join(parts, "") + "$", nil
}
