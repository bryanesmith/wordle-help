package main

import "testing"

func TestValidateGuesses_NoGuesses(t *testing.T) {
	_, err := ValidateGuesses(nil)
	if err != ErrNoGuesses {
		t.Fatalf("expected %v, got %v", ErrNoGuesses, err)
	}
}

func TestValidateGuesses_TooManyGuesses(t *testing.T) {
	_, err := ValidateGuesses([]string{"a", "b", "c", "d", "e", "f", "g"})
	if err != ErrTooManyGuesses {
		t.Fatalf("expected %v, got %v", ErrTooManyGuesses, err)
	}
}

func TestValidateGuesses_WrongLength(t *testing.T) {
	_, err := ValidateGuesses([]string{"abcd"})
	if err != ErrGuessWrongLength {
		t.Fatalf("expected %v, got %v", ErrGuessWrongLength, err)
	}
}

func TestValidateGuesses_InvalidChar(t *testing.T) {
	_, err := ValidateGuesses([]string{"ab#de"})
	if err != ErrGuessInvalidChar {
		t.Fatalf("expected %v, got %v", ErrGuessInvalidChar, err)
	}
}

func TestValidateGuesses_InvalidMarkers_NotSingleLetter(t *testing.T) {
	_, err := ValidateGuesses([]string{"[ab]cd"})
	if err != ErrInvalidMarkers {
		t.Fatalf("expected %v, got %v", ErrInvalidMarkers, err)
	}
}

func TestValidateGuesses_ContradictingGreens(t *testing.T) {
	_, err := ValidateGuesses([]string{"[a]bcde", "[b]cdef"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestValidateGuesses_TooManyMarkedLetters(t *testing.T) {
	_, err := ValidateGuesses([]string{"[a](b)[c](d)[e]", "(f)ghij"})
	if err != ErrTooManyMarked {
		t.Fatalf("expected %v, got %v", ErrTooManyMarked, err)
	}
}
