package container

import "fmt"

var factory = map[string]func() any{}

// Register 注册
func Provide(kind string, fn func() any) { factory[kind] = fn }

func Create(kind string) (any, error) {
	constructor := factory[kind]
	if constructor == nil {
		return nil, fmt.Errorf("kind: %s not registered", kind)
	}
	return constructor(), nil
}
