# wordle_help
CLI tool that helps players decide their next guess while playing Wordle.

## Usage

Print help:

`wordle_help -h` (or `wordle_help --help`)

Provide one or more guesses:

`wordle_help -g "slate"`

`wordle_help -g "slate" -g "[p][r](o)ps" -g "[p][r]i[o]n"`

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

## Building

`./build.sh`

Or:

```
go build -o wordle_help main.go candidates.go regex.go validations.go
```

## Running tests

`go test ./...`

## How results are sorted

Instead of printing candidates alphabetically, `wordle_help` prints them in ascending order of `E_remaining` (the estimated number of remaining candidates after making that guess).

Let:

- `C` be the current set of remaining candidate answers (size `N`)
- `g` be a potential guess
- `f(g, a)` be the 5-tile Wordle feedback pattern (gray/yellow/green, including duplicate-letter rules) you would see if the true answer were `a` and you guessed `g`

For a fixed guess `g`, the mapping `a -> f(g, a)` partitions `C` into buckets by feedback pattern. If bucket `p` has size `n_p`, then:

```
E_remaining(g) = (1/N) * Σ_p n_p^2
E_eliminated(g) = N - E_remaining(g)
```

The tool prints each candidate guess along with these two values.

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
