package case18

import (
	"github.com/rcpqc/odi"
)

func init() {
	odi.Provide("case18_a", func() any { return &A{} })
	odi.Provide("case18_b", func() any { return &B{} })
	odi.Provide("case18_c", func() any { return &C{} })
}

type Node interface{}

type A struct {
	X int `odi:",required"`
}

type B struct {
	N Node `odi:"iter,required"`
}

type C struct {
	IA A `odi:",inline"`
}
