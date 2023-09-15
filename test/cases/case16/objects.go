package case16

import (
	"context"
	"log"

	"github.com/rcpqc/odi"
)

func init() {
	odi.Provide("case16_map_objects", func() any { return &MapObjects{} })
	odi.Provide("case16_slc_objects", func() any { return &SlcObjects{} })
	odi.Provide("case16_a", func() any { return &A{} })
	odi.Provide("case16_b", func() any { return &B{} })
	odi.Provide("case16_c", func() any { return &C{} })
}

type Object interface {
	Execute(ctx context.Context)
}

type MapObjects struct {
	Objects map[string]Object
}

type SlcObjects struct {
	Objects []Object
}

type A struct {
	S string
}

func (o *A) Execute(ctx context.Context) { log.Print(o.S) }

type B struct {
	I int
}

func (o *B) Execute(ctx context.Context) { log.Print(o.I) }

type C struct {
	F float64
}

func (o *C) Execute(ctx context.Context) { log.Print(o.F) }
