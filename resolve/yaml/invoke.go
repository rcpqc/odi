package yaml

import (
	"fmt"
	"reflect"

	"github.com/rcpqc/odi/container"
	"github.com/rcpqc/odi/types"
	"gopkg.in/yaml.v3"
)

func (o *Resolver) invoke(node *yaml.Node) (any, error) {
	kind, err := o.classify(node)
	if err != nil {
		return nil, err
	}
	object, err := container.Create(kind)
	if err != nil {
		return nil, err
	}
	if err := o.inject(reflect.ValueOf(object), node); err != nil {
		return nil, err
	}
	if iface, ok := object.(types.Resolvable); ok {
		if err := iface.OnResolve(); err != nil {
			return nil, fmt.Errorf("resolve err(%v)", err)
		}
	}
	return object, nil
}

func (o *Resolver) classify(node *yaml.Node) (string, error) {
	st := reflect.StructOf([]reflect.StructField{{
		Name: "Key",
		Type: types.String,
		Tag:  reflect.StructTag(fmt.Sprintf(`yaml:"%s"`, o.Key)),
	}})
	base := reflect.New(st)
	if err := node.Decode(base.Interface()); err != nil {
		return "", err
	}
	return base.Elem().Field(0).String(), nil
}
