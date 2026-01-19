package recommendations

import (
	"os"
	"path/filepath"
	"testing"

	sdkutils "github.com/bryanesmith/wordle-help/go-sdk/utils"
)

func TestRecommendNextGuess_EndToEnd(t *testing.T) {
	dir := t.TempDir()
	dict := filepath.Join(dir, "words")

	contents := "slate\n" +
		"stale\n" +
		"plate\n" +
		"apple\n"
	if err := os.WriteFile(dict, []byte(contents), 0o644); err != nil {
		t.Fatalf("write dict: %v", err)
	}

	rated, err := NextGuessRecommendations([]string{"[s](t)[a](l)[e]"}, dict)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(rated) != 1 {
		t.Fatalf("expected 1 recommendation, got %d (%v)", len(rated), rated)
	}
	if rated[0].Guess != "slate" {
		t.Fatalf("expected top recommendation to be %q, got %q", "slate", rated[0].Guess)
	}
}

func TestRecommendNextGuess_ErrorOnNoGuesses(t *testing.T) {
	_, err := NextGuessRecommendations(nil, "/does/not/matter")
	if err != sdkutils.ErrNoGuesses {
		t.Fatalf("expected %v, got %v", sdkutils.ErrNoGuesses, err)
	}
}
