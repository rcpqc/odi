package case7

import (
	"github.com/rcpqc/odi"
)

func init() {
	odi.Provide("case7_b", func() any { return &B{} })
	odi.Provide("case7_h", func() any { return &H{} })
	odi.Provide("case7_e", func() any { return &E{} })
	odi.Provide("case7_f", func() any { return &F{} })
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
		B *B `odi:"b,inline"`
	} `odi:"h1,inline"`
	H2 struct {
		C int
	}
}

type Iface interface {
	Foo() error
}

type D struct {
	AA string
	YY string
	NN Iface `odi:",inline"`
}

type E struct {
	Ds []D
}

type F struct {
	Kind string         `odi:"object"`
	LL   map[string]int `odi:"ll,inline"`
}
