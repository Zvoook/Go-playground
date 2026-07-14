package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type fileLength struct {
	File   string
	Lenght int
}

func countStringsInFiles(filenames []string) (map[string]int, chan error) {
	cnt := len(filenames)
	lenghts := make(map[string]int)
	ansCh := make(chan fileLength)
	errCh := make(chan error)

	workersWg, managerWg := sync.WaitGroup{}, sync.WaitGroup{}

	managerWg.Add(1)
	go func() {
		defer managerWg.Done()
		for ans := range ansCh {
			lenghts[ans.File] = ans.Lenght
		}
	}()

	workersWg.Add(cnt)
	for i := range cnt {
		go func(filename string) {
			defer workersWg.Done()
			len := 0

			file, err := os.Open(filename)
			if err != nil {
				errCh <- err
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				len++
			}
			if err := scanner.Err(); err != nil {
				errCh <- err
				return
			}
			ansCh <- fileLength{filename, len}
		}(filenames[i])
	}

	workersWg.Wait()
	close(ansCh)
	close(errCh)

	managerWg.Wait()

	return lenghts, errCh
}

func calculateLinesInFiles(filesAmount int) {
	paths, err := prepareFilesForProcessing(filesAmount)
	if err != nil {
		fmt.Println(err)
		return
	}
	lenghts, errCh := countStringsInFiles(paths)
	for err := range errCh {
		fmt.Println(err)
	}

	for name, cnt := range lenghts {
		fmt.Printf("%s - %d strings\n", name, cnt)
	}
}
