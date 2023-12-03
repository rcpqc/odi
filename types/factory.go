package types

import (
	"reflect"
	"sync"
)

// Factory factory of profile
type Factory struct {
	TagKey string
	Cache  sync.Map
}

// GetProfile get or create the profile of type
func (o *Factory) GetProfile(t reflect.Type) *Profile {
	if f, ok := o.Cache.Load(t); ok {
		return f.(func() any)().(*Profile)
	}
	var once sync.Once
	var res any
	f, _ := o.Cache.LoadOrStore(t, func() any {
		once.Do(func() {
			res = (&Profile{}).init(t, o.TagKey)
			o.Cache.Store(t, func() any { return res })
		})
		return res
	})
	return f.(func() any)().(*Profile)
}
