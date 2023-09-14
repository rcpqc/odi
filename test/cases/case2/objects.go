package case2

import (
	"github.com/rcpqc/odi"
	"github.com/rcpqc/odi/resolve"
)

func init() {
	odi.Provide("case2_a", func() any { return &A{} })
	odi.Provide("case2_e", func() any { return &E{} })
}

type Interface interface {
	Foo() error
}

type A struct {
	Ifaces []Interface `json:"ifaces"`
}

type E struct {
	DFG string `json:"dfg"`
	CX  int    `json:"cx"`
	FF  *struct {
		VC []int `json:"vc"`
	} `json:"ff"`
}

func (o *E) Resolve(src any) error {
	resolve.Struct(o, nil, resolve.WithObjectKey("obj"), resolve.WithTagKey("json"))
	data := []any{map[any]any{"cx": 321}}
	resolve.Struct(o, data[0], resolve.WithObjectKey("obj"), resolve.WithTagKey("json"))
	if err := resolve.Struct(o, src, resolve.WithObjectKey("obj"), resolve.WithTagKey("json")); err != nil {
		return err
	}
	o.DFG = "[" + o.DFG + "]"
	return nil
}

func (o *E) Foo() error {
	return nil
}
