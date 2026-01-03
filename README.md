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
go build -o bin/wordle_help ./wordlehelp
go build -o bin/wordle_sims ./wordlesims
```

## Running tests

`go test ./...`

## Simulating games

After building, you can simulate one or more games using `bin/wordle_sims`.

Provide exactly one starting word and one or more answers:

```
./build.sh
bin/wordle_sims -s slate -a stale -a apple
```

The tool prints the `bin/wordle_help` commands it runs (without printing their output), then prints a summary table and the average guesses.

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

- `wordlehelp/`
  - `main.go`: CLI entrypoint for `wordle_help`; parses flags, validates guesses, builds regex, and prints sorted candidates
- `wordlesims/`
  - `main.go`: CLI entrypoint for `wordle_sims`; runs multiple simulations by repeatedly calling `bin/wordle_help`
- `utils/`
  - `validations.go`: Guess parsing and validation logic
  - `regex.go`: Converts validated guesses into a single regular expression
  - `candidates.go`: Loads candidate words and scores/sorts them by expected remaining candidates
- `bin/`
  - `wordle_help`: executable that helps players decide their next guess while playing Wordle
  - `wordle_sims`: executable that simulates games of Wordle and reports aggregate results
