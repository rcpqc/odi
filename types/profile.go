package types

import (
	"reflect"
	"regexp"
	"strings"
)

// Profile type's Profile
type Profile struct {
	Fields []*Field
}

// Field
type Field struct {
	Index  int
	Type   reflect.Type
	Name   string
	Inline bool
}

// NewProfile construct type's profile
func NewProfile(t reflect.Type, tagKey string) *Profile {
	val, _ := LoadOrCreate(t, tagKey, func() interface{} {
		return (&Profile{}).init(t, tagKey)
	})
	return val.(*Profile)
}

var firstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var allCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// snake translate to snake case
func snake(s string) string {
	snake := firstCap.ReplaceAllString(s, "${1}_${2}")
	snake = allCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// init initialize profile
func (o *Profile) init(t reflect.Type, tagKey string) *Profile {
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
	o.Fields = orders
	return o
}

func (o *Profile) GetFields() []*Field {
	return o.Fields
}
