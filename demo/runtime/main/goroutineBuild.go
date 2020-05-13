package main

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	count int32
	wg    sync.WaitGroup
)

func main() {
	data := []string{"test", "build"}
	dataChan := make(chan int, 1)
	for i, s := range data {
		go func(s string, i int) {
			dataChan <- i
			fmt.Printf("data %d, %s", i, s)
		}(s, i)
		<-dataChan
	}
	fmt.Println()
	wg.Add(2)
	go incCount()
	go incCount()
	wg.Wait()
	fmt.Println(count)
}
func incCount() {
	defer wg.Done()
	for i := 0; i < 2; i++ {
		value := count
		runtime.Gosched()
		value++
		count = value
	}
}

// var chan 定./义之后只是空
