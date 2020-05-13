package reflect

import (
	"fmt"
	"reflect"
)

type Person struct {
	name string
	age  int
}

func Method(in interface{}) (ok bool) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Slice {
		ok = true
	} else {
		//panic
		return false
	}

	num := v.Len()
	for i := 0; i < num; i++ {
		fmt.Println(v.Index(i).Interface())
	}
	return ok
}

func main() {
	s := []int{1, 3, 5, 7, 9}
	b := []float64{1.2, 3.4, 5.6, 7.8}
	Method(s)
	Method(b)

	list := new([]int)
	list2 := make([]int, 5)
	list2 = append(list2, 1)
	fmt.Println(list)
	m := get()
	a := len(m)
	fmt.Println(a)
}

func get() (m []map[string]int) {
	return nil
}
