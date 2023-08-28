package objects

import (
	"fmt"
)

type C struct {
	C any
	d any
	E map[string]int
	F map[bool]string
}

func (o *C) Foo() error {
	return nil
}

func (o *C) Dispose() error {
	return fmt.Errorf("ObjectC Dispose")
}
