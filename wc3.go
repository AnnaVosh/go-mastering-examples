package main

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"os"
	"runtime"
	"strings"
)

var (
	workers = runtime.GOMAXPROCS(0)
	sem     = semaphore.NewWeighted(int64(workers))
)

type Result3 struct {
	lines      int
	characters int
	words      int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: wc3 <file>")
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	outputChan := make(chan Result3)

	ctx := context.Background()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = sem.Acquire(ctx, 1)
		if err != nil {
			fmt.Println(err)
			break
		}

		go func(line string) {
			defer sem.Release(1)

			outputChan <- Result3{
				lines:      1,
				characters: len(line),
				words:      strings.Count(line, " "),
			}

		}(scanner.Text())
	}

	err = sem.Acquire(ctx, int64(workers))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	close(outputChan)

	result := Result3{}

	for lineResult := range outputChan {
		result.lines += lineResult.lines
		result.characters += lineResult.characters
		result.words += lineResult.words
	}

	fmt.Printf("%d lines, %d characters, %d words\n", result.lines, result.characters, result.words)
}
