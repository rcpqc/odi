package container

import "fmt"

var factory = map[string]func() any{}

// Bind bind kinds and their constructors
func Bind(kind string, fn func() any) { factory[kind] = fn }

// Make instantize a kind
func Make(kind string) (any, error) {
	constructor := factory[kind]
	if constructor == nil {
		return nil, fmt.Errorf("kind(%s) not registered", kind)
	}
	return constructor(), nil
}

// List list all kinds
func List() []string {
	kinds := make([]string, 0, len(factory))
	for k := range factory {
		kinds = append(kinds, k)
	}
	return kinds
}
