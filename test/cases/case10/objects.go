package case10

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case10_uints", func() any { return &Uints{} })
}

type Uints struct {
	ArrUint []uint
}
