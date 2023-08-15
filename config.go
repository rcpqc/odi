package odi

type Config struct {
	Resolver string
	Key      string
}

type Option func(o *Config)

func NewDefaultConfig() *Config {
	return &Config{Resolver: "yaml", Key: "object"}
}

func WithKey(key string) Option {
	return func(o *Config) { o.Key = key }
}

func WithResolver(resolver string) Option {
	return func(o *Config) { o.Resolver = resolver }
}
