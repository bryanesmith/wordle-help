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

func Candidates(re *regexp.Regexp, dictPath string) ([]string, error) {
	f, err := os.Open(dictPath)
	if err != nil {
		return nil, ErrDictOpen
	}
	defer f.Close()

	out := []string{}
	seen := map[string]struct{}{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		w := strings.TrimSpace(scanner.Text())
		if w == "" {
			continue
		}
		w = strings.ToLower(w)
		if len([]rune(w)) != 5 {
			continue
		}
		allLetters := true
		for _, r := range w {
			if !unicode.IsLetter(r) {
				allLetters = false
				break
			}
		}
		if !allLetters {
			continue
		}
		if !re.MatchString(w) {
			continue
		}
		if _, ok := seen[w]; ok {
			continue
		}
		seen[w] = struct{}{}
		out = append(out, w)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
