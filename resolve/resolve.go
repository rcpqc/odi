package resolve

import (
	"github.com/rcpqc/odi/resolve/yaml"
	"github.com/rcpqc/odi/types"
)

var Resolve = resolve

// NewResolver 创建resolver
func NewResolver(name string, key string) types.Resolver {
	switch name {
	case "yaml":
		return yaml.NewResolver(key)
	default:
		return nil
	}
}

// resolve 解析配置
func resolve(resolver types.Resolver, data []byte) (any, error) {
	return resolver.Resolve(data)
}
