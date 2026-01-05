package utils

import (
	"regexp"
	"testing"
)

func mustParse(t *testing.T, raw string) ParsedGuess {
	t.Helper()
	pg, err := ParseGuess(raw)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	return pg
}

func TestBuildRegex_Example(t *testing.T) {
	guesses := []ParsedGuess{
		mustParse(t, "slate"),
		mustParse(t, "[p][r](o)ps"),
		mustParse(t, "[p][r]i[o]n"),
	}

	pattern, err := BuildRegex(guesses)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	expected := "^pr[^aeilnost]o[^aeilnst]$"
	if pattern != expected {
		t.Fatalf("expected %q, got %q", expected, pattern)
	}
}

func TestBuildRegex_DuplicateLetterYellowAndGray(t *testing.T) {
	guess := mustParse(t, "(a)bcad")

	pattern, err := BuildRegex([]ParsedGuess{guess})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	re := regexp.MustCompile(pattern)
	if !re.MatchString("eaaee") {
		t.Fatalf("expected regex %q to match %q", pattern, "eaaee")
	}
}
