package assertion

import "fmt"

type data interface {
	getName() string
}

type A struct {
	name string
}

func (a A) getName() string {
	return a.name
}

type B struct {
	name string
}

func (b B) getName() string {
	return b.name
}

func execute(a data) {
	switch data := a.(type) {
	case A:
		fmt.Println("A" + data.name)
	case B:
		fmt.Println("B" + data.name)

	}
}

func InterfaceAssert() {
	var a = A{"a"}
	execute(a)
	var b = B{"b"}
	execute(b)
}
