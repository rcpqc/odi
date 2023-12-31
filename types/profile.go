package types

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var firstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var allCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// snake translate to snake case
func snake(s string) string {
	snake := firstCap.ReplaceAllString(s, "${1}_${2}")
	snake = allCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// Profile type's Profile
type Profile struct {
	Fields []*Field
	Names  map[any]struct{}
}

// Field
type Field struct {
	Index       []int
	Name        any
	Router      string
	InlineMap   bool
	InlineIface bool
	Error       error
}

func extractFields(index []int, router string, st reflect.StructField, tagKey string) []*Field {
	tags := strings.Split(st.Tag.Get(tagKey), ",")
	name := tags[0]
	if name == "-" {
		return nil
	}
	if name == "" {
		name = snake(st.Name)
	}
	inline := false
	for i := 1; i < len(tags); i++ {
		if tags[i] == "inline" {
			inline = true
		}
	}
	if !inline {
		return []*Field{{Index: index, Name: name, Router: router + "." + name}}
	}
	if tags[0] != "" {
		router = router + "." + name
	}
	t := st.Type
	if t.Kind() == reflect.Map && t.Key() == String {
		return []*Field{{Index: index, Name: name, Router: router, InlineMap: true}}
	}
	if t.Kind() == reflect.Interface {
		return []*Field{{Index: index, Name: name, Router: router, InlineIface: true}}
	}
	if t.Kind() == reflect.Struct {
		fields := []*Field{}
		for i := 0; i < t.NumField(); i++ {
			fields = append(fields, extractFields(append(index, i), router, t.Field(i), tagKey)...)
		}
		return fields
	}
	err := fmt.Errorf("illegal inline type(%v) expect type(struct, map[string]any, interface)", st.Type)
	return []*Field{{Index: index, Name: name, Router: router, Error: err}}
}

// init initialize profile
func (o *Profile) init(t reflect.Type, tagkey string) *Profile {
	for i := 0; i < t.NumField(); i++ {
		o.Fields = append(o.Fields, extractFields([]int{i}, "", t.Field(i), tagkey)...)
	}
	o.Names = map[any]struct{}{}
	for _, f := range o.Fields {
		if !f.InlineMap && !f.InlineIface && f.Error == nil {
			o.Names[f.Name] = struct{}{}
		}
	}
	return o
}
