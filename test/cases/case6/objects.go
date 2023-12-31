package case6

import (
	"github.com/rcpqc/odi"
)

func init() {
	odi.Provide("case6_g", func() any { return &G{} })
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

func (o *C) Dispose() {}

type G struct {
	MC    map[string]*C
	True  int
	False float32
}
