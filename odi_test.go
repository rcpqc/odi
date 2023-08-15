package odi

import (
	"os"
	"reflect"
	"testing"
)

func init() {
	Provide("object_a", func() any { return &ObjectA{} })
	Provide("object_b", func() any { return &ObjectB{} })
	Provide("object_c", func() any { return &ObjectC{} })
}

type Interface1 interface {
	Foo() error
}

type ObjectA struct {
	Arg0   int64
	Arg1   string
	Arg2   []uint
	Ifaces []Interface1 `yaml:"ifaces"`
}

type ObjectB struct {
	XX int64
	YY string
	ZZ []uint
}

func (o *ObjectB) Foo() error {
	return nil
}

type ObjectC struct {
	C any
	d any
}

func (o *ObjectC) Foo() error {
	return nil
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
			file: "data/1.yaml",
			want: &ObjectA{
				Arg0:   123,
				Arg1:   "fafdsa",
				Arg2:   []uint{1, 2, 3},
				Ifaces: []Interface1{&ObjectB{XX: 123, YY: "aaf", ZZ: []uint{4, 5, 6}}, &ObjectC{C: "abcde"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := os.Open(tt.file)
			got, err := Resolve(f)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClone(t *testing.T) {
	tests := []struct {
		name string
		want any
	}{
		{
			name: "case_1",
			want: &ObjectA{
				Arg0:   123,
				Arg1:   "fafdsa",
				Arg2:   []uint{1, 2, 3},
				Ifaces: []Interface1{&ObjectB{XX: 123, YY: "aaf", ZZ: []uint{4, 5, 6}}, &ObjectC{C: "abcde"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clone(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}
