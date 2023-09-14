package case7

import (
	"github.com/rcpqc/odi"
)

func init() {
	odi.Provide("case7_b", func() any { return &B{} })
	odi.Provide("case7_h", func() any { return &H{} })
}

type B struct {
	XX int64
	YY string
	ZZ []uint
	WW [2]float32
}

func (o *B) Foo() error {
	return nil
}

type H struct {
	H1 struct {
		A string
		B *B `odi:",inline"`
	} `odi:",inline"`
	H2 struct {
		C int
	}
}
