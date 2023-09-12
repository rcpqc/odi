package case4

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case4_a", func() any { return &A{} })
}

type Interface interface {
	Foo() error
}

type A struct {
	Ifaces []Interface `yaml:"ifaces"`
}
