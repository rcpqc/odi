package types

// Resolvable 可解析
type Resolvable interface{ OnResolve() error }

// Disposable 可处置
type Disposable interface{ OnDispose() error }

// Clonable 可克隆
type Clonable interface{ OnClone() any }

// Resolver 解析器
type Resolver interface {
	Resolve(data []byte) (any, error)
}
