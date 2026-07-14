package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type User string

var db = map[int]User{
	1: "Petya",
	2: "Alex",
	3: "Vasya",
}

func parseKey(key string) (int, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return 0, errors.New("Invalid key")
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("Invalid key")
	}
	return id, nil
}

func GetUser(id int) (User, error) {
	fmt.Println("Go in database...")
	time.Sleep(2 * time.Second)
	// id, err := parseKey(key)
	// if err != nil {
	// 	return "", err
	// }
	user := db[id]
	if user == "" {
		return "", errors.New("Not found")
	}
	return user, nil
}

type Call struct {
	wg     sync.WaitGroup
	result any
	err    error
}

func NewCall(result any, err error) *Call {
	return &Call{sync.WaitGroup{}, result, err}
}

type Group struct {
	requests map[any]*Call
	mtx      sync.Mutex
}

func NewGroup() *Group {
	return &Group{make(map[any]*Call), sync.Mutex{}}
}

func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mtx.Lock()
	if call, ok := g.requests[key]; ok {
		g.mtx.Unlock()
		call.wg.Wait()
		return call.result, call.err
	}

	call := NewCall(nil, nil)
	call.wg.Add(1)
	g.requests[key] = call
	g.mtx.Unlock()

	call.result, call.err = fn()
	call.wg.Done()

	g.mtx.Lock()
	delete(g.requests, key)
	g.mtx.Unlock()
	return call.result, call.err
}

func RunGroup(count int) {
	g := NewGroup()
	wg := sync.WaitGroup{}

	for range count {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := g.Do("user:1", func() (any, error) {
				return GetUser(1)
			})
			fmt.Println(res, err)
		}()
	}
	wg.Wait()
}
