package main

import "testing"

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

	expected := "^pr[^aeilnpst]o[^aeilnpst]$"
	if pattern != expected {
		t.Fatalf("expected %q, got %q", expected, pattern)
	}
}
