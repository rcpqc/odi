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

func analyzeFields(t reflect.Type) []*Field {
	fields := []*Field{}
	for i := 0; i < t.NumField(); i++ {
		tags := strings.Split(t.Field(i).Tag.Get("odi"), ",")
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

func injectStruct(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() != reflect.Map || src.Type().Key() != types.String {
		return errs.Newf("expect map[string] but %v", src.Kind())
	}

	ignore := ctxIgnore(ctx)
	if ignore == nil {
		ignore = map[string]struct{}{}
	}

	fields := analyzeFields(dst.Type())
	for _, field := range fields {
		vfield := dst.Field(field.Index)
		tfield := field.Type
		if !vfield.CanSet() {
			continue
		}
		// 已注入
		if _, ok := ignore[field.Name]; ok {
			continue
		}
		if field.Inline {
			if !field.InlineTypeValid {
				return errs.Newf("illegal inline type(%v) expect struct or map[string]", tfield).Prefix("." + field.Name)
			}
			if err := inject(ctxWithIgnore(ctx, ignore), vfield, src); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
			continue
		}
		val := src.MapIndex(reflect.ValueOf(field.Name))
		if !val.IsValid() {
			continue
		}
		if err := inject(ctx, vfield, val); err != nil {
			return errs.New(err).Prefix("." + field.Name)
		}
		ignore[field.Name] = struct{}{}
	}
	return nil
}
