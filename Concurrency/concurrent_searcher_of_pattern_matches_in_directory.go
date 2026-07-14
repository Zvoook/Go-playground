package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
)

type Pair struct {
	File  string
	Count int
}

func SearchPattern(paths []string, pattern string) ([]Pair, chan error) {
	var result []Pair
	errCh := make(chan error)
	ansCh := make(chan Pair)
	workerWg := sync.WaitGroup{}
	managerWg := sync.WaitGroup{}

	managerWg.Add(1)
	go func() {
		defer managerWg.Done()
		for ans := range ansCh {
			if ans.Count > 0 {
				result = append(result, Pair{ans.File, ans.Count})
			}
		}
	}()

	for i := range len(paths) {
		workerWg.Add(1)
		go func(path string) {
			defer workerWg.Done()
			cnt, err := countPatternInFile(path, pattern)
			if err != nil {
				errCh <- err
			} else {
				ansCh <- Pair{path, cnt}
			}
		}(paths[i])
	}
	workerWg.Wait()
	close(ansCh)
	close(errCh)
	managerWg.Wait()
	return result, errCh
}

func countPatternInFile(path string, pattern string) (int, error) {
	file, err := os.Open(path)
	res := 0
	if err != nil {
		return res, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Err() != nil {
		return res, err
	}
	for scanner.Scan() {
		line := scanner.Text()
		res += strings.Count(line, pattern)
	}
	return res, nil
}

func Search(filesCount int, pattern string) {
	paths, err := prepareFilesForProcessing(filesCount)
	if err != nil {
		fmt.Println(err)
		return
	}
	pairs, errCh := SearchPattern(paths, pattern)

	slices.SortFunc(pairs, func(a, b Pair) int {
		return b.Count - a.Count
	})

	for err := range errCh {
		fmt.Println(err)
	}

	fmt.Printf("Pattern: %s\n", pattern)

	for _, pair := range pairs {
		name := pair.File
		cnt := pair.Count
		fmt.Printf("%v\t %d matches\n", name, cnt)
	}
}
