package yaml

import (
	"fmt"
	"reflect"

	"github.com/rcpqc/odi/types"
	"gopkg.in/yaml.v3"
)

type Unmarshaler interface{ UnmarshalYAML(value *yaml.Node) error }

var UnmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

var injectors [types.MaxKinds]func(o *Resolver, rv reflect.Value, node *yaml.Node, injecteds map[string]struct{}) error

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
	injectors[reflect.Array] = (*Resolver).injectArray
	injectors[reflect.Interface] = (*Resolver).injectInterface
	injectors[reflect.Map] = (*Resolver).injectMap
	injectors[reflect.Pointer] = (*Resolver).injectPointer
	injectors[reflect.Slice] = (*Resolver).injectSlice
	injectors[reflect.String] = (*Resolver).injectBasic
	injectors[reflect.Struct] = (*Resolver).injectStruct
}

func (o *Resolver) inject(rv reflect.Value, node *yaml.Node, injecteds map[string]struct{}) error {
	if reflect.PointerTo(rv.Type()).Implements(UnmarshalerType) && rv.CanAddr() {
		return rv.Addr().Interface().(Unmarshaler).UnmarshalYAML(node)
	}
	injector := injectors[rv.Kind()]
	if injector == nil {
		return fmt.Errorf("not support kind: %v", rv.Kind())
	}
	return injector(o, rv, node, injecteds)
}

func (o *Resolver) injectPointer(rv reflect.Value, node *yaml.Node, injecteds map[string]struct{}) error {
	if rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	return o.inject(rv.Elem(), node, nil)
}

func (o *Resolver) injectArray(rv reflect.Value, node *yaml.Node, injecteds map[string]struct{}) error {
	nodes := []yaml.Node{}
	if err := node.Decode(&nodes); err != nil {
		return err
	}
	rv.Set(reflect.New(rv.Type()).Elem())
	for i, node := range nodes {
		elem := reflect.New(rv.Type().Elem()).Elem()
		if err := o.inject(elem, &node, nil); err != nil {
			return err
		}
		rv.Index(i).Set(elem)
	}
	return nil
}

func (o *Resolver) injectSlice(rv reflect.Value, node *yaml.Node, injecteds map[string]struct{}) error {
	nodes := []yaml.Node{}
	if err := node.Decode(&nodes); err != nil {
		return err
	}
	rv.Set(reflect.MakeSlice(rv.Type(), len(nodes), len(nodes)))
	for i, node := range nodes {
		elem := reflect.New(rv.Type().Elem()).Elem()
		if err := o.inject(elem, &node, nil); err != nil {
			return err
		}
		rv.Index(i).Set(elem)
	}
	return nil
}

func (o *Resolver) injectMap(rv reflect.Value, node *yaml.Node, injecteds map[string]struct{}) error {
	nodes := reflect.MakeMapWithSize(reflect.MapOf(rv.Type().Key(), reflect.TypeOf(yaml.Node{})), rv.Len())
	if err := node.Decode(nodes.Interface()); err != nil {
		return err
	}
	rv.Set(reflect.MakeMap(rv.Type()))
	iter := nodes.MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value().Interface().(yaml.Node)
		if injecteds != nil && key.Type() == types.String {
			if _, ok := injecteds[key.String()]; ok {
				continue
			}
		}
		elem := reflect.New(rv.Type().Elem()).Elem()
		if err := o.inject(elem, &val, nil); err != nil {
			return err
		}
		rv.SetMapIndex(key, elem)
	}
	return nil
}

func (o *Resolver) injectStruct(rv reflect.Value, node *yaml.Node, injecteds map[string]struct{}) error {
	nodes := map[string]yaml.Node{}
	if err := node.Decode(&nodes); err != nil {
		return err
	}
	if injecteds == nil {
		injecteds = map[string]struct{}{}
	}

	var stuInlines []reflect.Value
	var mapInline reflect.Value

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
		// 已注入
		if _, ok := injecteds[alias]; ok {
			continue
		}
		if inline {
			t := tfield.Type
			if t.Kind() == reflect.Pointer {
				t = t.Elem()
			}
			if t.Kind() == reflect.Map && t.Key() == types.String {
				mapInline = vfield
			} else if t.Kind() == reflect.Struct {
				stuInlines = append(stuInlines, vfield)
			} else {
				return fmt.Errorf("unsupported inline type=%v", t)
			}
			continue
		}
		node, ok := nodes[alias]
		if !ok {
			continue
		}
		if err := o.inject(vfield, &node, nil); err != nil {
			return err
		}
		injecteds[alias] = struct{}{}
	}

	for _, vfield := range stuInlines {
		if err := o.inject(vfield, node, injecteds); err != nil {
			return err
		}
	}

	if mapInline.IsValid() {
		if err := o.inject(mapInline, node, injecteds); err != nil {
			return err
		}
	}
	return nil
}

func (o *Resolver) injectBasic(rv reflect.Value, node *yaml.Node, injecteds map[string]struct{}) error {
	elemPtr := reflect.New(rv.Type())
	if err := node.Decode(elemPtr.Interface()); err != nil {
		return err
	}
	rv.Set(elemPtr.Elem())
	return nil
}

func (o *Resolver) injectInterface(rv reflect.Value, node *yaml.Node, injecteds map[string]struct{}) error {
	if rv.Type() == types.Any {
		return o.injectBasic(rv, node, nil)
	}
	object, err := o.invoke(node)
	if err != nil {
		return err
	}
	if !reflect.ValueOf(object).Type().Implements(rv.Type()) {
		return fmt.Errorf("the injected object does not implement the interface(%v)", rv.Type())
	}
	rv.Set(reflect.ValueOf(object))
	return nil
}
