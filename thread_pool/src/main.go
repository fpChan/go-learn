package main

import (
	"fmt"
)

func main() {
	pool := NewPool(30, 50)
	defer pool.Release()
	for i := 0; i < 10; i++ {
		job := SimpleExample{num: i}
		pool.jobQueue <- job
	}
	pool.WaitAll()
}

type SimpleExample struct{ num int }

func (s SimpleExample) handler() error {
	fmt.Printf("Simple function Test: the worker is %d .\n", s.num)
	return nil
}
