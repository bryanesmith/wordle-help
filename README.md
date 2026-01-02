# wordle-help
CLI tool that helps players decide their next guess while playing Wordle.

## Usage

Print help:

`wordle-help -h` (or `wordle-help --help`)

Provide one or more guesses:

`wordle-help -g "slate"`

`wordle-help -g "slate" -g "[p][r](o)ps" -g "[p][r]i[o]n"`

### Guess encoding

Each guess is a 5-letter word. You can optionally annotate letters with:

- **Green letters**: surround the letter with `[` `]` (correct letter, correct position)
- **Yellow letters**: surround the letter with `(` `)` (correct letter, wrong position)

Examples:

- `slate` means none of `s`, `l`, `a`, `t`, `e` appear in the answer.
- `[p][r](o)ps` means:
  - `p` and `r` are green in positions 1 and 2
  - `o` is present but not in position 3
  - the second `p` is unmarked, so `p` is excluded elsewhere

## Running tests

`go test ./...`

## Project layout

- `main.go`
  - Parses flags (`-g/--guess`, `-h/--help`)
  - Validates guesses
  - Builds a regular expression
  - Loads candidate words from `/usr/share/dict/words` and prints matches
- `validations.go`
  - Guess parsing and validation logic
  - Unit tests in `validations_test.go`
- `regex.go`
  - Converts validated guesses into a single regular expression
  - Unit tests in `regex_test.go`
- `candidates.go`
  - Reads a dictionary file and returns all 5-letter candidates matching a regex
  - Unit tests in `candidates_test.go`
