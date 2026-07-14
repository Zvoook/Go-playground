package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const goroutinsNum = 10

func sumArray(arr []int) int {
	ch := make(chan int)
	wg := new(sync.WaitGroup)
	sum := 0

	wg.Add(1)
	go func(sum *int) {
		defer wg.Done()
		for range goroutinsNum {
			x := <-ch
			*sum += x
		}
	}(&sum)

	for i := range goroutinsNum {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			sum := 0
			amount := len(arr) / goroutinsNum
			start := i * amount

			for id := range (amount) + start {
				sum += arr[id]
			}

			ch <- sum
		}(i)
	}

	wg.Wait()
	return sum
}

func sumThousandNumbers(cnt int) {
	arr := make([]int, cnt)
	for i := range cnt {
		arr[i] = rand.Intn(33)
	}
	for id, num := range arr {
		fmt.Println(id, " : ", num)
	}
	sum := sumArray(arr)
	fmt.Println("\nResult: ", sum)
}
