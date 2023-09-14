package case3

import (
	"github.com/rcpqc/odi"
)

func init() {
	odi.Provide("case3_a", func() any { return &A{} })
}

type A struct {
	Arg0 int64
}
