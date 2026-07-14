package main

import (
	"fmt"
	"sync"
)

func logInConsole() {
	ch := make(chan int, 1)
	wg := sync.WaitGroup{}
	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-ch
			for j := range 100 {
				fmt.Printf("Worker %d: message %d\n", i, j)
			}
			ch <- 1
		}(i)
	}
	ch <- 1
	wg.Wait()
}
