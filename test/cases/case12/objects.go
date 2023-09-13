package case12

import (
	"github.com/rcpqc/odi/odi"
)

func init() {
	odi.Provide("case12_strings", func() any { return &Strings{} })
}

type Strings struct {
	ArrString []string
}
