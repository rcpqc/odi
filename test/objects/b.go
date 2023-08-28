package objects

type B struct {
	XX int64
	YY string
	ZZ []uint
	WW [2]float32
}

func (o *B) Foo() error {
	return nil
}
