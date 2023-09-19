package case17

import (
	"fmt"

	"github.com/rcpqc/odi"
)

func init() {
	odi.Provide("case17_a", func() any { return &A{} })
	odi.Provide("case17_b", func() any { return &B{} })
	odi.Provide("case17_c", func() any { return &C{} })
	odi.Provide("case17_d", func() any { return &D{} })
	odi.Provide("case17_e", func() any { return &E{} })
	odi.Provide("case17_f", func() any { return &F{} })
	odi.Provide("case17_g", func() any { return &G{} })

}

type Node interface{}

type A struct {
	Arr  [3]Node
	Slc  []Node
	Map  map[any]Node
	A    *A
	Func func() error
	xxx  int
}

type B struct{}

func (o *B) Iter() (string, error) { return "B", nil }

type C struct{}

func (o *C) Iter() (string, error) { return "C", nil }

type D struct{}

func (o *D) Iter() (string, error) { return "D", nil }

type E struct{}

func (o *E) Iter() (string, error) { return "E", nil }

type F struct{}

func (o *F) Iter() (string, error) { return "F", nil }

type G struct {
	Error string
}

func (o *G) Iter() (string, error) { return "", fmt.Errorf(o.Error) }
