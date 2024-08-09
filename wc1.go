package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	lines      int
	characters int
	words      int
}

func writeToResult(in chan<- Result, scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()
		res := Result{
			characters: len(line),
			lines:      1,
			words:      strings.Count(line, " "),
		}
		in <- res
	}
	close(in)
}

func main() {
	if len(os.Args) != 2 {
		println("usage: wc1 <filename>")
		os.Exit(1)
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	ch := make(chan Result, 10)

	go writeToResult(ch, bufio.NewScanner(file))

	result := Result{}

	for lineResult := range ch {
		result.lines += lineResult.lines
		result.characters += lineResult.characters
		result.words += lineResult.words
	}

	fmt.Printf("%d lines, %d characters, %d words\n", result.lines, result.characters, result.words)
}
