package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func init() {
	numCPUs := runtime.NumCPU()
	fmt.Printf("the num of CPU in this machine is %d .\n", numCPUs)
	runtime.GOMAXPROCS(numCPUs)
}

// 目前机器CPU 为 8 核, 理论上就可以同时支持 8 个Machine 和 Process，目前最多能同时支持 8128 个Goruntinue
// 测试输入为 10 万个 Job，计算 hash，但是两种测试结果差别并不明显
// 原生 Goruntinue 就需要 10  万个 Goruntinue
// Goruntinue pool 如果是 100 个 worker, 大约需要每个 Goruntinue 调度执行一千次
// BenchmarkPoolGoruntinueFileIO-8               1        1611840458 ns/op        54192496 B/op    1200483 allocs/op
//BenchmarkNativeGoruntinueFileIO-8              1        1589631208 ns/op        71328408 B/op    1181652 allocs/op

const JobAmount = 100000

type calcHash struct {
	path string
	data string
	lock sync.Mutex
}

func (c calcHash) handler() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	f, err := os.OpenFile(c.path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Fatal("an error occurred by", err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s hash calc result :  %x\n", c.data, sha256.Sum256([]byte(c.data))))
	return nil
}

// 使用协程池
func BenchmarkPoolGoruntinueFileIO(b *testing.B) {
	pool := NewPool(20, 20)
	defer pool.Release()

	for n := 0; n < JobAmount; n++ {
		job := calcHash{path: "./poolhash", data: fmt.Sprintf("index_%d", n)}
		pool.jobQueue <- job
	}
	pool.WaitAll()
}

// 使用原生 Goruntinue
func BenchmarkNativeGoruntinueFileIO(b *testing.B) {
	var wg = sync.WaitGroup{}
	for n := 0; n < JobAmount; n++ {
		job := calcHash{path: "./nativehash", data: fmt.Sprintf("index_%d", n)}
		wg.Add(1)
		go func() {
			_ = job.handler()
			wg.Done()
		}()
	}
	wg.Wait()
}

// 累加 一亿 次，这个对比就比较明显了
//BenchmarkPoolGoroutineAtomicAdd-8              1        3838245709 ns/op           65088 B/op        545 allocs/op
//BenchmarkNativeGoroutineAtomicAdd-8            1        4061791417 ns/op        16280488 B/op      38729 allocs/op
var sum int64

type atomicAdd struct{}

func (d atomicAdd) handler() error {
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
	return nil
}

const AddTimes = 1000000

// 使用协程池
func BenchmarkPoolGoroutineAtomicAdd(b *testing.B) {
	pool := NewPool(20, 20)
	defer pool.Release()
	demoTask := atomicAdd{}

	for i := 0; i < AddTimes; i++ {
		pool.jobQueue <- demoTask
	}
	pool.WaitAll()
	fmt.Printf("native sum is %d\n", sum)
}

// 原生 goroutine
func BenchmarkNativeGoroutineAtomicAdd(b *testing.B) {
	var wg = sync.WaitGroup{}
	demoTask := atomicAdd{}
	for i := 0; i < AddTimes; i++ {
		wg.Add(1)
		go func() {
			demoTask.handler()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("native sum is %d\n", sum)
}
