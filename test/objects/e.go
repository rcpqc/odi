package objects

type E struct {
	DFG string `json:"dfg"`
	CX  int    `json:"cx"`
	FF  *struct {
		VC []int `json:"vc"`
	} `json:"ff"`
}

func (o *E) Resolve(src any) error {
	o.DFG = "[" + o.DFG + "]"
	return nil
}

func (o *E) Foo() error {
	return nil
}
