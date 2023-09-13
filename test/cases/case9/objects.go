package case9

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case9_ints", func() any { return &Ints{} })
}

type Ints struct {
	ArrInt []int
}
