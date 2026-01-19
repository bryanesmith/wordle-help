package utils

import (
	"regexp"
)

const DefaultDictionaryPath = "/usr/share/dict/words"

func NextGuessRecommendations(rawGuesses []string, dictPath string) ([]RatedGuess, error) {
	parsed, err := ValidateGuesses(rawGuesses)
	if err != nil {
		return nil, err
	}

	pattern, err := BuildRegex(parsed)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	cands, err := Candidates(re, parsed, dictPath)
	if err != nil {
		return nil, err
	}

	rated := SortCandidates(cands)
	return rated, nil
}
