package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Result5 struct {
	lines      int
	characters int
	words      int
}

type Data struct {
	result Result5
	mx     sync.Mutex
}

var data = Data{}

var size = runtime.GOMAXPROCS(0)
var linesInfo = make(chan Result5, size)

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range linesInfo {
		data.mx.Lock()
		data.result.lines += line.lines
		data.result.characters += line.characters
		data.result.words += line.words
		data.mx.Unlock()
	}

}

func create(scanner *bufio.Scanner) {
	defer close(linesInfo)
	for scanner.Scan() {
		line := scanner.Text()
		res := Result5{
			characters: len(line),
			lines:      1,
			words:      strings.Count(line, " "),
		}
		linesInfo <- res
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], " <file> <workers>")
		os.Exit(1)
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	nWorkers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	scanner := bufio.NewScanner(file)
	go create(scanner)

	var wg sync.WaitGroup
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()

	fmt.Printf("%d lines, %d characters, %d words\n", data.result.lines, data.result.characters, data.result.words)
}
