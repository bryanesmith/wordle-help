package recommendations

import (
	"regexp"

	sdkutils "github.com/bryanesmith/wordle-help/go-sdk/utils"
)

const DefaultDictionaryPath = "/usr/share/dict/words"

type RatedGuess = sdkutils.RatedGuess

func NextGuessRecommendations(rawGuesses []string, dictPath string) ([]RatedGuess, error) {
	parsed, err := sdkutils.ValidateGuesses(rawGuesses)
	if err != nil {
		return nil, err
	}

	pattern, err := sdkutils.BuildRegex(parsed)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	cands, err := sdkutils.Candidates(re, parsed, dictPath)
	if err != nil {
		return nil, err
	}

	rated := sdkutils.SortCandidates(cands)
	return rated, nil
}
