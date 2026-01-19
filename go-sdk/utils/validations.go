package utils

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrNoGuesses        = errors.New("at least one -g/--guess must be provided")
	ErrTooManyGuesses   = errors.New("no more than six -g/--guess values may be provided")
	ErrGuessWrongLength = errors.New("each guess must be exactly five letters (ignoring (), [])")
	ErrGuessInvalidChar = errors.New("guesses must consist only of letters (ignoring (), [])")
	ErrInvalidMarkers   = errors.New("invalid guess marker syntax")
	ErrContradiction    = errors.New("one guess contradicts another")
	ErrTooManyMarked    = errors.New("there cannot be more than five different marked letters")
)

type LetterState int

const (
	LetterStateNone LetterState = iota
	LetterStateYellow
	LetterStateGreen
)

type ParsedGuess struct {
	Raw         string
	Letters     [5]rune
	States      [5]LetterState
	MarkedSet   map[rune]struct{}
	YellowByPos [5]map[rune]struct{}
	GreenByPos  [5]rune
}

func ParseGuess(raw string) (ParsedGuess, error) {
	var out ParsedGuess
	out.Raw = raw
	out.MarkedSet = map[rune]struct{}{}
	for i := 0; i < 5; i++ {
		out.YellowByPos[i] = map[rune]struct{}{}
	}

	lettersIndex := 0
	runes := []rune(raw)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch r {
		case '[', '(':
			if i+2 >= len(runes) {
				return ParsedGuess{}, ErrInvalidMarkers
			}
			closer := ']'
			state := LetterStateGreen
			if r == '(' {
				closer = ')'
				state = LetterStateYellow
			}
			letter := runes[i+1]
			if runes[i+2] != closer {
				return ParsedGuess{}, ErrInvalidMarkers
			}
			if !unicode.IsLetter(letter) {
				return ParsedGuess{}, ErrGuessInvalidChar
			}
			if lettersIndex >= 5 {
				return ParsedGuess{}, ErrGuessWrongLength
			}
			lower := unicode.ToLower(letter)
			out.Letters[lettersIndex] = lower
			out.States[lettersIndex] = state
			out.MarkedSet[lower] = struct{}{}
			if state == LetterStateYellow {
				out.YellowByPos[lettersIndex][lower] = struct{}{}
			} else {
				out.GreenByPos[lettersIndex] = lower
			}
			lettersIndex++
			i += 2
		default:
			if !unicode.IsLetter(r) {
				return ParsedGuess{}, ErrGuessInvalidChar
			}
			if lettersIndex >= 5 {
				return ParsedGuess{}, ErrGuessWrongLength
			}
			out.Letters[lettersIndex] = unicode.ToLower(r)
			out.States[lettersIndex] = LetterStateNone
			lettersIndex++
		}
	}

	if lettersIndex != 5 {
		return ParsedGuess{}, ErrGuessWrongLength
	}

	return out, nil
}

// ValidateGuesses parses and validates 1-6 user-provided guess strings.
//
// On success, it returns the parsed guesses that can be fed into regex construction.
func ValidateGuesses(rawGuesses []string) ([]ParsedGuess, error) {
	if len(rawGuesses) == 0 {
		return nil, ErrNoGuesses
	}
	if len(rawGuesses) > 6 {
		return nil, ErrTooManyGuesses
	}

	parsed := make([]ParsedGuess, 0, len(rawGuesses))
	greens := [5]rune{}
	markedLetters := map[rune]struct{}{}

	for _, g := range rawGuesses {
		pg, err := ParseGuess(g)
		if err != nil {
			return nil, err
		}

		for m := range pg.MarkedSet {
			markedLetters[m] = struct{}{}
		}

		for i := 0; i < 5; i++ {
			if pg.States[i] == LetterStateGreen {
				if greens[i] != 0 && greens[i] != pg.Letters[i] {
					return nil, fmt.Errorf("%w: conflicting green letter at position %d", ErrContradiction, i+1)
				}
				greens[i] = pg.Letters[i]
			}
		}

		parsed = append(parsed, pg)
	}

	if len(markedLetters) > 5 {
		return nil, ErrTooManyMarked
	}

	return parsed, nil
}
