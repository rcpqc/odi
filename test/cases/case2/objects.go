package case2

import (
	"github.com/rcpqc/odi"
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

func (o *E) Resolve(src any, def func() error) error {
	o.CX = 321
	if err := def(); err != nil {
		return err
	}
	return nil
}

func (o *E) Prepare() error {
	o.DFG = "[" + o.DFG + "]"
	return nil
}

func (o *E) Foo() error {
	return nil
}
