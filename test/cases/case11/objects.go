package case11

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case11_floats", func() any { return &Floats{} })
}

type Floats struct {
	ArrFloat []float64
}
