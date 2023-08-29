package odi

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/rcpqc/odi/resolve"
	"github.com/rcpqc/odi/test/objects"
	"gopkg.in/yaml.v3"
)

func init() {
	Provide("object_a", func() any { return &objects.A{} })
	Provide("object_b", func() any { return &objects.B{} })
	Provide("object_c", func() any { return &objects.C{} })
	Provide("object_d", func() any { return &objects.D{} })
}

func TestResolve(t *testing.T) {
	tests := []struct {
		name string
		file string
		want any
		err  error
	}{
		{
			name: "case_1",
			file: "../test/cases/1.yaml",
			want: &objects.A{
				Other:   map[string]interface{}{"object": "object_a", "arg9": 6, "arg10": "t3"},
				Arg0:    123,
				Arg1:    "fafdsa",
				Arg2:    []uint{1, 2, 3},
				ObjectD: objects.D{KK: "kk123"},
				Ifaces: []objects.Interface1{
					&objects.B{XX: 123, YY: "aaf", ZZ: []uint{4, 5, 6}, WW: [2]float32{1.1, 4.3}},
					&objects.C{C: "abcde", E: map[string]int{"a": 3}, F: map[bool]string{true: "T"}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatal(err)
			}
			var data any
			if err := yaml.Unmarshal(bytes, &data); err != nil {
				t.Fatal(err)
			}
			got, err := Resolve(data, resolve.WithTagKey("yaml"))
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolve() \ngot:  %v\nwant: %v", got, tt.want)
			}
		})
	}
}

func TestDispose(t *testing.T) {
	tests := []struct {
		name string
		file string
		err  error
	}{
		{
			name: "case_1",
			file: "../test/cases/1.yaml",
			err:  fmt.Errorf("ObjectC Dispose"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatal(err)
			}
			var data any
			if err := yaml.Unmarshal(bytes, &data); err != nil {
				t.Fatal(err)
			}
			obj, err := Resolve(data)
			if err != nil {
				t.Fatal(err)
			}
			log.Print(obj)
			err = Dispose(obj)
			if (err == nil && tt.err != nil) ||
				(err != nil && tt.err == nil) ||
				(err != nil && tt.err != nil && err.Error() != tt.err.Error()) {
				t.Errorf("Dispose() \ngot:  %v\nwant: %v", err, tt.err)
			}
		})
	}
}
