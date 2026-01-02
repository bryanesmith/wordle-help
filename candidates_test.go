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
	cands, err := Candidates(re, dict)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if len(cands) != 1 || cands[0] != "prboy" {
		t.Fatalf("expected [prboy], got %v", cands)
	}
}

func TestCandidates_DictOpenError(t *testing.T) {
	re := regexp.MustCompile("^abcde$")
	_, err := Candidates(re, "/path/does/not/exist")
	if err != ErrDictOpen {
		t.Fatalf("expected %v, got %v", ErrDictOpen, err)
	}
}
