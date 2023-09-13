package case13

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case13_a", func() any { return &A{} })
	odi.Provide("case13_b", func() any { return &B{} })
	odi.Provide("case13_c", func() any { return &C{} })
}

type A struct {
	Kind string
}

type B struct {
	Channel chan int
}

type Iter interface {
	Foo(s string) error
}

type C struct {
	X [3]string
	Y []float32
	Z map[string]int
	I Iter
}
