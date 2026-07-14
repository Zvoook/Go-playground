package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var timeoutErr = errors.New("Time is out!")

func play(n int, outTime time.Duration) (string, error) {
	var winner string
	var err error

	ch := make(chan int)
	chWin := make(chan string)
	chErr := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), outTime)
	defer cancel()
	go func() {
		select {
		case num := <-ch:
			chWin <- fmt.Sprintf("Worker %d", num)

		case <-ctx.Done():
			chErr <- timeoutErr
		}
	}()
	for i := range n {
		go func(i int, ctx context.Context) {
			time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}(i, ctx)
	}
	select {
	case winner = <-chWin:
	case err = <-chErr:
	}
	return winner, err
}

func startRace(num int, outTime time.Duration) {
	now := time.Now()
	winner, err := play(num, outTime)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Winner: %s\n", winner)
	}
	duration := time.Since(now)
	fmt.Printf("Total time: %s", duration)
}
