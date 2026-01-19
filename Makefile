BIN_DIR := bin

.PHONY: build test clean

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/wordle_help ./wordlehelp
	go build -o $(BIN_DIR)/wordle_sims ./wordlesims

test:
	go test ./...

clean:
	rm -f $(BIN_DIR)/wordle_help $(BIN_DIR)/wordle_sims
