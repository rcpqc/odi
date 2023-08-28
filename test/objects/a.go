package objects

type Interface1 interface {
	Foo() error
}

type A struct {
	Other   map[string]interface{} `yaml:",inline" odi:",inline"`
	Arg0    int64
	Arg1    string
	Arg2    []uint
	ObjectD D            `yaml:",inline" odi:",inline"`
	Ifaces  []Interface1 `yaml:"ifaces"`
}
