package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type dictionary map[rune]int

const chunkSize = 10

func countFreqs(str []rune) dictionary {
	freqs := make(dictionary)

	wg := &sync.WaitGroup{}
	ch := make(chan dictionary)

	num := len(str)/chunkSize + 1

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range num {
			dict := <-ch
			for sym, cnt := range dict {
				freqs[sym] += cnt
			}
		}
	}()

	for i := range num {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			dict := make(dictionary)
			left := i * chunkSize
			right := (i + 1) * chunkSize
			if right > len(str) {
				right = len(str)
			}

			substr := str[left:right]
			for _, sym := range substr {
				dict[sym]++
			}
			ch <- dict
		}(i)
	}

	wg.Wait()
	return freqs
}

func countFrequencies(cnt int) {
	runes := make([]rune, cnt)
	for i := range cnt {
		runes[i] = rune(rand.Intn(33) + 'а')
	}
	freqs := countFreqs(runes)
	for sym, cnt := range freqs {
		fmt.Printf("%c - %d\n", sym, cnt)
	}
}
