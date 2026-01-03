#!/bin/sh
mkdir -p bin
if [ -f wordle_help ] && [ ! -f bin/wordle_help ]; then
	mv wordle_help bin/wordle_help
fi
go build -o bin/wordle_helper ./cli
