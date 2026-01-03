#!/bin/sh
mkdir -p bin
go build -o bin/wordle_help ./wordlehelp
go build -o bin/wordle_sims ./wordlesims
