package resolve

import (
	"context"
	"reflect"
	"regexp"
	"strings"

	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
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
	Index           int
	Type            reflect.Type
	Name            string
	Inline          bool
	InlineTypeValid bool
}

func inlineTypeCheck(t reflect.Type) bool {
	return (t.Kind() == reflect.Struct) ||
		(t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct) ||
		(t.Kind() == reflect.Map && t.Key() == types.String)
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
				field.InlineTypeValid = inlineTypeCheck(field.Type)
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

func cloneMap(src reflect.Value, ignores map[string]struct{}) reflect.Value {
	m := reflect.MakeMap(src.Type())
	for iter := src.MapRange(); iter.Next(); {
		m.SetMapIndex(iter.Key(), iter.Value())
	}
	for name := range ignores {
		m.SetMapIndex(reflect.ValueOf(name), reflect.Value{})
	}
	return m
}

func injectStructInternal(ctx context.Context, dst, src reflect.Value, ignores map[string]struct{}) error {
	if ignores == nil {
		ignores = map[string]struct{}{}
	}
	tagKey := ctxGetTagKey(ctx)
	fields := analyzeFields(tagKey, dst.Type())
	for _, field := range fields {
		if _, ok := ignores[field.Name]; ok {
			continue
		}
		vfield := dst.Field(field.Index)
		tfield := field.Type
		if !vfield.CanSet() {
			continue
		}
		if field.Inline {
			if field.Type.Kind() == reflect.Map && field.Type.Key().Kind() == reflect.String {
				if err := inject(ctx, vfield, cloneMap(src, ignores)); err != nil {
					return errs.New(err).Prefix("." + field.Name)
				}
			} else if field.Type.Kind() == reflect.Struct {
				if err := injectStructInternal(ctx, vfield, src, ignores); err != nil {
					return errs.New(err).Prefix("." + field.Name)
				}
			} else if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct {
				if err := injectStructInternal(ctx, vfield, src, ignores); err != nil {
					return errs.New(err).Prefix("." + field.Name)
				}
			} else {
				return errs.Newf("illegal inline type(%v) expect struct or map[string]", tfield).Prefix("." + field.Name)
			}
			continue
		}
		key := reflect.ValueOf(field.Name)
		val := src.MapIndex(key)
		if !val.IsValid() {
			continue
		}
		if err := inject(ctx, vfield, val); err != nil {
			return errs.New(err).Prefix("." + field.Name)
		}
		ignores[field.Name] = struct{}{}
	}
	return nil
}

func injectStruct(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() != reflect.Map {
		return errs.Newf("expect map but %v", src.Kind())
	}
	injectStructInternal(ctx, dst, src, nil)
	return nil
}
