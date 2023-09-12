package case8

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case8_j", func() any { return &J{} })
}

type J struct {
	ArrBool   []bool
	ArrInt    []int
	ArrUint   []uint
	ArrFloat  []float64
	ArrString []string
}
