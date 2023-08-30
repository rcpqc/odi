package types

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Profile type's Profile
type Profile struct {
	Fields []*Field
	Names  map[string]struct{}
}

// Field
type Field struct {
	Index     []int
	Name      string
	InlineMap bool
	Error     error
}

// GetProfile get or construct type's profile
func GetProfile(t reflect.Type, tagKey string) *Profile {
	val, _ := LoadOrCreate(t, tagKey, func() interface{} {
		return NewProfile(t, tagKey)
	})
	return val.(*Profile)
}

func extractFields(index []int, st reflect.StructField, tagKey string) []*Field {
	tags := strings.Split(st.Tag.Get(tagKey), ",")
	if tags[0] == "-" {
		return nil
	}
	if tags[0] == "" {
		tags[0] = snake(st.Name)
	}
	inline := false
	for i := 1; i < len(tags); i++ {
		if tags[i] == "inline" {
			inline = true
		}
	}
	if !inline {
		return []*Field{{Index: index, Name: tags[0]}}
	}
	t := st.Type
	if t.Kind() == reflect.Map && t.Key() == String {
		return []*Field{{Index: index, Name: tags[0], InlineMap: true}}
	}
	if t.Kind() == reflect.Struct {
		fields := []*Field{}
		for i := 0; i < t.NumField(); i++ {
			fields = append(fields, extractFields(append(index, i), t.Field(i), tagKey)...)
		}
		return fields
	}
	err := fmt.Errorf("illegal inline type(%v) expect struct/map[string]", st.Type)
	return []*Field{{Index: index, Name: tags[0], Error: err}}
}

// NewProfile construct type's profile
func NewProfile(t reflect.Type, tagKey string) *Profile {
	o := &Profile{}
	for i := 0; i < t.NumField(); i++ {
		o.Fields = append(o.Fields, extractFields([]int{i}, t.Field(i), tagKey)...)
	}
	o.Names = map[string]struct{}{}
	for _, f := range o.Fields {
		if !f.InlineMap && f.Error == nil {
			o.Names[f.Name] = struct{}{}
		}
	}
	return o
}

var firstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var allCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// snake translate to snake case
func snake(s string) string {
	snake := firstCap.ReplaceAllString(s, "${1}_${2}")
	snake = allCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
