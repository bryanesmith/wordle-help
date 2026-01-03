package main

import (
	"os"
	"path/filepath"
	"regexp"
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

	contents := "syrup\nswathe\nswathh\n"
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

	if len(cands) != 1 || cands[0] != "swathh" {
		t.Fatalf("expected [swathh], got %v", cands)
	}
}
