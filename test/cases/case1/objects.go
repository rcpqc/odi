package case1

import (
	"fmt"

	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case1_a", func() any { return &A{} })
	odi.Provide("case1_b", func() any { return &B{} })
	odi.Provide("case1_c", func() any { return &C{} })
	odi.Provide("case1_d", func() any { return &D{} })
}

type Interface interface {
	Foo() error
}

type A struct {
	Other   map[string]interface{} `yaml:",inline"`
	Arg0    int64
	Arg1    string
	Arg2    []uint
	Arg3    float32     `yaml:"arg3"`
	ObjectD D           `yaml:",inline"`
	Ifaces  []Interface `yaml:"ifaces"`
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

type C struct {
	C any
	d any
	E map[string]int `odi:",inline"`
	F map[bool]string
}

func (o *C) Foo() error {
	return nil
}

func (o *C) Dispose() error {
	return fmt.Errorf("ObjectC Dispose")
}

type D struct {
	KK string `yaml:"kk"`
	HH string `yaml:"-"`
	B  B      `yaml:",inline"`
}
