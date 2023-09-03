package objects

type H struct {
	H1 struct {
		A string
		B *B `odi:",inline"`
	} `odi:",inline"`
	H2 struct {
		C int
	}
}
