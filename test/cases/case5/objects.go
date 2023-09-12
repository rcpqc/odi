package case5

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case5_b", func() any { return &B{} })
	odi.Provide("case5_d", func() any { return &D{} })
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

type D struct {
	KK string `yaml:"kk"`
	HH string `yaml:"-"`
	B  B      `yaml:",inline"`
}
