package main

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"testing"
)

func TestCandidates_FiltersAndMatches(t *testing.T) {
	dir := t.TempDir()
	dict := filepath.Join(dir, "words")

	contents := "prone\nprboy\npr00n\nPRONE\nstone\nprion\n"
	if err := os.WriteFile(dict, []byte(contents), 0o644); err != nil {
		t.Fatalf("write dict: %v", err)
	}

	re := regexp.MustCompile("^pr[^aeilnpst]o[^aeilnpst]$")
	cands, err := Candidates(re, nil, dict)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if len(cands) != 1 || cands[0] != "prboy" {
		t.Fatalf("expected [prboy], got %v", cands)
	}
}

func TestCandidates_DictOpenError(t *testing.T) {
	re := regexp.MustCompile("^abcde$")
	_, err := Candidates(re, nil, "/path/does/not/exist")
	if err != ErrDictOpen {
		t.Fatalf("expected %v, got %v", ErrDictOpen, err)
	}
}

func TestCandidates_RequiresYellowLetters(t *testing.T) {
	dir := t.TempDir()
	dict := filepath.Join(dir, "words")

	contents := "syrup\nspore\nsitar\n"
	if err := os.WriteFile(dict, []byte(contents), 0o644); err != nil {
		t.Fatalf("write dict: %v", err)
	}

	guess, err := ParseGuess("[s]w(a)(t)h")
	if err != nil {
		t.Fatalf("parse guess: %v", err)
	}

	pattern, err := BuildRegex([]ParsedGuess{guess})
	if err != nil {
		t.Fatalf("build regex: %v", err)
	}
	re := regexp.MustCompile(pattern)

	cands, err := Candidates(re, []ParsedGuess{guess}, dict)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if len(cands) != 1 || cands[0] != "sitar" {
		t.Fatalf("expected [sitar], got %v", cands)
	}
}

func TestBuildYellowConstraints(t *testing.T) {
	g1, err := ParseGuess("ab(c)de")
	if err != nil {
		t.Fatalf("parse guess: %v", err)
	}
	g2, err := ParseGuess("(a)bcde")
	if err != nil {
		t.Fatalf("parse guess: %v", err)
	}

	req, forbid := buildYellowConstraints([]ParsedGuess{g1, g2})

	if _, ok := req['a']; !ok {
		t.Fatalf("expected required yellow to include 'a'")
	}
	if _, ok := req['c']; !ok {
		t.Fatalf("expected required yellow to include 'c'")
	}

	if _, ok := forbid[0]['a']; !ok {
		t.Fatalf("expected forbidden pos 1 to include 'a'")
	}
	if _, ok := forbid[2]['c']; !ok {
		t.Fatalf("expected forbidden pos 3 to include 'c'")
	}
}

func TestIsValidGuessWord(t *testing.T) {
	if isValidGuessWord("abcde") != true {
		t.Fatalf("expected abcde to be valid")
	}
	if isValidGuessWord("Abcde") != true {
		t.Fatalf("expected Abcde to be valid")
	}
	if isValidGuessWord("abcd") != false {
		t.Fatalf("expected abcd to be invalid")
	}
	if isValidGuessWord("abcde!") != false {
		t.Fatalf("expected abcde! to be invalid")
	}
}

func TestSortCandidates_MatchesManualScoring(t *testing.T) {
	candidates := []string{"cigar", "rebut", "sissy", "humph"}

	manual := make([]RatedGuess, 0, len(candidates))
	N := len(candidates)
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
		manual = append(manual, RatedGuess{Guess: g, ERemaining: eRemaining, EEliminated: eEliminated})
	}

	sort.Slice(manual, func(i, j int) bool {
		if manual[i].ERemaining == manual[j].ERemaining {
			return manual[i].Guess < manual[j].Guess
		}
		return manual[i].ERemaining < manual[j].ERemaining
	})

	got := SortCandidates(candidates)
	if len(got) != len(manual) {
		t.Fatalf("expected %d results, got %d", len(manual), len(got))
	}

	for i := range manual {
		if got[i].Guess != manual[i].Guess {
			t.Fatalf("at %d expected guess %q, got %q", i, manual[i].Guess, got[i].Guess)
		}
		if got[i].ERemaining != manual[i].ERemaining {
			t.Fatalf("at %d expected E_remaining %v, got %v", i, manual[i].ERemaining, got[i].ERemaining)
		}
		if got[i].EEliminated != manual[i].EEliminated {
			t.Fatalf("at %d expected E_eliminated %v, got %v", i, manual[i].EEliminated, got[i].EEliminated)
		}
	}
}
