package method

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"testing"
)

func test(p *int) {
	go func() {
		println(p)
	}()
}
func main() {
	x := 100
	p := &x
	test(p)
}

func getFileName() {
	_, fullFilename, _, _ := runtime.Caller(0)
	fmt.Println(fullFilename)
	f, err := os.Open(fullFilename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if bs, err := ioutil.ReadAll(f); err == nil {
		fmt.Println(string(bs))
	}
}

var m sync.Mutex

func call() {
	m.Lock()
	m.Unlock()
}
func deferCall() {
	m.Lock()
	defer m.Unlock()
}
func BenchmarkCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		call()
	}
}
func BenchmarkDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		deferCall()
	}
}

func TestParams(t *testing.T) {
	main()
}
