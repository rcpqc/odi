package objects

type D struct {
	KK string `yaml:"kk"`
	HH string `yaml:"-"`
	B  B      `yaml:",inline"`
}
