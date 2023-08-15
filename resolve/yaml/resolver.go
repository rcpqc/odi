package yaml

import (
	"gopkg.in/yaml.v3"
)

// Resolver YAML解析器
type Resolver struct {
	Key string
}

func NewResolver(key string) *Resolver {
	o := &Resolver{Key: key}
	return o
}

func (o *Resolver) Resolve(data []byte) (any, error) {
	node := &yaml.Node{}
	if err := yaml.Unmarshal(data, node); err != nil {
		return nil, err
	}
	return o.invoke(node)
}
