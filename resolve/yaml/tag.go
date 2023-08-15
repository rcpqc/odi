package yaml

import (
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

func (o *Resolver) analyzeTag(field reflect.StructField) (alias string, inline bool) {
	fields := strings.Split(field.Tag.Get("yaml"), ",")
	alias = fields[0]
	inline = false
	for i := 1; i < len(fields); i++ {
		if fields[i] == "inline" {
			inline = true
		}
	}
	if alias == "" {
		alias = snake(field.Name)
	}
	return
}
