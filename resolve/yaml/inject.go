package yaml

import (
	"fmt"
	"reflect"

	"github.com/rcpqc/odi/types"
	"gopkg.in/yaml.v3"
)

type Unmarshaler interface{ UnmarshalYAML(value *yaml.Node) error }

var UnmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

var injectors [types.MaxKinds]func(o *Resolver, rv reflect.Value, node *yaml.Node) error

func init() {
	injectors[reflect.Bool] = (*Resolver).injectBasic
	injectors[reflect.Int] = (*Resolver).injectBasic
	injectors[reflect.Int8] = (*Resolver).injectBasic
	injectors[reflect.Int16] = (*Resolver).injectBasic
	injectors[reflect.Int32] = (*Resolver).injectBasic
	injectors[reflect.Int64] = (*Resolver).injectBasic
	injectors[reflect.Uint] = (*Resolver).injectBasic
	injectors[reflect.Uint8] = (*Resolver).injectBasic
	injectors[reflect.Uint16] = (*Resolver).injectBasic
	injectors[reflect.Uint32] = (*Resolver).injectBasic
	injectors[reflect.Uint64] = (*Resolver).injectBasic
	injectors[reflect.Uintptr] = (*Resolver).injectBasic
	injectors[reflect.Float32] = (*Resolver).injectBasic
	injectors[reflect.Float64] = (*Resolver).injectBasic
	injectors[reflect.Complex64] = (*Resolver).injectBasic
	injectors[reflect.Complex128] = (*Resolver).injectBasic
	injectors[reflect.Array] = (*Resolver).injectSlice
	injectors[reflect.Interface] = (*Resolver).injectInterface
	injectors[reflect.Map] = (*Resolver).injectMap
	injectors[reflect.Pointer] = (*Resolver).injectPointer
	injectors[reflect.Slice] = (*Resolver).injectSlice
	injectors[reflect.String] = (*Resolver).injectBasic
	injectors[reflect.Struct] = (*Resolver).injectStruct
}

func (o *Resolver) inject(rv reflect.Value, node *yaml.Node) error {
	injector := injectors[rv.Kind()]
	if injector == nil {
		return fmt.Errorf("not support kind: %v", rv.Kind())
	}
	return injector(o, rv, node)
}

func (o *Resolver) injectPointer(rv reflect.Value, node *yaml.Node) error {
	if rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	return o.inject(rv.Elem(), node)
}

func (o *Resolver) injectSlice(rv reflect.Value, node *yaml.Node) error {
	nodes := []yaml.Node{}
	if err := node.Decode(&nodes); err != nil {
		return err
	}
	rv.Set(reflect.MakeSlice(rv.Type(), len(nodes), len(nodes)))
	for i, node := range nodes {
		elem := reflect.New(rv.Type().Elem()).Elem()
		if err := o.inject(elem, &node); err != nil {
			return err
		}
		rv.Index(i).Set(elem)
	}
	return nil
}

func (o *Resolver) injectMap(rv reflect.Value, node *yaml.Node) error {
	nodes := map[string]yaml.Node{}
	if err := node.Decode(&nodes); err != nil {
		return err
	}
	rv.Set(reflect.MakeMap(rv.Type()))
	for key, node := range nodes {
		elem := reflect.New(rv.Type().Elem()).Elem()
		if err := o.inject(elem, &node); err != nil {
			return err
		}
		rv.SetMapIndex(reflect.ValueOf(key), elem)
	}
	return nil
}

func (o *Resolver) injectStruct(rv reflect.Value, node *yaml.Node) error {

	if reflect.PointerTo(rv.Type()).Implements(UnmarshalerType) && rv.CanAddr() {
		return rv.Addr().Interface().(Unmarshaler).UnmarshalYAML(node)
	}

	nodes := map[string]yaml.Node{}
	if err := node.Decode(&nodes); err != nil {
		return err
	}
	for i := 0; i < rv.NumField(); i++ {
		vfield := rv.Field(i)
		tfield := rv.Type().Field(i)
		if !vfield.CanSet() {
			continue
		}
		alias, inline := o.analyzeTag(tfield)
		if alias == "-" {
			continue
		}
		if inline {
			if err := o.inject(vfield, node); err != nil {
				return err
			}
			continue
		}
		node, ok := nodes[alias]
		if !ok {
			continue
		}
		if err := o.inject(vfield, &node); err != nil {
			return err
		}
	}
	return nil
}

func (o *Resolver) injectBasic(rv reflect.Value, node *yaml.Node) error {
	elemPtr := reflect.New(rv.Type())
	if err := node.Decode(elemPtr.Interface()); err != nil {
		return err
	}
	rv.Set(elemPtr.Elem())
	return nil
}

func (o *Resolver) injectInterface(rv reflect.Value, node *yaml.Node) error {
	if rv.Type() == types.Any {
		return o.injectBasic(rv, node)
	}
	object, err := o.invoke(node)
	if err != nil {
		return err
	}
	if reflect.ValueOf(object).Type().Implements(rv.Type()) {
		rv.Set(reflect.ValueOf(object))
	}
	return nil
}
