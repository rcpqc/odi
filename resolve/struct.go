package resolve

import (
	"context"
	"reflect"
	"regexp"
	"strings"

	"github.com/rcpqc/odi/convert"
	"github.com/rcpqc/odi/errs"
)

var firstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var allCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// snake translate to snake case
func snake(s string) string {
	snake := firstCap.ReplaceAllString(s, "${1}_${2}")
	snake = allCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

type Field struct {
	Index  int
	Type   reflect.Type
	Name   string
	Inline bool
}

func analyzeFields(tagKey string, t reflect.Type) []*Field {
	fields := []*Field{}
	for i := 0; i < t.NumField(); i++ {
		tags := strings.Split(t.Field(i).Tag.Get(tagKey), ",")
		if tags[0] == "-" {
			continue
		}
		if tags[0] == "" {
			tags[0] = snake(t.Field(i).Name)
		}
		field := &Field{Index: i, Type: t.Field(i).Type, Name: tags[0]}
		for i := 1; i < len(tags); i++ {
			if tags[i] == "inline" {
				field.Inline = true
			}
		}
		fields = append(fields, field)
	}
	orders := make([]*Field, 0, len(fields))
	for _, field := range fields {
		if !field.Inline {
			orders = append(orders, field)
		}
	}
	for _, field := range fields {
		if field.Inline && field.Type.Kind() != reflect.Map {
			orders = append(orders, field)
		}
	}
	for _, field := range fields {
		if field.Inline && field.Type.Kind() == reflect.Map {
			orders = append(orders, field)
		}
	}
	return orders
}

// convertSource convert map key to string
func convertSource(src reflect.Value) reflect.Value {
	m := map[string]interface{}{}
	for iter := src.MapRange(); iter.Next(); {
		if s, err := convert.String(iter.Key()); err == nil {
			m[s] = iter.Value().Interface()
		}
	}
	return reflect.ValueOf(m)
}

func injectStructInternal(ctx context.Context, dst, src reflect.Value) error {
	tagKey := ctxGetTagKey(ctx)
	fields := analyzeFields(tagKey, dst.Type())
	for _, field := range fields {
		vfield := dst.Field(field.Index)
		tfield := field.Type
		if !vfield.CanSet() {
			continue
		}
		// direct injection of non-inline fields
		if !field.Inline {
			key := reflect.ValueOf(field.Name)
			val := src.MapIndex(key)
			if !val.IsValid() {
				continue
			}
			if err := inject(ctx, vfield, val); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
			// delete injected fields to prevent repeated injection of inline fields
			src.SetMapIndex(key, reflect.Value{})
			continue
		}
		// inline fields are injected according to type (map[string]any/struct/*struct)
		if field.Type.Kind() == reflect.Map && field.Type.Key().Kind() == reflect.String {
			if err := inject(ctx, vfield, src); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
		} else if field.Type.Kind() == reflect.Struct {
			if err := injectStructInternal(ctx, vfield, src); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
		} else if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct {
			if vfield.IsNil() {
				vfield.Set(reflect.New(vfield.Type().Elem()))
			}
			if err := injectStructInternal(ctx, vfield.Elem(), src); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
		} else {
			return errs.Newf("illegal inline type(%v) expect struct or map[string]", tfield).Prefix("." + field.Name)
		}
		continue
	}
	return nil
}

func injectStruct(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() != reflect.Map {
		return errs.Newf("expect map but %v", src.Kind())
	}
	return injectStructInternal(ctx, dst, convertSource(src))
}
