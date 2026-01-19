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

## Simulating games

After building, you can simulate one or more games using `bin/wordle_sims`. The tool will run `bin/wordle_help` multiple times to simulate the game, and gather statistics about the number of guesses it takes to solve each game if it uses the suggested guess at the top of the list.

Provide exactly one starting word and optionally one or more answers:

```
./build.sh

# Simulate 50 games with slate as the starting word and random answers:
bin/wordle_sims -s slate 

# Simulate 2 games with slate as the starting word and stale and apple as the answers:
bin/wordle_sims -s slate -a stale -a apple 
```

If you omit all `-a/--answer` flags, the tool will randomly sample 50 answers from `/usr/share/dict/words`.

The tool prints the `bin/wordle_help` commands it runs (without printing their output), then prints a summary table of the game results along with the average number of guesses it took to solve each game. E.g.,:

```
| Starting Word | Answer | Total Guesses |
| ------------- | ------ | ------------- |
| slate | fifie | 4 |
| slate | smith | 3 |
| slate | blare | 6 |
...
| slate | unfur | 4 |

Average guesses: 4.07

| Guesses | Count |
| ------- | ----- |
| 1       | 0 |
| 2       | 1 |
| 3       | 8 |
| 4       | 23 |
| 5       | 9 |
| 6       | 2 |
```

## Project layout

- `wordlehelp/`
  - `main.go`: CLI entrypoint for `wordle_help`; parses flags, validates guesses, builds regex, and prints sorted candidates
- `wordlesims/`
  - `main.go`: CLI entrypoint for `wordle_sims`; runs multiple simulations by repeatedly calling `bin/wordle_help`
- `utils/`
  - `validations.go`: Guess parsing and validation logic
  - `regex.go`: Converts validated guesses into a single regular expression
  - `recommendations.go`: End-to-end helper for recommending the next guess from raw guesses (`utils.NextGuessRecommendations`)
  - `candidates.go`: Loads candidate words and scores/sorts them by expected remaining candidates
- `bin/`
  - `wordle_help`: executable that helps players decide their next guess while playing Wordle
  - `wordle_sims`: executable that simulates games of Wordle and reports aggregate results
