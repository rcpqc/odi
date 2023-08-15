package di

import (
	"fmt"
	"io"
	"reflect"

	"github.com/rcpqc/odi/clone"
	"github.com/rcpqc/odi/container"
	"github.com/rcpqc/odi/dispose"
	"github.com/rcpqc/odi/resolve"
)

// Provide 注册
func Provide(kind string, constructor func() any) {
	container.Provide(kind, constructor)
}

// Resolve 解析
func Resolve(r io.Reader, opts ...Option) (any, error) {
	cfg := NewDefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	resolver := resolve.NewResolver(cfg.Resolver, cfg.Key)
	if r == nil {
		return nil, fmt.Errorf("unknown resolver=%s", cfg.Key)
	}
	bytes, _ := io.ReadAll(r)
	return resolve.Resolve(resolver, bytes)
}

// Dispose 释放
func Dispose(object any) {
	dispose.Dispose(reflect.ValueOf(object))
}

// Clone 克隆
func Clone(object any) any {
	return clone.Clone(reflect.ValueOf(object)).Interface()
}
