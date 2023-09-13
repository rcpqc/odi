package case14

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case14_a", func() any { return &A{} })
}

type A struct {
	M1 map[string]string
	M2 map[string]string
	S1 []int
	S2 []int
	B1 *B
	B2 *B
}

type B struct {
	XX int
}
