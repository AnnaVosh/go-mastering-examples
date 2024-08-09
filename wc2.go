package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Result2 struct {
	lines      int
	characters int
	words      int
	mx         sync.Mutex
}

var wg sync.WaitGroup

func write(line string, result *Result2) {
	defer wg.Done()
	result.mx.Lock()
	result.lines += 1
	result.characters += len(line)
	result.words += strings.Count(line, " ")
	result.mx.Unlock()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: wc2 <file>")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	result := Result2{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		go write(scanner.Text(), &result)
	}

	wg.Wait()
	fmt.Printf("%d lines, %d characters, %d words\n", result.lines, result.characters, result.words)
}
