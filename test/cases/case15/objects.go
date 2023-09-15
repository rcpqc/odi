package case15

import (
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/rcpqc/odi"
	"github.com/rcpqc/odi/resolve"
)

func init() {
	odi.Provide("case15_expr1", func() any { return &Expr1{} })
	odi.Provide("case15_expr2", func() any { return &Expr2{} })
	odi.Provide("case15_expr3", func() any { return &Expr3{} })
	odi.Provide("case15_expr4", func() any { return &Expr4{} })
	odi.Provide("case15_component", func() any { return &Component{} })
}

type Expr1 struct {
	Expr    string
	Timeout string

	e ast.Expr
}

func (o *Expr1) Resolve(src any) error {
	if err := resolve.Struct(o, src, resolve.WithObjectKey("object")); err != nil {
		if err := resolve.Object(&o.Expr, src, resolve.WithObjectKey("object")); err != nil {
			return fmt.Errorf("illegal expr(%v)", src)
		}
	}
	if o.Expr != "" {
		e, err := parser.ParseExpr(o.Expr)
		if err != nil {
			return err
		}
		o.e = e
	}
	return nil
}

type Expr2 struct {
	Expr    string
	Timeout string
}

func (o *Expr2) Resolve(src any) error {
	return resolve.Struct(*o, src)
}

type Expr3 struct {
	Expr    string
	Timeout string
}

func (o *Expr3) Resolve(src any) error {
	return resolve.Struct(new(int), src)
}

type Expr4 struct {
	Expr    string
	Timeout string
}

func (o *Expr4) Resolve(src any) error {
	return resolve.Object(&o.Expr, src, resolve.WithObjectKey("object"))
}

type Component struct {
	X Expr1
	Y Expr2
	Z Expr3
	W Expr4
}
