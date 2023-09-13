package case8

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case8_bools", func() any { return &Bools{} })
}

type Bools struct {
	ArrBool []bool
}
